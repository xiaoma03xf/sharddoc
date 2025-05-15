package storage

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/cockroachdb/pebble"
)

// Document defination
type Document struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Tags      []string          `json:"tags,omitempty"`
	Author    string            `json:"author,omitempty"`
	CreatedAt int64             `json:"created_at"`
	UpdatedAt int64             `json:"updated_at"`
	Version   int               `json:"version"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

type DB struct {
	index    int
	mu       sync.RWMutex
	pebbleDB *pebble.DB
	dbfile   string
	docCount atomic.Int64
	// generate distr id
	snowflake      *snowflake.Node
	collections    map[string]struct{}
	ttlMap         map[string]int64
	insertCallback func(key string)
	deleteCallback func(key string)
}

func makeDB(dirpath string, index int) (*DB, error) {
	var _err error
	db := &DB{}
	db.index = index
	db.snowflake, _err = snowflake.NewNode(int64(index))
	if _err != nil {
		return nil, _err
	}
	db.collections = make(map[string]struct{})
	db.ttlMap = make(map[string]int64)
	db.insertCallback = func(key string) {}
	db.deleteCallback = func(key string) {}
	db.pebbleDB, _err = pebble.Open(dirpath, &pebble.Options{})
	if _err != nil {
		return nil, _err
	}

	return db, nil
}

func (db *DB) Insert(coll string, doc Document) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// check collection is exist
	if _, ok := db.collections[coll]; !ok {
		db.collections[coll] = struct{}{}
	}

	// insert <collection>:<documentid>
	docid := clllcetionBindID(coll, doc.ID)
	timeunix := time.Now().UnixNano()
	doc.CreatedAt = timeunix
	doc.UpdatedAt = timeunix
	data, err := json.Marshal(doc)
	if err != nil {
		return doc.ID, ErrJsonMarshal
	}
	db.docCount.Add(1)

	return doc.ID, db.pebbleDB.Set([]byte(docid), data, nil)
}
func (db *DB) BatchInsert(coll string, docs []Document, batchCount int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// check collection is exist
	if _, ok := db.collections[coll]; !ok {
		db.collections[coll] = struct{}{}
	}

	total := len(docs)
	timeunix := time.Now().UnixNano()

	for i := 0; i < total; i += batchCount {
		end := i + batchCount
		if end > total {
			end = total
		}

		batch := db.pebbleDB.NewBatch()
		for _, doc := range docs[i:end] {
			if doc.ID == "" {
				doc.ID = db.snowflake.Generate().String()
			}

			key := clllcetionBindID(coll, doc.ID)
			doc.CreatedAt = timeunix
			doc.UpdatedAt = timeunix
			data, err := json.Marshal(doc)
			if err != nil {
				batch.Close()
				return ErrJsonMarshal
			}
			if err = batch.Set([]byte(key), data, nil); err != nil {
				batch.Close()
				return err
			}
		}
		if err := batch.Commit(pebble.Sync); err != nil {
			batch.Close()
			return err
		}
		batch.Close()
		db.docCount.Add(int64(end - i))
	}
	return nil
}

// GetDocument document
func (db *DB) GetByID(coll string, id string) (*Document, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if _, ok := db.collections[coll]; !ok {
		return nil, ErrCollNotFound
	}

	docid := clllcetionBindID(coll, id)
	val, closer, err := db.pebbleDB.Get([]byte(docid))
	if err != nil {
		return nil, nil
	}
	defer closer.Close()

	var doc Document
	err = json.Unmarshal(val, &doc)
	if err != nil {
		return nil, ErrJsonUnMarshal
	}
	return &doc, nil
}

func (db *DB) DeleteByID(coll string, id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.docCount.Add(-1)

	docid := clllcetionBindID(coll, id)
	return db.pebbleDB.Delete([]byte(docid), nil)
}

func clllcetionBindID(coll string, docID string) string {
	return fmt.Sprintf("%v:%v", coll, docID)
}

func (db *DB) GetDocCount() int64 {
	return db.docCount.Load()
}

func (db *DB) Close() {
	_ = db.pebbleDB.Close()
}

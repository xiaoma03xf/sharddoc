package storage

type StorageEngine interface {
	// base trie collection: <collection_name>:<document_id>
	Insert(coll string, doc Document) (string, error)
	// TODO
	BatchInsert(coll string, docs []Document, batchCount int) error
	GetByID(coll string, id string) (*Document, error)
	// Delete(coll string, query Query) (int, error)
	DeleteByID(coll string, id string) error
}

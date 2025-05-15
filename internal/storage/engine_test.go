package storage

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	// 初始化数据库
	db, err := makeDB("pebble-data", 1)
	if err != nil {
		panic(err)
	}

	total := 100_000
	batchSize := 1000 // 可选，用于每N条做一次提示
	collName := "articles"

	start := time.Now()

	for i := 0; i < total; i++ {
		id := db.snowflake.Generate().String()
		doc := Document{
			ID:        id,
			Title:     "Document " + strconv.Itoa(i),
			Content:   randomString(100),
			Tags:      []string{"tag1", "tag2"},
			Author:    "Author_" + strconv.Itoa(rand.Intn(1000)),
			Metadata:  map[string]string{"type": "test"},
			Version:   1,
			CreatedAt: time.Now().UnixNano(),
			UpdatedAt: time.Now().UnixNano(),
		}

		_, err := db.Insert(collName, doc)
		if err != nil {
			fmt.Println("插入失败:", err)
		}

		if i%batchSize == 0 {
			fmt.Printf("插入第 %d 条...\n", i)
		}
	}

	fmt.Printf("插入完成，总用时: %v\n", time.Since(start))
	fmt.Printf("总文档数: %d\n", db.GetDocCount())

	// 随机读取一条数据
	testID := db.snowflake.Generate().String()
	doc := Document{
		ID:      testID,
		Title:   "Test Read",
		Content: "This is a test read",
	}
	db.Insert(collName, doc)

	readDoc, err := db.GetByID(collName, testID)
	if err != nil {
		fmt.Println("读取失败:", err)
	} else {
		fmt.Printf("读取成功: %+v\n", readDoc)
	}

	// 删除测试数据
	err = db.DeleteByID(collName, testID)
	if err != nil {
		fmt.Println("删除失败:", err)
	} else {
		fmt.Println("删除成功")
	}
}

func TestBatchInsert(t *testing.T) {
	db, err := makeDB("pebble_data", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	const total = 5000000
	const batchSize = 5000
	const coll = "test_coll"

	// 构造数据
	docs := make([]Document, total)
	docIDs := make([]string, total)
	for i := 0; i < total; i++ {
		id := fmt.Sprintf("docid_%d", i)
		docs[i] = Document{
			ID:      id,
			Title:   fmt.Sprintf("Title %d", i),
			Content: "some content",
			Author:  "tester",
		}
		docIDs[i] = id
	}

	// 插入测试
	start := time.Now()
	err = db.BatchInsert(coll, docs, batchSize)
	if err != nil {
		log.Fatalf("batch insert failed: %v", err)
	}
	fmt.Printf("✅ Inserted %d docs in %v\n", total, time.Since(start))

	// 查询测试
	rand.Seed(time.Now().UnixNano())
	sampleCount := 10
	fmt.Printf("🔍 Start %d random GetByID tests\n", sampleCount)
	for i := 0; i < sampleCount; i++ {
		idx := rand.Intn(total)
		id := docIDs[idx]
		doc, err := db.GetByID(coll, id)
		if err != nil || doc == nil {
			log.Fatalf("❌ GetByID failed for id %s: %v", id, err)
		}
		fmt.Printf("📄 Found: ID=%s Title=%s\n", doc.ID, doc.Title)
	}

	// 文档总数检查
	fmt.Printf("📊 Total documents in DB: %d\n", db.GetDocCount())
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

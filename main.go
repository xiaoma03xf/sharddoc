package main

import (
	"encoding/json"
	"fmt"
	"log"
   
	"github.com/cockroachdb/pebble"
)

// Document 结构体定义
type Document struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// StoreDocument 插入文档
func StoreDocument(db *pebble.DB, doc Document) error {
	// 序列化 Document 为 JSON 格式
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("序列化文档失败: %v", err)
	}

	// 使用 Pebble 的 Set 方法将文档存入数据库
	return db.Set([]byte(doc.ID), data, nil)
}

// GetDocument 查询文档
func GetDocument(db *pebble.DB, id string) (*Document, error) {
	// 获取文档数据
	val, closer, err := db.Get([]byte(id))
	if err != nil {
		return nil, fmt.Errorf("查询文档失败: %v", err)
	}
	defer closer.Close()

	// 反序列化 JSON 数据为 Document 结构体
	var doc Document
	err = json.Unmarshal(val, &doc)
	if err != nil {
		return nil, fmt.Errorf("反序列化文档失败: %v", err)
	}
	return &doc, nil
}

// UpdateDocument 更新文档
func UpdateDocument(db *pebble.DB, doc Document) error {
	// 序列化 Document 为 JSON 格式
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("序列化文档失败: %v", err)
	}

	// 更新文档内容
	return db.Set([]byte(doc.ID), data, nil)
}

// DeleteDocument 删除文档
func DeleteDocument(db *pebble.DB, id string) error {
	// 使用 Pebble 的 Delete 方法删除文档
	return db.Delete([]byte(id), nil)
}

func main() {
	// 打开数据库
	db, err := pebble.Open("document_store.db", &pebble.Options{})
	if err != nil {
		log.Fatalf("无法打开数据库: %v", err)
	}
	defer db.Close()

	// 创建一个新的文档
	doc := Document{
		ID:      "123",
		Content: "This is a test document.",
	}

	// 存储文档
	err = StoreDocument(db, doc)
	if err != nil {
		log.Fatalf("存储文档失败: %v", err)
	}
	fmt.Println("文档存储成功!")

	// 查询文档
	retrievedDoc, err := GetDocument(db, "123")
	if err != nil {
		log.Fatalf("查询文档失败: %v", err)
	}
	fmt.Printf("查询文档: %+v\n", retrievedDoc)

	// 更新文档
	doc.Content = "This is an updated document."
	err = UpdateDocument(db, doc)
	if err != nil {
		log.Fatalf("更新文档失败: %v", err)
	}
	fmt.Println("文档更新成功!")

	// 查询更新后的文档
	updatedDoc, err := GetDocument(db, "123")
	if err != nil {
		log.Fatalf("查询更新后的文档失败: %v", err)
	}
	fmt.Printf("查询更新后的文档: %+v\n", updatedDoc)

	// 删除文档
	err = DeleteDocument(db, "123")
	if err != nil {
		log.Fatalf("删除文档失败: %v", err)
	}
	fmt.Println("文档删除成功!")

	// 确认文档是否已删除
	_, err = GetDocument(db, "123")
	if err != nil {
		fmt.Println("文档已删除或不存在")
	}
}

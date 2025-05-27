package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/cockroachdb/pebble"
)

// 创建倒排索引
func AddToInvertedIndex(db *pebble.DB, word string, docID string) error {
	key := []byte("idx:" + word)

	// 获取已有索引
	val, closer, err := db.Get(key)
	var docIDs []string
	if err == nil {
		defer closer.Close()
		if err := json.Unmarshal(val, &docIDs); err != nil {
			return fmt.Errorf("解析倒排索引失败: %v", err)
		}
	}

	// 避免重复添加
	for _, id := range docIDs {
		if id == docID {
			return nil // 已存在
		}
	}
	docIDs = append(docIDs, docID)

	// 更新索引
	newVal, _ := json.Marshal(docIDs)
	return db.Set(key, newVal, nil)
}

// 查询包含某个词的所有文档ID
func QueryInvertedIndex(db *pebble.DB, word string) ([]string, error) {
	key := []byte("idx:" + word)
	val, closer, err := db.Get(key)
	if err != nil {
		return nil, fmt.Errorf("未找到关键词 %s 的索引", word)
	}
	defer closer.Close()

	var docIDs []string
	if err := json.Unmarshal(val, &docIDs); err != nil {
		return nil, fmt.Errorf("解析索引失败: %v", err)
	}
	return docIDs, nil
}

// 对文档内容进行简单分词（按空格切）
func indexDocument(db *pebble.DB, docID string, content string) error {
	words := strings.Fields(strings.ToLower(content))
	for _, word := range words {
		if err := AddToInvertedIndex(db, word, docID); err != nil {
			return err
		}
	}
	return nil
}

func TestPostring(t *testing.T) {
	// 打开 Pebble 数据库
	db, err := pebble.Open("invindex.db", &pebble.Options{})
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	defer db.Close()

	// 模拟文档
	docs := map[string]string{
		"doc1": "Go is a programming language",
		"doc2": "Python is a programming language",
		"doc3": "Go and Python are popular languages",
	}

	// 索引所有文档
	for id, content := range docs {
		if err := indexDocument(db, id, content); err != nil {
			log.Fatalf("文档索引失败: %v", err)
		}
	}

	// 查询关键词
	keywords := []string{"go", "python", "programming", "java"}
	for _, word := range keywords {
		docIDs, err := QueryInvertedIndex(db, word)
		if err != nil {
			fmt.Printf("查询词 '%s' 失败: %v\n", word, err)
		} else {
			fmt.Printf("词 '%s' 出现在文档: %v\n", word, docIDs)
		}
	}
}

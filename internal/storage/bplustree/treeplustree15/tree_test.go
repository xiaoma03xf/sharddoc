package bplustree15

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	data := kvs[int, string]{
		{key: 5, value: "five"},
		{key: 2, value: "two"},
		{key: 8, value: "eight"},
		{key: 1, value: "one"},
		{key: 3, value: "three"},
		{key: 4, value: "four"},
		{key: 9, value: "nine"},
		{key: 7, value: "seven"},
	}

	fmt.Println("Before sort:")
	for _, kv := range data {
		fmt.Printf("key: %d, value: %s\n", kv.key, kv.value)
	}

	// 调用 sort.Sort
	sort.Sort(data)

	fmt.Println("\nAfter sort:")
	for _, kv := range data {
		fmt.Printf("key: %d, value: %s\n", kv.key, kv.value)
	}

	i := sort.Search(data.Len(), func(i int) bool {
		return data[i].key >= 5
	})

	fmt.Println("ceiling 5: answer", data[i].key)

}

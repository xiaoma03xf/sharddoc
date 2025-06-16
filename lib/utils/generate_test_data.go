package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

type RecordTestData struct {
	ID     int64
	Name   string
	Age    int64
	Height int64
}

func GenerateTestData(filepath string, count int) error {
	generateRandomString := func() string {
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		b := make([]byte, 6)
		b[0] = letters[rand.Intn(26)+26]
		for i := 1; i < 6; i++ {
			b[i] = letters[rand.Intn(26)]
		}
		return string(b)
	}

	records := make([]RecordTestData, count)
	for i := 0; i < count; i++ {
		records[i] = RecordTestData{
			ID:     int64(i + 1),
			Name:   generateRandomString(),
			Age:    rand.Int63n(30) + 15,
			Height: rand.Int63n(40) + 150,
		}
	}

	// 写入文件部分保持不变...
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		return fmt.Errorf("could not write data to file: %v", err)
	}

	if err := os.Chmod(filepath, 0644); err != nil {
		return fmt.Errorf("could not set file permissions: %v", err)
	}

	return nil
}
func ReadTestDataFromFile(filePath string) ([]RecordTestData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()
	var records []RecordTestData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&records); err != nil {
		return nil, fmt.Errorf("could not read data from file: %v", err)
	}

	return records, nil
}

func GenerateData(filepath string, count int) ([]RecordTestData, error) {
	if !FileExists(filepath) {
		_ = GenerateTestData(filepath, count)
	}

	records, err := ReadTestDataFromFile(filepath)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return nil, err
	}
	return records, nil
}

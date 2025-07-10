package tcp

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xiaoma03xf/sharddoc/discarded/storage"
	"github.com/xiaoma03xf/sharddoc/lib/utils"
)

func TestClientBasic(t *testing.T) {
	defer func() {
		os.RemoveAll("../clusterdb")
		os.Remove("./test_data.json")
	}()
	bootCluster(t)

	client, err := Open("127.0.0.1:29001")
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Exec(`
	CREATE TABLE users (
        id INT64,
        name BYTES,
        age INT64,
		height INT64,
		PRIMARY KEY (id),
        INDEX (age, height)
    );
	`)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("create table response:", resp)

	// 生成测试数据
	filePath := "test_data.json"
	datacnt := 2000
	recs, _ := utils.GenerateData(filePath, datacnt)
	for _, rec_data := range recs {
		cols := []string{"id", "name", "age", "height"}
		args := []interface{}{rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height}

		resp, err := client.Exec(storage.BuildInsertSQL("users", cols, args))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, resp, "OK")
	}

	// 查询插入后的数据
	sqlselect := "SELECT name, id FROM users WHERE age >= 0"
	records, err := client.Raw(sqlselect)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, len(records), datacnt)
	// fmt.Println("data after insert:")
	// PrintlnRec(records)

	// 执行数据更新
	updatasql := "UPDATE users SET age=10086 WHERE age>22"
	resp, err = client.Exec(updatasql)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("update data info", resp)

	// 查询修改更新后的数量
	updatedCnt := 0
	sqlselect = "SELECT name, id FROM users WHERE age=10086"
	records, err = client.Raw(sqlselect)
	if err != nil {
		t.Error(err)
	}
	for _, rec := range records {
		for i, x := range rec.Cols {
			if x == "age" && rec.Vals[i].I64 == 10086 {
				updatedCnt++
			}
		}
	}
	fmt.Println("updated data cnt: ", updatedCnt)
	// fmt.Println("data after update:")
	// PrintlnRec(records)

	// 查询要删除的数据
	sqlselect = "SELECT name, id FROM users WHERE age=10086 AND height > 175"
	records, err = client.Raw(sqlselect)
	if err != nil {
		t.Error(err)
	}
	deleteCnt := len(records)

	// 执行数据删除
	delsql := "DELETE FROM users WHERE age=10086 AND height > 175"
	resp, err = client.Exec(delsql)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("after delete info", resp)
	fmt.Println("delete data cnt: ", deleteCnt)

	// 执行删除后数据量
	sqlselect = "SELECT name, id FROM users WHERE age=10086"
	records, err = client.Raw(sqlselect)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("after delete count: ", len(records))
	assert.Equal(t, deleteCnt+len(records), updatedCnt)
}

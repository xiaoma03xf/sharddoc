package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdBasic(t *testing.T) {
	defer func() {
		os.RemoveAll("../clusterdb")
	}()
	conf := "../cluster1.yaml"

	go raft.BootstrapCluster(conf, "node1")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node2")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node3")
	time.Sleep(3 * time.Second)

	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1"})
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 5; i++ {
		client, _, err := db.getLeader("cluster1")
		if err != nil {
			t.Error(err)
		}
		resp, err := client.Status(context.Background(), &pb.StatusRequest{})
		if err != nil {
			t.Error(err)
		}
		fmt.Println(resp.Me)
		time.Sleep(1 * time.Second)
	}

	// default CachedTimeOut = 10*time.Second
	time.Sleep(10 * time.Second)
	client, _, err := db.getLeader("cluster1")
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Status(context.Background(), &pb.StatusRequest{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.Me)
}
func ClearEtcd(endpoints []string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create etcd client: %w", err))
	}
	_, err = client.Delete(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
}
func TestTableNew(t *testing.T) {
	ClearEtcd([]string{"118.89.66.104:2379"})

	defer func() {
		os.RemoveAll("../../clusterdb")
	}()
	conf := "../../cluster1.yaml"

	go raft.BootstrapCluster(conf, "node1")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node2")
	time.Sleep(3 * time.Second)
	go raft.BootstrapCluster(conf, "node3")
	time.Sleep(3 * time.Second)

	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1"})
	if err != nil {
		t.Error(err)
	}

	tdef1 := &kv.TableDef{
		Name:  "TableNew",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{kv.TYPE_INT64, kv.TYPE_BYTES, kv.TYPE_INT64, kv.TYPE_INT64},
		Indexes: [][]string{
			{"id"}, // 主键索引
			{"name"},
			{"age", "height"}, // 二级索引（复合索引）
		},
	}
	meta1, err := db.TablesDefDiscovery.GetMetaKey(NEXT_PREFIX)
	if err != nil {
		t.Error(err)
	}
	if meta1 != nil {
		fmt.Println("meta before", bytesToUint32LE(meta1))
	}
	err = db.TableNew(tdef1)
	if err != nil {
		t.Error(err)
	}
	meta2, err := db.TablesDefDiscovery.GetMetaKey(NEXT_PREFIX)
	if err != nil {
		t.Error(err)
	}
	if meta2 != nil {
		fmt.Println("meta after", bytesToUint32LE(meta2))
	}
}

func BootCluster(conf string, nodes []string) {
	for _, node := range nodes {
		go raft.BootstrapCluster(conf, node)
		time.Sleep(3 * time.Second)
	}
}

type RecordTestData struct {
	ID     int64
	Name   string
	Age    int64
	Height int64
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
func TestDBBasic(t *testing.T) {
	// go BootCluster("../../cluster1.yaml", []string{"node1", "node2", "node3"})
	// go BootCluster("../../cluster2.yaml", []string{"node4", "node5", "node6"})
	// go BootCluster("../../cluster3.yaml", []string{"node7", "node8", "node9"})
	// time.Sleep(10 * time.Second)

	ClearEtcd([]string{"118.89.66.104:2379"})
	db, err := NewDB([]string{"118.89.66.104:2379"}, []string{"cluster1", "cluster2", "cluster3"})
	if err != nil {
		t.Error(err)
	}
	tdef := &kv.TableDef{
		Name:  "user",
		Cols:  []string{"id", "name", "age", "height"},
		Types: []uint32{kv.TYPE_INT64, kv.TYPE_BYTES, kv.TYPE_INT64, kv.TYPE_INT64},
		Indexes: [][]string{
			{"id"}, // 主键索引
			{"name"},
			{"age", "height"}, // 二级索引（复合索引）
		},
	}
	err = db.TableNew(tdef)
	if err != nil {
		t.Error(err)
	}

	// 插入数据
	record := func(id int64, name string, age int64, height int64) kv.Record {
		rec := kv.Record{}
		rec.AddInt64("id", id).AddStr("name", []byte(name))
		rec.AddInt64("age", age).AddInt64("height", height)
		return rec
	}
	fmt.Printf("Adding test records...\n")

	// 持久化测试数据,便于观察
	filePath := "./test_data.json"
	records, err := ReadTestDataFromFile(filePath)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	for _, rec_data := range records {
		rec := record(rec_data.ID, rec_data.Name, rec_data.Age, rec_data.Height)
		ok, err := db.Insert("user", rec)
		if !ok || err != nil {
			t.Error(err)
			return
		}
	}
	fmt.Printf("Added test records successfully\n")
}

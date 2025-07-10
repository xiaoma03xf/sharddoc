package raft

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/raftpb"
	"github.com/xiaoma03xf/sharddoc/server/etcd"
	"google.golang.org/protobuf/proto"
)

func TestStoreInterface(t *testing.T) {
	var _ raftpb.KVStoreServer = &Store{}
	kv2 := kv.KV{Path: "./test.db"}
	if err := kv2.Open(); err != nil {
		panic(err)
	}
	// Generate key-value pairs
	type kvPair struct {
		key   string
		value string
	}
	kvPairs := make([]kvPair, 1000)
	for i := 0; i < 1000; i++ {
		kvPairs[i] = kvPair{
			key:   fmt.Sprintf("key_%d", i),
			value: fmt.Sprintf("val_%d", i),
		}
	}

}
func TestBootStrap(t *testing.T) {
	// defer func() {
	// 	os.RemoveAll("../clusterdb")
	// }()
	// conf := "../cluster1.yaml"

	// go BootstrapCluster(conf, "node1")
	// time.Sleep(2 * time.Second)
	// go BootstrapCluster(conf, "node2")
	// time.Sleep(2 * time.Second)
	// go BootstrapCluster(conf, "node3")
	// time.Sleep(2 * time.Second)

	// High-concurrency read/write test
	const (
		numOperations = 1000 // Total Put/Get operations
		concurrency   = 50   // Concurrent workers
		leaderAddr    = "127.0.0.1:29001"
	)
	var (
		wg          sync.WaitGroup
		putSuccess  int32
		getSuccess  int32
		getNotFound int32
		putErrors   int32
		getErrors   int32
		startTime   = time.Now()
	)
	// Generate key-value pairs
	type kvPair struct {
		key   string
		value string
	}
	kvPairs := make([]kvPair, numOperations)
	for i := 0; i < numOperations; i++ {
		kvPairs[i] = kvPair{
			key:   fmt.Sprintf("key_%d", i),
			value: fmt.Sprintf("val_%d", i),
		}
	}
	// Worker for concurrent Put and Get operations
	worker := func(start, end int) {
		defer wg.Done()
		// Create a separate gRPC connection per worker to avoid contention
		localClient, localConn, err := BuildGrpcConn(leaderAddr)
		if err != nil {
			atomic.AddInt32(&putErrors, 1)
			atomic.AddInt32(&getErrors, 1)
			t.Errorf("failed to connect to leader %s: %v", leaderAddr, err)
			return
		}
		defer localConn.Close()

		// Perform Put and Get concurrently
		for i := start; i < end; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			// Put operation
			putResp, err := localClient.Put(ctx, &raftpb.PutRequest{
				Key:   []byte(kvPairs[i].key),
				Value: []byte(kvPairs[i].value),
			})
			if err != nil {
				atomic.AddInt32(&putErrors, 1)
				t.Errorf("Put failed for key %s: %v", kvPairs[i].key, err)
			} else if putResp.Success {
				atomic.AddInt32(&putSuccess, 1)
			}

			// Get operation
			getResp, err := localClient.Get(ctx, &raftpb.GetRequest{
				Key: []byte(kvPairs[i].key),
			})
			if err != nil {
				atomic.AddInt32(&getErrors, 1)
				t.Errorf("Get failed for key %s: %v", kvPairs[i].key, err)
			} else if getResp.Found {
				if string(getResp.Value) != kvPairs[i].value {
					atomic.AddInt32(&getErrors, 1)
					t.Errorf("Get value mismatch for key %s: expected %s, got %s", kvPairs[i].key, kvPairs[i].value, getResp.Value)
				} else {
					atomic.AddInt32(&getSuccess, 1)
				}
			} else {
				atomic.AddInt32(&getNotFound, 1)
				t.Errorf("Key %s not found", kvPairs[i].key)
			}
			cancel()
		}
	}

	// Start workers
	operationsPerWorker := numOperations / concurrency
	for i := 0; i < concurrency; i++ {
		start := i * operationsPerWorker
		end := start + operationsPerWorker
		if i == concurrency-1 {
			end = numOperations // Handle remainder
		}
		wg.Add(1)
		go worker(start, end)
	}
	wg.Wait()

	// Performance metrics
	duration := time.Since(startTime)
	t.Logf("Test completed: %d Put successes, %d Put errors, %d Get successes, %d Get not found, %d Get errors", putSuccess, putErrors, getSuccess, getNotFound, getErrors)
	t.Logf("Test duration: %v, Throughput: %.2f ops/sec", duration, float64(numOperations*2)/duration.Seconds())

	// Assertions
	assert.Equal(t, int32(numOperations), putSuccess, "Put success count mismatch")
	assert.Equal(t, int32(0), putErrors, "Put errors occurred")
	assert.Equal(t, int32(numOperations), getSuccess, "Get success count mismatch")
	assert.Equal(t, int32(0), getNotFound, "Get not found errors occurred")
	assert.Equal(t, int32(0), getErrors, "Get errors occurred")
}

func TestBatchInsert(t *testing.T) {
	// defer func() {
	// 	os.RemoveAll("../clusterdb")
	// }()
	conf := "../cluster1.yaml"

	go BootstrapCluster(conf, "node1")
	time.Sleep(2 * time.Second)
	go BootstrapCluster(conf, "node2")
	time.Sleep(2 * time.Second)
	go BootstrapCluster(conf, "node3")
	time.Sleep(2 * time.Second)

	// Generate key-value pairs
	kvPairs := make([]*raftpb.KeyValue, 1000)
	for i := 0; i < 1000; i++ {
		kvPairs[i] = &raftpb.KeyValue{
			Key:   []byte(fmt.Sprintf("key_%d", i)),
			Value: []byte(fmt.Sprintf("val_%d", i)),
		}
	}
	leaderAddr := "127.0.0.1:29001"
	localClient, localConn, err := BuildGrpcConn(leaderAddr)
	Assert(err == nil)
	defer localConn.Close()

	resp, err := localClient.BatchPut(context.Background(), &raftpb.BatchPutRequest{Pairs: kvPairs})
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		Assert(resp.Success)
	}

	for _, pair := range kvPairs {
		resp, err := localClient.Get(context.Background(), &raftpb.GetRequest{Key: pair.Key})
		Assert(err == nil)
		if string(resp.Value) != string(pair.Value) {
			t.Errorf("value not expected, want: %v, got: %v", string(pair.Value), string(resp.Value))
		}
		fmt.Printf("value want: %v, got: %v\n", string(pair.Value), string(resp.Value))
	}
}

func TestEtcdConf(t *testing.T) {
	defer func() {
		os.RemoveAll("../clusterdb")
	}()
	conf := "../cluster1.yaml"

	go BootstrapCluster(conf, "node1")
	time.Sleep(2 * time.Second)
	go BootstrapCluster(conf, "node2")
	time.Sleep(2 * time.Second)
	go BootstrapCluster(conf, "node3")
	time.Sleep(2 * time.Second)

	// 启动后检测当前集群状态，询问node3当前集群状态
	leaderAddr := "127.0.0.1:29002"
	localClient, localConn, err := BuildGrpcConn(leaderAddr)
	Assert(err == nil)
	defer localConn.Close()
	resp, err := localClient.Status(context.Background(), &raftpb.StatusRequest{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Leader:", resp.Leader)
	fmt.Println("Me:", resp.Me)
	fmt.Println("Follower:", resp.Follower)

	endpoints := []string{"118.89.66.104:2379"}
	clusterID := []string{"cluster1"}

	// 检测此时etcd服务
	sd, err := etcd.NewServiceDiscovery(endpoints, clusterID)
	if err != nil {
		log.Fatalf("创建服务发现失败: %v", err)
	}
	defer sd.Close()

	sd.Start()
	// 获取当前etcd中的leader
	res := sd.GetServiceByClusterID("cluster1")
	fmt.Println("当前leader信息", res.Addr)
}

func TestProtobuf(t *testing.T) {
	// 构造 PutRequest，基于日志数据
	req := &raftpb.PutRequest{
		Key:   []byte{0, 0, 0, 100, 2, 128, 0, 0, 0, 0, 0, 0, 1},                                           // Key: [0 0 0 100 2 128 0 0 0 0 0 0 1]
		Value: []byte{1, 106, 97, 99, 107, 0, 2, 128, 0, 0, 0, 0, 0, 0, 23, 2, 128, 0, 0, 0, 0, 0, 0, 175}, // Value: [1 106 97 99 107 0 2 128 0 0 0 0 0 0 23 2 128 0 0 0 0 0 0 175]
		Mode:  1,                                                                                           // Mode: 1
	}

	// 序列化
	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal PutRequest: %v", err)
	}
	fmt.Printf("Serialized data: %x\n", data)

	// 验证序列化数据是否匹配日志中的 reqData
	expectedData := []byte{0x0a, 0x0d, 0x00, 0x00, 0x00, 0x64, 0x02, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x12, 0x18, 0x01, 0x6a, 0x61, 0x63, 0x6b, 0x00, 0x02, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17, 0x02, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaf, 0x18, 0x01}
	fmt.Printf("Expected data: %x\n", expectedData)
	if string(data) != string(expectedData) {
		log.Fatalf("Serialized data does not match expected: got %x, want %x", data, expectedData)
	}

	// 反序列化
	var putreq raftpb.PutRequest
	err = proto.Unmarshal(data, &putreq)
	if err != nil {
		log.Fatalf("Failed to unmarshal PutRequest: %v", err)
	}
	fmt.Printf("Unmarshaled PutRequest: key=%x, value=%x, mode=%d\n", putreq.Key, putreq.Value, putreq.Mode)

	// 验证反序列化结果
	if string(putreq.Key) != string(req.Key) || string(putreq.Value) != string(req.Value) || putreq.Mode != req.Mode {
		log.Fatalf("Unmarshaled data does not match: got key=%x, value=%x, mode=%d; want key=%x, value=%x, mode=%d",
			putreq.Key, putreq.Value, putreq.Mode, req.Key, req.Value, req.Mode)
	}

	fmt.Println("Serialization and deserialization successful!")
}

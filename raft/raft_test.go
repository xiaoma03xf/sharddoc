package raft

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xiaoma03xf/sharddoc/kv"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
)

func TestStoreInterface(t *testing.T) {
	var _ pb.KVStoreServer = &Store{}
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
	defer func() {
		os.RemoveAll("../clusterdb")
	}()
	conf1 := "../node1.yaml"
	conf2 := "../node2.yaml"
	conf3 := "../node3.yaml"

	go BootstrapCluster(conf1)
	time.Sleep(2 * time.Second)
	go BootstrapCluster(conf2)
	time.Sleep(2 * time.Second)
	go BootstrapCluster(conf3)
	time.Sleep(2 * time.Second)

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
			putResp, err := localClient.Put(ctx, &pb.PutRequest{
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
			getResp, err := localClient.Get(ctx, &pb.GetRequest{
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
	// Generate key-value pairs
	kvPairs := make([]*pb.KeyValue, 1000)
	for i := 0; i < 1000; i++ {
		kvPairs[i] = &pb.KeyValue{
			Key:   []byte(fmt.Sprintf("key_%d", i)),
			Value: []byte(fmt.Sprintf("val_%d", i)),
		}
	}
	leaderAddr := "127.0.0.1:29001"
	localClient, localConn, err := BuildGrpcConn(leaderAddr)
	Assert(err == nil)
	defer localConn.Close()

	resp, err := localClient.BatchPut(context.Background(), &pb.BatchPutRequest{Pairs: kvPairs})
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		Assert(resp.Success)
	}

	for _, pair := range kvPairs {
		resp, err := localClient.Get(context.Background(), &pb.GetRequest{Key: pair.Key})
		Assert(err == nil)
		if string(resp.Value) != string(pair.Value) {
			t.Errorf("value not expected, want: %v, got: %v", string(pair.Value), string(resp.Value))
		}
		fmt.Printf("value want: %v, got: %v\n", string(pair.Value), string(resp.Value))
	}
}

package server

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/xiaoma03xf/sharddoc/raft"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
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

	db, err := NewDB("cluster1", []string{"118.89.66.104:2379"}, []string{"cluster1"})
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

package raft

import (
	"context"
	"fmt"
	"testing"

	"github.com/xiaoma03xf/sharddoc/raft/pb"
)

func TestStoreInterface(t *testing.T) {
	var _ pb.KVStoreServer = &Store{}
	// 也可以使用 `assert` 来进一步验证
	s := Store{}
	usegrpcGet := func(s pb.KVStoreServer) {
		res, err := s.Join(context.Background(), &pb.JoinRequest{})
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
	usegrpcGet(&s)
}
func TestBootStrap(t *testing.T) {
	// defer func() {
	// 	os.RemoveAll("./clusterdb")
	// }()
	// conf1 := "../node1.yaml"
	// conf2 := "../node2.yaml"
	// conf3 := "../node3.yaml"

	// go BootstrapCluster(conf1)
	// time.Sleep(2 * time.Second)
	// go BootstrapCluster(conf2)
	// time.Sleep(2 * time.Second)

	// go BootstrapCluster(conf3)
	// time.Sleep(2 * time.Second)

	client, conn, err := BuildGrpcConn("127.0.0.1:29001")
	defer conn.Close()
	assert(err == nil)
	{
		resp, err := client.Status(context.Background(), &pb.StatusRequest{})
		if err != nil {
			t.Error(err)
		}
		fmt.Println(resp.Me)
		fmt.Println(resp.Leader)
		fmt.Println(resp.Follower)

	}

	{
		cnt := 5
		for i := 0; i < cnt; i++ {
			key := fmt.Sprintf("key_%d", i)
			val := fmt.Sprintf("val_%d", i)
			resp, err := client.Put(context.Background(), &pb.PutRequest{
				Key: []byte(key), Value: []byte(val),
			})
			if err != nil {
				t.Error(err)
			}
			if resp != nil {
				fmt.Println(resp.Success)
			}
		}
	}
}

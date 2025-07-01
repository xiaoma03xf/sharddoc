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

package raft

import (
	"context"
	"fmt"

	"github.com/xiaoma03xf/sharddoc/raft/pb"
)

func (s *Store) Join(context.Context, *pb.JoinRequest) (*pb.JoinResponse, error) {
	fmt.Println("join")
	return &pb.JoinResponse{}, nil
}

func (s *Store) Put(context.Context, *pb.PutRequest) (*pb.PutResponse, error) {
	return &pb.PutResponse{}, nil
}
func (s *Store) Get(context.Context, *pb.GetRequest) (*pb.GetResponse, error) {
	return &pb.GetResponse{}, nil
}
func (s *Store) Delete(context.Context, *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	return &pb.DeleteResponse{}, nil
}
func (s *Store) BatchPut(context.Context, *pb.BatchPutRequest) (*pb.BatchPutResponse, error) {
	return &pb.BatchPutResponse{}, nil
}
func (s *Store) Scan(context.Context, *pb.ScanRequest) (*pb.ScanResponse, error) {
	return &pb.ScanResponse{}, nil
}
func (s *Store) Status(context.Context, *pb.StatusRequest) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{}, nil
}

package etcd

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceRegistry struct {
	client    *clientv3.Client
	lease     clientv3.Lease
	leaseID   clientv3.LeaseID
	keepAlive <-chan *clientv3.LeaseKeepAliveResponse
	key       string
	value     string
	ttl       int64
}

func NewServiceRegistry(endpoints []string, clusterID, serviceAddr string, ttl int64) (*ServiceRegistry, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	// 创建lease对象, 并创建租约
	// lease := clientv3.NewLease(client)
	// leaseResp, err := lease.Grant(context.Background(), ttl)
	// if err != nil {
	// 	return nil, err
	// }
	key := fmt.Sprintf("/services/%s/%s", clusterID, serviceAddr)
	value := fmt.Sprintf(`{
		"addr": "%s",
		"timestamp": "%s",
		"version": "1.0.0"
	}`, serviceAddr, time.Now().Format(time.RFC3339))
	sr := &ServiceRegistry{
		client: client,
		// lease:   lease,
		// leaseID: leaseResp.ID,
		key:   key,
		value: value,
		ttl:   ttl,
	}
	return sr, nil
}
func (sr *ServiceRegistry) Register(isLeader func() bool, onDeregister func()) error {
	// 初始化租约
	lease := clientv3.NewLease(sr.client)
	leaseResp, err := lease.Grant(context.Background(), sr.ttl)
	if err != nil {
		return err
	}
	sr.lease = lease
	sr.leaseID = leaseResp.ID
	// 发送租约请求
	_, err = sr.client.Put(context.Background(), sr.key, sr.value, clientv3.WithLease(sr.leaseID))
	if err != nil {
		return err
	}
	sr.keepAlive, err = sr.lease.KeepAlive(context.Background(), sr.leaseID)
	if err != nil {
		return err
	}
	log.Printf("服务已注册: %s -> %s", sr.key, sr.value)
	go sr.watchKeepAlive(isLeader, onDeregister)
	return nil
}

func (sr *ServiceRegistry) watchKeepAlive(isLeader func() bool, onDeregister func()) {
	for ka := range sr.keepAlive {
		if !isLeader() {
			log.Println("检测到非 Leader，注销服务")
			_ = sr.Deregister()
			onDeregister()
			break
		}
		if ka != nil {
			log.Printf("续约成功, TTL: %d", ka.TTL)
		}
	}
	log.Println("续约channel已关闭")
}
func (sr *ServiceRegistry) Deregister() error {
	_, err := sr.lease.Revoke(context.Background(), sr.leaseID)
	if err != nil {
		log.Printf("撤销租约失败: %v", err)
		_, delErr := sr.client.Delete(context.Background(), sr.key)
		return delErr
	}

	log.Printf("服务已注销: %s", sr.key)
	return nil
}

func (sr *ServiceRegistry) Close() error {
	return sr.client.Close()
}

package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/xiaoma03xf/sharddoc/kv"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceRegistry struct {
	client      *clientv3.Client
	lease       clientv3.Lease
	leaseID     clientv3.LeaseID
	keepAlive   <-chan *clientv3.LeaseKeepAliveResponse
	clusterID   string
	serviceAddr string
	key         string
	value       string
	ttl         int64
}

func NewServiceRegistry(endpoints []string, clusterID, serviceAddr string, ttl int64) (*ServiceRegistry, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

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
		key:         key,
		value:       value,
		ttl:         ttl,
		clusterID:   clusterID,
		serviceAddr: serviceAddr,
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
	for range sr.keepAlive {
		if !isLeader() {
			log.Println("检测到非 Leader，注销服务")
			if err := sr.Deregister(); err != nil {
				log.Printf("注销服务失败: %v", err)
			}
			onDeregister()
			break
		}
		// 不打印续约成功日志，默默消费 channel 数据
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

/*
	// 表信息存储
*/
// TableDefRegistry manages persistent table definition registration in etcd.
type TableDefRegistry struct {
	client     *clientv3.Client
	ctx        context.Context
	cancel     context.CancelFunc
	mutex      sync.Mutex
	prefix     string // e.g., "/services/tables/"
	metaPrefix string // e.g., "/services/meta/"
}

// NewTableDefRegistry initializes a TableDefRegistry without lease.
func NewTableDefRegistry(endpoints []string) (*TableDefRegistry, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	reg := &TableDefRegistry{
		client:     client,
		ctx:        ctx,
		cancel:     cancel,
		prefix:     "/services/tables/",
		metaPrefix: "/services/meta/",
	}
	return reg, nil
}

// RegisterTable registers a TableDef in etcd persistently.
func (r *TableDefRegistry) RegisterTable(tdef *kv.TableDef) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	key := r.prefix + tdef.Name
	value, err := json.Marshal(tdef)
	if err != nil {
		return fmt.Errorf("failed to marshal TableDef: %w", err)
	}
	_, err = r.client.Put(r.ctx, key, string(value))
	if err != nil {
		return fmt.Errorf("failed to put TableDef to etcd: %w", err)
	}
	return nil
}

// GetTable retrieves a TableDef from etcd.
func (r *TableDefRegistry) GetTable(tableName string) (*kv.TableDef, error) {
	key := r.prefix + tableName
	resp, err := r.client.Get(r.ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get table %s: %w", tableName, err)
	}
	if len(resp.Kvs) == 0 {
		return nil, nil // Table not found
	}

	var tdef kv.TableDef
	if err := json.Unmarshal(resp.Kvs[0].Value, &tdef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TableDef: %w", err)
	}
	return &tdef, nil
}

// DeleteTable removes a TableDef from etcd.
func (r *TableDefRegistry) DeleteTable(tableName string) error {
	key := r.prefix + tableName
	resp, err := r.client.Delete(r.ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete table %s: %w", tableName, err)
	}
	if resp.Deleted == 0 {
		return fmt.Errorf("table %s not found", tableName)
	}
	return nil
}

// PutMetaKey stores or updates a key-value pair in /services/meta/.
func (r *TableDefRegistry) PutMetaKey(key string, value []byte) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	metaKey := r.metaPrefix + key
	_, err := r.client.Put(r.ctx, metaKey, string(value))
	if err != nil {
		return fmt.Errorf("failed to put key %s to etcd: %w", key, err)
	}
	return nil
}

// GetMetaKey retrieves a value for a key from /services/meta/.
func (r *TableDefRegistry) GetMetaKey(key string) ([]byte, error) {
	metaKey := r.metaPrefix + key
	resp, err := r.client.Get(r.ctx, metaKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get key %s: %w", key, err)
	}
	if len(resp.Kvs) == 0 {
		return nil, nil // Key not found
	}
	return resp.Kvs[0].Value, nil
}

// DeleteMetaKey removes a key from /services/meta/.
func (r *TableDefRegistry) DeleteMetaKey(key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	metaKey := r.metaPrefix + key
	resp, err := r.client.Delete(r.ctx, metaKey)
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	if resp.Deleted == 0 {
		return fmt.Errorf("key %s not found", key)
	}
	return nil
}

// UpdateMetaKey updates an existing key-value pair in /services/meta/.
func (r *TableDefRegistry) UpdateMetaKey(key string, value []byte) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	metaKey := r.metaPrefix + key
	resp, err := r.client.Txn(r.ctx).
		If(clientv3.Compare(clientv3.CreateRevision(metaKey), "!=", 0)).
		Then(clientv3.OpPut(metaKey, string(value))).
		Commit()
	if err != nil {
		return fmt.Errorf("etcd transaction failed for key %s: %w", key, err)
	}
	if !resp.Succeeded {
		return fmt.Errorf("key %s does not exist", key)
	}
	return nil
}

// Close shuts down the TableDefRegistry.
func (r *TableDefRegistry) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.cancel()
	return r.client.Close()
}

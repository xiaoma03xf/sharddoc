package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceInfo struct {
	Addr      string `json:"addr"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// ServiceDiscovery etcd服务发现器
type ServiceDiscovery struct {
	client *clientv3.Client
	mutex  sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	clusterIDs []string
	services   map[string]*ServiceInfo // 所有clusterID对应的Info
	watchChans map[string]chan struct{}
}

func NewServiceDiscovery(endpoints []string, clusterIDs []string) (*ServiceDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	sd := &ServiceDiscovery{
		client:     client,
		services:   make(map[string]*ServiceInfo),
		clusterIDs: clusterIDs,
		watchChans: make(map[string]chan struct{}),
		ctx:        ctx,
		cancel:     cancel,
	}
	for _, cid := range clusterIDs {
		sd.services[cid] = &ServiceInfo{}
		sd.watchChans[cid] = make(chan struct{}, 1)
	}
	return sd, nil
}

func (sd *ServiceDiscovery) Start() error {
	for _, clusterID := range sd.clusterIDs {
		if err := sd.loadServices(clusterID); err != nil {
			log.Printf("加载集群[%s]服务失败: %v", clusterID, err)
		}
		go sd.watch(clusterID)
	}
	log.Printf("服务发现已启动，监听多个集群: %v", sd.clusterIDs)
	return nil
}

func (sd *ServiceDiscovery) loadServices(clusterID string) error {
	prefix := fmt.Sprintf("/services/%s/", clusterID)

	resp, err := sd.client.Get(sd.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("etcd 获取失败: %w", err)
	}

	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	// 清空现有服务
	sd.services[clusterID] = nil
	for _, kv := range resp.Kvs {
		// addr := sd.extractAddrFromKey(string(kv.Key))
		var serviceInfo ServiceInfo
		if err := json.Unmarshal(kv.Value, &serviceInfo); err != nil {
			log.Printf("解析服务信息失败: %s, err: %v", string(kv.Value), err)
			continue
		}
		sd.services[clusterID] = &serviceInfo
		log.Printf("发现服务: %s -> %+v", clusterID, serviceInfo)
	}

	// 通知服务列表更新
	select {
	case sd.watchChans[clusterID] <- struct{}{}:
	default:
	}

	return nil
}

func (sd *ServiceDiscovery) extractAddrFromKey(key string) string {
	// 键格式: /services/{clusterID}/{addr}
	parts := strings.Split(key, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}

func (sd *ServiceDiscovery) watch(clusterID string) {
	prefix := fmt.Sprintf("/services/%s/", clusterID)
	for {
		watchChan := sd.client.Watch(sd.ctx, prefix, clientv3.WithPrefix())
		for resp := range watchChan {
			for _, event := range resp.Events {
				// addr := sd.extractAddrFromKey(string(event.Kv.Key))
				switch event.Type {
				case clientv3.EventTypePut:
					var serviceInfo ServiceInfo
					if err := json.Unmarshal(event.Kv.Value, &serviceInfo); err != nil {
						log.Printf("解析服务信息失败: %s, err: %v", string(event.Kv.Value), err)
						continue
					}
					sd.mutex.Lock()
					sd.services[clusterID] = &serviceInfo
					sd.mutex.Unlock()
					log.Printf("服务上线: %s -> %+v", clusterID, serviceInfo)

				case clientv3.EventTypeDelete:
					sd.mutex.Lock()
					delete(sd.services, clusterID)
					sd.mutex.Unlock()
					log.Printf("服务下线: %s", clusterID)
				}

				// 通知服务列表更新
				select {
				case sd.watchChans[clusterID] <- struct{}{}:
				default:
				}
			}
		}

		select {
		case <-sd.ctx.Done():
			log.Printf("停止监听服务变化: %s", clusterID)
			return
		default:
		}
	}
}

// GetServiceByAddr 根据地址获取特定服务
func (sd *ServiceDiscovery) GetServiceByClusterID(clusterID string) *ServiceInfo {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()

	return sd.services[clusterID]
}

func (sd *ServiceDiscovery) WatchServices(clusterID string) <-chan struct{} {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	return sd.watchChans[clusterID]
}

func (sd *ServiceDiscovery) Close() error {
	sd.mutex.Lock()
	defer sd.mutex.Unlock()

	sd.cancel()
	for _, clusterid := range sd.clusterIDs {
		close(sd.watchChans[clusterid])
	}
	return sd.client.Close()
}

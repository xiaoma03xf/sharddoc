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
	client     *clientv3.Client
	services   map[string]*ServiceInfo // 所有clusterID对应的Info
	clusterIDs []string
	mutex      sync.RWMutex
	watchChan  chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
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
		watchChan:  make(chan struct{}, 1),
		ctx:        ctx,
		cancel:     cancel,
	}
	for _, cid := range clusterIDs {
		sd.services[cid] = &ServiceInfo{}
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
	for attempt := 1; attempt <= 3; attempt++ {
		resp, err := sd.client.Get(sd.ctx, prefix, clientv3.WithPrefix())
		if err == nil {
			sd.mutex.Lock()
			defer sd.mutex.Unlock()

			// 清空现有服务
			sd.services = make(map[string]*ServiceInfo)

			for _, kv := range resp.Kvs {
				addr := sd.extractAddrFromKey(string(kv.Key))
				var serviceInfo ServiceInfo
				if err := json.Unmarshal(kv.Value, &serviceInfo); err != nil {
					log.Printf("解析服务信息失败: %s, err: %v", string(kv.Value), err)
					continue
				}
				sd.services[addr] = &serviceInfo
				log.Printf("发现服务: %s -> %+v", addr, serviceInfo)
			}

			// 通知服务列表更新
			select {
			case sd.watchChan <- struct{}{}:
			default:
			}
			return nil
		}
		log.Printf("加载服务失败，第%d次重试: %v", attempt, err)
		time.Sleep(time.Second * time.Duration(attempt)) // 指数退避
	}
	return fmt.Errorf("加载服务失败，达到最大重试次数")
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
	prefix := fmt.Sprintf("/services/%s/", sd.clusterIDs)
	for {
		watchChan := sd.client.Watch(sd.ctx, prefix, clientv3.WithPrefix())
		for resp := range watchChan {
			for _, event := range resp.Events {
				addr := sd.extractAddrFromKey(string(event.Kv.Key))
				switch event.Type {
				case clientv3.EventTypePut:
					var serviceInfo ServiceInfo
					if err := json.Unmarshal(event.Kv.Value, &serviceInfo); err != nil {
						log.Printf("解析服务信息失败: %s, err: %v", string(event.Kv.Value), err)
						continue
					}
					sd.mutex.Lock()
					sd.services[addr] = &serviceInfo
					sd.mutex.Unlock()
					log.Printf("服务上线: %s -> %+v", addr, serviceInfo)

				case clientv3.EventTypeDelete:
					sd.mutex.Lock()
					delete(sd.services, addr)
					sd.mutex.Unlock()
					log.Printf("服务下线: %s", addr)
				}

				// 通知服务列表更新
				select {
				case sd.watchChan <- struct{}{}:
				default:
				}
			}
		}

		select {
		case <-sd.ctx.Done():
			log.Printf("停止监听服务变化: %s", clusterID)
			return
		default:
			log.Printf("Watch通道关闭，尝试重新监听: %s", clusterID)
			time.Sleep(time.Second) // 避免快速重试
		}
	}
}

// GetServiceByAddr 根据地址获取特定服务
func (sd *ServiceDiscovery) GetServiceByAddr(clusterID string) *ServiceInfo {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()
	return sd.services[clusterID]
}

func (sd *ServiceDiscovery) WatchServices() <-chan struct{} {
	return sd.watchChan
}

func (sd *ServiceDiscovery) Close() error {
	sd.cancel()
	close(sd.watchChan)
	return sd.client.Close()
}

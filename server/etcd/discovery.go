package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	client    *clientv3.Client
	clusterID string
	services  map[string]*ServiceInfo
	mutex     sync.RWMutex
	watchChan chan struct{}
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewServiceDiscovery(endpoints []string, clusterID string) (*ServiceDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	sd := &ServiceDiscovery{
		client:    client,
		clusterID: clusterID,
		services:  make(map[string]*ServiceInfo),
		watchChan: make(chan struct{}, 1),
		ctx:       ctx,
		cancel:    cancel,
	}

	return sd, nil
}

func (sd *ServiceDiscovery) Start() error {
	if err := sd.loadServices(); err != nil {
		return fmt.Errorf("加载服务失败: %v", err)
	}
	go sd.watch()
	log.Printf("服务发现已启动，监听集群: %s", sd.clusterID)
	return nil
}

func (sd *ServiceDiscovery) loadServices() error {
	prefix := fmt.Sprintf("/services/%s/", sd.clusterID)
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

// extractAddrFromKey 从键中提取地址
func (sd *ServiceDiscovery) extractAddrFromKey(key string) string {
	// 键格式: /services/{clusterID}/{addr}
	parts := strings.Split(key, "/")
	if len(parts) >= 4 {
		return parts[3] // 提取 serviceAddr
	}
	return ""
}

// watch 监听服务变化
func (sd *ServiceDiscovery) watch() {
	prefix := fmt.Sprintf("/services/%s/", sd.clusterID)
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
		// Watch 通道关闭，检查上下文是否取消
		select {
		case <-sd.ctx.Done():
			log.Printf("停止监听服务变化: %s", sd.clusterID)
			return
		default:
			log.Printf("Watch通道关闭，尝试重新监听: %s", sd.clusterID)
			time.Sleep(time.Second) // 避免快速重试
		}
	}
}

// GetServices 获取所有可用服务
func (sd *ServiceDiscovery) GetServices() map[string]*ServiceInfo {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()

	result := make(map[string]*ServiceInfo)
	for k, v := range sd.services {
		result[k] = v
	}
	return result
}

// GetService 获取一个可用服务（随机负载均衡）
func (sd *ServiceDiscovery) GetService() *ServiceInfo {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()

	if len(sd.services) == 0 {
		return nil
	}

	keys := make([]string, 0, len(sd.services))
	for k := range sd.services {
		keys = append(keys, k)
	}
	return sd.services[keys[rand.Intn(len(keys))]]
}

// GetServiceByAddr 根据地址获取特定服务
func (sd *ServiceDiscovery) GetServiceByAddr(addr string) *ServiceInfo {
	sd.mutex.RLock()
	defer sd.mutex.RUnlock()
	return sd.services[addr]
}

// WatchServices 监听服务变化
func (sd *ServiceDiscovery) WatchServices() <-chan struct{} {
	return sd.watchChan
}

// Close 关闭服务发现
func (sd *ServiceDiscovery) Close() error {
	sd.cancel()
	close(sd.watchChan)
	return sd.client.Close()
}

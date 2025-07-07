package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/xiaoma03xf/sharddoc/raft"
	"github.com/xiaoma03xf/sharddoc/raft/etcd"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"google.golang.org/grpc"
)

const (
	CachedTimeOut            = 8 * time.Second
	LeaderHealthCheckTimeOut = 2 * time.Second
)

// 直接以键所在表的主键分片
type DBServer struct {
	mu sync.RWMutex
	// 缓存每一个cluster中leader信息, kvclients 表示和leader的grpc client
	// conn 连接, leaderCache 为缓存的leader信息
	leaders          map[string]*LeaderManager
	serviceDiscovery *etcd.ServiceDiscovery
	cacheTimeout     time.Duration

	metaClusterID string // 元数据, 系统信息放在哪个集群
	stopCh        chan struct{}
}

type LeaderManager struct {
	mu          sync.RWMutex
	GrpcAddress string
	Client      pb.KVStoreClient
	Conn        *grpc.ClientConn
	LastActive  time.Time
}

func NewDB(metaNodeID string, endpoints, clusterAddrs []string) (*DBServer, error) {
	if metaNodeID == "" {
		return nil, fmt.Errorf("metaNodeID is nil")
	}
	if len(clusterAddrs) == 0 {
		return nil, fmt.Errorf("clusterAddrs is nil")
	}
	// grpc serve Discovery
	serverDiscovery, err := etcd.NewServiceDiscovery(endpoints, clusterAddrs)
	if err != nil {
		return nil, err
	}
	if err := serverDiscovery.Start(); err != nil {
		return nil, fmt.Errorf("启动服务发现失败:%w", err)
	}

	db := &DBServer{}
	db.leaders = make(map[string]*LeaderManager)
	db.serviceDiscovery = serverDiscovery
	db.cacheTimeout = CachedTimeOut
	db.metaClusterID = metaNodeID
	db.stopCh = make(chan struct{}, 1)

	for _, addr := range clusterAddrs {
		go db.RunServiceListener(addr)
		go db.StartLeaderHealthCheck(addr, LeaderHealthCheckTimeOut)
	}
	return db, nil
}

func (db *DBServer) getLeader(clusterID string) (pb.KVStoreClient, *grpc.ClientConn, error) {
	client, conn, err := db.getLeaderFromCache(clusterID)
	if err == nil {
		log.Printf("[集群 %s] 缓存获取服务成功！", clusterID)
		return client, conn, nil
	}
	// leader cache 不可用或第一次初始化
	// 如果从 leadercache 中获取失败, 会清理旧连接的
	serveInfo := db.serviceDiscovery.GetServiceByClusterID(clusterID)
	leaderClient, leaderConn, err := raft.BuildGrpcConn(serveInfo.Addr)
	if err != nil {
		log.Printf("[监听器] 创建连接失败: %v", err)
		return nil, nil, err
	}
	leaderManager := &LeaderManager{}
	leaderManager.GrpcAddress = serveInfo.Addr
	leaderManager.Client = leaderClient
	leaderManager.Conn = leaderConn
	leaderManager.LastActive = time.Now()

	db.mu.Lock()
	db.leaders[clusterID] = leaderManager
	db.mu.Unlock()

	return leaderClient, leaderConn, nil
}
func (db *DBServer) RunServiceListener(clusterID string) {
	ch := db.serviceDiscovery.WatchServices(clusterID)
	for {
		select {
		case <-ch:
			// 服务发现器 ServiceDiscovery 发现有变更 update或者delete
			serviceInfo := db.serviceDiscovery.GetServiceByClusterID(clusterID)
			if serviceInfo == nil || serviceInfo.Addr == "" {
				log.Printf("[监听器] %s 当前没有可用服务，跳过连接更新", clusterID)
				continue
			}
			db.mu.Lock()
			current := db.leaders[clusterID]
			if current != nil && current.GrpcAddress == serviceInfo.Addr {
				log.Printf("[监听器] %s 服务地址未变，跳过更新", clusterID)
				db.mu.Unlock()
				continue
			}
			// 关闭旧连接 并且 建立新连接
			if current != nil && current.Conn != nil {
				log.Printf("[监听器] 关闭旧连接: %s", current.GrpcAddress)
				current.Conn.Close()
			}
			leaderClient, leaderConn, err := raft.BuildGrpcConn(serviceInfo.Addr)
			if err != nil {
				log.Printf("[监听器] 创建连接失败: %v", err)
				db.mu.Unlock()
				continue
			}
			db.leaders[clusterID] = &LeaderManager{
				GrpcAddress: serviceInfo.Addr,
				Client:      leaderClient,
				Conn:        leaderConn,
				LastActive:  time.Now(),
			}
			db.mu.Unlock()
			log.Printf("[监听器] %s 已更新为新 Leader: %s", clusterID, serviceInfo.Addr)

		case <-db.stopCh:
			return
		}
	}
}
func (db *DBServer) StartLeaderHealthCheck(clusterID string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-db.stopCh:
			return
		case <-ticker.C:
			db.mu.RLock()
			leader := db.leaders[clusterID]
			db.mu.RUnlock()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := leader.Client.Status(ctx, &pb.StatusRequest{})
			cancel()
			if err == nil {
				log.Printf("[心跳检测] %s ping 成功!", clusterID)

				leader.mu.Lock()
				leader.LastActive = time.Now()
				leader.mu.Unlock()
			} else {
				log.Printf("[心跳检测] %s ping 失败: %v", clusterID, err)
			}
		}
	}
}

func (db *DBServer) getLeaderFromCache(clusterID string) (pb.KVStoreClient, *grpc.ClientConn, error) {
	// 检查缓存是否失效
	db.mu.Lock()
	lm := db.leaders[clusterID]
	db.mu.Unlock()
	if lm == nil {
		return nil, nil, fmt.Errorf("leader [%s] 不存在", clusterID)
	}
	if lm.Client == nil || time.Since(lm.LastActive) > CachedTimeOut {
		log.Printf("[缓存] leader %s 缓存过期，关闭连接", clusterID)
		lm.mu.Lock()
		if lm.Conn != nil {
			_ = lm.Conn.Close()
			lm.Conn = nil
			lm.Client = nil
			lm.GrpcAddress = ""
		}
		lm.mu.Unlock()

		db.mu.Lock()
		delete(db.leaders, clusterID)
		db.mu.Unlock()

		return nil, nil, fmt.Errorf("leader [%s] 缓存过期", clusterID)
	}

	lm.mu.Lock()
	lm.LastActive = time.Now()
	lm.mu.Unlock()
	return lm.Client, lm.Conn, nil
}

func (db *DBServer) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()

	close(db.stopCh)
	for clusterID, leader := range db.leaders {
		leader.mu.Lock()
		if leader.Conn != nil {
			leader.Conn.Close()
		}
		leader.mu.Unlock()
		delete(db.leaders, clusterID)
	}
	err := db.serviceDiscovery.Close()
	if err != nil {
		panic(err)
	}
}

package server

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/xiaoma03xf/sharddoc/lib/hash"
	"github.com/xiaoma03xf/sharddoc/raft"
	"github.com/xiaoma03xf/sharddoc/raft/raftpb"
	"github.com/xiaoma03xf/sharddoc/server/etcd"
	"google.golang.org/grpc"
)

func assert(cond bool) {
	if !cond {
		panic("assertion failure")
	}
}

const (
	CachedTimeOut            = 8 * time.Second
	LeaderHealthCheckTimeOut = 2 * time.Second
)

// 直接以键所在表的主键分片
type DBServer struct {
	mu sync.RWMutex
	// 缓存每一个cluster中leader信息, kvclients 表示和leader的grpc client
	// conn 连接, leaderCache 为缓存的leader信息
	Leaders            map[string]*LeaderManager
	ServiceDiscovery   *etcd.ServiceDiscovery
	TablesDefDiscovery *etcd.TableDefRegistry
	CacheTimeout       time.Duration
	CHash              *hash.Map // key -> cluster
	StopCh             chan struct{}
	SQLParser          *SQLParser // 解释器
}

type LeaderManager struct {
	mu          sync.RWMutex
	GrpcAddress string
	Client      raftpb.KVStoreClient
	Conn        *grpc.ClientConn
	LastActive  time.Time
}

func NewDB(endpoints, clusterAddrs []string) (*DBServer, error) {
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
	// tabledef discovery
	tablesDefDiscovery, err := etcd.NewTableDefRegistry(endpoints)
	if err != nil {
		return nil, err
	}
	cHash := hash.NewMap(3, nil)

	db := &DBServer{}
	db.Leaders = make(map[string]*LeaderManager)
	db.ServiceDiscovery = serverDiscovery
	db.CacheTimeout = CachedTimeOut
	db.TablesDefDiscovery = tablesDefDiscovery
	db.StopCh = make(chan struct{}, 1)
	db.CHash = cHash
	db.SQLParser = new(SQLParser)

	for _, addr := range clusterAddrs {
		go db.RunServiceListener(addr)
		go db.StartLeaderHealthCheck(addr, LeaderHealthCheckTimeOut)
		db.CHash.Add(addr)
	}
	return db, nil
}

func (db *DBServer) getLeader(clusterID string) (raftpb.KVStoreClient, *grpc.ClientConn, error) {
	client, conn, err := db.getLeaderFromCache(clusterID)
	if err == nil {
		return client, conn, nil
	}
	// leader cache 不可用或第一次初始化
	// 如果从 leadercache 中获取失败, 会清理旧连接的
	serveInfo := db.ServiceDiscovery.GetServiceByClusterID(clusterID)
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
	db.Leaders[clusterID] = leaderManager
	db.mu.Unlock()

	return leaderClient, leaderConn, nil
}
func (db *DBServer) RunServiceListener(clusterID string) {
	ch := db.ServiceDiscovery.WatchServices(clusterID)
	for {
		select {
		case <-ch:
			// 服务发现器 ServiceDiscovery 发现有变更 update或者delete
			serviceInfo := db.ServiceDiscovery.GetServiceByClusterID(clusterID)
			if serviceInfo == nil || serviceInfo.Addr == "" {
				log.Printf("[监听器] %s 当前没有可用服务，跳过连接更新", clusterID)
				continue
			}
			db.mu.Lock()
			current := db.Leaders[clusterID]
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
			db.Leaders[clusterID] = &LeaderManager{
				GrpcAddress: serviceInfo.Addr,
				Client:      leaderClient,
				Conn:        leaderConn,
				LastActive:  time.Now(),
			}
			db.mu.Unlock()
			log.Printf("[监听器] %s 已更新为新 Leader: %s", clusterID, serviceInfo.Addr)

		case <-db.StopCh:
			return
		}
	}
}
func (db *DBServer) StartLeaderHealthCheck(clusterID string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-db.StopCh:
			return
		case <-ticker.C:
			db.mu.RLock()
			leader, ok := db.Leaders[clusterID]
			db.mu.RUnlock()
			if !ok || leader == nil || leader.Client == nil {
				log.Printf("[心跳检测] %s 无法执行 ping，因为 leader 或 client 为 nil", clusterID)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := leader.Client.Status(ctx, &raftpb.StatusRequest{})
			cancel()
			if err == nil {
				leader.mu.Lock()
				leader.LastActive = time.Now()
				leader.mu.Unlock()
			} else {
				log.Printf("[心跳检测] %s ping 失败: %v", clusterID, err)
			}
		}
	}
}

func (db *DBServer) getLeaderFromCache(clusterID string) (raftpb.KVStoreClient, *grpc.ClientConn, error) {
	// 检查缓存是否失效
	db.mu.Lock()
	lm := db.Leaders[clusterID]
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
		delete(db.Leaders, clusterID)
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

	close(db.StopCh)
	for clusterID, leader := range db.Leaders {
		leader.mu.Lock()
		if leader.Conn != nil {
			leader.Conn.Close()
		}
		leader.mu.Unlock()
		delete(db.Leaders, clusterID)
	}
	err := db.ServiceDiscovery.Close()
	if err != nil {
		panic(err)
	}
}

func (db *DBServer) getGrpcClientForPrimaryKey(primaryKey []string) (raftpb.KVStoreClient, *grpc.ClientConn, error) {
	key := ""
	for _, k := range primaryKey {
		key += k
	}
	cluster := db.CHash.Get(key)
	return db.getLeader(cluster)
}

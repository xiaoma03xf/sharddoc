package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xiaoma03xf/sharddoc/lib/logger"
	"github.com/xiaoma03xf/sharddoc/raft"
	"github.com/xiaoma03xf/sharddoc/raft/pb"
	"google.golang.org/grpc"
)

const (
	CachedTimeOut = 5 * time.Minute
)

type DB struct {
	mu sync.RWMutex
	// cluster1 ->[node1, node2, node2]
	// cluster2 ->[node1, node2, node2]
	clusterAddrs map[string][]ClusterNode

	// 缓存每一个cluster中leader信息, kvclients 表示和leader的grpc client
	// conn 连接, leaderCache 为缓存的leader信息
	leaders      map[string]*LeaderManager
	cacheTimeout time.Duration

	metaClusterID string // 元数据, 系统信息放在哪个集群
	stopCh        chan struct{}
}
type ClusterNode struct {
	ID       string
	RaftAddr string
	GrpcAddr string
}
type LeaderManager struct {
	mu          sync.RWMutex
	ID          string
	GrpcAddress string
	Client      pb.KVStoreClient
	Conn        *grpc.ClientConn
	ChangeChan  chan struct{}
	LastActive  time.Time
}

func NewDB(metaNodeID string, clusterAddrs map[string][]ClusterNode) (*DB, error) {
	if metaNodeID == "" {
		return nil, fmt.Errorf("metaNodeID is nil")
	}
	if len(clusterAddrs) == 0 {
		return nil, fmt.Errorf("clusterAddrs is nil")
	}
	if _, ok := clusterAddrs[metaNodeID]; !ok {
		return nil, fmt.Errorf("metaClusterID %s is not contained in clusterAddrs", metaNodeID)
	}
	db := &DB{
		clusterAddrs:  make(map[string][]ClusterNode),
		leaders:       make(map[string]*LeaderManager),
		cacheTimeout:  CachedTimeOut,
		metaClusterID: metaNodeID,
		stopCh:        make(chan struct{}, 1),
	}
	return db, nil
}

func (db *DB) Close() {
	db.mu.Lock()
	defer db.mu.Unlock()
	close(db.stopCh)
	for clusterID, leader := range db.leaders {
		leader.mu.Lock()
		if leader.Conn != nil {
			leader.Conn.Close()
		}
		close(leader.ChangeChan) // 统一关闭 ChangeChan
		leader.mu.Unlock()
		delete(db.leaders, clusterID)
	}
}
func (db *DB) watchClusterChanges(clusterID string) {
	for {
		select {
		case <-db.stopCh:
			return
		case <-db.leaders[clusterID].ChangeChan:
			// 集群成员变更，重新查询 leader
		}
	}
}
func (db *DB) getLeader(clusterID string) (pb.KVStoreClient, *grpc.ClientConn, error) {
	client, conn, err := db.getLeaderFromCache(clusterID)
	if err == nil {
		return client, conn, nil
	}
	db.mu.RLock()
	nodes, ok := db.clusterAddrs[clusterID]
	db.mu.RUnlock()
	if !ok || len(nodes) == 0 {
		return nil, nil, fmt.Errorf("集群 %s 无可用节点", clusterID)
	}

	// 获取leader信息
	leaderManager := db.syncGetClusterLeader(clusterID)
	if leaderManager == nil {
		return nil, nil, fmt.Errorf("无法找到集群%s的leader", clusterID)
	}

	var leaderClient pb.KVStoreClient
	var leaderConn *grpc.ClientConn
	if !containsNode(nodes, leaderManager.ID) {
		// leader 不在已知节点列表，触发成员变更通知
		if lm, ok := db.leaders[clusterID]; ok {
			lm.mu.Lock()
			select {
			case lm.ChangeChan <- struct{}{}:
			default:
			}
			lm.mu.Unlock()
		}
		return nil, nil, fmt.Errorf("leader %s not in known nodes for cluster %s", leaderManager.ID, clusterID)
	}
	leaderClient = leaderManager.Client
	leaderConn = leaderManager.Conn

	// 5. 更新 leader 缓存
	db.mu.Lock()
	cachedleader, exists := db.leaders[clusterID]
	if !exists {
		cachedleader = &LeaderManager{
			ChangeChan: make(chan struct{}, 1),
		}
		db.leaders[clusterID] = cachedleader
		go db.watchClusterChanges(clusterID)
	}
	cachedleader.mu.Lock()
	if cachedleader.Conn != nil && cachedleader.Conn != leaderConn {
		cachedleader.Conn.Close()
	}
	cachedleader.Client = leaderClient
	cachedleader.Conn = leaderConn
	cachedleader.ID = leaderManager.ID
	cachedleader.GrpcAddress = leaderManager.GrpcAddress
	cachedleader.LastActive = time.Now()
	cachedleader.mu.Unlock()
	db.mu.Unlock()

	return leaderClient, leaderConn, nil
}
func (db *DB) getLeaderFromCache(clusterID string) (pb.KVStoreClient, *grpc.ClientConn, error) {
	// 检查缓存是否失效
	var cachedleader *LeaderManager
	db.mu.Lock()
	if lm, ok := db.leaders[clusterID]; ok && lm.Client != nil && time.Since(lm.LastActive) < db.cacheTimeout {
		cachedleader = lm
	}
	db.mu.Unlock()
	if cachedleader == nil {
		return nil, nil, fmt.Errorf("cache Leader is nil")
	}

	// 快速验证缓存的 leader
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := cachedleader.Client.Status(ctx, &pb.StatusRequest{})
	if err == nil && resp.Leader.Id == cachedleader.ID {
		go func() {
			cachedleader.mu.Lock()
			cachedleader.LastActive = time.Now()
			cachedleader.mu.Unlock()
		}()
		return cachedleader.Client, cachedleader.Conn, nil
	}
	// 缓存失效，清理旧连接
	cachedleader.mu.Lock()
	if cachedleader.Conn != nil {
		cachedleader.Conn.Close()
		cachedleader.Client = nil
		cachedleader.Conn = nil
		cachedleader.ID = ""
		cachedleader.GrpcAddress = ""
		close(cachedleader.ChangeChan)
	}
	cachedleader.mu.Unlock()
	return nil, nil, fmt.Errorf("cached leader invalid for cluster %s", clusterID)
}
func (db *DB) syncGetClusterLeader(clusterID string) *LeaderManager {
	// 当前ClusterID集群下所有节点
	nodes := db.clusterAddrs[clusterID]
	results := make(chan *pb.StatusResponse, len(nodes))
	go func() {
		var wg sync.WaitGroup
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, node := range nodes {
			wg.Add(1)
			go func(nodeinfo ClusterNode) {
				defer wg.Done()
				db.askStatus(nodeinfo.GrpcAddr, ctx, results)
			}(node)
		}
		wg.Wait()
		close(results)
	}()

	// 收集 leader 信息
	var leaderInfo *pb.Node
	var leaderClient pb.KVStoreClient
	var leaderConn *grpc.ClientConn
	for resp := range results {
		if resp != nil && resp.Leader != nil {
			if leaderInfo == nil || resp.Leader.Id == leaderInfo.Id {
				leaderInfo = resp.Leader
				// 获取 Client 和 Conn
				client, conn, err := raft.BuildGrpcConn(leaderInfo.Grpcaddress)
				if err != nil {
					logger.Warn(fmt.Errorf("连接 leader %s 失败: %v", leaderInfo.Grpcaddress, err))
					continue
				}
				leaderClient = client
				leaderConn = conn
			}
		}
	}
	if leaderInfo == nil {
		return nil
	}
	return &LeaderManager{
		ID:          leaderInfo.Id,
		GrpcAddress: leaderInfo.Grpcaddress,
		Client:      leaderClient,
		Conn:        leaderConn,
		ChangeChan:  make(chan struct{}, 1),
		LastActive:  time.Now(),
	}
}

// askStatus 访问任意一个raft集群中节点,当前status
func (db *DB) askStatus(grpcaddr string, ctx context.Context, res chan<- *pb.StatusResponse) {
	client, conn, err := raft.BuildGrpcConn(grpcaddr)
	defer conn.Close()
	if err != nil {
		logger.Warn(fmt.Errorf("连接节点 %s 失败 :%v", grpcaddr, err))
		return
	}
	r, err := client.Status(ctx, &pb.StatusRequest{})
	if err != nil {
		logger.Warn(fmt.Errorf("Status %s 请求失败 :%v", grpcaddr, err))
		return
	}
	res <- r
}
func containsNode(nodes []ClusterNode, id string) bool {
	for _, node := range nodes {
		if node.ID == id {
			return true
		}
	}
	return false
}

package tcc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/xiaoma03xf/sharddoc/tcc/log"
)

type TXManager struct {
	ctx            context.Context
	stop           context.CancelFunc
	opts           *Options
	txStore        TXStore         //内置的事务日志存储模块，需要由使用方实现并完成注入
	registryCenter *registryCenter //TCC 组件的注册管理中心
}

func NewTXManager(txStore TXStore, opts ...Option) *TXManager {
	ctx, cancel := context.WithCancel(context.Background())
	txmanager := &TXManager{}
	txmanager.opts = &Options{}
	txmanager.txStore = txStore
	txmanager.registryCenter = newRegistryCenter()
	txmanager.ctx = ctx
	txmanager.stop = cancel

	for _, opt := range opts {
		opt(txmanager.opts)
	}
	repair(txmanager.opts)

	go txmanager.run()
	return txmanager
}

func (t *TXManager) Stop() {
	t.stop()
}

func (t *TXManager) Register(component TccComponent) error {
	return t.registryCenter.register(component)
}

// 启动事务
func (t *TXManager) Transaction(ctx context.Context, reqs ...*RequestEntity) (string, bool, error) {
	// 1 限制分布式事务执行时长
	tctx, cancel := context.WithTimeout(ctx, t.opts.Timeout)
	defer cancel()

	// 2 获得所有的涉及使用的tcc组件
	componentEntities, err := t.getComponents(tctx, reqs...)
	if err != nil {
		return "", false, err
	}

	// 3 调用 txStore 模块, 创建新的事务明细记录, 并取得全局唯一的事务 id
	txID, err := t.txStore.CreateTX(tctx, componentEntities.ToComponents()...)
	if err != nil {
		return "", false, nil
	}

	// 4 两阶段提交, try-confirm/cancel
	return txID, t.twoPhaseCommit(ctx, txID, componentEntities), nil
}

func (t *TXManager) backOffTick(tick time.Duration) time.Duration {
	tick <<= 1 // tick = tick * 2
	// 就是 t.opts.MonitorTick 乘以 2^3 = 8
	if threshold := t.opts.MonitorTick << 3; tick > threshold {
		return threshold
	}
	return tick
}

func (t *TXManager) run() {
	var tick time.Duration
	var err error
	for {
		// 如果出现了失败, tick 需要避让， 遵循退避策略增大 tick 间隔时长
		if err != nil {
			tick = t.opts.MonitorTick
		} else {
			tick = t.backOffTick(tick)
		}
		select {
		//倘若 txManager.ctx 被终止，则异步轮询任务退出
		case <-t.ctx.Done():
			return
		//等待tick对应时长后，开始执行任务
		case <-time.After(tick):
			// 加锁，避免多个分布式多个节点的监控任务重复执行
			if err := t.txStore.Lock(t.ctx, t.opts.MonitorTick); err != nil {
				// 取消失败时（大概率被其他节点占有）, 不对tick进行退避升级
				err = nil
				continue
			}

			// 获取任然处于 hanging 状态的事务
			var txs []*Transaction
			if txs, err = t.txStore.GetHangingTXs(t.ctx); err != nil {
				_ = t.txStore.Unlock(t.ctx)
				continue
			}
			//批量推进事务进度
			err = t.batchAdvanceProgress(txs)
			_ = t.txStore.Unlock(t.ctx)
		}
	}
}

func (t *TXManager) batchAdvanceProgress(txs []*Transaction) error {
	// 对每笔事务进行状态推进
	errCh := make(chan error)
	go func() {
		// 并发执行, 推进各比事务的进度
		var wg sync.WaitGroup
		for _, tx := range txs {
			// shadow
			tx := tx
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 每个goroutine 负责处理一笔事务
				if err := t.advanceProgress(tx); err != nil {
					// 遇到错误则投递到 errch
					errCh <- err
				}
			}()
		}
		wg.Wait()
		close(errCh)
	}()

	var firstErr error
	// 通过 chan 阻塞在这里, 直到所有 goroutine 执行完成, chan 被 close 才能往下
	for err := range errCh {
		// 记录遇到的第一个错误
		if firstErr != nil {
			continue
		}
		firstErr = err
	}
	return firstErr
}

// 传入一个事务id推进进度
func (t *TXManager) advanceProgressByTXID(txID string) error {
	// 获取事务日志记录
	tx, err := t.txStore.GetTX(t.ctx, txID)
	if err != nil {
		return err
	}
	// 推进进度
	return t.advanceProgress(tx)
}

// 传入一个事务id推进其进度
func (t *TXManager) advanceProgress(tx *Transaction) error {
	// 根据各个 component try 请求的情况, 推断出事务当前的状态
	// 1.1 倘若所有组件 try 都成功, 则为 successful
	// 1.2 倘若所有组件 try 失败, 则为 failure
	// 1.3 倘若事务超时了， 则为failure
	// 1.4 否则事务为hanging
	txStatus := tx.getStatus(time.Now().Add(-t.opts.Timeout))
	// hanging 状态的暂时不处理
	if txStatus == TXHanging {
		return nil
	}

	// 根据事务是否成功, 定制不同的处理函数
	success := txStatus == TXSuccessful
	var confirmOrCancel func(ctx context.Context, component TccComponent) (*TCCResp, error)
	var txAdvanceProgress func(ctx context.Context) error
	if success {
		confirmOrCancel = func(ctx context.Context, component TccComponent) (*TCCResp, error) {
			// 对 component 进行第二阶段的 confirm 操作
			return component.Confirm(ctx, tx.TXID)
		}
		txAdvanceProgress = func(ctx context.Context) error {
			// 更新事务日志记录的状态为成功
			return t.txStore.TXSubmit(ctx, tx.TXID, true)
		}
	} else {
		confirmOrCancel = func(ctx context.Context, component TccComponent) (*TCCResp, error) {
			// 对 component 进行第二阶段的 cancel 操作
			return component.Cancel(ctx, tx.TXID)
		}

		txAdvanceProgress = func(ctx context.Context) error {
			// 更新事务日志记录的状态为失败
			return t.txStore.TXSubmit(ctx, tx.TXID, false)
		}
	}

	for _, component := range tx.Components {
		// 获取对应的 tcc component
		components, err := t.registryCenter.getComponents(component.ComponentID)
		if err != nil || len(components) == 0 {
			return errors.New("get tcc component faild")
		}
		// 执行二阶段的confirm 或者cancel 操作
		resp, err := confirmOrCancel(t.ctx, components[0])
		if err != nil {
			return err
		}
		if !resp.ACK {
			return fmt.Errorf("component: %s ack failed", component.ComponentID)
		}
	}

	// 两阶段都执行完成后,对事务状态进行提交
	return txAdvanceProgress(t.ctx)
}

func (t *TXManager) twoPhaseCommit(ctx context.Context, txID string, componentEntities ComponentEntities) bool {
	// 1 创建子 context 用于管理子goroutine生命周期
	// 手握 cancel 终止器, 能保证在需要的时候终止所有子 goroutine 生命周期
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 2 创建一个chan,用于接收子 goroutine 传递的错误
	errCh := make(chan error)
	// 3 并发启动, 批量执行各 tcc 组件的 try 流程
	go func() {
		// 通过 waitGroup 进行多个子 goroutine 的汇总
		var wg sync.WaitGroup
		for _, componentEntity := range componentEntities {
			// shadow
			componentEntity := componentEntity
			wg.Add(1)
			// 并发执行各组件的 try 流程
			go func() {
				defer wg.Done()
				resp, err := componentEntity.Component.Try(cctx, &TCCReq{
					ComponentID: componentEntity.Component.ID(),
					TXID:        txID,
					Data:        componentEntity.Request,
				})
				// 但凡有一个 component try报错或者拒绝， 都是需要进行cancel的，但会放在 advanceProgressByTXID 流程处理
				if err != nil || !resp.ACK {
					log.ErrorContextf(cctx, "tx try failed, tx id: %s, comonent id: %s, err: %v", txID, componentEntity.Component.ID(), err)
					// 对对应的事务进行更新
					if _err := t.txStore.TXUpdate(cctx, txID, componentEntity.Component.ID(), false); _err != nil {
						log.ErrorContextf(cctx, "tx updated failed, tx id: %s, component id: %s, err: %v", txID, componentEntity.Component.ID(), _err)
					}
					errCh <- fmt.Errorf("component: %s try failed", componentEntity.Component.ID())
					return
				}
				// try 请求成功, 但是请求结果更新到事务日志失败时,也需要视为处理失败
				if err = t.txStore.TXUpdate(cctx, txID, componentEntity.Component.ID(), true); err != nil {
					log.ErrorContextf(ctx, "tx updated failed, tx id: %s, component id: %s, err: %v", txID, componentEntity.Component.ID(), err)
					errCh <- err
				}
			}()
		}
		wg.Wait()
		close(errCh)
	}()

	successful := true
	if err := <-errCh; err != nil {
		// 只要有一笔 try 请求出现问题, 其他的都进行终止
		cancel()
		successful = false
	}

	// 执行二阶段，即便第二阶段执行失败也无妨，可以通过轮询任务进行兜底处理
	if err := t.advanceProgressByTXID(txID); err != nil {
		log.ErrorContextf(ctx, "advance tx progress fail, txid: %s, err: %v", txID, err)
	}
	return successful
}

func (t *TXManager) getComponents(ctx context.Context, reqs ...*RequestEntity) (ComponentEntities, error) {
	if len(reqs) == 0 {
		return nil, errors.New("emtpy task")
	}
	// 调一下接口， 确认这些都是合法的
	idToReq := make(map[string]*RequestEntity, len(reqs))
	componentIDs := make([]string, 0, len(reqs))
	for _, req := range reqs {
		if _, ok := idToReq[req.ComponentID]; ok {
			return nil, fmt.Errorf("repeat component: %s", req.ComponentID)
		}
		idToReq[req.ComponentID] = req
		componentIDs = append(componentIDs, req.ComponentID)
	}

	// 校验合法性
	components, err := t.registryCenter.getComponents(componentIDs...)
	if err != nil {
		return nil, err
	}
	if len(componentIDs) != len(components) {
		return nil, errors.New("invalid componentIDs")
	}

	entities := make(ComponentEntities, 0, len(components))
	for _, component := range components {
		entities = append(entities, &ComponentEntity{
			Request:   idToReq[component.ID()].Request,
			Component: component,
		})
	}
	return entities, nil
}

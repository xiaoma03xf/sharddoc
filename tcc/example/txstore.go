package example

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/demdxx/gocast"
	"github.com/xiaoma03xf/sharddoc/lib/redislock"
	"github.com/xiaoma03xf/sharddoc/tcc"
	expdao "github.com/xiaoma03xf/sharddoc/tcc/example/dao"
	"github.com/xiaoma03xf/sharddoc/tcc/example/pkg"
)

type MockTXStore struct {
	client *redislock.Client
	dao    TXRecordDAO
}

func NewMockTXStore(dao TXRecordDAO, client *redislock.Client) *MockTXStore {
	return &MockTXStore{
		dao:    dao,
		client: client,
	}
}

func (m *MockTXStore) CreateTX(ctx context.Context, components ...tcc.TccComponent) (string, error) {
	// 创建一项内容，里面以唯一事务 id 为 key
	componentTryStatuses := make(map[string]*expdao.ComponentTryStatus, len(components))
	for _, component := range components {
		componentTryStatuses[component.ID()] = &expdao.ComponentTryStatus{
			ComponentID: component.ID(),
			TryStatus:   tcc.TryHanging.String(),
		}
	}

	statusesBody, _ := json.Marshal(componentTryStatuses)
	txID, err := m.dao.CreateTXRecord(ctx, &expdao.TXRecordPO{
		Status:               tcc.TXHanging.String(),
		ComponentTryStatuses: string(statusesBody),
	})
	if err != nil {
		return "", err
	}

	return gocast.ToString(txID), nil
}

func (m *MockTXStore) TXUpdate(ctx context.Context, txID string, componentID string, accept bool) error {
	_txID := gocast.ToUint(txID)
	status := tcc.TXFailure.String()
	if accept {
		status = tcc.TXSuccessful.String()
	}
	return m.dao.UpdateComponentStatus(ctx, _txID, componentID, status)
}

func (m *MockTXStore) GetHangingTXs(ctx context.Context) ([]*tcc.Transaction, error) {
	records, err := m.dao.GetTXRecords(ctx, expdao.WithStatus(tcc.TryHanging))
	if err != nil {
		return nil, err
	}

	txs := make([]*tcc.Transaction, 0, len(records))
	for _, record := range records {
		componentTryStatuses := make(map[string]*expdao.ComponentTryStatus)
		_ = json.Unmarshal([]byte(record.ComponentTryStatuses), &componentTryStatuses)
		components := make([]*tcc.ComponentTryEntity, 0, len(componentTryStatuses))
		for _, component := range componentTryStatuses {
			components = append(components, &tcc.ComponentTryEntity{
				ComponentID: component.ComponentID,
				TryStatus:   tcc.ComponentTryStatus(component.TryStatus),
			})
		}

		txs = append(txs, &tcc.Transaction{
			TXID:       gocast.ToString(record.ID),
			Status:     tcc.TXHanging,
			CreatedAt:  record.CreatedAt,
			Components: components,
		})
	}

	return txs, nil
}

func (m *MockTXStore) Lock(ctx context.Context, expireDuration time.Duration) error {
	lock := redislock.NewRedisLock(pkg.BuildTXRecordLockKey(), m.client, redislock.WithExpireSeconds(int64(expireDuration.Seconds())))
	return lock.Lock(ctx)
}

func (m *MockTXStore) Unlock(ctx context.Context) error {
	lock := redislock.NewRedisLock(pkg.BuildTXRecordLockKey(), m.client)
	return lock.Unlock(ctx)
}

// 提交事务的最终状态
func (m *MockTXStore) TXSubmit(ctx context.Context, txID string, success bool) error {
	do := func(ctx context.Context, dao *expdao.TXRecordDAO, record *expdao.TXRecordPO) error {
		if success {
			if record.Status == tcc.TXFailure.String() {
				return fmt.Errorf("invalid tx status: %s, txid: %s", record.Status, txID)
			}
			record.Status = tcc.TXSuccessful.String()
		} else {
			if record.Status == tcc.TXSuccessful.String() {
				return fmt.Errorf("invalid tx status: %s, txid: %s", record.Status, txID)
			}
			record.Status = tcc.TXFailure.String()
		}
		return dao.UpdateTXRecord(ctx, record)
	}
	return m.dao.LockAndDo(ctx, gocast.ToUint(txID), do)
}

// 获取指定的一笔事务
func (m *MockTXStore) GetTX(ctx context.Context, txID string) (*tcc.Transaction, error) {
	records, err := m.dao.GetTXRecords(ctx, expdao.WithID(gocast.ToUint(txID)))
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, errors.New("get tx failed")
	}

	componentTryStatuses := make(map[string]*expdao.ComponentTryStatus)
	_ = json.Unmarshal([]byte(records[0].ComponentTryStatuses), &componentTryStatuses)

	components := make([]*tcc.ComponentTryEntity, 0, len(componentTryStatuses))
	for _, tryItem := range componentTryStatuses {
		components = append(components, &tcc.ComponentTryEntity{
			ComponentID: tryItem.ComponentID,
			TryStatus:   tcc.ComponentTryStatus(tryItem.TryStatus),
		})
	}
	return &tcc.Transaction{
		TXID:       txID,
		Status:     tcc.TXStatus(records[0].Status),
		Components: components,
		CreatedAt:  records[0].CreatedAt,
	}, nil
}

type TXRecordDAO interface {
	GetTXRecords(ctx context.Context, opts ...expdao.QueryOption) ([]*expdao.TXRecordPO, error)
	CreateTXRecord(ctx context.Context, record *expdao.TXRecordPO) (uint, error)
	UpdateComponentStatus(ctx context.Context, id uint, componentID string, status string) error
	UpdateTXRecord(ctx context.Context, record *expdao.TXRecordPO) error
	LockAndDo(ctx context.Context, id uint, do func(ctx context.Context, dao *expdao.TXRecordDAO, record *expdao.TXRecordPO) error) error
}

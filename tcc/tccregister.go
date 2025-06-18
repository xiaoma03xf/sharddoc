package tcc

import (
	"errors"
	"fmt"
	"sync"
)

type registryCenter struct {
	mux        sync.RWMutex
	components map[string]TccComponent
}

func newRegistryCenter() *registryCenter {
	return &registryCenter{
		components: make(map[string]TccComponent),
	}
}

func (r *registryCenter) register(component TccComponent) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	if _, ok := r.components[component.ID()]; ok {
		return errors.New("repeat component id")
	}
	r.components[component.ID()] = component
	return nil
}

// 上游 TXManager 可以通过 component id，进行 TCC 组件的查询.
// 倘若某个 component id 不存在，则会抛出错误
func (r *registryCenter) getComponents(componentIDs ...string) ([]TccComponent, error) {
	components := make([]TccComponent, 0, len(componentIDs))

	r.mux.RLock()
	defer r.mux.RUnlock()

	for _, componentID := range componentIDs {
		component, ok := r.components[componentID]
		if !ok {
			return nil, fmt.Errorf("component id :%s not existed", componentID)
		}
		components = append(components, component)
	}
	return components, nil
}

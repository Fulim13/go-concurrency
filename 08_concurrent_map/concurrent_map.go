package concurrentmap

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewConcurrentMap[K comparable, V any](cap int) *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data: make(map[K]V, cap),
	}
}

func (cm *ConcurrentMap[K, V]) Store(key K, value V) {
	// the reason use *ConccurrentMap is that we also need to change the state of the RLock and RUnlock
	// AND we need to change the cm.data property
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.data[key] = value
}

func (cm *ConcurrentMap[K, V]) Load(key K) (value V, exists bool) {
	// the reason use *ConccurrentMap is that we also need to change the state of the RLock and RUnlock
	// Where there is write lock, block it until release, so it can read the updated value
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	value, exists = cm.data[key]
	return
}

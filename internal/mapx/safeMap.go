package mapx

import (
	"sync"
)

type ISafeMapKey interface {
}

// SafeMap provides a map alternative to avoid memory leak.
// This implementation is not needed until issue below fixed.
// https://github.com/golang/go/issues/20135

var defaultMaxDeleteCount int64 = 1000

func SetDefaultMaxDeleteCount(newDefaultMaxDeleteCount int64) {
	defaultMaxDeleteCount = newDefaultMaxDeleteCount
}

type SafeMap[Key comparable, Value any] struct {
	maxDeleteCount int64
	lock           sync.RWMutex
	m              map[Key]Value
	deleteCount    int64
}

func NewSafeMap[Key comparable, Value any](opts ...SafeMapOption) *SafeMap[Key, Value] {
	// handle options
	var opt safeMapOption
	for _, o := range opts {
		o(&opt)
	}

	maxDeleteCount := opt.maxDeleteCount
	if maxDeleteCount == 0 {
		maxDeleteCount = defaultMaxDeleteCount
	}

	return &SafeMap[Key, Value]{
		maxDeleteCount: maxDeleteCount,
		m:              make(map[Key]Value),
	}
}

// Set 插入或更新
func (s *SafeMap[Key, Value]) Set(key Key, val Value) {
	s.lock.Lock()
	s.m[key] = val
	s.lock.Unlock()
}

// Get 读取
func (s *SafeMap[Key, Value]) Get(key Key) (Value, bool) {
	s.lock.RLock()
	v, ok := s.m[key]
	s.lock.RUnlock()
	return v, ok
}

// Del 删除，并在阈值时重建底层 map
func (s *SafeMap[Key, Value]) Del(key Key) {
	s.lock.Lock()
	// 删除并计数
	if _, ok := s.m[key]; ok {
		delete(s.m, key)
		s.deleteCount++
	}
	// 超过阈值，做一次 shrink
	if s.deleteCount >= s.maxDeleteCount {
		newMap := make(map[Key]Value, len(s.m))
		for k, v := range s.m {
			newMap[k] = v
		}
		s.m = newMap
		s.deleteCount = 0
	}
	s.lock.Unlock()
}

// Size 返回当前元素数量
func (s *SafeMap[Key, Value]) Size() int {
	s.lock.RLock()
	n := len(s.m)
	s.lock.RUnlock()
	return n
}

// Range 迭代所有 kv
func (s *SafeMap[Key, Value]) Range(f func(key Key, val Value) bool) {
	s.lock.RLock()
	for k, v := range s.m {
		if !f(k, v) {
			break
		}
	}
	s.lock.RUnlock()
}

func (s *SafeMap[Key, Value]) GetOrSet(key Key, init func() Value) (val Value, alreadyInit bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if v, ok := s.m[key]; ok {
		return v, true
	}
	s.m[key] = init()
	return s.m[key], false
}

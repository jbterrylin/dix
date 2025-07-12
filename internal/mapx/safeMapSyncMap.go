package mapx

// type ISafeMapKey interface {
// }

// // SafeMap provides a map alternative to avoid memory leak.
// // This implementation is not needed until issue below fixed.
// // https://github.com/golang/go/issues/20135

// var defaultMaxDeleteCount int64 = 1000

// func SetDefaultMaxDeleteCount(newDefaultMaxDeleteCount int64) {
// 	defaultMaxDeleteCount = newDefaultMaxDeleteCount
// }

// type SafeMap[Key comparable, Value any] struct {
// 	maxDeleteCount int64
// 	m              *sync.Map
// 	deleteCount    int64
// }

// func NewSafeMap[Key comparable, Value any](opts ...SafeMapOption) *SafeMap[Key, Value] {
// 	// handle options
// 	var opt safeMapOption
// 	for _, o := range opts {
// 		o(&opt)
// 	}

// 	maxDeleteCount := opt.maxDeleteCount
// 	if maxDeleteCount == 0 {
// 		maxDeleteCount = defaultMaxDeleteCount
// 	}

// 	return &SafeMap[Key, Value]{
// 		maxDeleteCount: maxDeleteCount,
// 		m:              new(sync.Map),
// 	}
// }

// // Set 插入或更新
// func (s *SafeMap[Key, Value]) Set(key Key, val Value) {
// 	s.m.Store(key, val)
// }

// // Get 读取
// func (s *SafeMap[Key, Value]) Get(key Key) (Value, bool) {
// 	v, ok := s.m.Load(key)
// 	if !ok {
// 		var zero Value
// 		return zero, ok
// 	}
// 	return v.(Value), ok
// }

// // Del 删除，并在阈值时重建底层 map
// func (s *SafeMap[Key, Value]) Del(key Key) {
// 	// 删除并计数
// 	if _, ok := s.m.Load(key); ok {
// 		s.m.Delete(key)
// 		s.deleteCount++
// 	}
// 	// 超过阈值，做一次 shrink
// 	if s.deleteCount >= s.maxDeleteCount {
// 		var newMap sync.Map
// 		s.m.Range(func(key, value any) bool {
// 			newMap.Store(key, value)
// 			return true
// 		})
// 		s.m = &newMap
// 		s.deleteCount = 0
// 	}
// }

// // Size 返回当前元素数量
// func (s *SafeMap[Key, Value]) Size() int {
// 	n := 0
// 	s.m.Range(func(_, _ any) bool {
// 		n++
// 		return true
// 	})
// 	return n
// }

// // Range 迭代所有 kv
// func (s *SafeMap[Key, Value]) Range(f func(key Key, val Value) bool) {
// 	s.m.Range(func(key, value any) bool {
// 		return f(key.(Key), value.(Value))
// 	})
// }

// func (s *SafeMap[Key, Value]) GetOrSet(key Key, init func() Value) (val Value, alreadyInit bool) {
// 	tmp, loaded := s.m.LoadOrStore(key, init())
// 	val = tmp.(Value)
// 	return val, loaded
// }

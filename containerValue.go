package dix

import (
	"sync"
	"sync/atomic"
	"time"
)

var _ iContainerData = &containerValue{}

type containerValue struct {
	mu          sync.RWMutex
	value       any
	onCloseHook func()
	isAccessed  atomic.Bool

	refCounter     int64
	refCounterCond *sync.Cond

	createdAt  time.Time
	accessedAt atomic.Int64
	tagMap     map[string]any
}

func newContainerValue(
	value any,
	onCloseHook func(),
	tagMap map[string]any,
) *containerValue {
	return &containerValue{
		value:       value,
		onCloseHook: onCloseHook,

		refCounterCond: sync.NewCond(&sync.Mutex{}),

		createdAt: time.Now(),
		tagMap:    tagMap,
	}
}

func (c *containerValue) setAccessed() (isFirstAccess bool) {
	isFirstAccess = c.isAccessed.CompareAndSwap(false, true)
	if isFirstAccess {
		c.accessedAt.Store(time.Now().UnixMicro())
	}
	return
}

func (c *containerValue) lock() {
	c.mu.Lock()
}

func (c *containerValue) unlock() {
	c.mu.Unlock()
}

func (c *containerValue) triggerOnCloseHook() {
	if c.onCloseHook != nil {
		c.waitUntilRefZero()
		c.onCloseHook()
	}
}

func (c *containerValue) refCounterIncr() {
	if !globalContainer.safeDelete {
		return
	}
	atomic.AddInt64(&c.refCounter, 1)
}

func (c *containerValue) refCounterDecr() {
	if !globalContainer.safeDelete {
		return
	}
	newRefCounter := atomic.AddInt64(&c.refCounter, -1)
	if newRefCounter < 0 {
		panic(ErrRefCounterBelowZero)
	}
	if newRefCounter == 0 {
		c.refCounterCond.L.Lock()
		c.refCounterCond.Broadcast()
		c.refCounterCond.L.Unlock()
	}
}

func (c *containerValue) waitUntilRefZero() {
	if !globalContainer.safeDelete {
		return
	}
	c.refCounterCond.L.Lock()
	for atomic.LoadInt64(&c.refCounter) > 0 {
		c.refCounterCond.Wait()
	}
	c.refCounterCond.L.Unlock()
}

// func (c *containerValue) GetValue() any             { return c.value }
// func (c *containerValue) GetOnCloseHook() func()    { return c.onCloseHook }
func (c *containerValue) GetIsAccessed() bool       { return c.isAccessed.Load() }
func (c *containerValue) GetRefCounter() int64      { return atomic.LoadInt64(&c.refCounter) }
func (c *containerValue) GetCreatedAt() time.Time   { return c.createdAt }
func (c *containerValue) GetAccessedAt() time.Time  { return time.UnixMicro(c.accessedAt.Load()) }
func (c *containerValue) GetTagMap() map[string]any { return copyMap(c.tagMap) }

func (c *containerValue) OnCloseHookExist() bool {
	return c.onCloseHook != nil
}
func (c *containerValue) TriggerOnCloseHook() {
	if c.onCloseHook != nil {
		c.onCloseHook()
	}
}

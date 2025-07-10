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
	isAccessed  bool

	refCounter     int64
	refCounterCond *sync.Cond

	createdAt  time.Time
	accessedAt time.Time
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

func (c *containerValue) setAccessed() {
	c.isAccessed = true
	c.accessedAt = time.Now()
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
func (c *containerValue) GetIsAccessed() bool       { return c.isAccessed }
func (c *containerValue) GetRefCounter() int64      { return c.refCounter }
func (c *containerValue) GetCreatedAt() time.Time   { return c.createdAt }
func (c *containerValue) GetAccessedAt() time.Time  { return c.accessedAt }
func (c *containerValue) GetTagMap() map[string]any { return copyMap(c.tagMap) }

func (c *containerValue) OnCloseHookExist() bool {
	return c.onCloseHook != nil
}
func (c *containerValue) TriggerOnCloseHook() {
	if c.onCloseHook != nil {
		c.onCloseHook()
	}
}

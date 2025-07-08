package dix

import (
	"context"
	"sync"
	"time"
)

var _ iContainerData = &containerProvider{}

type containerProvider struct {
	mu             sync.RWMutex
	value          func() (any, error)
	valueWithCtx   func(context.Context) (any, error)
	isValueWithCtx bool
	noCache        bool
	cacheValue     any
	isAccessed     bool

	createdAt  time.Time
	accessedAt time.Time
	tagMap     map[string]any
}

func newContainerProvider(
	value func() (any, error),
	noCache bool,
	tagMap map[string]any,
) *containerProvider {
	return &containerProvider{
		value:     value,
		noCache:   noCache,
		createdAt: time.Now(),
		tagMap:    tagMap,
	}
}

func newCtxContainerProvider(
	valueWithCtx func(context.Context) (any, error),
	noCache bool,
	tagMap map[string]any,
) *containerProvider {
	return &containerProvider{
		valueWithCtx:   valueWithCtx,
		isValueWithCtx: true,
		noCache:        noCache,
		createdAt:      time.Now(),
		tagMap:         tagMap,
	}
}

func (c *containerProvider) setAccessed() {
	c.isAccessed = true
	c.accessedAt = time.Now()
}

func (c *containerProvider) lock() {
	c.mu.Lock()
}

func (c *containerProvider) unlock() {
	c.mu.Unlock()
}

func (c *containerProvider) triggerOnCloseHook() {
}

// func (c *containerProvider) GetValue() func() (any, error) { return c.value }
//
//	func (c *containerProvider) GetValueWithCtx() func(context.Context) (any, error) {
//		return c.valueWithCtx
//	}
func (c *containerProvider) GetIsValueWithCtx() bool   { return c.isValueWithCtx }
func (c *containerProvider) GetNoCache() bool          { return c.noCache }
func (c *containerProvider) GetCacheValue() any        { return c.cacheValue }
func (c *containerProvider) GetIsAccessed() bool       { return c.isAccessed }
func (c *containerProvider) GetCreatedAt() time.Time   { return c.createdAt }
func (c *containerProvider) GetAccessedAt() time.Time  { return c.accessedAt }
func (c *containerProvider) GetTagMap() map[string]any { return copyMap(c.tagMap) }

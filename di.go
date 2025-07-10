package dix

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/jbterrylin/dix/internal/mapx"
)

func SetDefaultValueKey(defaultValueKey string) {
	DefaultValueKey = ValueKey(defaultValueKey)
}

func SetDefaultProviderKey(defaultProviderKey string) {
	DefaultProviderKey = ProviderKey(defaultProviderKey)
}

func SetSafeDelete(safeDelete bool) {
	globalContainer.safeDelete = safeDelete
}

func SetResetMaxConcurrent(resetMaxConcurrent int) {
	if resetMaxConcurrent == 0 {
		globalContainer.resetMaxConcurrent = 100
		return
	}
	globalContainer.resetMaxConcurrent = resetMaxConcurrent
}

var globalContainer = newContainer()

type (
	container struct {
		typeKeyValueMap    *mapx.SafeMap[reflect.Type, *mapx.SafeMap[ValueKey, *containerValue]]
		typeKeyProviderMap *mapx.SafeMap[reflect.Type, *mapx.SafeMap[ProviderKey, *containerProvider]]

		afterAdd                AfterAddFunc
		afterProviderRun        AfterProviderRunFunc
		afterFirstAccess        AfterFirstAccessFunc
		beforeDuplicateRegister BeforeDuplicateRegisterFunc

		safeDelete bool

		resetMaxConcurrent int
	}
)

func newContainer() *container {
	return &container{
		typeKeyValueMap:    mapx.NewSafeMap[reflect.Type, *mapx.SafeMap[ValueKey, *containerValue]](),
		typeKeyProviderMap: mapx.NewSafeMap[reflect.Type, *mapx.SafeMap[ProviderKey, *containerProvider]](),

		resetMaxConcurrent: 100,
	}
}

func Reset(opts ...ResetOption) []error {
	// handle options
	var opt resetOption
	for _, o := range opts {
		o(&opt)
	}

	errs := reset(opt.skipOnClose, globalContainer.typeKeyValueMap, DefaultValueKey)
	errs = append(errs, reset(opt.skipOnClose, globalContainer.typeKeyProviderMap, DefaultProviderKey)...)

	return errs
}

func reset[Key ~string, Value iContainerData](
	skipOnClose bool,
	typeKeyValueMap *mapx.SafeMap[reflect.Type, *mapx.SafeMap[Key, Value]],
	defaultKey Key,
) []error {
	typeKeysMap := getTypeKeysMap(typeKeyValueMap)

	var (
		errs     = make([]error, 0)
		errsLock sync.Mutex
		wg       sync.WaitGroup
	)

	sem := make(chan struct{}, globalContainer.resetMaxConcurrent)

	for typ, keys := range typeKeysMap {
		for _, key := range keys {
			wg.Add(1)
			go func(typ reflect.Type, key Key) {
				defer wg.Done()

				// 获取并发许可（阻塞）
				sem <- struct{}{}
				defer func() { <-sem }()

				val, err := getContainerNestedMapValue(typeKeyValueMap, typ, key)
				if err != nil {
					errsLock.Lock()
					errs = append(errs, fmt.Errorf("failed at type=%v, key=%v: %w", typ, key, err))
					errsLock.Unlock()
					return
				}

				// delete first so others won't see it
				_, _ = deleteContainerNestedMapValue(typeKeyValueMap, typ, key)

				val.lock()
				defer val.unlock()

				if !skipOnClose && key != defaultKey {
					val.triggerOnCloseHook()
				}
			}(typ, key)
		}
	}

	wg.Wait()
	return errs
}

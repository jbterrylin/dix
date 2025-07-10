package dix

import (
	"context"
	"reflect"
)

func AddProvider[T any](key ProviderKey, value func() (T, error), opts ...ProviderAddOption) error {
	if value == nil {
		return ErrValueIsNil
	}
	return addProvider(key, value, nil, opts...)
}

func AddCtxProvider[T any](key ProviderKey, valueWithCtx func(context.Context) (T, error), opts ...ProviderAddOption) error {
	if valueWithCtx == nil {
		return ErrValueIsNil
	}
	return addProvider(key, nil, valueWithCtx, opts...)
}

func addProvider[T any](key ProviderKey, value func() (T, error), valueWithCtx func(context.Context) (T, error), opts ...ProviderAddOption) error {
	if !key.IsValid() {
		return ErrInvalidKey
	}

	// handle options
	var opt providerAddOption
	for _, o := range opts {
		o(&opt)
	}

	t := reflect.TypeOf((*T)(nil)).Elem()

	var tmp *containerProvider
	if value != nil {
		tmp = newContainerProvider(
			func() (any, error) {
				return value()
			},
			opt.noCache,
			opt.tagMap,
		)
	}
	if valueWithCtx != nil {
		tmp = newCtxContainerProvider(
			func(ctx context.Context) (any, error) {
				return valueWithCtx(ctx)
			},
			opt.noCache,
			opt.tagMap,
		)
	}

	oldValue, _ := getContainerNestedMapValue(globalContainer.typeKeyProviderMap, t, key)
	if oldValue != nil && globalContainer.beforeDuplicateRegister != nil {
		oldValue.mu.RLock()
		err := globalContainer.beforeDuplicateRegister(NewBeforeDuplicateRegisterCtx(t, nil, nil, nil, &key, oldValue, tmp, false))
		oldValue.mu.RUnlock()
		if err != nil {
			return err
		}
	}

	if opt.setDefault {
		oldValue, _ := getContainerNestedMapValue(globalContainer.typeKeyProviderMap, t, DefaultProviderKey)
		if oldValue != nil && globalContainer.beforeDuplicateRegister != nil {
			oldValue.mu.RLock()
			err := globalContainer.beforeDuplicateRegister(NewBeforeDuplicateRegisterCtx(t, nil, nil, nil, &key, oldValue, tmp, true))
			oldValue.mu.RUnlock()
			if err != nil {
				return err
			}
		}

		setValueToContainerNestedMap(globalContainer.typeKeyProviderMap, t, DefaultProviderKey, tmp)
	}

	// to avoid replace DefaultProviderKey value fail but key value get set
	setValueToContainerNestedMap(globalContainer.typeKeyProviderMap, t, key, tmp)

	if globalContainer.afterAdd != nil {
		globalContainer.afterAdd(NewAfterAddCtx(t, nil, nil, &key, tmp))
	}

	return nil
}

func GetProvider[T any](opts ...ProviderGetOption) (T, error) {
	return GetProviderByKey[T](DefaultProviderKey, opts...)
}

func GetProviderWithCtx[T any](ctx context.Context, opts ...ProviderGetOption) (T, error) {
	return GetProviderByKeyWithCtx[T](ctx, DefaultProviderKey, opts...)
}

func GetProviderByKey[T any](key ProviderKey, opts ...ProviderGetOption) (T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	_, val, err := getProviderByTypeKey(context.Background(), t, key, opts...)
	if err != nil {
		var zero T
		return zero, err
	}

	return val.(T), nil
}

func GetProviderByKeyWithCtx[T any](ctx context.Context, key ProviderKey, opts ...ProviderGetOption) (T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	_, val, err := getProviderByTypeKey(ctx, t, key, opts...)
	if err != nil {
		var zero T
		return zero, err
	}

	return val.(T), nil
}

func getProviderByTypeKey(ctx context.Context, t reflect.Type, key ProviderKey, opts ...ProviderGetOption) (*containerProvider, any, error) {
	// handle options
	var opt providerGetOption
	for _, o := range opts {
		o(&opt)
	}

	provider, err := getContainerNestedMapValue(globalContainer.typeKeyProviderMap, t, key)
	if err != nil {
		return nil, nil, err
	}

	provider.mu.Lock()
	defer provider.mu.Unlock()

	if !opt.reload && provider.cacheValue != nil {
		return provider, provider.cacheValue, nil
	}

	done := make(chan struct{})
	var tmp any

	go func() {
		if provider.isValueWithCtx {
			tmp, err = provider.valueWithCtx(ctx)
		} else {
			tmp, err = provider.value()
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err() // timeout, canceled
	case <-done:
		// continue
	}
	if err != nil {
		return nil, nil, err
	}

	if !provider.noCache {
		provider.cacheValue = tmp
	}

	isFirstAccess := false
	if !provider.isAccessed {
		isFirstAccess = true
		provider.setAccessed()
	}

	if globalContainer.afterProviderRun != nil {
		globalContainer.afterProviderRun(NewAfterProviderRunCtx(t, key, provider, tmp))
	}
	if isFirstAccess && globalContainer.afterFirstAccess != nil {
		globalContainer.afterFirstAccess(NewAfterFirstAccessCtx(t, nil, nil, &key, provider))
	}

	return provider, tmp, nil
}

func MustGetProvider[T any](opts ...ProviderGetOption) T {
	return MustGetProviderByKey[T](DefaultProviderKey, opts...)
}

func MustGetProviderWithCtx[T any](ctx context.Context, opts ...ProviderGetOption) T {
	return MustGetProviderByKeyWithCtx[T](ctx, DefaultProviderKey, opts...)
}

func MustGetProviderByKey[T any](key ProviderKey, opts ...ProviderGetOption) T {
	v, err := GetProviderByKey[T](key, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

func MustGetProviderByKeyWithCtx[T any](ctx context.Context, key ProviderKey, opts ...ProviderGetOption) T {
	v, err := GetProviderByKeyWithCtx[T](ctx, key, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

func ProviderExist[T any]() bool {
	return ProviderExistByKey[T](DefaultProviderKey)
}

func ProviderExistByKey[T any](key ProviderKey) bool {
	_, err := GetProviderByKey[T](key)
	return err == nil
}

func DeleteProvider[T any]() {
	DeleteProviderByKey[T](DefaultProviderKey)
}

func DeleteProviderByKey[T any](key ProviderKey) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	_, err := deleteContainerNestedMapValue(globalContainer.typeKeyProviderMap, t, key)
	if err != nil {
		return err
	}
	return nil
}

func ListProviderKeys[T any]() []ProviderKey {
	t := reflect.TypeOf((*T)(nil)).Elem()
	keyProviderMap, exist := globalContainer.typeKeyProviderMap.Get(t)
	if !exist {
		return []ProviderKey{}
	}
	keys := make([]ProviderKey, 0, keyProviderMap.Size())
	keyProviderMap.Range(func(key ProviderKey, val *containerProvider) bool {
		if key == DefaultProviderKey {
			return true
		}

		keys = append(keys, key)
		return true
	})
	return keys
}

func GetAllProvider[T any](opts ...ProviderGetOption) ([]T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	keyProviderMap, exist := globalContainer.typeKeyProviderMap.Get(t)
	if !exist {
		return []T{}, nil
	}
	var err error
	values := make([]T, 0, keyProviderMap.Size())
	keyProviderMap.Range(func(key ProviderKey, val *containerProvider) bool {
		if key == DefaultProviderKey {
			return true
		}

		var tmp T
		tmp, err = GetProviderByKey[T](key, opts...)
		if err != nil {
			return false
		}
		values = append(values, tmp)
		return true
	})
	if err != nil {
		return []T{}, nil
	}
	return values, nil
}

package dix

import (
	"reflect"
)

func Add[T any](key ValueKey, val T, opts ...ValueAddOption) error {
	if !key.IsValid() {
		return ErrInvalidKey
	}

	// handle options
	var opt valueAddOption
	for _, o := range opts {
		o(&opt)
	}

	t := reflect.TypeOf((*T)(nil)).Elem()

	tmp := newContainerValue(val, opt.onCloseHook, opt.tagMap)

	oldValue, _ := getContainerNestedMapValue(globalContainer.typeKeyValueMap, t, key)
	if oldValue != nil && globalContainer.beforeDuplicateRegister != nil {
		oldValue.mu.RLock()
		err := globalContainer.beforeDuplicateRegister(NewBeforeDuplicateRegisterCtx(t, &key, oldValue, tmp, nil, nil, nil, false))
		oldValue.mu.RUnlock()
		if err != nil {
			return err
		}
	}

	if opt.setDefault {
		oldValue, _ := getContainerNestedMapValue(globalContainer.typeKeyValueMap, t, DefaultValueKey)
		if oldValue != nil && globalContainer.beforeDuplicateRegister != nil {
			oldValue.mu.RLock()
			err := globalContainer.beforeDuplicateRegister(NewBeforeDuplicateRegisterCtx(t, &key, oldValue, tmp, nil, nil, nil, true))
			oldValue.mu.RUnlock()
			if err != nil {
				return err
			}
		}

		setValueToContainerNestedMap(globalContainer.typeKeyValueMap, t, DefaultValueKey, tmp)
	}

	// to avoid replace DefaultValueKey value fail but key value get set
	setValueToContainerNestedMap(globalContainer.typeKeyValueMap, t, key, tmp)

	if globalContainer.afterAdd != nil {
		globalContainer.afterAdd(NewAfterAddCtx(t, &key, tmp, nil, nil))
	}

	return nil
}

func Get[T any]() (T, error) {
	return GetByKey[T](DefaultValueKey)
}

func GetByKey[T any](key ValueKey) (T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()

	tmp, err := getByTypeKey(t, key)
	if err != nil {
		var zero T
		return zero, err
	}

	return tmp.value.(T), nil
}

func getByTypeKey(t reflect.Type, key ValueKey) (*containerValue, error) {
	val, err := getContainerNestedMapValue(globalContainer.typeKeyValueMap, t, key)
	if err != nil {
		return nil, err
	}

	val.mu.Lock()
	defer val.mu.Unlock()

	if val.isAccessed {
		val.refCounterIncr()
		return val, nil
	}

	val.setAccessed()

	if globalContainer.afterFirstAccess != nil {
		globalContainer.afterFirstAccess(NewAfterFirstAccessCtx(t, &key, val, nil, nil))
	}

	val.refCounterIncr()
	return val, nil
}

func MustGet[T any]() T {
	return MustGetByKey[T](DefaultValueKey)
}

func MustGetByKey[T any](key ValueKey) T {
	v, err := GetByKey[T](key)
	if err != nil {
		panic(err)
	}
	return v
}

func Exist[T any]() bool {
	return ExistByKey[T](DefaultValueKey)
}

func ExistByKey[T any](key ValueKey) bool {
	_, err := GetByKey[T](key)
	return err == nil
}

func Delete[T any](opts ...ValueDeleteOption) {
	DeleteByKey[T](DefaultValueKey, opts...)
}

func DeleteByKey[T any](key ValueKey, opts ...ValueDeleteOption) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	return deleteByTypeKey(t, key, opts...)
}

func deleteByTypeKey(t reflect.Type, key ValueKey, opts ...ValueDeleteOption) error {
	// handle options
	var opt valueDeleteOption
	for _, o := range opts {
		o(&opt)
	}

	// delete first so others won't see it
	value, err := deleteContainerNestedMapValue(globalContainer.typeKeyValueMap, t, key)
	if err != nil {
		return err
	}
	value.mu.Lock()
	defer value.mu.Unlock()
	if !opt.skipOnClose {
		value.triggerOnCloseHook()
	}
	return nil
}

func ListKeys[T any]() []ValueKey {
	t := reflect.TypeOf((*T)(nil)).Elem()
	keyValueMap, exist := globalContainer.typeKeyValueMap.Get(t)
	if !exist {
		return []ValueKey{}
	}
	keys := make([]ValueKey, 0, keyValueMap.Size())
	keyValueMap.Range(func(key ValueKey, val *containerValue) bool {
		if key == DefaultValueKey {
			return true
		}

		keys = append(keys, key)
		return true
	})
	return keys
}

func GetAll[T any]() []T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	keyValueMap, exist := globalContainer.typeKeyValueMap.Get(t)
	if !exist {
		return []T{}
	}
	values := make([]T, 0, keyValueMap.Size())
	keyValueMap.Range(func(key ValueKey, val *containerValue) bool {
		if key == DefaultValueKey {
			return true
		}

		values = append(values, val.value.(T))
		return true
	})
	return values
}

func DeductRefCount[T any]() error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	val, err := getContainerNestedMapValue(globalContainer.typeKeyValueMap, t, DefaultValueKey)
	if err != nil {
		return err
	}

	val.refCounterDecr()
	return nil
}

func DeductRefCountByKey[T any](key ValueKey) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	val, err := getContainerNestedMapValue(globalContainer.typeKeyValueMap, t, key)
	if err != nil {
		return err
	}

	val.refCounterDecr()
	return nil
}

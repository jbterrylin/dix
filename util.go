package dix

import (
	"reflect"

	"github.com/jbterrylin/dix/internal/mapx"
)

func getContainerNestedMapValue[Key ~string, Value any](
	typeKeyValueMap *mapx.SafeMap[reflect.Type, *mapx.SafeMap[Key, Value]],
	t reflect.Type,
	key Key,
) (Value, error) {
	keyValueMap, exist := typeKeyValueMap.Get(t)
	if !exist {
		var zero Value
		return zero, ErrValueNotFound
	}
	val, exist := keyValueMap.Get(key)
	if !exist {
		var zero Value
		return zero, ErrValueNotFound
	}

	return val, nil
}

func setValueToContainerNestedMap[Key ~string, Value any](
	typeKeyValueMap *mapx.SafeMap[reflect.Type, *mapx.SafeMap[Key, Value]],
	t reflect.Type,
	key Key,
	value Value,
) {
	keyValueMap, _ := typeKeyValueMap.GetOrSet(t, func() *mapx.SafeMap[Key, Value] {
		return mapx.NewSafeMap[Key, Value]()
	})
	keyValueMap.Set(key, value)
}

func deleteContainerNestedMapValue[Key ~string, Value any](
	typeKeyValueMap *mapx.SafeMap[reflect.Type, *mapx.SafeMap[Key, Value]],
	t reflect.Type,
	key Key,
) (Value, error) {
	keyValueMap, exist := typeKeyValueMap.Get(t)
	if !exist {
		var zero Value
		return zero, ErrValueNotFound
	}
	val, exist := keyValueMap.Get(key)
	if !exist {
		var zero Value
		return zero, ErrValueNotFound
	}
	keyValueMap.Del(key)

	return val, nil
}

func getTypeKeysMap[Key ~string, Value iContainerData](
	typeKeyValueMap *mapx.SafeMap[reflect.Type, *mapx.SafeMap[Key, Value]],
) map[reflect.Type][]Key {
	typeKeysMap := make(map[reflect.Type][]Key, typeKeyValueMap.Size())
	typeKeyValueMap.Range(func(typ reflect.Type, keyValueMap *mapx.SafeMap[Key, Value]) bool {
		keys := make([]Key, 0, keyValueMap.Size())
		keyValueMap.Range(func(key Key, _ Value) bool {
			keys = append(keys, key)
			return true
		})
		typeKeysMap[typ] = keys
		return true
	})
	return typeKeysMap
}

func copyMap[Key comparable, Value any](tmp map[Key]Value) map[Key]Value {
	copyMap := make(map[Key]Value, len(tmp))
	for k, v := range tmp {
		copyMap[k] = v
	}
	return copyMap
}

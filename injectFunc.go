package dix

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// If the function returns two or more values and the last one is an error, it will be used as the returned error.
func InjectFunc(fn any, opts ...InjectFuncOption) error {
	return InjectFuncWithCtx(context.Background(), fn, opts...)
}

// If the function returns two or more values and the last one is an error, it will be used as the returned error.
func InjectFuncWithCtx(ctx context.Context, fn any, opts ...InjectFuncOption) error {
	v := reflect.ValueOf(fn)
	t := v.Type()

	if t.Kind() != reflect.Func {
		return ErrInjectFuncMustBeFunc
	}

	variableNameIndexMap := make(map[string]int, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		variableNameIndexMap[paramType.Name()] = i
	}

	optMap := make(map[int]*injectFuncOption)
	for _, o := range opts {
		tmp := &injectFuncOption{}
		o(tmp)

		index := -1
		if i, err := strconv.Atoi(tmp.variable); err == nil {
			if i > t.NumIn()-1 || i < 0 {
				return ErrInvalidVariable
			}
			index = i
		} else if idx, ok := variableNameIndexMap[tmp.variable]; ok {
			index = idx
		} else {
			return ErrInvalidVariable
		}

		if _, exist := optMap[int(index)]; !exist {
			optMap[int(index)] = tmp
		} else {
			mergeInjectFuncOpt(optMap[int(index)], tmp)
		}
	}

	in := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		var tag injectTag
		if opt, exist := optMap[i]; exist {
			tag = newInjectTag(opt.valType, opt.key, opt.reload, opt.optional)
		}

		var err error

		var tmp *reflect.Value
		switch tag.valType {
		case injectTagFlagTypeOptProvider.Value():
			tmp, err = getFromProvider(ctx, paramType, tag)
		default:
			tmp, err = getFromValue(paramType, tag)
		}

		if err != nil {
			if errors.Is(err, ErrValueNotFound) && tag.optional {
				continue
			}
			return fmt.Errorf("failed at param type name=%v, param type=%v: %w", paramType.Name(), paramType.String(), err)
		}

		in[i] = *tmp
	}

	v.Call(in)

	return nil
}

func mergeInjectFuncOpt(dst, src *injectFuncOption) {
	if src.valType != "" {
		dst.valType = src.valType
	}
	if src.key != "" {
		dst.key = src.key
	}
	if src.reload {
		dst.reload = true
	}
	if src.optional {
		dst.optional = true
	}
}

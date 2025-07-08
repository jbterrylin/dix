package dix

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const injectStructTag = "di"

var injectTagFlagType injectTagFlag = "type"
var injectTagFlagKey injectTagFlag = "key"
var injectTagFlagReload injectTagFlag = "reload"
var injectTagFlagOptional injectTagFlag = "optional"

var injectTagFlagTypeOptProvider injectTagFlagTypeOpt = "provider"

type injectTag struct {
	valType  string
	key      string
	reload   bool
	optional bool
}

func newInjectTag(valType string, key string, reload bool, optional bool) injectTag {
	return injectTag{
		valType:  valType,
		key:      key,
		reload:   reload,
		optional: optional,
	}
}

func InjectStruct(target any) error {
	return InjectStructWithCtx(context.Background(), target)
}

func InjectStructWithCtx(ctx context.Context, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return ErrInjectStructMustBePointerStruct
	}

	t := v.Elem().Type()
	v = v.Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(injectStructTag)
		if tag == "-" {
			continue
		}

		// parse tag
		injectTag := parseDITag(tag)

		fieldVal := v.Field(i)
		if !fieldVal.CanSet() {
			return fmt.Errorf("failed at field name=%v, field type=%v: %w", field.Name, field.Type, ErrFieldCannotBeSet)
		}

		typ := field.Type

		// determine injection path
		var err error

		var tmp *reflect.Value
		switch injectTag.valType {
		case injectTagFlagTypeOptProvider.Value():
			tmp, err = getFromProvider(ctx, typ, injectTag)
		default:
			tmp, err = getFromValue(typ, injectTag)
		}

		if err != nil {
			if errors.Is(err, ErrValueNotFound) && injectTag.optional {
				continue
			}
			return fmt.Errorf("failed at field name=%v, field type=%v: %w", field.Name, field.Type, err)
		}
		fieldVal.Set(*tmp)
	}
	return nil
}

func parseDITag(tag string) injectTag {
	opts := make(map[string]string)
	for _, part := range strings.Split(tag, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			opts[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else {
			switch part {
			case injectTagFlagReload.Value(),
				injectTagFlagOptional.Value():
				opts[part] = "true"
			}
		}
	}

	return newInjectTag(
		opts[injectTagFlagType.Value()],
		opts[injectTagFlagKey.Value()],
		strings.ToLower(opts[injectTagFlagReload.Value()]) == "true",
		strings.ToLower(opts[injectTagFlagOptional.Value()]) == "true",
	)
}

func getFromProvider(ctx context.Context, typ reflect.Type, opts injectTag) (*reflect.Value, error) {
	key := ProviderKey(opts.key)
	reload := opts.reload

	providerGetOptions := []ProviderGetOption{}

	if reload {
		providerGetOptions = append(providerGetOptions, WithProviderReload())
	}

	_, val, err := getProviderByTypeKey(ctx, typ, key, providerGetOptions...)
	if err != nil {
		return nil, err
	}

	tmp := reflect.ValueOf(val)
	return &tmp, nil
}

func getFromValue(typ reflect.Type, opts injectTag) (*reflect.Value, error) {
	key := ValueKey(opts.key)
	value, err := getByTypeKey(typ, key)
	if err != nil {
		return nil, err
	}
	tmp := reflect.ValueOf(value.value)
	return &tmp, nil
}

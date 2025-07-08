package dix

import "reflect"

type (
	AfterAddCtx struct {
		Type              reflect.Type
		ValueKey          *ValueKey
		ContainerValue    *containerValue
		ProviderKey       *ProviderKey
		ContainerProvider *containerProvider
	}

	AfterAddFunc func(ctx AfterAddCtx)

	AfterProviderRunCtx struct {
		Type              reflect.Type
		Key               ProviderKey
		ContainerProvider *containerProvider
		Value             any
	}

	AfterProviderRunFunc func(ctx AfterProviderRunCtx)

	AfterFirstAccessCtx struct {
		Type              reflect.Type
		ValueKey          *ValueKey
		ContainerValue    *containerValue
		ProviderKey       *ProviderKey
		ContainerProvider *containerProvider
	}

	AfterFirstAccessFunc func(ctx AfterFirstAccessCtx)

	BeforeDuplicateRegisterCtx struct {
		Type                 reflect.Type
		ValueKey             *ValueKey
		OldContainerValue    *containerValue
		NewContainerValue    *containerValue
		ProviderKey          *ProviderKey
		OldContainerProvider *containerProvider
		NewContainerProvider *containerProvider
		IsDefault            bool
	}

	BeforeDuplicateRegisterFunc func(ctx BeforeDuplicateRegisterCtx) error
)

func NewAfterAddCtx(
	typ reflect.Type,
	valueKey *ValueKey, containerValue *containerValue,
	providerKey *ProviderKey, containerProvider *containerProvider,
) AfterAddCtx {
	return AfterAddCtx{
		Type:              typ,
		ValueKey:          valueKey,
		ContainerValue:    containerValue,
		ProviderKey:       providerKey,
		ContainerProvider: containerProvider,
	}
}

func AfterAdd(f AfterAddFunc) {
	Container.afterAdd = f
}

func NewAfterProviderRunCtx(
	typ reflect.Type,
	key ProviderKey, containerProvider *containerProvider,
	value any,
) AfterProviderRunCtx {
	return AfterProviderRunCtx{
		Type:              typ,
		Key:               key,
		ContainerProvider: containerProvider,
		Value:             value,
	}
}

func AfterProviderRun(f AfterProviderRunFunc) {
	Container.afterProviderRun = f
}

func NewAfterFirstAccessCtx(
	typ reflect.Type,
	valueKey *ValueKey, containerValue *containerValue,
	providerKey *ProviderKey, containerProvider *containerProvider,
) AfterFirstAccessCtx {
	return AfterFirstAccessCtx{
		Type:              typ,
		ValueKey:          valueKey,
		ContainerValue:    containerValue,
		ProviderKey:       providerKey,
		ContainerProvider: containerProvider,
	}
}

func AfterFirstAccess(f AfterFirstAccessFunc) {
	Container.afterFirstAccess = f
}

func NewBeforeDuplicateRegisterCtx(
	typ reflect.Type,
	valueKey *ValueKey, oldContainerValue *containerValue, newContainerValue *containerValue,
	providerKey *ProviderKey, oldContainerProvider *containerProvider, newContainerProvider *containerProvider,
	isDefault bool,
) BeforeDuplicateRegisterCtx {
	return BeforeDuplicateRegisterCtx{
		Type:                 typ,
		ValueKey:             valueKey,
		OldContainerValue:    oldContainerValue,
		NewContainerValue:    newContainerValue,
		ProviderKey:          providerKey,
		OldContainerProvider: oldContainerProvider,
		NewContainerProvider: newContainerProvider,
		IsDefault:            isDefault,
	}
}

func BeforeDuplicateRegister(f BeforeDuplicateRegisterFunc) {
	Container.beforeDuplicateRegister = f
}

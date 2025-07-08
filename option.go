package dix

type valueAddOption struct {
	onCloseHook func()
	setDefault  bool
	tagMap      map[string]any
}

type ValueAddOption func(*valueAddOption)

func WithValueOnClose(f func()) ValueAddOption {
	return func(o *valueAddOption) {
		o.onCloseHook = f
	}
}

func WithValueSetDefault() ValueAddOption {
	return func(o *valueAddOption) {
		o.setDefault = true
	}
}

func WithValueTag(tagMap map[string]any) ValueAddOption {
	return func(o *valueAddOption) {
		o.tagMap = tagMap
	}
}

type valueDeleteOption struct {
	skipOnClose bool
}

type ValueDeleteOption func(*valueDeleteOption)

func WithValueSkipOnClose() ValueDeleteOption {
	return func(o *valueDeleteOption) {
		o.skipOnClose = true
	}
}

type providerAddOption struct {
	setDefault bool
	noCache    bool
	tagMap     map[string]any
}

type ProviderAddOption func(*providerAddOption)

func WithProviderSetDefault() ProviderAddOption {
	return func(o *providerAddOption) {
		o.setDefault = true
	}
}

func WithProviderNoCache() ProviderAddOption {
	return func(o *providerAddOption) {
		o.noCache = true
	}
}

func WithProviderTag(tagMap map[string]any) ProviderAddOption {
	return func(o *providerAddOption) {
		o.tagMap = tagMap
	}
}

type providerGetOption struct {
	reload bool
}

type ProviderGetOption func(*providerGetOption)

func WithProviderReload() ProviderGetOption {
	return func(o *providerGetOption) {
		o.reload = true
	}
}

type injectFuncOption struct {
	variable string // can be variable name / index

	valType  string // provider / value
	key      string
	reload   bool
	optional bool
}

type InjectFuncOption func(*injectFuncOption)

func WithInjectFuncProvider(variable string) InjectFuncOption {
	return func(o *injectFuncOption) {
		o.variable = variable
		o.valType = "provider"
	}
}

func WithInjectFuncKey(variable string, key string) InjectFuncOption {
	return func(o *injectFuncOption) {
		o.variable = variable
		o.key = key
	}
}

func WithInjectFuncReload(variable string) InjectFuncOption {
	return func(o *injectFuncOption) {
		o.variable = variable
		o.reload = true
	}
}

func WithInjectFuncOptional(variable string) InjectFuncOption {
	return func(o *injectFuncOption) {
		o.variable = variable
		o.optional = true
	}
}

type resetOption struct {
	skipOnClose bool
}

type ResetOption func(*resetOption)

func WithResetSkipOnClose() ResetOption {
	return func(o *resetOption) {
		o.skipOnClose = true
	}
}

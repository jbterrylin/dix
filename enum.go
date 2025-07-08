package dix

type ValueKey string

func (c ValueKey) IsValid() bool {
	return c.Value() != DefaultValueKey.Value()
}

func (c ValueKey) Value() string {
	return string(c)
}

type ProviderKey string

func (c ProviderKey) IsValid() bool {
	return c.Value() != DefaultProviderKey.Value()
}

func (c ProviderKey) Value() string {
	return string(c)
}

type injectTagFlag string

func (c injectTagFlag) Value() string {
	return string(c)
}

type injectTagFlagTypeOpt string

func (c injectTagFlagTypeOpt) Value() string {
	return string(c)
}

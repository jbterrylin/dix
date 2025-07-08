package dix_test

import (
	"github.com/jbterrylin/dix"
)

const TestKey = dix.ValueKey("test")
const TestInterfaceKey = dix.ValueKey("test interface")

const TestProviderKey = dix.ProviderKey("test")

var _ ITestInterface = NewTestInterface("test interface")

type ITestInterface interface {
	Name() string
	Count() int
}

type Test struct {
	count int
	name  string
}

func NewTest(name string) *Test {
	return &Test{
		name: name,
	}
}

func NewTestInterface(name string) ITestInterface {
	return &Test{
		name: name,
	}
}

func (v *Test) Name() string {
	return v.name
}

func (v *Test) Count() int {
	v.count += 1
	return v.count
}

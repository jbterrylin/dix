package dix_test

import (
	"testing"

	"github.com/jbterrylin/dix"
)

func TestInjectStruct(t *testing.T) {
	test := NewTest("test")
	dix.Add(TestKey, test, dix.WithValueSetDefault())

	testInterface := NewTestInterface("test interface")
	dix.Add(TestInterfaceKey, testInterface, dix.WithValueSetDefault())

	testFromC, err := dix.Get[*Test]()
	if err != nil {
		t.Errorf("unexpected Get() err: got %v, want %v", err, nil)
	}

	testInterfaceFromC, err := dix.Get[ITestInterface]()
	if err != nil {
		t.Errorf("unexpected Get() err: got %v, want %v", err, nil)
	}

	testFromC.Count()          // count + 1
	testInterfaceFromC.Count() // count + 1

	var tmp struct {
		Val  *Test
		IVal ITestInterface
	}

	err = dix.InjectStruct(&tmp)
	if err != nil {
		t.Errorf("unexpected InjectStruct() err: got %v, want %v", err, nil)
	}

	count := tmp.Val.Count()
	if count != 2 {
		t.Errorf("unexpected Count(): got %v, want %v", count, 2)
	}

	count = tmp.IVal.Count()
	if count != 2 {
		t.Errorf("unexpected Count(): got %v, want %v", count, 2)
	}
}

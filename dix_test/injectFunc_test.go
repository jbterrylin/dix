package dix_test

import (
	"testing"

	"github.com/jbterrylin/dix"
)

func TestInjectFunc(t *testing.T) {
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

	err = dix.InjectFunc(func(val *Test, iVal ITestInterface) {
		count := val.Count()
		if count != 2 {
			t.Errorf("unexpected Count(): got %v, want %v", count, 2)
		}

		count = iVal.Count()
		if count != 2 {
			t.Errorf("unexpected Count(): got %v, want %v", count, 2)
		}

	})
	if err != nil {
		t.Errorf("unexpected InjectFunc() err: got %v, want %v", err, nil)
	}

}

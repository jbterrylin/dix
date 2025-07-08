package dix_test

import (
	"testing"

	"github.com/jbterrylin/dix"
)

func TestAddProvider(t *testing.T) {
	err := dix.AddProvider(TestProviderKey, func() (*Test, error) {
		return NewTest("test"), nil
	}, dix.WithProviderSetDefault())
	if err != nil {
		t.Errorf("unexpected AddProvider() err: got %v, want %v", err, nil)
	}

	testFromC, err := dix.GetProvider[*Test]()
	if err != nil {
		t.Errorf("unexpected GetProvider() err: got %v, want %v", err, nil)
	}

	name := testFromC.Name()
	count := testFromC.Count()

	if name != "test" {
		t.Errorf("unexpected Name(): got %v, want %v", name, "test")
	}

	if count != 1 {
		t.Errorf("unexpected Count(): got %v, want %v", count, 1)
	}
}

func TestProviderReload(t *testing.T) {
	err := dix.AddProvider(TestProviderKey, func() (*Test, error) {
		return NewTest("test"), nil
	},
		dix.WithProviderSetDefault(),
		dix.WithProviderNoCache(),
	)
	if err != nil {
		t.Errorf("unexpected AddProvider() err: got %v, want %v", err, nil)
	}

	testFromC, err := dix.GetProvider[*Test]()
	if err != nil {
		t.Errorf("unexpected GetProvider() err: got %v, want %v", err, nil)
	}
	testFromC.Count()
	testFromC.Count()
	testFromC.Count()

	testFromC, err = dix.GetProvider[*Test]()
	if err != nil {
		t.Errorf("unexpected GetProvider() err: got %v, want %v", err, nil)
	}
	count := testFromC.Count()

	if count != 1 {
		t.Errorf("unexpected Count(): got %v, want %v", count, 1)
	}
}

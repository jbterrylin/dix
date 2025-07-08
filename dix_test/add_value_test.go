package dix_test

import (
	"errors"
	"testing"

	"github.com/jbterrylin/dix"
)

func TestAddValue(t *testing.T) {
	test := NewTest("test")
	err := dix.Add(TestKey, test, dix.WithValueSetDefault())
	if err != nil {
		t.Errorf("unexpected Add() err: got %v, want %v", err, nil)
	}

	testFromC, err := dix.Get[*Test]()
	if err != nil {
		t.Errorf("unexpected Get() err: got %v, want %v", err, nil)
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

func TestAddInterface(t *testing.T) {
	testInterface := NewTestInterface("test interface")
	dix.Add(TestInterfaceKey, testInterface, dix.WithValueSetDefault())

	testInterfaceFromC, err := dix.Get[ITestInterface]()
	if err != nil {
		t.Errorf("unexpected Get() err: got %v, want %v", err, nil)
	}

	name := testInterfaceFromC.Name()
	count := testInterfaceFromC.Count()

	if name != "test interface" {
		t.Errorf("unexpected Name(): got %v, want %v", name, "test interface")
	}

	if count != 1 {
		t.Errorf("unexpected Count(): got %v, want %v", count, 1)
	}
}

func TestAddBothAndGetCorrect(t *testing.T) {
	test := NewTest("test")
	dix.Add(TestKey, test, dix.WithValueSetDefault())

	testInterface := NewTestInterface("test interface")
	dix.Add(TestInterfaceKey, testInterface, dix.WithValueSetDefault())

	testFromC, err := dix.Get[*Test]()
	if err != nil {
		t.Errorf("unexpected Get() err: got %v, want %v", err, nil)
	}

	name := testFromC.Name()
	count := testFromC.Count()

	if name != "test" {
		t.Errorf("unexpected Name(): got %v, want %v", name, "test")
	}

	if count != 1 {
		t.Errorf("unexpected Count(): got %v, want %v", count, 1)
	}

	testInterfaceFromC, err := dix.Get[ITestInterface]()
	if err != nil {
		t.Errorf("unexpected Get() err: got %v, want %v", err, nil)
	}
	name = testInterfaceFromC.Name()
	count = testInterfaceFromC.Count()

	if name != "test interface" {
		t.Errorf("unexpected Name(): got %v, want %v", name, "test interface")
	}

	if count != 1 {
		t.Errorf("unexpected Count(): got %v, want %v", count, 1)
	}
}

func TestNoAdd(t *testing.T) {
	_, err := dix.Get[*Test]()
	if !errors.Is(err, dix.ErrValueNotFound) {
		t.Errorf("unexpected Get() err: got %v, want %v", err, dix.ErrValueNotFound)
	}

	_, err = dix.Get[ITestInterface]()
	if !errors.Is(err, dix.ErrValueNotFound) {
		t.Errorf("unexpected Get() err: got %v, want %v", err, dix.ErrValueNotFound)
	}
}

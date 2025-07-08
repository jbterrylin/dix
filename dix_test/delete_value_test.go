package dix_test

import (
	"errors"
	"testing"

	"github.com/jbterrylin/dix"
)

func TestDeleteValue(t *testing.T) {
	test := NewTest("test")
	err := dix.Add(TestKey, test, dix.WithValueSetDefault())
	if err != nil {
		t.Errorf("unexpected Add() err: got %v, want %v", err, nil)
	}

	_, err = dix.Get[*Test]()
	if err != nil {
		t.Errorf("unexpected Get() err: got %v, want %v", err, nil)
	}

	dix.Delete[*Test]()

	_, err = dix.Get[*Test]()
	if !errors.Is(err, dix.ErrValueNotFound) {
		t.Errorf("unexpected Get() err: got %v, want %v", err, dix.ErrValueNotFound)
	}
}

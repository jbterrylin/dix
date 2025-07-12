package dix_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

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

func TestGetConcurrentGetProvider(t *testing.T) {
	err := dix.AddProvider(TestProviderKey, func() (*Test, error) {
		return NewTest("test"), nil
	}, dix.WithProviderSetDefault())
	if err != nil {
		t.Errorf("unexpected AddProvider() err: got %v, want %v", err, nil)
	}

	runCount := 1000000
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < runCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testFromC, err := dix.GetProvider[*Test]()
			if err != nil {
				t.Errorf("unexpected GetProvider() err: got %v, want %v", err, nil)
			}

			testFromC.Name()
			testFromC.Count()

		}()
	}
	wg.Wait()
	end := time.Now()
	fmt.Println("dix used time", end.Sub(start))

	test := NewTest("test")

	start = time.Now()
	for i := 0; i < runCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			test.Name()
			test.Count()
		}()
	}
	wg.Wait()
	end = time.Now()
	fmt.Println("no dix used time", end.Sub(start))
}

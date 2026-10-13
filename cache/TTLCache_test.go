package cache

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDisabledWhenTTLZero(t *testing.T) {
	c := New[string](0, 100)
	if c.Enabled() {
		t.Fatal("TTL 0 must disable the cache")
	}
	calls := 0
	for i := 0; i < 3; i++ {
		v, hit, err := c.GetOrLoad("k", func() (string, error) { calls++; return "v", nil })
		if v != "v" || hit || err != nil {
			t.Fatalf("got %q hit=%v err=%v", v, hit, err)
		}
	}
	if calls != 3 {
		t.Fatalf("disabled cache must call load every time, got %d calls", calls)
	}
}

func TestHitWithinTTLAndExpiry(t *testing.T) {
	c := New[string](50*time.Millisecond, 100)
	calls := 0
	load := func() (string, error) { calls++; return "v", nil }

	if _, hit, _ := c.GetOrLoad("k", load); hit {
		t.Fatal("first call must miss")
	}
	if _, hit, _ := c.GetOrLoad("k", load); !hit {
		t.Fatal("second call within TTL must hit")
	}
	time.Sleep(60 * time.Millisecond)
	if _, hit, _ := c.GetOrLoad("k", load); hit {
		t.Fatal("call after TTL must miss")
	}
	if calls != 2 {
		t.Fatalf("want 2 loads, got %d", calls)
	}
}

func TestErrorsNotCached(t *testing.T) {
	c := New[string](time.Minute, 100)
	calls := 0
	fail := func() (string, error) { calls++; return "", errors.New("rate_limit_error") }

	for i := 0; i < 3; i++ {
		if _, _, err := c.GetOrLoad("k", fail); err == nil {
			t.Fatal("error must be returned")
		}
	}
	if calls != 3 || c.Len() != 0 {
		t.Fatalf("errors must not be cached: calls=%d len=%d", calls, c.Len())
	}
}

func TestBoundedLRU(t *testing.T) {
	c := New[int](time.Minute, 3)
	for i := 0; i < 10; i++ {
		c.Set(fmt.Sprint(i), i)
	}
	if c.Len() != 3 {
		t.Fatalf("want 3 entries, got %d", c.Len())
	}
	if _, ok := c.Get("0"); ok {
		t.Fatal("oldest entry must be evicted")
	}
	// Touch 7 so it becomes most recent, then add one more: 8 is evicted.
	c.Get("7")
	c.Set("10", 10)
	if _, ok := c.Get("8"); ok {
		t.Fatal("least recently used entry must be evicted")
	}
	if _, ok := c.Get("7"); !ok {
		t.Fatal("recently used entry must survive")
	}
}

func TestConcurrentMissesCollapse(t *testing.T) {
	c := New[string](time.Minute, 100)
	var calls int32
	release := make(chan struct{})
	load := func() (string, error) {
		atomic.AddInt32(&calls, 1)
		<-release
		return "v", nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if v, _, err := c.GetOrLoad("k", load); v != "v" || err != nil {
				t.Errorf("got %q err=%v", v, err)
			}
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(release)
	wg.Wait()

	if calls != 1 {
		t.Fatalf("20 concurrent misses must trigger 1 load, got %d", calls)
	}
}

package cache

import (
	"container/list"
	"errors"
	"math/rand"
	"sync"
	"time"
)

// TTLCache is a small in-memory cache with per-entry expiry and a hard
// entry limit (least-recently-used entries are evicted first), so it can
// never grow without bound. A nil *TTLCache is valid and caches nothing,
// which is how a TTL of 0 disables caching.
//
// GetOrLoad also collapses concurrent misses for the same key into a
// single load: when several callers ask for an uncached key at once
// (e.g. multiple anonymous page views landing together right after an
// entry expired), only one of them calls the loader and the rest wait
// for its result instead of each hitting the backend.
type TTLCache[V any] struct {
	ttl        time.Duration
	maxEntries int

	mu       sync.Mutex
	order    *list.List // front = most recently used
	entries  map[string]*list.Element
	inflight map[string]*call[V]
}

var errLoadFailed = errors.New("cache: load did not complete")

type entry[V any] struct {
	key     string
	value   V
	expires time.Time
}

type call[V any] struct {
	done  chan struct{}
	value V
	err   error
}

// New returns a cache, or nil (caching disabled) if ttl <= 0 or
// maxEntries <= 0.
func New[V any](ttl time.Duration, maxEntries int) *TTLCache[V] {
	if ttl <= 0 || maxEntries <= 0 {
		return nil
	}
	return &TTLCache[V]{
		ttl:        ttl,
		maxEntries: maxEntries,
		order:      list.New(),
		entries:    make(map[string]*list.Element),
		inflight:   make(map[string]*call[V]),
	}
}

// Enabled reports whether this cache actually stores anything.
func (c *TTLCache[V]) Enabled() bool {
	return c != nil
}

// Get returns a live (unexpired) entry for key.
func (c *TTLCache[V]) Get(key string) (V, bool) {
	var zero V
	if c == nil {
		return zero, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.getLocked(key)
}

func (c *TTLCache[V]) getLocked(key string) (V, bool) {
	var zero V
	el, ok := c.entries[key]
	if !ok {
		return zero, false
	}
	e := el.Value.(*entry[V])
	if time.Now().After(e.expires) {
		c.order.Remove(el)
		delete(c.entries, key)
		return zero, false
	}
	c.order.MoveToFront(el)
	return e.value, true
}

// Set stores value under key. Expiry gets up to 10% random jitter so
// entries loaded together (lemmyBB resolves every community on a page in
// one burst) don't all expire in the same instant and come back as one
// burst again.
func (c *TTLCache[V]) Set(key string, value V) {
	if c == nil {
		return
	}
	jitter := time.Duration(rand.Int63n(int64(c.ttl)/10 + 1))
	expires := time.Now().Add(c.ttl - jitter)

	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.entries[key]; ok {
		e := el.Value.(*entry[V])
		e.value = value
		e.expires = expires
		c.order.MoveToFront(el)
		return
	}

	c.entries[key] = c.order.PushFront(&entry[V]{key: key, value: value, expires: expires})
	for c.order.Len() > c.maxEntries {
		oldest := c.order.Back()
		c.order.Remove(oldest)
		delete(c.entries, oldest.Value.(*entry[V]).key)
	}
}

// Delete removes key if present.
func (c *TTLCache[V]) Delete(key string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.entries[key]; ok {
		c.order.Remove(el)
		delete(c.entries, key)
	}
}

// Len returns the number of stored entries (including any not yet
// lazily evicted after expiry).
func (c *TTLCache[V]) Len() int {
	if c == nil {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

// GetOrLoad returns the cached value for key, or calls load to produce
// it. Only successful loads (err == nil) are stored — errors are returned
// to every waiting caller but never cached, so the next request retries.
// hit reports whether the value came from the cache (or from another
// caller's in-flight load) rather than from this call's own load.
func (c *TTLCache[V]) GetOrLoad(key string, load func() (V, error)) (value V, hit bool, err error) {
	if c == nil {
		value, err = load()
		return value, false, err
	}

	c.mu.Lock()
	if v, ok := c.getLocked(key); ok {
		c.mu.Unlock()
		return v, true, nil
	}
	if inflight, ok := c.inflight[key]; ok {
		c.mu.Unlock()
		<-inflight.done
		return inflight.value, true, inflight.err
	}
	own := &call[V]{done: make(chan struct{})}
	c.inflight[key] = own
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.inflight, key)
		c.mu.Unlock()
		close(own.done)
	}()

	// Pre-set an error so that if load panics, callers waiting on this
	// in-flight call get an error rather than a zero value with err == nil.
	own.err = errLoadFailed
	own.value, own.err = load()
	if own.err == nil {
		c.Set(key, own.value)
	}
	return own.value, false, own.err
}

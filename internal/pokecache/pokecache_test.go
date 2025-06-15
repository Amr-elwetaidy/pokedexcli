package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)

	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://example.com",
			value: []byte("testdata"),
		},
		{
			key:   "https://example.com/path",
			value: []byte("moredatatest"),
		},
	}

	for _, c := range cases {
		t.Run("Test adding and getting from cache", func(t *testing.T) {
			cache.Add(c.key, c.value)
			actual, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key %s in cache", c.key)
			}
			if string(actual) != string(c.value) {
				t.Errorf("expected value %s, but got %s", string(c.value), string(actual))
			}
		})
	}
}

func TestCacheReap(t *testing.T) {
	const reapInterval = 100 * time.Millisecond
	const testEntryLifetime = 50 * time.Millisecond

	cache := NewCache(reapInterval)

	key := "https://example.com/reap"
	cache.Add(key, []byte("testreap"))

	// Ensure the entry is there initially
	_, ok := cache.Get(key)
	if !ok {
		t.Fatalf("expected to find key %s right after adding", key)
	}

	// Wait long enough for the entry to be reaped
	time.Sleep(reapInterval + testEntryLifetime)

	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected key %s to be reaped from the cache", key)
	}
}

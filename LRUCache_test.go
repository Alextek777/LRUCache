package LRUCache

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cache := New[int, string](10, time.Minute)
	if cache == nil {
		t.Error("Expected new cache instance, got nil")
	}
	if cache.limit != 10 {
		t.Errorf("Expected limit 10, got %d", cache.limit)
	}
	if cache.ttl != time.Minute {
		t.Errorf("Expected TTL 1m, got %v", cache.ttl)
	}
}

func TestPushAndGet(t *testing.T) {
	cache := New[int, string](3, time.Minute)

	// Test basic push and get
	cache.Push(1, "one")
	val, ok := cache.Get(1)
	if !ok {
		t.Error("Expected to find key 1")
	}
	if val != "one" {
		t.Errorf("Expected 'one', got '%s'", val)
	}

	// Test non-existent key
	_, ok = cache.Get(2)
	if ok {
		t.Error("Expected not to find key 2")
	}
}

func TestLRUEviction(t *testing.T) {
	cache := New[int, string](2, time.Minute)

	cache.Push(1, "one")
	cache.Push(2, "two")
	cache.Push(3, "three") // This should evict key 1

	_, ok := cache.Get(1)
	if ok {
		t.Error("Expected key 1 to be evicted")
	}

	// Access key 2 to make it recently used
	cache.Get(2)
	cache.Push(4, "four") // This should evict key 3, not 2

	_, ok = cache.Get(3)
	if ok {
		t.Error("Expected key 3 to be evicted")
	}
	_, ok = cache.Get(2)
	if !ok {
		t.Error("Expected key 2 to still be in cache")
	}
}

func TestTTLEviction(t *testing.T) {
	shortTTL := time.Millisecond * 100
	cache := New[int, string](10, shortTTL)

	cache.Push(1, "one")
	time.Sleep(shortTTL * 2)

	_, ok := cache.Get(1)
	if ok {
		t.Error("Expected key 1 to be expired")
	}
}

func TestRemove(t *testing.T) {
	cache := New[int, string](3, time.Minute)

	cache.Push(1, "one")
	removed := cache.Remove(1)
	if !removed {
		t.Error("Expected to remove key 1")
	}

	_, ok := cache.Get(1)
	if ok {
		t.Error("Expected key 1 to be removed")
	}

	removed = cache.Remove(2) // non-existent key
	if removed {
		t.Error("Expected not to remove non-existent key 2")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cacheSize := 100
	cache := New[int, int](cacheSize, time.Minute)
	done := make(chan bool)

	// Concurrent writers
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < cacheSize; j++ {
				cache.Push(j, j)
				cache.Get(j)
			}
			done <- true
		}()
	}

	// Concurrent readers
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < cacheSize; j++ {
				cache.Get(j)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to finish
	for i := 0; i < 20; i++ {
		<-done
	}

	// Verify cache size hasn't exceeded limit
	count := 0
	for j := 0; j < cacheSize; j++ {
		if val, ok := cache.Get(j); ok && val == j {
			count++
		}
	}

	if count > cacheSize {
		t.Errorf("Cache size exceeded limit: got %d, expected max %d", count, cacheSize)
	}

	// Verify no corrupted entries
	for j := 0; j < cacheSize; j++ {
		if val, ok := cache.Get(j); ok && val != j {
			t.Errorf("Cache corruption detected for key %d: got %d, expected %d", j, val, j)
		}
	}
}

func TestReplaceExisting(t *testing.T) {
	cache := New[int, string](3, time.Minute)

	cache.Push(1, "one")
	cache.Push(1, "new one") // Replace existing

	val, ok := cache.Get(1)
	if !ok {
		t.Error("Expected to find key 1")
	}
	if val != "new one" {
		t.Errorf("Expected 'new one', got '%s'", val)
	}
}

func TestCacheSizeLimit(t *testing.T) {
	cache := New[int, int](3, time.Minute)

	for i := 0; i < 5; i++ {
		cache.Push(i, i)
	}

	// Only the last 3 items should be in cache
	for i := 0; i < 2; i++ {
		_, ok := cache.Get(i)
		if ok {
			t.Errorf("Expected key %d to be evicted", i)
		}
	}
	for i := 2; i < 5; i++ {
		_, ok := cache.Get(i)
		if !ok {
			t.Errorf("Expected key %d to be in cache", i)
		}
	}
}

func TestTTLRefreshOnGet(t *testing.T) {
	shortTTL := time.Millisecond * 100
	cache := New[int, string](3, shortTTL)

	cache.Push(1, "one")
	time.Sleep(shortTTL / 2) // Halfway through TTL

	// Access should refresh TTL
	cache.Get(1)
	time.Sleep(shortTTL / 2) // Now original TTL would be expired

	// Item should still be there because we refreshed it
	_, ok := cache.Get(1)
	if !ok {
		t.Error("Expected key 1 to still be in cache after refresh")
	}
}

func TestClearExpired(t *testing.T) {
	mixedTTL := time.Millisecond * 100
	cache := New[int, string](10, mixedTTL)

	cache.Push(1, "one") // Will expire
	cache.Push(2, "two") // Will expire
	time.Sleep(mixedTTL * 2)

	cache.Push(3, "three") // This should clear expired items

	_, ok := cache.Get(1)
	if ok {
		t.Error("Expected key 1 to be expired")
	}
	_, ok = cache.Get(2)
	if ok {
		t.Error("Expected key 2 to be expired")
	}
	_, ok = cache.Get(3)
	if !ok {
		t.Error("Expected key 3 to be in cache")
	}
}

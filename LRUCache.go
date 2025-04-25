package LRUCache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type listItem[K comparable] struct {
	Key       K
	ExpiresAt time.Time
}

type cacheItem[V any] struct {
	Value V
	Node  *list.Element
}

type LRUCache[K comparable, V any] struct {
	ttlList *list.List
	cache   map[K]cacheItem[V]
	ttl     time.Duration
	limit   int

	mu sync.RWMutex
}

func (l *listItem[K]) isExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

func New[K comparable, V any](limit int, ttl time.Duration) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		limit:   limit,
		ttlList: list.New(),
		cache:   make(map[K]cacheItem[V], limit),
		mu:      sync.RWMutex{},
		ttl:     ttl,
	}
}

func (s *LRUCache[K, V]) Push(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// clear old elements TTL expired
	s.unsafeClear()

	// clear least recently used if limit exceeded
	if s.ttlList.Len() >= s.limit {
		s.unsafeRemove(s.ttlList.Front().Value.(*listItem[K]).Key)
	}

	// push elements (replace if exists)
	_, ok := s.cache[key]
	if ok {
		s.unsafeRemove(key)
	}

	node := s.ttlList.PushBack(&listItem[K]{
		Key:       key,
		ExpiresAt: time.Now().Add(s.ttl),
	})
	s.cache[key] = cacheItem[V]{
		Value: value,
		Node:  node,
	}
}

func (s *LRUCache[K, V]) Remove(key K) bool {
	s.mu.Lock()
	s.unsafeClear()
	s.mu.Unlock()

	s.mu.RLock()

	elem, ok := s.cache[key]
	if !ok {
		s.mu.RUnlock()
		return false
	}

	s.mu.RUnlock()
	s.mu.Lock()

	s.unsafeClear()
	s.ttlList.Remove(elem.Node)
	delete(s.cache, key)

	s.mu.Unlock()
	return true
}

func (s *LRUCache[K, V]) Get(key K) (V, bool) {
	s.mu.Lock()
	s.unsafeClear()
	s.mu.Unlock()

	s.mu.RLock()

	v, ok := s.cache[key]
	if !ok {
		s.mu.RUnlock()
		return v.Value, false
	}

	s.mu.RUnlock()
	s.mu.Lock()

	node, ok := v.Node.Value.(*listItem[K])
	if !ok {
		s.mu.Unlock()
		fmt.Println("LRUCache error! Could not cast to listItem on Get")
		return v.Value, false
	}

	node.ExpiresAt = time.Now().Add(s.ttl)
	s.ttlList.MoveToBack(v.Node)
	value := s.cache[key].Value

	s.mu.Unlock()
	return value, true
}

func (s *LRUCache[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ttlList.Len()
}

func (s *LRUCache[K, V]) unsafeClear() {
	e := s.ttlList.Front()
	for e != nil {
		item, ok := e.Value.(*listItem[K])
		if !ok {
			fmt.Println("LRUCache error! Could not cast to listItem on unsafeClear")
			continue
		}

		if item.isExpired() {
			next := e.Next()

			delete(s.cache, item.Key)
			s.ttlList.Remove(e)

			e = next
			continue
		}

		return
	}
}

func (s *LRUCache[K, V]) unsafeRemove(key K) {
	s.ttlList.Remove(s.cache[key].Node)
	delete(s.cache, key)
}

func (s *LRUCache[K, V]) Display() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("LRUCache")
	for k, v := range s.cache {
		fmt.Println("Key: ", k, " Value: ", v.Value)
	}

	fmt.Println("\nLRUCache TTL List")
	for e := s.ttlList.Front(); e != nil; e = e.Next() {
		fmt.Println("Key: ", e.Value.(*listItem[K]).Key, " TTL : ", e.Value.(*listItem[K]).ExpiresAt)
	}
}

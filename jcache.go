package jcache

import (
	"sync"
	"time"
)

// CleanupStrategy represents the strategy used for cache cleanup.
type CleanupStrategy int

const (
	cleanupStrategyOnTheFly CleanupStrategy = iota
	cleanupStrategyCollect
)

const (
	defaultInterval = 1 * time.Minute
	defaultCapacity = 100000
	defaultStrategy = cleanupStrategyOnTheFly
)

type item struct {
	value     any
	deadAfter time.Time
}

// Item represents a cache entry with its key and value.
type Item struct {
	Key   string
	Value any
}

// JustCache is an in-memory cache with expiration and cleanup features.
type JustCache struct {
	mu sync.RWMutex

	cleanupStrategy CleanupStrategy
	cleanupInterval time.Duration
	cleanupStopCh   chan struct{}

	items    map[string]item
	capacity int64
}

// New creates a new JustCache instance with the specified cleanup interval and capacity.
func New(interval time.Duration, capacity int64) *JustCache {
	jc := &JustCache{
		items:           make(map[string]item, capacity),
		cleanupInterval: interval,
		cleanupStrategy: cleanupStrategyOnTheFly,
		cleanupStopCh:   make(chan struct{}),
		capacity:        capacity,
	}

	go jc.cleanup()

	return jc
}

// Default returns a JustCache instance with default settings.
func Default() *JustCache {
	jc := New(defaultInterval, defaultCapacity)
	jc.WithStrategy(defaultStrategy)

	return jc
}

// WithStrategy sets the cleanup strategy and returns the modified JustCache.
func (jc *JustCache) WithStrategy(strategy CleanupStrategy) *JustCache {
	jc.cleanupStrategy = strategy
	return jc
}

// WithCleanupInterval sets the cleanup interval and returns the modified JustCache.
func (jc *JustCache) WithCleanupInterval(interval time.Duration) *JustCache {
	jc.cleanupInterval = interval
	return jc
}

// WithCapacity sets the cache capacity and returns the modified JustCache.
func (jc *JustCache) WithCapacity(capacity int64) *JustCache {
	jc.capacity = capacity
	return jc
}

// Get retrieves the value associated with the specified key.
func (jc *JustCache) Get(key string) any {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	return jc.items[key].value
}

// Has checks whether the specified key exists in the cache.
func (jc *JustCache) Has(key string) bool {
	jc.mu.RLock()
	defer jc.mu.RUnlock()
	_, ok := jc.items[key]
	return ok
}

// Item returns a cache Item for the specified key.
func (jc *JustCache) Item(key string) Item {
	return Item{Key: key, Value: jc.Get(key)}
}

// Set stores a key-value pair in the cache with a time-to-live duration.
// Returns an error if the cache is full.
func (jc *JustCache) Set(key string, value any, ttl time.Duration) error {
	if int64(len(jc.items)) >= jc.capacity {
		return ErrorCacheIsFull
	}

	jc.mu.Lock()
	defer jc.mu.Unlock()
	jc.items[key] = item{
		value:     value,
		deadAfter: time.Now().Add(ttl),
	}

	return nil
}

// Del deletes the specified key from the cache.
func (jc *JustCache) Del(key string) {
	jc.mu.Lock()
	defer jc.mu.Unlock()

	delete(jc.items, key)
}

// Keys returns a slice of all keys in the cache.
func (jc *JustCache) Keys() []string {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	keys := make([]string, 0, len(jc.items))
	for k := range jc.items {
		keys = append(keys, k)
	}

	return keys
}

// Items returns a slice of all cache items.
func (jc *JustCache) Items() []Item {
	jc.mu.RLock()
	defer jc.mu.RUnlock()
	items := make([]Item, 0, len(jc.items))
	for key, item := range jc.items {
		items = append(items, Item{Key: key, Value: item.value})
	}

	return items
}

// Clear removes all items from the cache.
func (jc *JustCache) Clear() {
	jc.mu.Lock()
	defer jc.mu.Unlock()
	jc.items = make(map[string]item)
}

// Clean removes expired items from the cache based on the provided time.
func (jc *JustCache) Clean(now time.Time) {
	switch jc.cleanupStrategy {
	case cleanupStrategyOnTheFly:
		for key, record := range jc.items {
			if now.After(record.deadAfter) {
				jc.mu.Lock()
				delete(jc.items, key)
				jc.mu.Unlock()
			}
		}
	case cleanupStrategyCollect:
		dead := jc.dead(now)

		jc.mu.Lock()
		for _, key := range dead {
			delete(jc.items, key)
		}
		jc.mu.Unlock()
	}
}

// Close stops the cleanup process and clears the cache.
func (jc *JustCache) Close() error {
	jc.stopCleanup()
	close(jc.cleanupStopCh)
	jc.Clear()

	return nil
}

func (jc *JustCache) dead(tick time.Time) []string {
	jc.mu.RLock()
	defer jc.mu.RUnlock()

	var dead []string
	for k, r := range jc.items {
		if tick.After(r.deadAfter) {
			dead = append(dead, k)
		}
	}

	return dead
}

func (jc *JustCache) stopCleanup() {
	jc.cleanupStopCh <- struct{}{}
}

func (jc *JustCache) cleanup() {
	ticker := time.NewTicker(jc.cleanupInterval)

	for {
		select {
		case <-jc.cleanupStopCh:
			return
		case tick := <-ticker.C:
			jc.Clean(tick)
		}
	}
}

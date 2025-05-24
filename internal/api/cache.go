package api

import (
	"log"
	"sync"
	"time"
)

type CachedData[T any] struct {
	Data       T
	CachedTime time.Time
}

type Cache[T any] struct {
	Cache     map[string]CachedData[T]
	mu        sync.RWMutex
	RenewTime time.Duration
}

func NewCache[T any](renewTime time.Duration) *Cache[T] {
	return &Cache[T]{
		Cache:     make(map[string]CachedData[T]),
		RenewTime: renewTime,
	}
}

func (cache *Cache[T]) ReadCache(steamID string) (T, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	item, found := cache.Cache[steamID]
	if !found {
		var zero T
		return zero, false
	}
	if time.Since(item.CachedTime) >= cache.RenewTime {
		var zero T
		return zero, false
	}

	return cache.Cache[steamID].Data, true
}

func (cache *Cache[T]) UpdateCache(steamID string, data T) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cacheData := cache.Cache[steamID]
	cacheData.Data = data
	cacheData.CachedTime = time.Now().UTC()
	cache.Cache[steamID] = cacheData
}

type Cleaner[T any] struct {
	Name     string
	ticker   *time.Ticker
	done     chan bool
	Cache    *Cache[T]
	Interval time.Duration
}

func (cleaner *Cleaner[T]) CacheCleanerStart() {
	cleaner.ticker = time.NewTicker(cleaner.Interval)

	go func() {
		for {
			select {
			case <-cleaner.done:
				return
			case <-cleaner.ticker.C:
				log.Printf("Cleaner being ran for cache '%s'\n", cleaner.Name)
				cleaner.Cache.CleanCache()
				log.Printf("Successfully removed expired entries from cache '%s'\n", cleaner.Name)
			}
		}
	}()
}

func (cache *Cache[T]) CleanCache() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	now := time.Now()

	for steamID, item := range cache.Cache {
		if now.Sub(item.CachedTime) > cache.RenewTime {
			log.Printf("Clearing expired cache entry for steamID: %s\n", steamID)
			delete(cache.Cache, steamID)
		}
	}
}

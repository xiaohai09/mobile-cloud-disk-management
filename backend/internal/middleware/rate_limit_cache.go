package middleware

import (
	"sync"
	"time"
)

// shardCache is a single shard of the bounded rate-limiter visitor cache.
type shardCache struct {
	entries    map[string]*rateLimiterCacheEntry
	mu         sync.Mutex
	maxEntries int
	ttl        time.Duration
}

type rateLimiterCacheEntry struct {
	v        *visitor
	expireAt time.Time
}

// newShardCache creates an empty shard.
func newShardCache(maxEntries int, ttl time.Duration) *shardCache {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	if maxEntries <= 0 {
		maxEntries = 10000
	}
	return &shardCache{
		entries:    make(map[string]*rateLimiterCacheEntry),
		maxEntries: maxEntries,
		ttl:        ttl,
	}
}

// get returns the visitor for key if present and not expired, and removes it if expired.
func (s *shardCache) get(key string) (*visitor, bool) {
	if s == nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expireAt) {
		delete(s.entries, key)
		return nil, false
	}
	return entry.v, true
}

// set stores a visitor for key with the configured TTL.
func (s *shardCache) set(key string, v *visitor) {
	if s == nil || v == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key] = &rateLimiterCacheEntry{
		v:        v,
		expireAt: time.Now().Add(s.ttl),
	}
}

// del removes a key from the shard.
func (s *shardCache) del(key string) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// len returns the number of entries in the shard.
func (s *shardCache) len() int {
	if s == nil {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}

// evictOldest removes the entry with the oldest lastSeen time to make room.
func (s *shardCache) evictOldest() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	var oldestKey string
	var oldestSeen time.Time
	for k, entry := range s.entries {
		if oldestKey == "" || entry.v.lastSeen.Before(oldestSeen) {
			oldestKey = k
			oldestSeen = entry.v.lastSeen
		}
	}
	if oldestKey != "" {
		delete(s.entries, oldestKey)
	}
}

// ShardedRateLimiterCache is a sharded bounded cache for rate-limiter visitors.
type ShardedRateLimiterCache struct {
	shards    []*shardCache
	numShards int
	ttl       time.Duration
	stopCh    chan struct{}
	stopOnce  sync.Once
}

// NewShardedRateLimiterCache creates a new sharded cache.
// numShards: number of shards (must be >0).
// maxEntriesPerShard: maximum entries per shard.
// ttl: TTL for each entry.
func NewShardedRateLimiterCache(numShards, maxEntriesPerShard int, ttl time.Duration) *ShardedRateLimiterCache {
	if numShards <= 0 {
		numShards = 16
	}
	shards := make([]*shardCache, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = newShardCache(maxEntriesPerShard, ttl)
	}
	cache := &ShardedRateLimiterCache{
		shards:    shards,
		numShards: numShards,
		ttl:       ttl,
		stopCh:    make(chan struct{}),
	}
	go cache.cleanupLoop()
	return cache
}

// shardIndex returns the shard index for a key using FNV-1a.
func (c *ShardedRateLimiterCache) shardIndex(key string) int {
	if c == nil || len(c.shards) == 0 || c.numShards <= 0 {
		return 0
	}
	h := fnv1a(key)
	return int(h % uint64(c.numShards))
}

// Compute applies fn to the value for key while holding the shard lock.
// If fn returns nil, the entry is deleted.
// If the entry is missing, fn is called with nil.
// The returned value is stored under key.
func (c *ShardedRateLimiterCache) Compute(key string, fn func(v *visitor) *visitor) {
	if c == nil || key == "" {
		return
	}
	idx := c.shardIndex(key)
	shard := c.shards[idx]
	shard.mu.Lock()
	defer shard.mu.Unlock()

	entry, exists := shard.entries[key]
	if exists && time.Now().After(entry.expireAt) {
		// Expired, treat as missing.
		exists = false
		delete(shard.entries, key)
	}

	var oldV *visitor
	if exists {
		oldV = entry.v
	}
	newV := fn(oldV)
	if newV == nil {
		delete(shard.entries, key)
		return
	}

	// If the shard is at capacity and this is an update (not a new insert),
	// allow it. If it's a new insert and at capacity, evict oldest first.
	if !exists && len(shard.entries) >= shard.maxEntries {
		// evictOldest needs the lock, but we already hold it.
		// We inline a simple eviction here to avoid deadlock.
		var oldestKey string
		var oldestSeen time.Time
		for k, e := range shard.entries {
			if oldestKey == "" || e.v.lastSeen.Before(oldestSeen) {
				oldestKey = k
				oldestSeen = e.v.lastSeen
			}
		}
		if oldestKey != "" {
			delete(shard.entries, oldestKey)
		}
	}

	shard.entries[key] = &rateLimiterCacheEntry{
		v:        newV,
		expireAt: time.Now().Add(c.ttl),
	}
}

// Delete removes a key from the cache.
func (c *ShardedRateLimiterCache) Delete(key string) {
	if c == nil || key == "" {
		return
	}
	idx := c.shardIndex(key)
	c.shards[idx].del(key)
}

// Len returns the total number of entries across all shards.
func (c *ShardedRateLimiterCache) Len() int {
	if c == nil {
		return 0
	}
	total := 0
	for _, shard := range c.shards {
		total += shard.len()
	}
	return total
}

// cleanupLoop periodically evicts expired entries across all shards.
func (c *ShardedRateLimiterCache) cleanupLoop() {
	interval := c.ttl / 2
	if interval < time.Minute {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCh:
			return
		}
	}
}

// cleanup scans all shards and removes expired entries.
func (c *ShardedRateLimiterCache) cleanup() {
	now := time.Now()
	for _, shard := range c.shards {
		shard.mu.Lock()
		for k, entry := range shard.entries {
			if now.After(entry.expireAt) {
				delete(shard.entries, k)
			}
		}
		shard.mu.Unlock()
	}
}

// Stop signals the cleanup goroutine to exit.
func (c *ShardedRateLimiterCache) Stop() {
	if c == nil {
		return
	}
	c.stopOnce.Do(func() {
		close(c.stopCh)
	})
}

// fnv1a implements the FNV-1a hash algorithm for strings.
func fnv1a(s string) uint64 {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	var h uint64 = offset64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}

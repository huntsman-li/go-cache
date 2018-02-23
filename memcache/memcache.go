package cache

import (
	"strings"

	"github.com/Unknwon/com"
	"github.com/bradfitz/gomemcache/memcache"

	"github.com/huntsman-li/go-cache"
)

// MemcacheCacher represents a memcache cache adapter implementation.
type MemcacheCacher struct {
	c *memcache.Client
}

func NewItem(key string, data []byte, expire int32) *memcache.Item {
	return &memcache.Item{
		Key:        key,
		Value:      data,
		Expiration: expire,
	}
}

// Put puts value into cache with key and expire time.
// If expired is 0, it lives forever.
func (c *MemcacheCacher) Put(key string, val interface{}, expire int64) error {
	return c.c.Set(NewItem(key, []byte(com.ToStr(val)), int32(expire)))
}

// Get gets cached value by given key.
func (c *MemcacheCacher) Get(key string) interface{} {
	item, err := c.c.Get(key)
	if err != nil {
		return nil
	}
	return string(item.Value)
}

// Delete deletes cached value by given key.
func (c *MemcacheCacher) Delete(key string) error {
	return c.c.Delete(key)
}

// Incr increases cached int-type value by given key as a counter.
func (c *MemcacheCacher) Incr(key string) error {
	_, err := c.c.Increment(key, 1)
	return err
}

// Decr decreases cached int-type value by given key as a counter.
func (c *MemcacheCacher) Decr(key string) error {
	_, err := c.c.Decrement(key, 1)
	return err
}

// IsExist returns true if cached value exists.
func (c *MemcacheCacher) IsExist(key string) bool {
	_, err := c.c.Get(key)
	return err == nil
}

// Flush deletes all cached data.
func (c *MemcacheCacher) Flush() error {
	return c.c.FlushAll()
}

// StartAndGC starts GC routine based on config string settings.
// AdapterConfig: 127.0.0.1:9090;127.0.0.1:9091
func (c *MemcacheCacher) StartAndGC(opt cache.Options) error {
	c.c = memcache.New(strings.Split(opt.AdapterConfig, ";")...)
	return nil
}

func init() {
	cache.Register("memcache", &MemcacheCacher{})
}

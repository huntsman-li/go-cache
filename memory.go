package cache

import (
	"errors"
	"sync"
	"time"
)

// MemoryItem represents a memory cache item.
type MemoryItem struct {
	val     interface{}
	created int64
	expire  int64
}

func (item *MemoryItem) hasExpired() bool {
	return item.expire > 0 &&
		(time.Now().Unix()-item.created) >= item.expire
}

// MemoryCacher represents a memory cache adapter implementation.
type MemoryCacher struct {
	lock     sync.RWMutex
	items    map[string]*MemoryItem
	interval int // GC interval.
}

// NewMemoryCacher creates and returns a new memory cacher.
func NewMemoryCacher() *MemoryCacher {
	return &MemoryCacher{items: make(map[string]*MemoryItem)}
}

// Put puts value into cache with key and expire time.
// If expired is 0, it will be deleted by next GC operation.
func (c *MemoryCacher) Put(key string, val interface{}, expire int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = &MemoryItem{
		val:     val,
		created: time.Now().Unix(),
		expire:  expire,
	}
	return nil
}

// Get gets cached value by given key.
func (c *MemoryCacher) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil
	}
	if item.hasExpired() {
		go c.Delete(key)
		return nil
	}
	return item.val
}

// Delete deletes cached value by given key.
func (c *MemoryCacher) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.items, key)
	return nil
}

// Incr increases cached int-type value by given key as a counter.
func (c *MemoryCacher) Incr(key string) (err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return errors.New("key not exist")
	}
	item.val, err = Incr(item.val)
	return err
}

// Decr decreases cached int-type value by given key as a counter.
func (c *MemoryCacher) Decr(key string) (err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return errors.New("key not exist")
	}

	item.val, err = Decr(item.val)
	return err
}

// IsExist returns true if cached value exists.
func (c *MemoryCacher) IsExist(key string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.items[key]
	return ok
}

// Flush deletes all cached data.
func (c *MemoryCacher) Flush() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items = make(map[string]*MemoryItem)
	return nil
}

func (c *MemoryCacher) checkRawExpiration(key string) {
	item, ok := c.items[key]
	if !ok {
		return
	}

	if item.hasExpired() {
		delete(c.items, key)
	}
}

func (c *MemoryCacher) checkExpiration(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.checkRawExpiration(key)
}

func (c *MemoryCacher) startGC() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.interval < 1 {
		return
	}

	if c.items != nil {
		for key, _ := range c.items {
			c.checkRawExpiration(key)
		}
	}

	time.AfterFunc(time.Duration(c.interval)*time.Second, func() { c.startGC() })
}

// StartAndGC starts GC routine based on config string settings.
func (c *MemoryCacher) StartAndGC(opt Options) error {
	c.lock.Lock()
	c.interval = opt.Interval
	c.lock.Unlock()

	go c.startGC()
	return nil
}

func init() {
	Register("memory", NewMemoryCacher())
}

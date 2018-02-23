package cache

import (
	"fmt"
	"gopkg.in/ini.v1"
)

const _VERSION = "0.0.1"

func Version() string {
	return _VERSION
}

// Cache is the interface that operates the cache data.
type Cache interface {
	// Put puts value into cache with key and expire time.
	Put(key string, val interface{}, timeout int64) error
	// Get gets cached value by given key.
	Get(key string) interface{}
	// Delete deletes cached value by given key.
	Delete(key string) error
	// Incr increases cached int-type value by given key as a counter.
	Incr(key string) error
	// Decr decreases cached int-type value by given key as a counter.
	Decr(key string) error
	// IsExist returns true if cached value exists.
	IsExist(key string) bool
	// Flush deletes all cached data.
	Flush() error
	// StartAndGC starts GC routine based on config string settings.
	StartAndGC(opt Options) error
}

// Options represents a struct for specifying configuration options for the cache middleware.
type Options struct {
	// Name of adapter. Default is "memory".
	Adapter string
	// Adapter configuration, it's corresponding to adapter.
	AdapterConfig string
	// GC interval time in seconds. Default is 60.
	Interval int
	// Occupy entire database. Default is false.
	OccupyMode bool
	// Configuration section name. Default is "cache".
	Section string
}

var cfg *ini.File

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Section) == 0 {
		opt.Section = "cache"
	}
	sec := Config().Section(opt.Section)

	if len(opt.Adapter) == 0 {
		opt.Adapter = sec.Key("ADAPTER").MustString("memory")
	}
	if opt.Interval == 0 {
		opt.Interval = sec.Key("INTERVAL").MustInt(60)
	}
	if len(opt.AdapterConfig) == 0 {
		opt.AdapterConfig = sec.Key("ADAPTER_CONFIG").MustString("data/caches")
	}

	return opt
}

// NewCacher creates and returns a new cacher by given adapter name and configuration.
// It panics when given adapter isn't registered and starts GC automatically.
func NewCacher(name string, opt Options) (Cache, error) {
	adapter, ok := adapters[name]
	if !ok {
		return nil, fmt.Errorf("cache: unknown adapter '%s'(forgot to import?)", name)
	}
	return adapter, adapter.StartAndGC(opt)
}

// Cacher is a middleware that maps a cache.Cache service into the Macaron handler chain.
// An single variadic cache.Options struct can be optionally provided to configure.
func Cacher(options ...Options) (Cache, error) {
	opt := prepareOptions(options)
	return NewCacher(opt.Adapter, opt)
}

var adapters = make(map[string]Cache)

// Register registers a adapter.
func Register(name string, adapter Cache) {
	if adapter == nil {
		panic("cache: cannot register adapter with nil value")
	}
	if _, dup := adapters[name]; dup {
		panic(fmt.Errorf("cache: cannot register adapter '%s' twice", name))
	}
	adapters[name] = adapter
}

func Config() *ini.File {
	if cfg == nil {
		return ini.Empty()
	}
	return cfg
}

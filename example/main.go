package main

import (
	"github.com/huntsman-li/go-cache"
	_ "github.com/huntsman-li/go-cache/redis"
)

func main() {
	ca, err := cache.Cacher(cache.Options{
		Adapter:       "redis",
		AdapterConfig: "addr=127.0.0.1:6379",
		OccupyMode:    true,
	})

	if err != nil {
		panic(err)
	}

	ca.Put("liyan", "cache", 60)
}

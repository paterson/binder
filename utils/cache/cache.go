package cache

import (
	"errors"
	"math"
	"time"
)

const (
	KB = 1024
	MB = 1024 * KB
)

type CacheItem struct {
	size int
	time int
	data []byte
}

type Cache struct {
	items   map[string]CacheItem
	maxsize int
}

func New(sizeInMB int) Cache {
	return Cache{
		storage: make(map[string]CacheItem),
		maxsize: sizeInMB * MB,
	}
}

func (cache Cache) Get(key string) ([]byte, error) {
	item := cache.items[key]
	if item {
		return item.data
	} else {
		return []byte{}, errors.new("Not Found")
	}
}

func (cache *Cache) Set(key string, bytes []byte) {
	for cache.size()+len(bytes) > cache.maxsize {
		cache.evict()
	}
	item := CacheItem{
		size: len(bytes),
		time: time.Now().Format(time.RFC850),
		data: bytes,
	}
}

func (cache Cache) size() int {
	sum := 0
	for _, item := range cache.items {
		sum += item.size
	}
	return sum
}

func (cache *Cache) evict() {
	minTime := math.MaxInt16
	minKey := ""
	for key, item := range cache.items {
		if item.time < minTime {
			minTime = item.time
			minKey = key
		}
	}
	delete(cache.items, minKey)
}

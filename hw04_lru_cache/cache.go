package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш// Place your code here
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	mutex    sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	existedItem, exists := cache.items[key]

	newItem := cacheItem{
		key:   key,
		value: value,
	}

	if exists {
		existedItem.Value = newItem
		cache.queue.MoveToFront(existedItem)

		return true
	}
	if cache.capacity == cache.queue.Len() {
		lastElement := cache.queue.Back()
		cache.queue.Remove(lastElement)
		delete(cache.items, lastElement.Value.(cacheItem).key)
	}

	frontItem := cache.queue.PushFront(newItem)
	cache.items[key] = frontItem

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if existedItem, exists := cache.items[key]; exists {
		cache.queue.MoveToFront(existedItem)
		return existedItem.Value.(cacheItem).value, true
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for key, existedItem := range cache.items {
		cache.queue.Remove(existedItem)
		delete(cache.items, key)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*listItem),
	}
}

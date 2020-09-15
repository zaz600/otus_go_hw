package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.RWMutex
	capacity int
	queue    List
	items    map[string]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	k := string(key)
	l.Lock()
	defer l.Unlock()

	if item, ok := l.items[k]; ok {
		item.Value = cacheItem{
			key:   k,
			value: value,
		}
		l.queue.MoveToFront(item)
		return true
	}

	item := l.queue.PushFront(cacheItem{
		key:   k,
		value: value,
	})
	l.items[k] = item

	if l.queue.Len() > l.capacity {
		back := l.queue.Back()
		delete(l.items, back.Value.(cacheItem).key)
		l.queue.Remove(back)
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	k := string(key)
	l.RLock()
	item, ok := l.items[k]
	l.RUnlock()

	if ok {
		l.Lock()
		defer l.Unlock()
		l.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()
	l.queue = NewList()
	l.items = make(map[string]*ListItem, l.capacity)
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[string]*ListItem, capacity),
	}
}

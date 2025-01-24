package hw04lrucache

type (
	Key        string
	CacheValue struct {
		key   Key
		value interface{}
	}
)

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lruCash *lruCache) Set(key Key, value interface{}) bool {
	cacheValue := CacheValue{
		key:   key,
		value: value,
	}
	oldValue, exists := lruCash.items[key]
	if exists {
		oldValue.Value = cacheValue
		lruCash.queue.MoveToFront(oldValue)
		return true
	}
	lruCash.queue.PushFront(cacheValue)
	lruCash.items[key] = lruCash.queue.Front()
	if lruCash.queue.Len() > lruCash.capacity {
		last := lruCash.queue.Back()
		lastValue := last.Value.(CacheValue)
		delete(lruCash.items, lastValue.key)
		lruCash.queue.Remove(last)
	}
	return false
}

func (lruCash *lruCache) Get(key Key) (interface{}, bool) {
	value, exists := lruCash.items[key]
	if exists {
		lruCash.queue.MoveToFront(value)
		return value.Value.(CacheValue).value, true
	}
	return nil, false
}

func (lruCash *lruCache) Clear() {
	lruCash.queue = NewList()
	lruCash.items = make(map[Key]*ListItem, lruCash.capacity)
}

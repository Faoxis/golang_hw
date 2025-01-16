package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

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
	oldValue, exists := lruCash.items[key]
	if exists {
		lruCash.queue.MoveToFront(oldValue)
		oldValue.Value = value
		return true
	}
	lruCash.queue.PushFront(value)
	lruCash.items[key] = lruCash.queue.Front()
	if lruCash.queue.Len() > lruCash.capacity {
		last := lruCash.queue.Back()
		lruCash.queue.Remove(last)

	}
	return false
}

func (lruCash *lruCache) Get(key Key) (interface{}, bool) {

}

func (lruCash *lruCache) Clear() {

}

package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   string
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	v, ok := lc.items[key]
	if ok {
		v.Value = cacheItem{value: value, key: string(key)}
		lc.queue.MoveToFront(v)
		return true
	}

	if lc.queue.Len() < lc.capacity {
		el := lc.queue.PushFront(cacheItem{value: value, key: string(key)})
		lc.items[key] = el
		return false
	}

	if lc.queue.Len() == lc.capacity {
		back := lc.queue.Back()
		cItem := back.Value.(cacheItem)
		delete(lc.items, Key(cItem.key))
		lc.queue.Remove(back)
	}

	el := lc.queue.PushFront(cacheItem{value: value, key: string(key)})
	lc.items[key] = el
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	v, ok := lc.items[key]
	if !ok {
		return nil, false
	}

	lc.queue.MoveToFront(v)
	cItem := v.Value.(cacheItem)
	return cItem.value, true
}

func (lc *lruCache) Clear() {
	if lc.queue.Len() == 0 {
		return
	}
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

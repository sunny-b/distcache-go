package lrucache

import (
	"distcache-go/list"
	"sync"
)

/*
API:
New(size int)
Set(key string, value interface{})
Get(key) -> value interface{}

LRUCache:
- List
- HashMap map[key string] -> *Node
- MaxSize
*/

type LRUCache struct {
	list    *list.List
	items   map[string]*list.Node
	maxSize int
	mutex   sync.RWMutex
}

type element struct {
	key   string
	value interface{}
}

func New(size int) *LRUCache {
	return &LRUCache{
		maxSize: size,
		items:   make(map[string]*list.Node),
		list:    list.New(),
	}
}

func (l *LRUCache) Get(key string) interface{} {
	node := l.get(key)
	if node == nil {
		return nil
	}

	defer func() {
		l.list.MoveFront(node)
	}()

	ele := node.Value.(*element)

	return ele.value
}

func (l *LRUCache) Set(key string, value interface{}) interface{} {
	node := l.get(key)
	if node != nil {
		l.mutex.Lock()
		defer l.mutex.Unlock()
		node.Value.(*element).value = value
		return nil
	}

	ele := &element{
		key:   key,
		value: value,
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	node = l.list.PushFront(ele)

	l.items[key] = node

	tail := new(list.Node)
	if l.list.Length() > l.maxSize {
		tail = l.list.Tail()
		l.list.Remove(tail)

		delete(l.items, key)
	}

	if tail.Value == nil {
		return nil
	}

	return tail.Value.(*element).value
}

func (l *LRUCache) get(key string) *list.Node {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	node, exists := l.items[key]
	if !exists {
		return nil
	}

	return node
}

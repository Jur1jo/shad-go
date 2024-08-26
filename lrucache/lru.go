//go:build !solution

package lrucache

import "container/list"

type elem struct {
	key   int
	value int
}

type LRUCache struct {
	val   map[int]*list.Element
	items *list.List
	cap   int
	size  int
}

func (c *LRUCache) Get(key int) (int, bool) {
	if val, ok := c.val[key]; ok {
		c.items.MoveToFront(val)
		return val.Value.(*elem).value, ok
	}
	return 0, false
}

func (c *LRUCache) updateElement(e *list.Element, value int) {
	e.Value.(*elem).value = value // TODO
	c.items.MoveToFront(e)
}

func (c *LRUCache) Set(key, value int) {
	if val, ok := c.val[key]; ok {
		c.updateElement(val, value)
		c.items.MoveToFront(val)
	} else {
		c.items.PushFront(&elem{key: key, value: value})
		c.val[key] = c.items.Front()
		c.size++
		if c.size > c.cap {
			delete(c.val, c.items.Back().Value.(*elem).key)
			c.items.Remove(c.items.Back())
		}
	}
}

func (c *LRUCache) Clear() {
	c.size = 0
	c.items = list.New()
	c.val = make(map[int]*list.Element)
}

func (c *LRUCache) Range(f func(key, value int) bool) {
	for v := c.items.Back(); v != nil; {
		ok := f(v.Value.(*elem).key, v.Value.(*elem).value)
		if !ok {
			return
		}
		v = v.Prev()
	}
}

func New(cap int) Cache {
	return &LRUCache{cap: cap, val: make(map[int]*list.Element), items: list.New()}
}

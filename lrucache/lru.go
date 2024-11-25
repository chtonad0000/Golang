//go:build !solution

package lrucache

import "container/list"

type CacheMap struct {
	capacity   int
	storage    map[int]int
	ptrStamps  map[int]*list.Element
	timeStamps *list.List
}

func (c *CacheMap) Get(key int) (int, bool) {
	v, ok := c.storage[key]
	if ok {
		c.timeStamps.MoveToFront(c.ptrStamps[key])
	}
	return v, ok
}

func (c *CacheMap) Set(key, value int) {
	if c.capacity == 0 {
		return
	}
	if v, ok := c.ptrStamps[key]; ok {
		c.timeStamps.MoveToFront(v)
		c.storage[key] = value
	} else {
		if len(c.storage) == c.capacity {
			delete(c.storage, c.timeStamps.Back().Value.(int))
			delete(c.ptrStamps, c.timeStamps.Back().Value.(int))
			c.timeStamps.Remove(c.timeStamps.Back())
		}
		c.ptrStamps[key] = c.timeStamps.PushFront(key)
		c.storage[key] = value
	}
}

func (c *CacheMap) Range(f func(key, value int) bool) {
	for i := c.timeStamps.Back(); i != nil; i = i.Prev() {
		if !f(i.Value.(int), c.storage[i.Value.(int)]) {
			break
		}
	}
}

func (c *CacheMap) Clear() {
	c.storage = make(map[int]int, c.capacity)
	c.ptrStamps = make(map[int]*list.Element, c.capacity)
	c.timeStamps.Init()
}

func New(cap int) Cache {
	return &CacheMap{capacity: cap, storage: make(map[int]int, cap), ptrStamps: make(map[int]*list.Element, cap), timeStamps: list.New()}
}

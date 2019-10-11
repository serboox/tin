package main

import (
	"fmt"
	"sync"
)

type (
	Cache struct {
		cap int
		len int

		firstItem *Item
		lastItem  *Item

		dict  map[string]*Item
		mutex *sync.Mutex
	}

	Item struct {
		key    string
		value  int
		parent *Item
		child  *Item
	}
)

func NewCache(cap int) *Cache {
	return &Cache{
		cap:       cap,
		len:       0,
		firstItem: nil,
		lastItem:  nil,
		dict:      make(map[string]*Item),
		mutex:     new(sync.Mutex),
	}
}

func NewItem(key string, value int) *Item {
	return &Item{
		key:    key,
		value:  value,
		parent: nil,
		child:  nil,
	}
}

func (c *Cache) Add(key string, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item := NewItem(key, value)
	if c.firstItem == nil {
		c.firstItem = item
		c.lastItem = item
		c.len++
		c.dict[key] = item
		return
	}
	if c.cap == c.len {
		lastKey := c.lastItem.key
		lastItem := c.lastItem
		c.lastItem = c.lastItem.child
		lastItem.parent = nil
		lastItem.child = nil
		c.lastItem.parent = nil


		item.parent = c.firstItem
		c.firstItem.child = item
		c.firstItem = item
		delete(c.dict, lastKey)
		c.dict[key] = item
		return
	}
	c.firstItem.child = item
	item.parent = c.firstItem
	c.firstItem = item
	c.len++
	c.dict[key] = item
}

func (c *Cache) Get(key string) (res int, val error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if item, ok := c.dict[key]; ok {
		return item.value, nil
	}
	return -1, fmt.Errorf("%s not be found", key)
}

func (c *Cache) GetOrderedKeySlice() (res []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.len == 0 || c.firstItem == nil {
		return nil
	}
	res = make([]string, 0, c.cap)
	for item := c.lastItem; item != nil; item = item.child {
		res = append(res, item.key)
	}
	return
}

func (c *Cache) GetOrderedValuesSlice() (res []int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.len == 0 || c.firstItem == nil {
		return nil
	}
	res = make([]int, 0, c.cap)
	for item := c.lastItem; item != nil; item = item.child {
		res = append(res, item.value)
	}
	return
}

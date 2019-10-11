package main

import (
	"github.com/hashicorp/packer/common/random"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache(3)
	cache.Add("Sergei", 14)
	cache.Add("Roma", 12)
	cache.Add("Peter", 3)
	val, err := cache.Get("Roma")
	if val != 12 || err != nil {
		t.Fatalf("Value is not equal: %#v != %#v ", val, 12)
	}
	cache.Add("Masha", 4)
	val, err = cache.Get("Masha")
	if val != 4 || err != nil {
		t.Fatalf("Value is not equal: %#v != %#v ", val, 12)
	}
	_, err = cache.Get("Sergei")
	if err == nil {
		t.Fatal("Error is not received")
	}
}

func TestGetOrderedKeySlice(t *testing.T) {
	data := []string{
		"Sergei",
		"Roma",
		"Peter",
		"Masha",
	}
	cache := NewCache(3)
	cache.Add(data[0], 14)
	cache.Add(data[1], 12)
	cache.Add(data[2], 3)
	cache.Add(data[3], 4)
	for i, key := range cache.GetOrderedKeySlice() {
		if data[i+1] != key {
			t.Fatalf("Key is not equal: %#v != %#v ", data[i], key)
		}
	}
}

func TestGetOrderedValuesSlice(t *testing.T) {
	data := []int{
		14,
		12,
		3,
		4,
	}
	cache := NewCache(3)
	cache.Add("Sergei", data[0])
	cache.Add("Sergei", data[1])
	cache.Add("Peter", data[2])
	cache.Add("Masha", data[3])
	for i, key := range cache.GetOrderedValuesSlice() {
		if data[i+1] != key {
			t.Fatalf("Value is not equal: %#v != %#v ", data[i], key)
		}
	}
}

func TestCacheConcurrent(t *testing.T) {
	capacity := 10
	length := 100
	cache := NewCache(capacity)
	mutex := new(sync.Mutex)
	keySlice := make([]string, 0, length)
	valSlice := make([]int, 0, length)
	wg := new(sync.WaitGroup)
	wg.Add(length)
	for i := 1; i <= length; i++ {
		go func(i int) {
			defer wg.Done()
			mutex.Lock()
			key := random.String(random.PossibleLowerCase, 5)
			cache.Add(key, i)
			keySlice = append(keySlice, key)
			valSlice = append(valSlice, i)
			mutex.Unlock()
		}(i)
	}
	wg.Wait()

	j := capacity - 1
	cacheKeySlice := cache.GetOrderedKeySlice()
	cacheValueSlice := cache.GetOrderedValuesSlice()
	for i := length - 1; i >= length-capacity; i-- {
		if cacheKeySlice[j] != keySlice[i] {
			t.Fatalf(
				"Key is not equal: %#v != %#v ",
				cacheKeySlice[i],
				keySlice[j],
			)
			return
		}
		if cacheValueSlice[j] != valSlice[i] {
			t.Fatalf(
				"Value is not equal: %#v != %#v ",
				cacheKeySlice[i],
				keySlice[j],
			)
			return
		}
		j--
	}
}

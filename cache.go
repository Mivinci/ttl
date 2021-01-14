package ttl

import (
	"errors"
	"time"
)

var (
	ErrDup              = errors.New("key exists")
	ErrExpire           = errors.New("key has expired")
	ErrNotFound         = errors.New("key does not exists")
	ErrNegativeDuration = errors.New("duration cannot be negative")
)

type entry struct {
	dl  int64
	val interface{}
}

type Cache struct {
	c     map[interface{}]*entry
	Evict func(k interface{}, v interface{})
}

func New() *Cache {
	return &Cache{
		c: make(map[interface{}]*entry),
	}
}

// Get gets a value from the cache
func (c *Cache) Get(key interface{}) (interface{}, error) {
	if ent, ok := c.c[key]; ok {
		if ent.dl == -1 || ent.dl > time.Now().Unix() {
			return ent.val, nil
		}
		delete(c.c, key)
		if c.Evict != nil {
			c.Evict(key, ent.val)
		}
		return nil, ErrExpire
	}
	return nil, ErrNotFound
}

// Add adds a new pair to the cache and an error will be returned if the key exists.
func (c *Cache) Add(key interface{}, value interface{}, d time.Duration) error {
	if d < 0 {
		return ErrNegativeDuration
	}
	if _, ok := c.c[key]; ok {
		return ErrDup
	}
	c.c[key] = &entry{
		dl:  time.Now().Add(d).Unix(),
		val: value,
	}
	return nil
}

// Set sets a new pair to the cache if the given key does not exist, if
// the given key does exist, updates its value and sets a new expiration
// interval if `d` is not smaller than zero whereby the key can be expired
// at once if `d` equals to zero.
func (c *Cache) Set(key interface{}, value interface{}, d time.Duration) error {
	if ent, ok := c.c[key]; ok {
		ent.val = value
		if d >= 0 {
			ent.dl = time.Now().Add(d).Unix()
		}
		return nil
	}
	c.c[key] = &entry{
		dl:  time.Now().Add(d).Unix(),
		val: value,
	}
	return nil
}

// Expire sets a new expiration interval for an existing key
func (c *Cache) Expire(key interface{}, d time.Duration) error {
	if ent, ok := c.c[key]; ok {
		ent.dl = time.Now().Add(d).Unix()
		return nil
	}
	return ErrNotFound
}

// Remove removes an existing pair from the cache
func (c *Cache) Remove(key interface{}) error {
	if ent, ok := c.c[key]; ok {
		delete(c.c, key)
		if c.Evict != nil {
			c.Evict(key, ent.val)
		}
		return nil
	}
	return ErrNotFound
}

// GetAndRemove gets a value from the cache and removes it immediately only
// if no error occurs.
func (c *Cache) GetAndRemove(key interface{}) (value interface{}, err error) {
	if value, err = c.Get(key); err != nil {
		return
	}
	delete(c.c, key)
	if c.Evict != nil {
		c.Evict(key, value)
	}
	return
}

package ttl

import "time"

var defaultCache = New()

// Get gets a value from the cache
func Get(key interface{}) (interface{}, error) {
	return defaultCache.Get(key)
}

// Add adds a new pair to the cache and an error will be returned if the key exists.
func Add(key interface{}, value interface{}, d time.Duration) error {
	return defaultCache.Add(key, value, d)
}

// Set sets a new pair to the cache if the given key does not exist, if
// the given key does exist, updates its value and sets a new expiration
// interval if `d` is not smaller than zero whereby the key can be expired
// at once if `d` equals to zero.
func Set(key interface{}, value interface{}, d time.Duration) error {
	return defaultCache.Set(key, value, d)
}

// Expire sets a new expiration interval for an existing key
func Expire(key interface{}, d time.Duration) error {
	return defaultCache.Expire(key, d)
}

// Remove removes an existing pair from the cache
func Remove(key interface{}) error {
	return defaultCache.Remove(key)
}

// GetAndRemove gets a value from the cache and removes it immediately only
// if no error occurs.
func GetAndRemove(key interface{}) (value interface{}, err error) {
	return defaultCache.GetAndRemove(key)
}

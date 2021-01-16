package ttl

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := New()
	v, e := c.Get("a")
	assert.ErrorIs(t, ErrNotFound, e)
	assert.Nil(t, v)
	assert.Nil(t, c.Add("a", 1, time.Second))
	v, e = c.Get("a")
	assert.Nil(t, e)
	assert.Equal(t, v, 1)
	time.Sleep(time.Second)
	v, e = c.Get("a")
	assert.Nil(t, v)
	assert.ErrorIs(t, ErrExpire, e)
	assert.Nil(t, c.Add("a", 1, time.Second))
	assert.Nil(t, c.Remove("a"))
	v, e = c.Get("a")
	assert.Nil(t, v)
	assert.ErrorIs(t, ErrNotFound, e)
	assert.Nil(t, c.Add("a", 1, -1))
	v, e = c.Get("a")
	assert.Nil(t, e)
	assert.Equal(t, v, 1)
}

func TestSet(t *testing.T) {
	c := New()
	assert.Nil(t, c.Add("a", 1, time.Second))
	v, e := c.Get("a")
	assert.Nil(t, e)
	assert.Equal(t, 1, v)
	assert.Nil(t, c.Set("a", 2, -1))
	v, e = c.Get("a")
	assert.Nil(t, e)
	assert.Equal(t, 2, v)
	assert.Nil(t, c.Set("a", 3, 0))
	v, e = c.Get("a")
	assert.Nil(t, v)
	assert.ErrorIs(t, ErrExpire, e)
}

func TestExpire(t *testing.T) {
	c := New()
	assert.Nil(t, c.Add("a", 1, time.Second))
	assert.Nil(t, c.Expire("a", 2*time.Second))
	time.Sleep(time.Second)
	v, e := c.Get("a")
	assert.Nil(t, e)
	assert.Equal(t, 1, v)
	time.Sleep(time.Second)
	v, e = c.Get("a")
	assert.Nil(t, v)
	assert.ErrorIs(t, ErrExpire, e)
}

func TestGetAndRemove(t *testing.T) {
	c := New()
	assert.Nil(t, c.Add("a", 1, time.Second))
	v, e := c.GetAndRemove("a")
	assert.Equal(t, 1, v)
	assert.Nil(t, e)
	v, e = c.Get("a")
	assert.Nil(t, v)
	assert.ErrorIs(t, ErrNotFound, e)
}

func Example() {
	c := New()
	c.Evict = func(k, v interface{}) {
		fmt.Println(k, v)
	}
	c.Add("a", 1, 0) // nolint: errcheck
	c.Remove("a")    // nolint: errcheck

	// Output:
	// a 1
}

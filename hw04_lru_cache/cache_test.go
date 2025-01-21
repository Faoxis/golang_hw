package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge element by Set method", func(t *testing.T) {
		// Write me
		cache := NewCache(2)

		wasInCache := cache.Set("first", 1)
		require.False(t, wasInCache)

		wasInCache = cache.Set("second", 2)
		require.False(t, wasInCache)

		wasInCache = cache.Set("third", 3)
		require.False(t, wasInCache)

		value, ok := cache.Get("first")
		require.False(t, ok)
		require.Nil(t, value)

		cache.Clear()
		value, ok = cache.Get("first")
		require.False(t, ok)
		require.Nil(t, value)

		value, ok = cache.Get("second")
		require.False(t, ok)
		require.Nil(t, value)

		value, ok = cache.Get("third")
		require.False(t, ok)
		require.Nil(t, value)
	})

	t.Run("purge element by Get method", func(t *testing.T) {
		cache := NewCache(2)

		cache.Set("first", 1)
		cache.Set("second", 2)

		cache.Get("first")

		cache.Set("third", 3)

		value, ok := cache.Get("second")
		require.False(t, ok)
		require.Nil(t, value)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

package validator

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type field struct {
	idx  int
	name string
	tags string
}

// fieldCache stores cached fields.
type fieldCache struct {
	value atomic.Value // map[reflect.Type][]field
	mu    sync.Mutex   // used only by writers
}

func (c *fieldCache) get(t reflect.Type) (f []field, ok bool) {
	m, _ := c.value.Load().(map[reflect.Type][]field)
	f, ok = m[t]
	return
}

func (c *fieldCache) save(t reflect.Type, f []field) {
	c.mu.Lock()
	m, _ := c.value.Load().(map[reflect.Type][]field)
	newM := make(map[reflect.Type][]field, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[t] = f
	c.value.Store(newM)
	c.mu.Unlock()
}

// intCache is a cache of integer with string key.
type intCache struct {
	value atomic.Value // map[string]int64
	mu    sync.Mutex
}

func (c *intCache) get(s string) (val int64, ok bool) {
	m, _ := c.value.Load().(map[string]int64)
	val, ok = m[s]
	return
}

func (c *intCache) save(s string, val int64) {
	c.mu.Lock()
	m, _ := c.value.Load().(map[string]int64)
	newM := make(map[string]int64, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[s] = val
	c.value.Store(newM)
	c.mu.Unlock()
}

// uintCache is a cache of unsigned integer with string key.
type uintCache struct {
	value atomic.Value // map[string]uint
	mu    sync.Mutex
}

func (c *uintCache) get(s string) (val uint64, ok bool) {
	m, _ := c.value.Load().(map[string]uint64)
	val, ok = m[s]
	return
}

func (c *uintCache) save(s string, val uint64) {
	c.mu.Lock()
	m, _ := c.value.Load().(map[string]uint64)
	newM := make(map[string]uint64, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[s] = val
	c.value.Store(newM)
	c.mu.Unlock()
}

// floatCache is a cache of float with string key.
type floatCache struct {
	value atomic.Value // map[string]float64
	mu    sync.Mutex
}

func (c *floatCache) get(s string) (val float64, ok bool) {
	m, _ := c.value.Load().(map[string]float64)
	val, ok = m[s]
	return
}

func (c *floatCache) save(s string, val float64) {
	c.mu.Lock()
	m, _ := c.value.Load().(map[string]float64)
	newM := make(map[string]float64, len(m)+1)
	for k, v := range m {
		newM[k] = v
	}
	newM[s] = val
	c.value.Store(newM)
	c.mu.Unlock()
}

package util

import (
	"encoding/json"
	"sync"
)

// Cell holds a value protected by an RWMutex.
type Cell[T any] struct {
	mu    sync.RWMutex
	v     T
	after func()
}

func New[T any](initial T) *Cell[T] {
	return &Cell[T]{v: initial}
}

// After sets a hook that is called after writes (optional).
func (c *Cell[T]) After(fn func()) {
	c.mu.Lock()
	c.after = fn
	c.mu.Unlock()
}

func (c *Cell[T]) Get() T {
	c.mu.RLock()
	v := c.v
	c.mu.RUnlock()
	return v
}

func (c *Cell[T]) Set(v T) {
	c.mu.Lock()
	c.v = v
	after := c.after
	c.mu.Unlock()

	if after != nil {
		after()
	}
}

// Read runs fn while holding a read lock.
// Use it when you want to read multiple fields of T consistently.
func (c *Cell[T]) Read(fn func(v T)) {
	c.mu.RLock()
	fn(c.v)
	c.mu.RUnlock()
}

// Write runs fn while holding a write lock, so you can update in place.
// The hook (if set) fires after the write completes.
func (c *Cell[T]) Write(fn func(v *T)) {
	c.mu.Lock()
	fn(&c.v)
	after := c.after
	c.mu.Unlock()

	if after != nil {
		after()
	}
}

// Optional: make the cell JSON-friendly if you want to embed it.
func (c *Cell[T]) MarshalJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return json.Marshal(c.v)
}

func (c *Cell[T]) UnmarshalJSON(b []byte) error {
	var tmp T
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	c.mu.Lock()
	c.v = tmp
	after := c.after
	c.mu.Unlock()

	if after != nil {
		after()
	}
	return nil
}

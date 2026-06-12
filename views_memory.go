package main

import (
	"context"
	"sync"
)

type memoryViewCounter struct {
	mu     sync.Mutex
	counts map[string]int64
}

func newMemoryViewCounter() *memoryViewCounter {
	return &memoryViewCounter{counts: make(map[string]int64)}
}

func (v *memoryViewCounter) Record(_ context.Context, pageKey string) (int64, error) {
	key := viewDocID(todayVNDate(), pageKey)
	v.mu.Lock()
	defer v.mu.Unlock()
	v.counts[key]++
	return v.counts[key], nil
}

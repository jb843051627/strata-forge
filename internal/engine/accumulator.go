package engine

import "sync"

type Accumulator struct {
	mu     sync.RWMutex
	values map[int64]float64
}

func NewAccumulator() *Accumulator {
	return &Accumulator{values: make(map[int64]float64)}
}

func (a *Accumulator) Add(key int64, value float64) {
	a.mu.Lock()
	a.values[key] += value
	a.mu.Unlock()
}

func (a *Accumulator) Get(key int64) float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.values[key]
}

func (a *Accumulator) Snapshot() map[int64]float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	result := make(map[int64]float64, len(a.values))
	for key, value := range a.values {
		result[key] = value
	}
	return result
}

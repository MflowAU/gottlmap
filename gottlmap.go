package gottlmap

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	// ErrKeyNotFound is returned when the key is not found
	ErrTickerNotSet = fmt.Errorf("cleanupTicker not set")
	ErrCtxNotSet    = fmt.Errorf("ctx not set")
)

type Hook func(string, Element) error

// TTLMap is a container of map[string]*Element but with expirable Items.
type TTLMap struct {
	data map[string]*Element
	mu   *sync.RWMutex
	// cleanupTicker signals how often the expired elements are cleaned up
	cleanupTicker  *time.Ticker
	cleanupPreHook Hook
	ctx            context.Context
}

// Element is a value that expires after a given Time
type Element struct {
	Value interface{}
	Ex    time.Time
}

// New return a new TTLMap with the given cleanupTicker and cleanupPreHook
func New(t *time.Ticker, f Hook, ctx context.Context) (*TTLMap, error) {
	if t == nil {
		return nil, ErrTickerNotSet
	}
	if ctx == nil {
		return nil, ErrCtxNotSet
	}
	tm := &TTLMap{
		data:          make(map[string]*Element),
		mu:            &sync.RWMutex{},
		cleanupTicker: t,
		//cleanupPreHook should only work on copy of the data map
		cleanupPreHook: f,
		ctx:            ctx,
	}

	tm.startCleanupRoutine()
	return tm, nil
}

// Set the value for the given key. If the key already exists, the value
// will be overwritten. If the key does not exist, it will be created.
// The key will expire after the given ttl.
// The duration must be greater than 0.
// ttl is in seconds.
func (m *TTLMap) Set(key string, value interface{}, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ex := time.Now().Add(ttl)
	m.data[key] = &Element{value, ex}
}

// Get returns the value for the given key. If the key does not exist,
// an empty Element is returned along with a false boolean.
func (m *TTLMap) Get(key string) (Element, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a := Element{}
	if a, ok := m.data[key]; ok {
		return *a, true
	}
	return a, false
}

// Keys will return all keys of the TTLMap.
// The keys are returned are a snapshot of the keys at the time of the call.
// If the keys are modified after the call, the returned keys will not be updated.
func (m *TTLMap) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys
}

// Values will return all Values of the TTLMap
// The values are returned are a snapshot of the values at the time of the call.
func (m *TTLMap) Values() []Element {
	m.mu.RLock()
	defer m.mu.RUnlock()
	values := make([]Element, 0, len(m.data))
	for _, v := range m.data {
		values = append(values, *v)
	}
	return values
}

// GetDataCopy returns a copy of the data map
func (m *TTLMap) GetDataCopy() map[string]Element {
	m.mu.Lock()
	defer m.mu.Unlock()
	x := make(map[string]Element)
	for k, v := range m.data {
		x[k] = *v
	}
	return x
}

// startCleanupRoutine starts the cleanup routine
func (m *TTLMap) startCleanupRoutine() {
	fmt.Println("Starting cleanup routine")
	go func() {
		for {
			select {
			case <-m.cleanupTicker.C:
				fmt.Println("Cleanup routine ticked")
				m.cleanup()
			case <-m.ctx.Done():
				fmt.Println("Cleanup routine stopped")
				return
			}
		}
	}()
}

// cleanup removes all expired Elements from the map
// and invokes the cleanupPreHook if it is set
func (m *TTLMap) cleanup() {
	// Below will prevent looking at entries in the current minute
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if v.Ex.Before(now) {
			if m.cleanupPreHook != nil {
				err := m.cleanupPreHook(k, *v)
				if err != nil {
					log.Println(err)
				}
			}
			delete(m.data, k)
		}
	}
}

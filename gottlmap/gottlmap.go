package goTTLMap

import (
	"log"
	"sync"
	"time"
)

type Hook func(string, Element) error

type TTLMap struct {
	data map[string]*Element
	mu   *sync.RWMutex
	// cleanupTicker signals how often the expired elements are cleaned up
	cleanupTicker  *time.Ticker
	cleanupPreHook Hook
}

type Element struct {
	Value interface{}
	Ex    int64
}

func NewTTLMap(t *time.Ticker, f Hook) *TTLMap {
	tm := &TTLMap{
		data:          make(map[string]*Element),
		mu:            &sync.RWMutex{},
		cleanupTicker: t,
		//cleanupPreHook should only work on copy of the data map
		cleanupPreHook: f,
	}

	tm.startCleanupRoutine()
	return tm
}

// Set the value for the given key. If the key already exists, the value
// will be overwritten. If the key does not exist, it will be created.
// The key will expire after the given ttl.
// The duration must be greater than 0.
// ttl is in seconds.
func (m *TTLMap) Set(key string, value interface{}, ttl int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ex := time.Now().Unix() + ttl
	m.data[key] = &Element{value, ex}
}

func (m *TTLMap) Get(key string) (Element, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a := Element{}
	if a, ok := m.data[key]; ok {
		return *a, true
	}
	return a, false
}

// Keys will return all keys of the TTLMap
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
func (m *TTLMap) Values() []Element {
	m.mu.RLock()
	defer m.mu.RUnlock()
	values := make([]Element, 0, len(m.data))
	for _, v := range m.data {
		values = append(values, *v)
	}
	return values
}

func (m *TTLMap) GetDataCopy() map[string]Element {
	m.mu.Lock()
	defer m.mu.Unlock()
	x := make(map[string]Element)
	for k, v := range m.data {
		x[k] = *v
	}
	return x
}

// Startcleanup starts the cleanup routine
func (m *TTLMap) startCleanupRoutine() {
	go func() {
		for range m.cleanupTicker.C {
			m.cleanup()
		}
	}()
}

// cleanup removes all expired Elements from the map
func (m *TTLMap) cleanup() {
	// Below will prevent looking at entries in the current minute
	now := time.Now().Truncate(time.Minute).Unix()
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if v.Ex < now {
			if m.cleanupPreHook != nil {
				err := m.cleanupPreHook(k, *v)
				if err != nil {
					log.Printf("Error while executing cleanup hook %s, %v, %s \n", k, *v, err)
					log.Println(err)
					// continue
				}
			}
			delete(m.data, k)
		}
	}
}

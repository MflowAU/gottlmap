package gottlmap_test

import (
	"context"
	"testing"
	"time"

	"github.com/MFlowAU/gottlmap"
)

func TestNew(t *testing.T) {
	type testScenario struct {
		Description string
		KeyValues   []map[string]interface{}
		Tickers     *time.Ticker
		Hook        gottlmap.Hook
		ctx         context.Context
		err         error
	}

	tests := []testScenario{
		{
			Description: "New() should return a new TTLMap with out Errors",
			KeyValues: []map[string]interface{}{
				{"key1": "value1"},
				{"key2": "value2"},
				{"key3": "value3"},
			},
			Tickers: time.NewTicker(1 * time.Second),
			Hook:    nil,
			ctx:     context.Background(),
			err:     nil,
		},
		{
			Description: "New() should return error ErrTickerNotSet",
			KeyValues:   []map[string]interface{}{},
			Tickers:     nil,
			Hook:        nil,
			ctx:         context.Background(),
			err:         gottlmap.ErrTickerNotSet,
		},
		{
			Description: "New() should return error ErrCtxNotSet",
			KeyValues:   []map[string]interface{}{},
			Tickers:     time.NewTicker(1 * time.Second),
			Hook:        nil,
			ctx:         nil,
			err:         gottlmap.ErrCtxNotSet,
		},
	}

	for _, test := range tests {
		ttl_map, err := gottlmap.New(test.Tickers, test.Hook, test.ctx)
		if err != test.err {
			t.Errorf("%s: %s", test.Description, err)
			t.Fail()
		}
		if ttl_map == nil {
			continue
		}
		for _, keyValue := range test.KeyValues {
			for key, value := range keyValue {
				ttl_map.Set(key, value, 10*time.Second)
			}
		}
		if len(ttl_map.Keys()) != len(test.KeyValues) {
			t.Errorf("%s: Expected %d keys, got %d", test.Description, len(test.KeyValues), len(ttl_map.Keys()))
		}
	}
}

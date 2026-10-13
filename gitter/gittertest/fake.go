// Package gittertest provides a fake gitter.Querier for tests.
package gittertest

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Fake answers every query by decoding Data (the JSON of a GraphQL "data"
// field) into the response, or by returning Err. It records each query.
type Fake struct {
	Data string
	Err  error

	mu      sync.Mutex
	queries []string
}

func (f *Fake) Query(ctx context.Context, query string, response any) error {
	f.mu.Lock()
	f.queries = append(f.queries, query)
	f.mu.Unlock()

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("fake query: %w", err)
	}

	if f.Err != nil {
		return f.Err
	}

	if f.Data == "" {
		return nil
	}

	err = json.Unmarshal([]byte(f.Data), response)
	if err != nil {
		return fmt.Errorf("fake query decoding: %w", err)
	}

	return nil
}

// Queries returns a copy of the queries received so far.
func (f *Fake) Queries() []string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]string(nil), f.queries...)
}

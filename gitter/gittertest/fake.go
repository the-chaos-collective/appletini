// Package gittertest provides a fake gitter.Querier for tests.
package gittertest

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Fake answers every query or mutation by decoding Data (the JSON of a
// GraphQL "data" field) into the response, or by returning Err. It records
// each document and the variables sent with it.
type Fake struct {
	Data string
	Err  error

	mu      sync.Mutex
	queries []string
	vars    []map[string]any
}

func (f *Fake) Query(ctx context.Context, query string, response any) error {
	return f.answer(ctx, query, nil, response)
}

func (f *Fake) Mutate(ctx context.Context, mutation string, vars map[string]any, response any) error {
	return f.answer(ctx, mutation, vars, response)
}

func (f *Fake) answer(ctx context.Context, document string, vars map[string]any, response any) error {
	f.mu.Lock()
	f.queries = append(f.queries, document)
	f.vars = append(f.vars, vars)
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

// Vars returns the variables sent with each document, in order (nil for queries).
func (f *Fake) Vars() []map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]map[string]any(nil), f.vars...)
}

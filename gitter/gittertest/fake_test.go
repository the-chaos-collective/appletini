package gittertest_test

import (
	"context"
	"errors"
	"testing"

	"appletini/gitter/gittertest"
)

func TestFake(t *testing.T) {
	errBoom := errors.New("boom")

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	type payload struct {
		Name string `json:"name"`
	}

	testCases := []struct {
		name         string
		fake         *gittertest.Fake
		ctx          context.Context
		expected     payload
		expectedErr  error
		expectAnyErr bool
	}{
		{name: "decodes data", fake: &gittertest.Fake{Data: `{"name":"x"}`}, ctx: context.Background(), expected: payload{Name: "x"}},
		{name: "empty data", fake: &gittertest.Fake{}, ctx: context.Background()},
		{name: "configured error", fake: &gittertest.Fake{Err: errBoom}, ctx: context.Background(), expectedErr: errBoom},
		{name: "cancelled", fake: &gittertest.Fake{Data: `{"name":"x"}`}, ctx: cancelled, expectedErr: context.Canceled},
		{name: "invalid json", fake: &gittertest.Fake{Data: `{`}, ctx: context.Background(), expectAnyErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got payload

			err := tc.fake.Query(tc.ctx, "query q", &got)

			switch {
			case tc.expectAnyErr:
				if err == nil {
					t.Fatal("expected an error")
				}
			case !errors.Is(err, tc.expectedErr):
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if got != tc.expected {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}

			if queries := tc.fake.Queries(); len(queries) != 1 || queries[0] != "query q" {
				t.Fatalf("recorded queries %v", queries)
			}
		})
	}
}

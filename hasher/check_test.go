package hasher_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/charmbracelet/log"

	"appletini/gitter"
	"appletini/hasher"
	"appletini/logging"
)

type prMap = map[string][]gitter.PullRequest

func newHasher() *hasher.Hasher {
	//exhaustruct:ignore
	return &hasher.Hasher{Logger: logging.NewCharm(log.New(io.Discard))}
}

func TestHash(t *testing.T) {
	//exhaustruct:ignore
	a := prMap{"repo_0": {{Number: 1}}, "label_0": {{Number: 2}}}
	//exhaustruct:ignore
	sameContent := prMap{"label_0": {{Number: 2}}, "repo_0": {{Number: 1}}}
	//exhaustruct:ignore
	different := prMap{"repo_0": {{Number: 1}}, "label_0": {{Number: 3}}}

	testCases := []struct {
		name  string
		left  prMap
		right prMap
		equal bool
	}{
		{name: "same map", left: a, right: a, equal: true},
		{name: "same content, different insertion order", left: a, right: sameContent, equal: true},
		{name: "different content", left: a, right: different, equal: false},
		{name: "empty vs nil", left: prMap{}, right: nil, equal: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Map iteration is random, so check repeatedly.
			for range 20 {
				if got := hasher.Hash(tc.left) == hasher.Hash(tc.right); got != tc.equal {
					t.Fatalf("hashes equal = %v, expected %v", got, tc.equal)
				}
			}
		})
	}
}

func TestCheck(t *testing.T) {
	//exhaustruct:ignore
	first := prMap{"repo_0": {{Number: 1}}}
	//exhaustruct:ignore
	second := prMap{"repo_0": {{Number: 2}}}

	type step struct {
		prs         prMap
		cancelled   bool
		expectSend  bool
		expectedErr error
	}

	testCases := []struct {
		name  string
		steps []step
	}{
		{
			name:  "first map is always sent",
			steps: []step{{prs: first, expectSend: true}},
		},
		{
			name:  "unchanged map is not resent",
			steps: []step{{prs: first, expectSend: true}, {prs: first, expectSend: false}},
		},
		{
			name:  "changed map is sent",
			steps: []step{{prs: first, expectSend: true}, {prs: second, expectSend: true}, {prs: first, expectSend: true}},
		},
		{
			name: "cancelled send is retried on the next check",
			steps: []step{
				{prs: first, cancelled: true, expectedErr: context.Canceled},
				{prs: first, expectSend: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newHasher()

			for i, s := range tc.steps {
				ctx, cancel := context.WithCancel(context.Background())

				var ch chan prMap
				if s.cancelled {
					ch = make(chan prMap) // unbuffered with no reader: only cancellation can unblock
					cancel()
				} else {
					ch = make(chan prMap, 1)
				}

				err := h.Check(ctx, s.prs, ch)
				cancel()

				if !errors.Is(err, s.expectedErr) {
					t.Fatalf("step %d: got error %v, expected %v", i, err, s.expectedErr)
				}

				sent := len(ch) == 1
				if sent != s.expectSend {
					t.Fatalf("step %d: sent = %v, expected %v", i, sent, s.expectSend)
				}
			}
		})
	}
}

func TestCheckNilChannel(t *testing.T) {
	err := newHasher().Check(context.Background(), prMap{}, nil)
	if !errors.Is(err, hasher.ErrNilChannel) {
		t.Fatalf("got error %v, expected %v", err, hasher.ErrNilChannel)
	}
}

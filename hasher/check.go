package hasher

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	"appletini/gitter"
	"appletini/logging"
)

var ErrNilChannel = errors.New("pull request channel must not be nil")

type Hasher struct {
	Logger      logging.Logger
	currentHash string
}

// Hash returns a stable digest of prMap. fmt prints map keys in sorted order,
// so equal maps always produce the same hash.
func Hash(prMap map[string][]gitter.PullRequest) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%v", prMap)))

	return fmt.Sprintf("%x", sum)
}

// Check sends prMap on prChannel only when it differs from the last map sent.
// The stored hash is updated only after a successful send, so a send aborted
// by ctx cancellation is retried on the next call.
func (hasher *Hasher) Check(
	ctx context.Context,
	prMap map[string][]gitter.PullRequest,
	prChannel chan<- map[string][]gitter.PullRequest,
) error {
	if prChannel == nil {
		return ErrNilChannel
	}

	newHash := Hash(prMap)
	hasher.Logger.Debug("Computed hash", "hash", newHash)

	if hasher.currentHash == newHash {
		return nil
	}

	select {
	case prChannel <- prMap:
		hasher.currentHash = newHash
		return nil
	case <-ctx.Done():
		return fmt.Errorf("sending pull requests: %w", ctx.Err())
	}
}

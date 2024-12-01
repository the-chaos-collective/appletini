package hasher

import (
	"crypto/sha256"
	"fmt"
	"log"

	"git_applet/gitter"
)

type Hasher struct {
	Logger      *log.Logger
	currentHash string
}

func (hasher *Hasher) Check(prMap map[string][]gitter.PullRequest, prChannel chan<- map[string][]gitter.PullRequest) {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", prMap)))
	newHash := fmt.Sprintf("%x", h.Sum(nil))
	hasher.Logger.Printf("HASH: %v", newHash)
	if hasher.currentHash != newHash {
		hasher.currentHash = newHash
		prChannel <- prMap
	}
}

package main

import (
	"crypto/sha256"
	"fmt"

	"git_applet/gitter"
)

func hashCheck(currentHash *string, prMap map[string][]gitter.PullRequest, prChannel chan<- map[string][]gitter.PullRequest) {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", prMap)))
	newHash := fmt.Sprintf("%x", h.Sum(nil))
	logger.Printf("HASH: %v", newHash)
	if currentHash != &newHash {
		currentHash = &newHash
		prChannel <- prMap
	}
}

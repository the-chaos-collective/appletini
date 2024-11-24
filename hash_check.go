package main

import (
	"crypto/sha256"
	"fmt"
)

func hashCheck(callback func(string)) {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", trackedPrs)))
	newHash := fmt.Sprintf("%x", h.Sum(nil))
	callback(newHash)
}

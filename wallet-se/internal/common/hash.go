package common

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hasher(word string) string {
	hasher := sha1.New()
	hasher.Write([]byte(word))
	return hex.EncodeToString(hasher.Sum(nil))
}

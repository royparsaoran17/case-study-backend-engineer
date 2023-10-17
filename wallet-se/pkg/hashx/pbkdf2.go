package hashx

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

func HashPBKDF2(plain, salt string, iter, keyLen int) string {
	hash := pbkdf2.Key([]byte(plain), []byte(salt), iter, keyLen, sha256.New)
	return hex.EncodeToString(hash)
}

package cryptox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func EncryptGCM(plain []byte, key string) ([]byte, error) {
	// Membuat blok cipher AES dari kunci
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed create chipher blok aes %v", err)
	}

	// Membuat Galois/Counter Mode (GCM) dengan blok cipher AES
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed create Galois/Counter Mode (GCM): %v", err)
	}

	// Membuat nonce acak (IV) dengan panjang yang disarankan oleh GCM (12 bytes)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed create nonce %v", err)
	}

	// Enkripsi pesan menggunakan mode GCM
	return gcm.Seal(nonce, nonce, plain, nil), nil

}

func DecryptGCM(ciphertext []byte, key string) (plain []byte, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed create cipher blok aes %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed create Galois/Counter Mode (GCM): %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Dekripsi pesan menggunakan mode GCM
	return gcm.Open(nil, nonce, ciphertext, nil)
}

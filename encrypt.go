package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

func generateKey(pin string, salt []byte) ([]byte, error) {
	if salt == nil {
		salt := make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, err
		}
	}

	key, err := scrypt.Key([]byte(pin), salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func Encrypt(pwd, pin string) ([]byte, error) {
	data := []byte(pwd)
	key, err := generateKey(pin, nil)
	if err != nil {
		return nil, err
	}

	// good resource: https://www.vpngeeks.com/advanced-encryption-explained/
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// good resource: https://www.cryptologie.net/article/277/what-is-gcm-galois-counter-mode/
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// nonce = number once used
	// piece of data that should not be repeated and only used once in combination with any particular key
	nonce := make([]byte, gcm.NonceSize())
	// populates nonce with a cryptographically secure random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func Decrypt(pin string, data []byte) (string, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, err := generateKey(pin, salt)
	if err != nil {
		return "", err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

var AESKey = "nullbeacon00000000000000000000000000000000000000000000000000000000"

func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

func xorBytes(data []byte, key []byte) []byte {
	keyLen := len(key)
	xored := make([]byte, len(data))
	for i := range data {
		xored[i] = data[i] ^ key[i%keyLen]
	}
	return xored
}

func xorKey() []byte {
	return []byte{0xDE, 0xAD, 0xBE, 0xEF, 0xCA, 0xFE, 0xBA, 0xBE}
}

func packPayload(payload []byte) ([]byte, error) {
	key := xorKey()
	xoredPayload := xorBytes(payload, key)
	encryptedPayload, err := encrypt(xoredPayload, aesKey())
	if err != nil {
		return nil, err
	}
	return encryptedPayload, nil
}

func unpackPayload(packedPayload []byte) ([]byte, error) {
	decryptedPayload, err := decrypt(packedPayload, aesKey())
	if err != nil {
		return nil, err
	}
	key := xorKey()
	unpackedPayload := xorBytes(decryptedPayload, key)
	return unpackedPayload, nil
}

func aesKey() []byte {
	key, _ := hex.DecodeString(AESKey)
	return key
}

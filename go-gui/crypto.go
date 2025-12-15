package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// Crypto keys for PS3 Syscon authentication.
var (
	sc2tb       = mustDecodeHex("71f03f184c01c5ebc3f6a22a42ba9525") // Syscon to TestBench Key
	tb2sc       = mustDecodeHex("907e730f4d4e0a0b7b75f030eb1d9d36") // TestBench to Syscon Key
	authValue   = mustDecodeHex("3350BD7820345C29056A223BA220B323") // 0x45B8
	zeroIV      = mustDecodeHex("00000000000000000000000000000000")
	auth1rHdr   = mustDecodeHex("10100000FFFFFFFF0000000000000000")
	auth2Header = mustDecodeHex("10010000000000000000000000000000")
)

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

// aesDecryptCBC decrypts data using AES-CBC.
func aesDecryptCBC(key, iv, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	return decrypted, nil
}

// aesEncryptCBC encrypts data using AES-CBC.
func aesEncryptCBC(key, iv, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("plaintext is not a multiple of block size")
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	mode.CryptBlocks(encrypted, data)

	return encrypted, nil
}

// bytesEqual compares two byte slices.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

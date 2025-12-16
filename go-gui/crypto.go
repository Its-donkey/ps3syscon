// Package main provides AES cryptographic functions for PS3 Syscon authentication.
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// Cryptographic keys and constants for PS3 Syscon authentication protocol.
var (
	// sc2tb is the Syscon to TestBench AES key.
	sc2tb = mustDecodeHex("71f03f184c01c5ebc3f6a22a42ba9525")

	// tb2sc is the TestBench to Syscon AES key.
	tb2sc = mustDecodeHex("907e730f4d4e0a0b7b75f030eb1d9d36")

	// authValue is the expected authentication value at offset 0x45B8.
	authValue = mustDecodeHex("3350BD7820345C29056A223BA220B323")

	// zeroIV is a zero initialization vector for AES-CBC.
	zeroIV = mustDecodeHex("00000000000000000000000000000000")

	// auth1rHdr is the expected header in AUTH1 response.
	auth1rHdr = mustDecodeHex("10100000FFFFFFFF0000000000000000")

	// auth2Header is the header prefix for AUTH2 command.
	auth2Header = mustDecodeHex("10010000000000000000000000000000")
)

// mustDecodeHex decodes a hex string and panics on error.
// This is used for compile-time constant initialization.
func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

// aesDecryptCBC decrypts data using AES-CBC mode.
// Returns ErrDecryptionFailed if the data length is not a multiple of block size.
func aesDecryptCBC(key, iv, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("%w: ciphertext is not a multiple of block size", ErrDecryptionFailed)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	return decrypted, nil
}

// aesEncryptCBC encrypts data using AES-CBC mode.
// Returns ErrEncryptionFailed if the data length is not a multiple of block size.
func aesEncryptCBC(key, iv, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEncryptionFailed, err)
	}

	if len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("%w: plaintext is not a multiple of block size", ErrEncryptionFailed)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	mode.CryptBlocks(encrypted, data)

	return encrypted, nil
}

// bytesEqual compares two byte slices for equality.
// This is a convenience wrapper around bytes.Equal.
func bytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}

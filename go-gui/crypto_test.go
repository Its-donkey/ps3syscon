package main

import (
	"encoding/hex"
	"testing"
)

func TestMustDecodeHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "valid hex string",
			input:    "deadbeef",
			expected: []byte{0xde, 0xad, 0xbe, 0xef},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []byte{},
		},
		{
			name:     "all zeros",
			input:    "00000000",
			expected: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "all ff",
			input:    "ffffffff",
			expected: []byte{0xff, 0xff, 0xff, 0xff},
		},
		{
			name:     "16 bytes (AES key size)",
			input:    "71f03f184c01c5ebc3f6a22a42ba9525",
			expected: []byte{0x71, 0xf0, 0x3f, 0x18, 0x4c, 0x01, 0xc5, 0xeb, 0xc3, 0xf6, 0xa2, 0x2a, 0x42, 0xba, 0x95, 0x25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mustDecodeHex(tt.input)
			if !bytesEqual(result, tt.expected) {
				t.Errorf("mustDecodeHex(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMustDecodeHexPanicsOnInvalidInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("mustDecodeHex did not panic on invalid hex")
		}
	}()
	mustDecodeHex("invalid_hex")
}

func TestMustDecodeHexPanicsOnOddLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("mustDecodeHex did not panic on odd length hex")
		}
	}()
	mustDecodeHex("abc")
}

func TestBytesEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        []byte
		b        []byte
		expected bool
	}{
		{
			name:     "equal slices",
			a:        []byte{1, 2, 3, 4},
			b:        []byte{1, 2, 3, 4},
			expected: true,
		},
		{
			name:     "different values",
			a:        []byte{1, 2, 3, 4},
			b:        []byte{1, 2, 3, 5},
			expected: false,
		},
		{
			name:     "different lengths",
			a:        []byte{1, 2, 3},
			b:        []byte{1, 2, 3, 4},
			expected: false,
		},
		{
			name:     "empty slices",
			a:        []byte{},
			b:        []byte{},
			expected: true,
		},
		{
			name:     "nil and empty",
			a:        nil,
			b:        []byte{},
			expected: true,
		},
		{
			name:     "both nil",
			a:        nil,
			b:        nil,
			expected: true,
		},
		{
			name:     "first byte different",
			a:        []byte{0, 2, 3, 4},
			b:        []byte{1, 2, 3, 4},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesEqual(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("bytesEqual(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestAesDecryptCBC(t *testing.T) {
	key := mustDecodeHex("00000000000000000000000000000000")
	iv := mustDecodeHex("00000000000000000000000000000000")

	// Test with known values - encrypt then decrypt should give original
	plaintext := mustDecodeHex("00000000000000000000000000000000")

	// First encrypt
	encrypted, err := aesEncryptCBC(key, iv, plaintext)
	if err != nil {
		t.Fatalf("aesEncryptCBC failed: %v", err)
	}

	// Then decrypt
	decrypted, err := aesDecryptCBC(key, iv, encrypted)
	if err != nil {
		t.Fatalf("aesDecryptCBC failed: %v", err)
	}

	if !bytesEqual(decrypted, plaintext) {
		t.Errorf("Round trip failed: got %x, want %x", decrypted, plaintext)
	}
}

func TestAesDecryptCBCWithRealKeys(t *testing.T) {
	// Use the actual keys from the code
	key := mustDecodeHex("71f03f184c01c5ebc3f6a22a42ba9525")
	iv := mustDecodeHex("00000000000000000000000000000000")

	// Test multiple blocks
	plaintext := mustDecodeHex("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff")

	encrypted, err := aesEncryptCBC(key, iv, plaintext)
	if err != nil {
		t.Fatalf("aesEncryptCBC failed: %v", err)
	}

	decrypted, err := aesDecryptCBC(key, iv, encrypted)
	if err != nil {
		t.Fatalf("aesDecryptCBC failed: %v", err)
	}

	if !bytesEqual(decrypted, plaintext) {
		t.Errorf("Round trip failed: got %x, want %x", decrypted, plaintext)
	}
}

func TestAesDecryptCBCInvalidBlockSize(t *testing.T) {
	key := mustDecodeHex("00000000000000000000000000000000")
	iv := mustDecodeHex("00000000000000000000000000000000")
	invalidData := []byte{1, 2, 3, 4, 5} // Not a multiple of block size

	_, err := aesDecryptCBC(key, iv, invalidData)
	if err == nil {
		t.Error("Expected error for invalid block size, got nil")
	}
}

func TestAesDecryptCBCInvalidKeySize(t *testing.T) {
	invalidKey := []byte{1, 2, 3, 4, 5} // Invalid key size
	iv := mustDecodeHex("00000000000000000000000000000000")
	data := mustDecodeHex("00000000000000000000000000000000")

	_, err := aesDecryptCBC(invalidKey, iv, data)
	if err == nil {
		t.Error("Expected error for invalid key size, got nil")
	}
}

func TestAesEncryptCBC(t *testing.T) {
	key := mustDecodeHex("00000000000000000000000000000000")
	iv := mustDecodeHex("00000000000000000000000000000000")
	plaintext := mustDecodeHex("00000000000000000000000000000000")

	encrypted, err := aesEncryptCBC(key, iv, plaintext)
	if err != nil {
		t.Fatalf("aesEncryptCBC failed: %v", err)
	}

	// Encrypted should be different from plaintext (unless plaintext is already "encrypted" form)
	if len(encrypted) != len(plaintext) {
		t.Errorf("Encrypted length %d != plaintext length %d", len(encrypted), len(plaintext))
	}
}

func TestAesEncryptCBCWithRealKeys(t *testing.T) {
	// Use the TB2SC key from the code
	key := mustDecodeHex("907e730f4d4e0a0b7b75f030eb1d9d36")
	iv := mustDecodeHex("00000000000000000000000000000000")

	// Test 48 bytes (3 blocks) like the auth code uses
	plaintext := mustDecodeHex("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	encrypted, err := aesEncryptCBC(key, iv, plaintext)
	if err != nil {
		t.Fatalf("aesEncryptCBC failed: %v", err)
	}

	if len(encrypted) != 48 {
		t.Errorf("Expected 48 bytes encrypted, got %d", len(encrypted))
	}
}

func TestAesEncryptCBCInvalidBlockSize(t *testing.T) {
	key := mustDecodeHex("00000000000000000000000000000000")
	iv := mustDecodeHex("00000000000000000000000000000000")
	invalidData := []byte{1, 2, 3, 4, 5} // Not a multiple of block size

	_, err := aesEncryptCBC(key, iv, invalidData)
	if err == nil {
		t.Error("Expected error for invalid block size, got nil")
	}
}

func TestAesEncryptCBCInvalidKeySize(t *testing.T) {
	invalidKey := []byte{1, 2, 3, 4, 5} // Invalid key size
	iv := mustDecodeHex("00000000000000000000000000000000")
	data := mustDecodeHex("00000000000000000000000000000000")

	_, err := aesEncryptCBC(invalidKey, iv, data)
	if err == nil {
		t.Error("Expected error for invalid key size, got nil")
	}
}

func TestCryptoKeyConstantsExist(t *testing.T) {
	// Test that all the crypto constants are properly initialized
	if len(sc2tb) != 16 {
		t.Errorf("sc2tb key should be 16 bytes, got %d", len(sc2tb))
	}
	if len(tb2sc) != 16 {
		t.Errorf("tb2sc key should be 16 bytes, got %d", len(tb2sc))
	}
	if len(authValue) != 16 {
		t.Errorf("authValue should be 16 bytes, got %d", len(authValue))
	}
	if len(zeroIV) != 16 {
		t.Errorf("zeroIV should be 16 bytes, got %d", len(zeroIV))
	}
	if len(auth1rHdr) != 16 {
		t.Errorf("auth1rHdr should be 16 bytes, got %d", len(auth1rHdr))
	}
	if len(auth2Header) != 16 {
		t.Errorf("auth2Header should be 16 bytes, got %d", len(auth2Header))
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	testCases := []struct {
		name      string
		key       []byte
		plaintext []byte
	}{
		{
			name:      "SC2TB key single block",
			key:       sc2tb,
			plaintext: mustDecodeHex("00112233445566778899aabbccddeeff"),
		},
		{
			name:      "TB2SC key single block",
			key:       tb2sc,
			plaintext: mustDecodeHex("ffeeddccbbaa99887766554433221100"),
		},
		{
			name:      "SC2TB key triple block",
			key:       sc2tb,
			plaintext: mustDecodeHex("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := aesEncryptCBC(tc.key, zeroIV, tc.plaintext)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			decrypted, err := aesDecryptCBC(tc.key, zeroIV, encrypted)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			if !bytesEqual(decrypted, tc.plaintext) {
				t.Errorf("Round trip failed:\noriginal:  %s\ndecrypted: %s",
					hex.EncodeToString(tc.plaintext),
					hex.EncodeToString(decrypted))
			}
		})
	}
}

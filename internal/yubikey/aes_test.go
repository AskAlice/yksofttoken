package yubikey

import (
	"testing"
)

func TestAESEncryptDecrypt(t *testing.T) {
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	plaintext := []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}

	// Encrypt
	ciphertext, err := AESEncrypt(key, plaintext)
	if err != nil {
		t.Fatalf("AESEncrypt returned error: %v", err)
	}

	// Ciphertext should be different from plaintext
	match := true
	for i := range plaintext {
		if plaintext[i] != ciphertext[i] {
			match = false
			break
		}
	}
	if match {
		t.Error("Ciphertext matches plaintext - encryption failed")
	}

	// Decrypt
	decrypted, err := AESDecrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("AESDecrypt returned error: %v", err)
	}

	// Decrypted should match original plaintext
	for i := range plaintext {
		if decrypted[i] != plaintext[i] {
			t.Errorf("Decryption mismatch at %d: got %02x, expected %02x",
				i, decrypted[i], plaintext[i])
		}
	}
}

func TestAESInvalidKeySize(t *testing.T) {
	invalidKey := []byte{0x00, 0x01, 0x02} // Too short
	plaintext := []byte{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}

	_, err := AESEncrypt(invalidKey, plaintext)
	if err == nil {
		t.Error("AESEncrypt should have returned error for invalid key size")
	}
}

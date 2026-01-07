package yubikey

import (
	"crypto/aes"
)

// AESEncrypt encrypts a 16-byte block with AES-128-ECB
func AESEncrypt(key, plaintext []byte) ([]byte, error) {
	if len(key) != 16 || len(plaintext) != 16 {
		return nil, aes.KeySizeError(len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, 16)
	block.Encrypt(ciphertext, plaintext)
	return ciphertext, nil
}

// AESDecrypt decrypts a 16-byte block with AES-128-ECB
func AESDecrypt(key, ciphertext []byte) ([]byte, error) {
	if len(key) != 16 || len(ciphertext) != 16 {
		return nil, aes.KeySizeError(len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, 16)
	block.Decrypt(plaintext, ciphertext)
	return plaintext, nil
}

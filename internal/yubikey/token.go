package yubikey

import (
	"encoding/binary"
)

const (
	// UIDSize is the size of the UID (private ID) in bytes
	UIDSize = 6
	// KeySize is the size of the AES key in bytes
	KeySize = 16
	// OTPSize is the size of the OTP block in bytes
	OTPSize = 16
	// PublicIDSize is the size of the public ID in bytes
	PublicIDSize = 6
)

// TokenBlock represents the internal structure of a Yubikey OTP
type TokenBlock struct {
	UID       [UIDSize]byte // Private ID (6 bytes)
	Counter   uint16        // Usage counter (2 bytes)
	Timestamp uint32        // 24-bit timestamp, stored in low, high order
	Session   uint8         // Session use counter (1 byte)
	Random    uint16        // Random value (2 bytes)
	CRC       uint16        // CRC16 (2 bytes)
}

// MarshalBinary encodes the token block to bytes
func (t *TokenBlock) MarshalBinary() []byte {
	data := make([]byte, OTPSize)

	// Copy UID
	copy(data[0:6], t.UID[:])

	// Counter (little-endian)
	binary.LittleEndian.PutUint16(data[6:8], t.Counter)

	// Timestamp low (16 bits, little-endian)
	binary.LittleEndian.PutUint16(data[8:10], uint16(t.Timestamp&0xffff))

	// Timestamp high (8 bits)
	data[10] = byte((t.Timestamp >> 16) & 0xff)

	// Session counter
	data[11] = t.Session

	// Random (little-endian)
	binary.LittleEndian.PutUint16(data[12:14], t.Random)

	// CRC (little-endian)
	binary.LittleEndian.PutUint16(data[14:16], t.CRC)

	return data
}

// UnmarshalBinary decodes bytes to the token block
func (t *TokenBlock) UnmarshalBinary(data []byte) error {
	if len(data) != OTPSize {
		return ErrInvalidLength
	}

	// Copy UID
	copy(t.UID[:], data[0:6])

	// Counter (little-endian)
	t.Counter = binary.LittleEndian.Uint16(data[6:8])

	// Timestamp (24 bits)
	t.Timestamp = uint32(binary.LittleEndian.Uint16(data[8:10])) |
		(uint32(data[10]) << 16)

	// Session counter
	t.Session = data[11]

	// Random (little-endian)
	t.Random = binary.LittleEndian.Uint16(data[12:14])

	// CRC (little-endian)
	t.CRC = binary.LittleEndian.Uint16(data[14:16])

	return nil
}

// ComputeCRC computes and sets the CRC for the token block
func (t *TokenBlock) ComputeCRC() {
	data := t.MarshalBinary()
	t.CRC = ^CRC16(data[:14])
}

// Generate generates an encrypted OTP from the token block and key
func (t *TokenBlock) Generate(key []byte) (string, error) {
	// Compute CRC
	t.ComputeCRC()

	// Marshal to bytes
	plaintext := t.MarshalBinary()

	// Encrypt with AES
	ciphertext, err := AESEncrypt(key, plaintext)
	if err != nil {
		return "", err
	}

	// Encode as modhex
	return ModHexEncode(ciphertext), nil
}

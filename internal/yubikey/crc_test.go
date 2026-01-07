package yubikey

import (
	"testing"
)

func TestCRC16Empty(t *testing.T) {
	// Empty data should return initial CRC value
	result := CRC16([]byte{})
	if result != 0xffff {
		t.Errorf("CRC16([]) = 0x%04x, expected 0xffff", result)
	}
}

func TestCRC16Residual(t *testing.T) {
	// When CRC is appended to data correctly, the CRC of the entire
	// block should equal CRCOKResidual
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e}
	crc := ^CRC16(data)

	// Append CRC (little-endian)
	fullData := append(data, byte(crc&0xff), byte(crc>>8))

	// Verify residual
	result := CRC16(fullData)
	if result != CRCOKResidual {
		t.Errorf("CRC residual = 0x%04x, expected 0x%04x", result, CRCOKResidual)
	}
}

func TestCRC16Consistency(t *testing.T) {
	// Same data should always produce the same CRC
	data := []byte{0x10, 0x20, 0x30, 0x40, 0x50}
	crc1 := CRC16(data)
	crc2 := CRC16(data)

	if crc1 != crc2 {
		t.Errorf("CRC16 inconsistent: 0x%04x != 0x%04x", crc1, crc2)
	}
}

func TestCRC16Different(t *testing.T) {
	// Different data should produce different CRCs (usually)
	data1 := []byte{0x10, 0x20, 0x30, 0x40, 0x50}
	data2 := []byte{0x10, 0x20, 0x30, 0x40, 0x51}
	crc1 := CRC16(data1)
	crc2 := CRC16(data2)

	if crc1 == crc2 {
		t.Errorf("CRC16 collision for different data: both 0x%04x", crc1)
	}
}

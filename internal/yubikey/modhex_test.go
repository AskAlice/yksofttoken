package yubikey

import (
	"testing"
)

func TestModHexEncode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0x00}, "cc"},
		{[]byte{0xff}, "vv"},
		{[]byte{0x22, 0x22}, "dddd"},
	}

	for _, tt := range tests {
		result := ModHexEncode(tt.input)
		if result != tt.expected {
			t.Errorf("ModHexEncode(%v) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestModHexDecode(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"cc", []byte{0x00}},
		{"vv", []byte{0xff}},
		{"dddd", []byte{0x22, 0x22}},
	}

	for _, tt := range tests {
		result, err := ModHexDecode(tt.input)
		if err != nil {
			t.Errorf("ModHexDecode(%s) returned error: %v", tt.input, err)
			continue
		}
		if len(result) != len(tt.expected) {
			t.Errorf("ModHexDecode(%s) length = %d, expected %d", tt.input, len(result), len(tt.expected))
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("ModHexDecode(%s)[%d] = %02x, expected %02x", tt.input, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestModHexDecodeInvalid(t *testing.T) {
	// Odd length
	_, err := ModHexDecode("ccc")
	if err == nil {
		t.Error("Expected error for odd-length modhex string")
	}

	// Invalid character
	_, err = ModHexDecode("cx")
	if err == nil {
		t.Error("Expected error for invalid modhex character")
	}
}

func TestHexEncode(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{0x00}, "00"},
		{[]byte{0xff}, "ff"},
		{[]byte{0x10, 0xa0, 0x35}, "10a035"},
	}

	for _, tt := range tests {
		result := HexEncode(tt.input)
		if result != tt.expected {
			t.Errorf("HexEncode(%v) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestHexDecode(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"00", []byte{0x00}},
		{"ff", []byte{0xff}},
		{"10a035", []byte{0x10, 0xa0, 0x35}},
		{"10A035", []byte{0x10, 0xa0, 0x35}}, // uppercase
	}

	for _, tt := range tests {
		result, err := HexDecode(tt.input)
		if err != nil {
			t.Errorf("HexDecode(%s) returned error: %v", tt.input, err)
			continue
		}
		if len(result) != len(tt.expected) {
			t.Errorf("HexDecode(%s) length = %d, expected %d", tt.input, len(result), len(tt.expected))
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("HexDecode(%s)[%d] = %02x, expected %02x", tt.input, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestHexDecodeInvalid(t *testing.T) {
	// Odd length
	_, err := HexDecode("abc")
	if err == nil {
		t.Error("Expected error for odd-length hex string")
	}

	// Invalid character
	_, err = HexDecode("zz")
	if err == nil {
		t.Error("Expected error for invalid hex character")
	}
}

func TestModHexRoundTrip(t *testing.T) {
	original := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}
	encoded := ModHexEncode(original)
	decoded, err := ModHexDecode(encoded)
	if err != nil {
		t.Fatalf("ModHexDecode returned error: %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("Round trip length mismatch: got %d, expected %d", len(decoded), len(original))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("Round trip mismatch at %d: got %02x, expected %02x", i, decoded[i], original[i])
		}
	}
}

func TestHexRoundTrip(t *testing.T) {
	original := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}
	encoded := HexEncode(original)
	decoded, err := HexDecode(encoded)
	if err != nil {
		t.Fatalf("HexDecode returned error: %v", err)
	}

	if len(decoded) != len(original) {
		t.Fatalf("Round trip length mismatch: got %d, expected %d", len(decoded), len(original))
	}

	for i := range original {
		if decoded[i] != original[i] {
			t.Errorf("Round trip mismatch at %d: got %02x, expected %02x", i, decoded[i], original[i])
		}
	}
}

// Package yubikey provides Yubikey encoding, encryption, and CRC functions
package yubikey

import (
	"errors"
	"strings"
)

// ModHex alphabet used by Yubikey
const modHexAlphabet = "cbdefghijklnrtuv"

var (
	modHexDecode = make(map[byte]byte)
	hexAlphabet  = "0123456789abcdef"
)

func init() {
	for i := 0; i < 16; i++ {
		modHexDecode[modHexAlphabet[i]] = byte(i)
	}
}

// ModHexEncode encodes a byte slice to modhex string
func ModHexEncode(data []byte) string {
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = modHexAlphabet[b>>4]
		result[i*2+1] = modHexAlphabet[b&0x0f]
	}
	return string(result)
}

// ModHexDecode decodes a modhex string to byte slice
func ModHexDecode(s string) ([]byte, error) {
	s = strings.ToLower(s)
	if len(s)%2 != 0 {
		return nil, errors.New("modhex string must have even length")
	}
	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		high, okHigh := modHexDecode[s[i]]
		low, okLow := modHexDecode[s[i+1]]
		if !okHigh || !okLow {
			return nil, errors.New("invalid modhex character")
		}
		result[i/2] = (high << 4) | low
	}
	return result, nil
}

// HexEncode encodes a byte slice to hex string
func HexEncode(data []byte) string {
	result := make([]byte, len(data)*2)
	for i, b := range data {
		result[i*2] = hexAlphabet[b>>4]
		result[i*2+1] = hexAlphabet[b&0x0f]
	}
	return string(result)
}

// HexDecode decodes a hex string to byte slice
func HexDecode(s string) ([]byte, error) {
	s = strings.ToLower(s)
	if len(s)%2 != 0 {
		return nil, errors.New("hex string must have even length")
	}
	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		var high, low byte
		if s[i] >= '0' && s[i] <= '9' {
			high = s[i] - '0'
		} else if s[i] >= 'a' && s[i] <= 'f' {
			high = s[i] - 'a' + 10
		} else {
			return nil, errors.New("invalid hex character")
		}
		if s[i+1] >= '0' && s[i+1] <= '9' {
			low = s[i+1] - '0'
		} else if s[i+1] >= 'a' && s[i+1] <= 'f' {
			low = s[i+1] - 'a' + 10
		} else {
			return nil, errors.New("invalid hex character")
		}
		result[i/2] = (high << 4) | low
	}
	return result, nil
}

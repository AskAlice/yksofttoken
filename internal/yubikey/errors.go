package yubikey

import "errors"

var (
	// ErrInvalidLength indicates the data has an invalid length
	ErrInvalidLength = errors.New("invalid data length")
	// ErrInvalidModHex indicates an invalid modhex character
	ErrInvalidModHex = errors.New("invalid modhex character")
	// ErrInvalidHex indicates an invalid hex character
	ErrInvalidHex = errors.New("invalid hex character")
	// ErrCRCMismatch indicates a CRC mismatch
	ErrCRCMismatch = errors.New("CRC mismatch")
)

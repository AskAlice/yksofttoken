// Package token provides soft token management and persistence
package token

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/arr2036/yksofttoken/internal/yubikey"
)

const (
	// Field names for persistence
	PublicIDField  = "public_id"
	PrivateIDField = "private_id"
	AESKeyField    = "aes_key"
	CounterField   = "counter"
	SessionField   = "session"
	CreatedField   = "created"
	LastUseField   = "lastuse"
	PonRandField   = "ponrand"
)

// SoftToken represents a software Yubikey token
type SoftToken struct {
	PublicID  [yubikey.PublicIDSize]byte // 6 byte public identifier
	PrivateID [yubikey.UIDSize]byte      // 6 byte private identifier
	AESKey    [yubikey.KeySize]byte      // 16 byte AES key
	Counter   uint16                     // Usage counter
	Session   uint8                      // Session use counter
	Created   int64                      // Unix timestamp of creation
	LastUse   int64                      // Unix timestamp of last use
	PonRand   uint32                     // Power-on random value
}

// New creates a new SoftToken with random values
func New() (*SoftToken, error) {
	t := &SoftToken{}

	// Generate random public ID with dddd prefix (0x2222 in modhex)
	t.PublicID[0] = 0x22
	t.PublicID[1] = 0x22
	if _, err := rand.Read(t.PublicID[2:]); err != nil {
		return nil, fmt.Errorf("failed to generate public ID: %w", err)
	}

	// Generate random private ID
	if _, err := rand.Read(t.PrivateID[:]); err != nil {
		return nil, fmt.Errorf("failed to generate private ID: %w", err)
	}

	// Generate random AES key
	if _, err := rand.Read(t.AESKey[:]); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}

	// Initialize counters
	t.Counter = 1 // First "power up"
	t.Session = 1 // First "session"

	// Record creation time
	now := time.Now().Unix()
	t.Created = now
	t.LastUse = now

	// Generate power-on random
	var ponRandBytes [4]byte
	if _, err := rand.Read(ponRandBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to generate ponrand: %w", err)
	}
	t.PonRand = binary.LittleEndian.Uint32(ponRandBytes[:]) & 0xfffffff0

	return t, nil
}

// NewWithOptions creates a new SoftToken with specified options
func NewWithOptions(publicID, privateID []byte, aesKey []byte, counter uint16) (*SoftToken, error) {
	t, err := New()
	if err != nil {
		return nil, err
	}

	if publicID != nil {
		copy(t.PublicID[:], publicID)
	}

	if privateID != nil {
		copy(t.PrivateID[:], privateID)
	}

	if aesKey != nil {
		copy(t.AESKey[:], aesKey)
	}

	t.Counter = counter + 1 // Always increment on first use

	return t, nil
}

// Load loads a token from a file
func Load(path string) (*SoftToken, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	t := &SoftToken{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case PublicIDField:
			decoded, err := yubikey.ModHexDecode(value)
			if err != nil {
				return nil, fmt.Errorf("invalid public_id: %w", err)
			}
			copy(t.PublicID[:], decoded)

		case PrivateIDField:
			decoded, err := yubikey.HexDecode(value)
			if err != nil {
				return nil, fmt.Errorf("invalid private_id: %w", err)
			}
			copy(t.PrivateID[:], decoded)

		case AESKeyField:
			decoded, err := yubikey.HexDecode(value)
			if err != nil {
				return nil, fmt.Errorf("invalid aes_key: %w", err)
			}
			copy(t.AESKey[:], decoded)

		case CounterField:
			v, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("invalid counter: %w", err)
			}
			t.Counter = uint16(v)

		case SessionField:
			v, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid session: %w", err)
			}
			t.Session = uint8(v)

		case CreatedField:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid created: %w", err)
			}
			t.Created = v

		case LastUseField:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid lastuse: %w", err)
			}
			t.LastUse = v
			// Check for time travel
			if t.LastUse > time.Now().Unix() {
				return nil, errors.New("lastuse time travel detected")
			}

		case PonRandField:
			v, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid ponrand: %w", err)
			}
			t.PonRand = uint32(v)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return t, nil
}

// Save saves the token to a file
func (t *SoftToken) Save(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	publicIDModHex := yubikey.ModHexEncode(t.PublicID[:])
	privateIDHex := yubikey.HexEncode(t.PrivateID[:])
	aesKeyHex := yubikey.HexEncode(t.AESKey[:])

	fmt.Fprintf(file, "%s: %s\n", PublicIDField, publicIDModHex)
	fmt.Fprintf(file, "%s: %s\n", PrivateIDField, privateIDHex)
	fmt.Fprintf(file, "%s: %s\n", AESKeyField, aesKeyHex)
	fmt.Fprintf(file, "%s: %d\n", CounterField, t.Counter)
	fmt.Fprintf(file, "%s: %d\n", SessionField, t.Session)
	fmt.Fprintf(file, "%s: %d\n", CreatedField, t.Created)
	fmt.Fprintf(file, "%s: %d\n", LastUseField, t.LastUse)
	fmt.Fprintf(file, "%s: %d\n", PonRandField, t.PonRand)

	return nil
}

// GenerateOTP generates a new OTP and updates the token state
func (t *SoftToken) GenerateOTP() (string, error) {
	// Update session counter
	if t.Session == 0xff {
		// Session counter wrapped, increment main counter
		if t.Counter >= 0x7fff {
			return "", errors.New("token counter at max, token must be regenerated")
		}
		t.Counter++

		// Generate new power-on random
		var ponRandBytes [4]byte
		if _, err := rand.Read(ponRandBytes[:]); err != nil {
			return "", fmt.Errorf("failed to generate ponrand: %w", err)
		}
		t.PonRand = binary.LittleEndian.Uint32(ponRandBytes[:]) & 0xfffffff0
		t.Session = 1
	} else {
		t.Session++
	}

	now := time.Now().Unix()

	// Handle rate limiting
	if now == t.LastUse {
		if (t.PonRand & 0x0000000f) > 6 {
			// Rate limit - wait 1 second
			time.Sleep(time.Second)
			now = time.Now().Unix()
			t.PonRand &= 0xfffffff0
		} else {
			t.PonRand++
		}
	} else {
		t.LastUse = now
		t.PonRand &= 0xfffffff0
	}

	// Calculate 8hz timestamp
	hzTime := uint32((now-t.Created)*8) + t.PonRand
	hzTime %= 0xffffff // 24-bit wrap

	// Generate random for this OTP
	var rndBytes [2]byte
	if _, err := rand.Read(rndBytes[:]); err != nil {
		return "", fmt.Errorf("failed to generate random: %w", err)
	}
	random := binary.LittleEndian.Uint16(rndBytes[:])

	// Create token block
	block := &yubikey.TokenBlock{
		Counter:   t.Counter,
		Timestamp: hzTime,
		Session:   t.Session,
		Random:    random,
	}
	copy(block.UID[:], t.PrivateID[:])

	// Generate encrypted OTP
	otp, err := block.Generate(t.AESKey[:])
	if err != nil {
		return "", err
	}

	// Prepend public ID
	publicIDModHex := yubikey.ModHexEncode(t.PublicID[:])

	return publicIDModHex + otp, nil
}

// RegistrationInfo returns the registration information for the token
func (t *SoftToken) RegistrationInfo() string {
	publicIDModHex := yubikey.ModHexEncode(t.PublicID[:])
	privateIDHex := yubikey.HexEncode(t.PrivateID[:])
	aesKeyHex := yubikey.HexEncode(t.AESKey[:])

	return fmt.Sprintf("%s, %s, %s", publicIDModHex, privateIDHex, aesKeyHex)
}

// GetDefaultTokenDir returns the default token directory
func GetDefaultTokenDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".yksoft"), nil
}

// GetTokenPath returns the full path for a token file
func GetTokenPath(tokenDir, tokenName string) string {
	if tokenName == "" {
		tokenName = "default"
	}
	return filepath.Join(tokenDir, tokenName)
}

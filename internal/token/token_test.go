package token

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewToken(t *testing.T) {
	tok, err := New()
	if err != nil {
		t.Fatalf("Failed to create new token: %v", err)
	}

	// Check public ID prefix (should be 0x2222 = dddd in modhex)
	if tok.PublicID[0] != 0x22 || tok.PublicID[1] != 0x22 {
		t.Errorf("Public ID prefix incorrect: got %02x%02x, expected 2222",
			tok.PublicID[0], tok.PublicID[1])
	}

	// Check counters are initialized
	if tok.Counter != 1 {
		t.Errorf("Counter = %d, expected 1", tok.Counter)
	}
	if tok.Session != 1 {
		t.Errorf("Session = %d, expected 1", tok.Session)
	}

	// Check timestamps are set
	if tok.Created == 0 {
		t.Error("Created timestamp is 0")
	}
	if tok.LastUse == 0 {
		t.Error("LastUse timestamp is 0")
	}
}

func TestTokenSaveLoad(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "yksoft-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new token
	tok, err := New()
	if err != nil {
		t.Fatalf("Failed to create new token: %v", err)
	}

	// Save token
	path := filepath.Join(tmpDir, "test-token")
	err = tok.Save(path)
	if err != nil {
		t.Fatalf("Failed to save token: %v", err)
	}

	// Load token
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Failed to load token: %v", err)
	}

	// Compare fields
	if tok.PublicID != loaded.PublicID {
		t.Error("PublicID mismatch")
	}
	if tok.PrivateID != loaded.PrivateID {
		t.Error("PrivateID mismatch")
	}
	if tok.AESKey != loaded.AESKey {
		t.Error("AESKey mismatch")
	}
	if tok.Counter != loaded.Counter {
		t.Errorf("Counter mismatch: %d != %d", tok.Counter, loaded.Counter)
	}
	if tok.Session != loaded.Session {
		t.Errorf("Session mismatch: %d != %d", tok.Session, loaded.Session)
	}
	if tok.Created != loaded.Created {
		t.Errorf("Created mismatch: %d != %d", tok.Created, loaded.Created)
	}
}

func TestGenerateOTP(t *testing.T) {
	tok, err := New()
	if err != nil {
		t.Fatalf("Failed to create new token: %v", err)
	}

	initialSession := tok.Session

	// Generate an OTP
	otp, err := tok.GenerateOTP()
	if err != nil {
		t.Fatalf("Failed to generate OTP: %v", err)
	}

	// OTP should be 44 characters (12 modhex public ID + 32 modhex encrypted block)
	if len(otp) != 44 {
		t.Errorf("OTP length = %d, expected 44", len(otp))
	}

	// Session should have incremented
	if tok.Session != initialSession+1 {
		t.Errorf("Session = %d, expected %d", tok.Session, initialSession+1)
	}

	// Generate another OTP - should be different
	otp2, err := tok.GenerateOTP()
	if err != nil {
		t.Fatalf("Failed to generate second OTP: %v", err)
	}

	if otp == otp2 {
		t.Error("Two OTPs should be different")
	}
}

func TestRegistrationInfo(t *testing.T) {
	tok, err := New()
	if err != nil {
		t.Fatalf("Failed to create new token: %v", err)
	}

	info := tok.RegistrationInfo()

	// Should contain 3 comma-separated values
	// Format: public_id_modhex, private_id_hex, aes_key_hex
	// 12 + 2 + 12 + 2 + 32 = 60 characters minimum
	if len(info) < 60 {
		t.Errorf("Registration info too short: %s", info)
	}
}

func TestNewWithOptions(t *testing.T) {
	publicID := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	privateID := []byte{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	aesKey := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}

	tok, err := NewWithOptions(publicID, privateID, aesKey, 100)
	if err != nil {
		t.Fatalf("Failed to create token with options: %v", err)
	}

	// Check public ID
	for i, b := range publicID {
		if tok.PublicID[i] != b {
			t.Errorf("PublicID[%d] = %02x, expected %02x", i, tok.PublicID[i], b)
		}
	}

	// Check private ID
	for i, b := range privateID {
		if tok.PrivateID[i] != b {
			t.Errorf("PrivateID[%d] = %02x, expected %02x", i, tok.PrivateID[i], b)
		}
	}

	// Check AES key
	for i, b := range aesKey {
		if tok.AESKey[i] != b {
			t.Errorf("AESKey[%d] = %02x, expected %02x", i, tok.AESKey[i], b)
		}
	}

	// Counter should be initial + 1
	if tok.Counter != 101 {
		t.Errorf("Counter = %d, expected 101", tok.Counter)
	}
}

func TestGetDefaultTokenDir(t *testing.T) {
	dir, err := GetDefaultTokenDir()
	if err != nil {
		t.Fatalf("GetDefaultTokenDir failed: %v", err)
	}

	if dir == "" {
		t.Error("GetDefaultTokenDir returned empty string")
	}

	// Should end with .yksoft
	if filepath.Base(dir) != ".yksoft" {
		t.Errorf("Default token dir = %s, expected to end with .yksoft", dir)
	}
}

func TestGetTokenPath(t *testing.T) {
	path := GetTokenPath("/home/user/.yksoft", "mytoken")
	expected := "/home/user/.yksoft/mytoken"
	if path != expected {
		t.Errorf("GetTokenPath = %s, expected %s", path, expected)
	}

	// Empty name should default to "default"
	path = GetTokenPath("/home/user/.yksoft", "")
	expected = "/home/user/.yksoft/default"
	if path != expected {
		t.Errorf("GetTokenPath with empty name = %s, expected %s", path, expected)
	}
}

package utility

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/mail"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Validate user password according to the following rules:
//
// - Minimum length: 8 characters
//
// - Maximum length: 32 characters
//
// - Format:
//
//   - At least 1 lowercase letter
//
//   - At least 1 uppercase letter
//
//   - At least 1 digit
func ValidatePassword(password string) bool {
	// Minimum length: 8 characters
	if len(password) < 8 {
		return false
	}

	// Maximum length: 32 characters
	if len(password) > 32 {
		return false
	}

	// At least 1 lowercase, 1 uppercase, and 1 digit
	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	for _, ch := range password {
		if ch >= 'A' && ch <= 'Z' {
			hasUppercase = true
		} else if ch >= 'a' && ch <= 'z' {
			hasLowercase = true
		} else if ch >= '0' && ch <= '9' {
			hasDigit = true
		}
	}
	if !hasUppercase || !hasLowercase || !hasDigit {
		return false
	}

	return true
}

// Generate a PHC string for the given password using Argon2id
func ComputePHC(password string) (string, error) {
	// Define parameters
	// These are the recommended parameters for Argon2id
	const time = uint32(1)
	const memory = uint32(64 * 1024)
	const threads = uint8(4)
	const keyLen = uint32(32)
	const saltLen = uint32(16)

	// Generate a secure random salt
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Generate the hash
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Encode the PHC string
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	phc := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, time, threads, b64Salt, b64Hash)

	return phc, nil
}

// Verify a password against a given PHC string
func VerifyPassword(password, phc string) (bool, error) {
	// Parse PHC string
	parts := strings.Split(phc, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid PHC format")
	}

	// Extract parameters
	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	// Decode salt and hash from base64
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Compute hash with the same parameters and salt
	keyLen := uint32(len(expectedHash))
	computedHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Compare hashes in constant time
	return subtle.ConstantTimeCompare(computedHash, expectedHash) == 1, nil
}

// Validate username according to the following rules:
//
// - Minimum length: 3 characters
//
// - Maximum length: 20 characters
//
// - Format:
//
//   - Only alphanumeric characters
//
//   - Must start with an alphabetic character
func ValidateUsername(username string) bool {
	// Minimum length: 3 characters
	if len(username) < 3 {
		return false
	}

	// Maximum length: 20 characters
	if len(username) > 20 {
		return false
	}

	// Format:
	// + Only alphanumeric characters
	// + Must start with an alphabetic character
	firstChar := username[0]
	if !((firstChar >= 'A' && firstChar <= 'Z') || (firstChar >= 'a' && firstChar <= 'z')) {
		return false
	}
	for _, ch := range username {
		if !(ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			return false
		}
	}

	return true
}

// Try to parse user ID from string to uint
func TryParseUserId(idStr string) (uint, bool) {
	num, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(num), true
}

// Validate email using net/mail package
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Validate avatar URL (stub function, always returns true)
func ValidateAvatarUrl(avatarUrl string) bool {
	return true
}

// Validate bio (stub function, always returns true)
func ValidateBio(bio string) bool {
	return true
}

// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //

// crypto defines
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Iterations is the number of iterations used when generating a key with PBKDF2.
	Iterations = 1000000

	// KeyLength is the length of the key generated by PBKDF2 for encrypting data.
	KeyLength = 32

	// SaltLength is the length of the salt used by PBKDF2 when generating a key.
	SaltLength = 32

	// MinPasswordLength is the minimum number of characters that must be contained in a
	// user-provided password.
	MinPasswordLength = 16
)

// Password is a user-provided string that has been validated and meets all criteria for a password.
type Password string

// Salt is the slice of saltlen bytes used when generating a key from a user-provided password.
type Salt []byte

// Key is the slice of keylen bytes used to encrypt data with AES-256.
type Key []byte

// NewPassword verifies the provided string value meets the criteria for a password and then wraps
// it in the Password type to indicate the string has been validated. If the provided string
// does not meet the password criteria for Kolob, then the function will return an error explaining
// which criteria failed.
//
// Note that although the function only returns a single error value, the message inside that error
// value is dynamic depending on which criteria for the password were not met.
func NewPassword(val string) (Password, error) {

	count := 0
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, c := range val {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		default:
		}
		count++
	}

	fails := make([]string, 0)

	if count < MinPasswordLength {
		fails = append(fails, fmt.Sprintf("at least %v characters", MinPasswordLength))
	}
	if !hasUpper {
		fails = append(fails, "one uppercase letter")
	}
	if !hasLower {
		fails = append(fails, "one lowercase letter")
	}
	if !hasNumber {
		fails = append(fails, "one number")
	}
	if !hasSpecial {
		fails = append(fails, "one special character")
	}

	if len(fails) > 0 {
		return "", fmt.Errorf("password must contain %v", strings.Join(fails, ", "))
	}

	return Password(val), nil
}

// NewSalt creates a new slice containing saltlen bytes. The resulting salt is used when generating
// a key from a user-provided password.
func NewSalt() (Salt, error) {
	salt := make([]byte, SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to create salt: %v", err)
	}

	return salt, nil
}

// LoadSalt verifies that an existing byte slice only contains saltlen bytes. This ensures that a
// bytes loaded from authentication information can be used to recreate the original key used to
// encrypt data.
func LoadSalt(val []byte) (Salt, error) {
	if len(val) != SaltLength {
		return nil, errors.New("Salt is not the correct size")
	}

	return val, nil
}

// NewKey uses the PBKDF2 key derivation algorithm to create a 256-bit key that can be used by the
// AES algorithm for encrypting and decrypting data.
func NewKey(pass Password, salt Salt) Key {
	return pbkdf2.Key([]byte(pass), salt, Iterations, KeyLength, sha256.New)
}

// Encrypt uses AES-256 to encrypt the provided plaintext and produce a newly allocated byte slice
// of ciphertext. The byte slice is only valid if err is nil.
func Encrypt(key Key, plaintext []byte) (ciphertext []byte, err error) {

	// Prepare the block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	// Create a new nonce and fill it with cryptographically strong random values.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return
	}

	// Encrypt the plaintext data in the same byte sequence as the nonce. This will make it easy for
	// us to extract the nonce from the ciphertext when we decrypt it later.
	ciphertext = gcm.Seal(nonce, nonce, plaintext, nil)
	return
}

// Decrypt uses AES-256 to decrypt the provided ciphertext and produce a newly allocated byte slice
// of the plaintext contents.
func Decrypt(key Key, ciphertext []byte) (plaintext []byte, err error) {

	// Prepare the block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	// The ciphertext contains both the nonce and the encrypted data. We need to split the slice
	// so that we can pass both to Open().
	delim := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:delim], ciphertext[delim:]

	// Decrypt
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	return
}

// String returns a string containing a hexadecimal representation of the Salt receiver.
func (s Salt) String() string {
	return hex.EncodeToString(s)
}

// String returns a string containing a hexadecimal representation of the Key receiver.
func (k Key) String() string {
	return hex.EncodeToString(k)
}

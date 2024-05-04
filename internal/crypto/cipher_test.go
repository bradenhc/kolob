// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package crypto_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
)

func TestNewPasswordValid(t *testing.T) {
	_, err := crypto.NewPassword("This!s@val1dPa$$w0rd!")
	if err != nil {
		t.Errorf("Should not get any errors: instead got %v", err)
	}
}

func TestNewPasswordTooShort(t *testing.T) {
	_, err := crypto.NewPassword("sH0rt!")
	if err == nil {
		t.Error("Missing expected error")
	}
	if !strings.Contains(err.Error(), "at least") {
		t.Error("The error message does not tell the user the password is too short")
	}
}

func TestNewPasswordMissingUppercase(t *testing.T) {
	_, err := crypto.NewPassword("m1$$inguppercasevalue")
	if err == nil {
		t.Error("Missing expected error")
	}
	if !strings.Contains(err.Error(), "one uppercase letter") {
		t.Error("The error message does not tell the user the password is missing uppercase")
	}
}

func TestNewPasswordMissingLowercase(t *testing.T) {
	_, err := crypto.NewPassword("M1$$INGLOWERCASECHARACTER")
	if err == nil {
		t.Error("Missing expected error")
	}
	if !strings.Contains(err.Error(), "one lowercase letter") {
		t.Error("The error message does not tell the user the password is missing lowercase")
	}
}

func TestNewPasswordMissingNumber(t *testing.T) {
	_, err := crypto.NewPassword("Mi$$ingNumberInPassword")
	if err == nil {
		t.Error("Missing expected error")
	}
	if !strings.Contains(err.Error(), "one number") {
		t.Error("The error message does not tell the user the password is missing a number")
	}
}

func TestNewPasswordMissingSpecial(t *testing.T) {
	_, err := crypto.NewPassword("M1issingSpecialCharacter")
	if err == nil {
		t.Error("Missing expected error")
	}
	if !strings.Contains(err.Error(), "one special character") {
		t.Error("The error message does not tell the user the password is missing a special")
	}
}

func TestNewSalt(t *testing.T) {
	salt, err := crypto.NewSalt()
	if err != nil {
		t.Errorf("Should have created a slice without error but got: %v", err)
	}
	if len(salt) != crypto.SaltLength {
		t.Errorf("len(salt) == %v failed, got %v", crypto.SaltLength, len(salt))
	}
	fmt.Printf("Salt (new): %v\n", salt)
}

func TestLoadSaltValid(t *testing.T) {
	buf := make([]byte, crypto.SaltLength)
	_, err := rand.Read(buf)
	if err != nil {
		t.Errorf("Failed to create salt for test: %v", err)
	}

	salt, err := crypto.LoadSalt(buf)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Printf("Salt (loaded): %v\n", salt)
}

func TestLoadSaltInvalid(t *testing.T) {
	buf := make([]byte, crypto.SaltLength+1)
	_, err := rand.Read(buf)
	if err != nil {
		t.Errorf("Failed to create bad salt for test: %v", err)
	}

	_, err = crypto.LoadSalt(buf)
	if err == nil {
		t.Error("LoadSalt() with a buffer that's too large should fail")
	}
}

func TestNewKey(t *testing.T) {
	salt, err := crypto.NewSalt()
	if err != nil {
		t.Errorf("Test setup failed to create salt: %v", err)
	}
	pass, err := crypto.NewPassword("This1s@validPassw0rd")
	if err != nil {
		t.Errorf("Test setup failed to create password: %v", err)
	}

	key := crypto.NewDerivedKey(pass, salt)
	if len(key) != crypto.KeyLength {
		t.Errorf("len(key) == %v, expected %v", len(key), crypto.KeyLength)
	}
	fmt.Printf("Key: %v\n", key)
}

func TestRandomKey(t *testing.T) {
	key1, err := crypto.NewRandomKey()
	if err != nil {
		t.Fatalf("failed to generate key1: %v", err)
	}
	key2, err := crypto.NewRandomKey()
	if err != nil {
		t.Fatalf("failed to generate key2: %v", err)
	}

	if bytes.Equal(key1, key2) {
		t.Errorf("key1 == key2, but they should be different")
	}
}

func TestHashData(t *testing.T) {
	data := []byte("some data to hash")
	hash := crypto.HashData(data)

	if !crypto.CheckDataHash(data, hash) {
		t.Errorf("check data hash should pass")
	}
}

func TestPassHash(t *testing.T) {
	pass, _ := crypto.NewPassword("Password12345!")
	hash, err := crypto.HashPassword(pass)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !crypto.CheckPasswordHash(pass, hash) {
		t.Errorf("check password hash should pass")
	}
}

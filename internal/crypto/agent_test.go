// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package crypto_test

import (
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
)

type TestPerson struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

func TestAgentEncryptDecrypt(t *testing.T) {
	dkey, err := crypto.NewRandomKey()
	if err != nil {
		t.Fatalf("failed to create random key: %v", err)
	}

	a := crypto.NewAgent[TestPerson](dkey)

	p1 := TestPerson{
		Name:    "Alice",
		Age:     30,
		Address: "1234 Fiver Ave. Six City, Seven 890123",
	}

	data, err := a.Encrypt(p1)
	if err != nil {
		t.Fatalf("failed to encrypt test user: %v", err)
	}

	p2, err := a.Decrypt(data)
	if err != nil {
		t.Fatalf("failed to decrypt test user: %v", err)
	}

	if p1 != p2 {
		t.Errorf("p1 != p2 : decrypting encrypted person resulted in differences")
	}
}

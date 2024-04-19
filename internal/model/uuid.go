// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package model

import (
	"crypto/rand"
	"fmt"
)

// Uuid is a string that conforms to the RFC 4122 format for a Universally Unique Identifier.
type Uuid string

// NewUuid creates a new RFC 4122 v4 UUID and returns its string representation. Creation of the
// UUID could fail while trying to fetch cryptographically strong random bytes for the ID.
func NewUuid() (Uuid, error) {
	var uuid Uuid

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return uuid, fmt.Errorf("failed to create UUID: %v", err)

	}

	uuid = Uuid(fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))

	return uuid, nil
}

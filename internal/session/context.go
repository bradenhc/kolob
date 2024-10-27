// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package session

import (
	"context"
	"errors"

	"github.com/bradenhc/kolob/internal/crypto"
)

type key int

var encryptionKeyContextKey key

var (
	ErrInvalidSessionContextValue = errors.New("invalid session context value")
)

func NewContext(ctx context.Context, key crypto.Key) context.Context {
	return context.WithValue(ctx, encryptionKeyContextKey, key)
}

func FromContext(ctx context.Context) (crypto.Key, error) {
	k, ok := ctx.Value(encryptionKeyContextKey).(crypto.Key)
	if !ok {
		return k, ErrInvalidSessionContextValue
	}
	return k, nil
}

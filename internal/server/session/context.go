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

var pkeyKey key

var (
	ErrInvalidSessionContextValue = errors.New("invalid session context value")
)

func NewContext(ctx context.Context, pkey crypto.Key) context.Context {
	return context.WithValue(ctx, pkeyKey, pkey)
}

func FromContext(ctx context.Context) (crypto.Key, error) {
	k, ok := ctx.Value(pkeyKey).(crypto.Key)
	if !ok {
		return k, ErrInvalidSessionContextValue
	}
	return k, nil
}

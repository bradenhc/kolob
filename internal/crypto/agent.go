// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package crypto

import (
	"encoding/json"
	"fmt"
)

type Agent[V any] struct {
	key Key
}

func NewAgent[V any](key Key) Agent[V] {
	return Agent[V]{key}
}

func (a Agent[V]) Encrypt(v V) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize %T value: %v", v, err)
	}

	edata, err := Encrypt(a.key, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt %T value data: %v", v, err)
	}

	return edata, nil
}

func (a Agent[V]) Decrypt(edata []byte) (V, error) {
	data, err := Decrypt(a.key, edata)
	if err != nil {
		var v V
		return v, fmt.Errorf("failed to decrypt %T data: %v", v, err)
	}

	var v V
	err = json.Unmarshal(data, v)
	if err != nil {
		return v, fmt.Errorf("failed to deserialize %T value: %v", v, err)
	}

	return v, nil
}

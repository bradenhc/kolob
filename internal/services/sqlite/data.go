// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

type EncryptedDataAccessor[V any] struct {
	db    *sql.DB
	table string
	agent crypto.Agent[V]
}

func NewEncryptedDataAccessor[V any](
	ctx context.Context, db *sql.DB, table string, pkey crypto.Key,
) (EncryptedDataAccessor[V], error) {
	var eda EncryptedDataAccessor[V]

	gs := NewGroupService(db)
	dkey, err := gs.GetGroupDataKey(ctx, model.GetGroupDataKeyParams{PassKey: pkey})
	if err != nil {
		return eda, fmt.Errorf("failed to get group data key for decryption: %v", err)
	}

	agent := crypto.NewAgent[V](dkey)

	eda = EncryptedDataAccessor[V]{
		db:    db,
		table: table,
		agent: agent,
	}

	return eda, nil
}

func (e *EncryptedDataAccessor[V]) Get(ctx context.Context, id model.Uuid) (V, error) {
	var data []byte
	query := fmt.Sprintf("SELECT data FROM [%s] WHERE id = ?", e.table)
	err := e.db.QueryRowContext(ctx, query, id).Scan(&data)
	if err != nil {
		var v V
		return v, fmt.Errorf("failed to get %s data from database: %v", e.table, err)
	}

	return e.Decrypt(data)
}

func (e *EncryptedDataAccessor[V]) GetByIdHash(ctx context.Context, h crypto.DataHash) (V, error) {
	var d []byte
	query := fmt.Sprintf("SELECT data FROM [%s] WHERE idhash = ?", e.table)
	err := e.db.QueryRowContext(ctx, query, h[:]).Scan(&d)
	if err != nil {
		var v V
		return v, fmt.Errorf("failed to get %s data from database by ID hash: %v", e.table, err)
	}

	return e.Decrypt(d)
}

func (e *EncryptedDataAccessor[V]) GetList(ctx context.Context) ([]V, error) {
	vs := make([]V, 0)

	query := fmt.Sprintf("SELECT data FROM [%s]", e.table)
	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s list from database: %v", e.table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		if err = rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan %s row: %v", e.table, err)
		}

		v, err := e.Decrypt(data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt %s data: %v", e.table, err)
		}

		vs = append(vs, v)
	}

	return vs, nil
}

func (e *EncryptedDataAccessor[V]) Set(ctx context.Context, id model.Uuid, v V) error {
	data, err := e.Encrypt(v)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated %s data: %v", e.table, err)
	}

	updated := time.Now().Format(time.RFC3339)

	query := fmt.Sprintf("UPDATE [%s] SET data = ?, updated = ?", e.table)
	_, err = e.db.ExecContext(ctx, query, data, updated)
	if err != nil {
		return fmt.Errorf("failed to store updated %s data in database: %v", e.table, err)
	}

	return nil
}

func (e *EncryptedDataAccessor[V]) SetWithIdHash(
	ctx context.Context, id model.Uuid, h crypto.DataHash, v V,
) error {
	data, err := e.Encrypt(v)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated %s data: %v", e.table, err)
	}

	updated := time.Now().Format(time.RFC3339)

	query := fmt.Sprintf("UPDATE [%s] SET data = ?, updated = ?, idhash = ?", e.table)
	_, err = e.db.ExecContext(ctx, query, data, updated, h[:])
	if err != nil {
		return fmt.Errorf("failed to store updated %s data in database: %v", e.table, err)
	}

	return nil
}

func (e *EncryptedDataAccessor[V]) Encrypt(v V) ([]byte, error) {
	return e.agent.Encrypt(v)
}

func (e *EncryptedDataAccessor[V]) Decrypt(d []byte) (V, error) {
	return e.agent.Decrypt(d)
}

// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type EncryptedDataAccessor struct {
	table string
	dkey  crypto.Key
}

func NewEncryptedDataAccessor(table string, dkey crypto.Key) EncryptedDataAccessor {
	return EncryptedDataAccessor{
		table,
		dkey,
	}
}

func (e *EncryptedDataAccessor) Get(
	ctx context.Context, db QueryExecutor, id model.Uuid,
) ([]byte, error) {
	var data []byte
	query := fmt.Sprintf("SELECT data FROM [%s] WHERE id = ?", e.table)
	err := db.QueryRowContext(ctx, query, id).Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s data from database: %v", e.table, err)
	}

	return e.Decrypt(data)
}

func (e *EncryptedDataAccessor) GetByIdHash(
	ctx context.Context, db QueryExecutor, h crypto.DataHash,
) ([]byte, error) {
	var d []byte
	query := fmt.Sprintf("SELECT data FROM [%s] WHERE idhash = ?", e.table)
	err := db.QueryRowContext(ctx, query, h[:]).Scan(&d)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s data from database by ID hash: %v", e.table, err)
	}

	return e.Decrypt(d)
}

func (e *EncryptedDataAccessor) GetList(ctx context.Context, db QueryExecutor) ([][]byte, error) {
	query := fmt.Sprintf("SELECT data FROM [%s]", e.table)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s list from database: %v", e.table, err)
	}
	defer rows.Close()

	return e.decryptList(rows)
}

func (e *EncryptedDataAccessor) GetListFilt(
	ctx context.Context, db QueryExecutor, pred map[string]any,
) ([][]byte, error) {
	where := make([]string, len(pred))
	parms := make([]any, len(pred))
	i := 0
	for k, v := range pred {
		where[i] = k
		parms[i] = v
		i++
	}
	query := fmt.Sprintf("SELECT data FROM [%s] WHERE %s", e.table, strings.Join(where, " AND "))
	rows, err := db.QueryContext(ctx, query, parms...)
	if err != nil {
		return nil, fmt.Errorf("failed to get filtered %s list from database: %v", e.table, err)
	}
	defer rows.Close()

	return e.decryptList(rows)
}

func (e *EncryptedDataAccessor) Set(
	ctx context.Context, db QueryExecutor, id model.Uuid, v []byte,
) error {
	data, err := e.Encrypt(v)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated %s data: %v", e.table, err)
	}

	updated := time.Now().Format(time.RFC3339)

	query := fmt.Sprintf("UPDATE [%s] SET data = ?, updated = ?", e.table)
	_, err = db.ExecContext(ctx, query, data, updated)
	if err != nil {
		return fmt.Errorf("failed to store updated %s data in database: %v", e.table, err)
	}

	return nil
}

func (e *EncryptedDataAccessor) SetWithIdHash(
	ctx context.Context, db QueryExecutor, id model.Uuid, h crypto.DataHash, v []byte,
) error {
	data, err := e.Encrypt(v)
	if err != nil {
		return fmt.Errorf("failed to encrypt updated %s data: %v", e.table, err)
	}

	updated := time.Now().Format(time.RFC3339)

	query := fmt.Sprintf("UPDATE [%s] SET data = ?, updated = ?, idhash = ?", e.table)
	_, err = db.ExecContext(ctx, query, data, updated, h[:])
	if err != nil {
		return fmt.Errorf("failed to store updated %s data in database: %v", e.table, err)
	}

	return nil
}

func (e *EncryptedDataAccessor) Encrypt(d []byte) ([]byte, error) {
	return crypto.Decrypt(e.dkey, d)
}

func (e *EncryptedDataAccessor) Decrypt(d []byte) ([]byte, error) {
	return crypto.Decrypt(e.dkey, d)
}

func (e *EncryptedDataAccessor) decryptList(rows *sql.Rows) ([][]byte, error) {
	vs := make([][]byte, 0)
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
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

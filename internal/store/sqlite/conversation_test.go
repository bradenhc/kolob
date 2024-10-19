// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite_test

import (
	"context"
	"database/sql"
	"path"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
	"github.com/bradenhc/kolob/internal/store/sqlite"
)

func TestConversationSqliteStore(t *testing.T) {
	// Setup the test
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "kolob.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Create a random key we will use for encryption/decryption
	key, err := crypto.NewRandomKey()
	if err != nil {
		t.Fatalf("failed to create encryption key: %v", err)
	}

	// Create a user id to use for a moderator
	moderator, err := model.NewUuid()
	if err != nil {
		t.Fatalf("failed to create member UUID: %v", err)
	}

	// Run the tests
	store := doTestConversationStoreSqliteCreate(t, db)
	entity := doTestConversationStoreSqliteInsert(t, store, key, moderator)
	doTestConversationStoreSqliteGet(t, store, key, entity)
	entity = doTestConversationStoreSqliteUpdate(t, store, key, entity)
	doTestConversationStoreSqliteList(t, store, key, moderator)
	doTestConversationStoreSqliteRemove(t, store, entity)
}

func doTestConversationStoreSqliteCreate(t *testing.T, db *sql.DB) store.ConversationStore {
	s, err := sqlite.NewConversationStore(db)
	if err != nil {
		t.Fatalf("failed to create sqlite conversation store: %v", err)
	}

	return s
}

func doTestConversationStoreSqliteInsert(
	t *testing.T, s store.ConversationStore, k crypto.Key, mod model.Uuid,
) store.ConversationEntity {
	c, err := model.NewConversation("TestName", "TestDesc", []model.Uuid{mod})
	if err != nil {
		t.Fatalf("failed to create new conversation: %v", err)
	}

	entity, err := store.NewConversationEntity(c, k)
	if err != nil {
		t.Fatalf("failed to create conversation entity: %v", err)
	}

	err = s.AddConversationEntity(context.Background(), entity)
	if err != nil {
		t.Fatalf("failed to add conversation to store: %v", err)
	}

	return entity
}

func doTestConversationStoreSqliteGet(
	t *testing.T, s store.ConversationStore, k crypto.Key, expected store.ConversationEntity,
) {
	actual, err := s.GetConversationEntity(context.Background(), expected.Id)
	if err != nil {
		t.Fatalf("failed to get conversation entity: %v", err)
	}

	if actual.CreatedAt != expected.CreatedAt {
		t.Errorf("actual and expected creation dates are different")
	}

	if actual.UpdatedAt != expected.UpdatedAt {
		t.Errorf("actual and expected updated dates are different")
	}

	a, err := actual.Decrypt(k)
	if err != nil {
		t.Fatalf("failed to decrypt actual conversation: %v", err)
	}

	e, err := expected.Decrypt(k)
	if err != nil {
		t.Fatalf("failed to decrypt expected conversation: %v", err)
	}

	if !model.ConversationEqual(a, e) {
		t.Errorf("decrypted actual/expected conversations are different")
	}
}

func doTestConversationStoreSqliteUpdate(
	t *testing.T, s store.ConversationStore, k crypto.Key, e store.ConversationEntity,
) store.ConversationEntity {
	expected, err := e.Update(k, []byte("UpdatedName"), []byte("UpdatedDescription"), nil)
	if err != nil {
		t.Fatalf("failed to update conversation entity: %v", err)
	}

	err = s.UpdateConversationEntity(context.Background(), e)
	if err != nil {
		t.Fatalf("failed to store updated conversation: %v", err)
	}

	e, err = s.GetConversationEntity(context.Background(), e.Id)
	if err != nil {
		t.Fatalf("failed to get updated conversation: %v", err)
	}

	actual, err := e.Decrypt(k)
	if err != nil {
		t.Fatalf("failed to decrypt updated conversation: %v", err)
	}

	if !model.ConversationEqual(actual, expected) {
		t.Errorf("expected stored updated conversation different than actual")
	}

	return e
}

func doTestConversationStoreSqliteList(
	t *testing.T, s store.ConversationStore, k crypto.Key, m model.Uuid,
) {
	for range 3 {
		doTestConversationStoreSqliteInsert(t, s, k, m)
	}

	entities, err := s.ListConversationEntities(context.Background())
	if err != nil {
		t.Fatalf("failed to list conversation entities: %v", err)
	}

	if len(entities) != 4 {
		t.Errorf("expected 4 conversation entities: got %d", len(entities))
	}
}

func doTestConversationStoreSqliteRemove(
	t *testing.T, s store.ConversationStore, e store.ConversationEntity,
) {
	err := s.RemoveConversationEntity(context.Background(), e.Id)
	if err != nil {
		t.Fatalf("failed to remove conversation entity from store: %v", err)
	}

	entities, err := s.ListConversationEntities(context.Background())
	if err != nil {
		t.Fatalf("failed to list conversation entities: %v", err)
	}

	if len(entities) != 3 {
		t.Errorf("expected 3 conversation entities after delete: got %d", len(entities))
	}
}

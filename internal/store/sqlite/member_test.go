// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"slices"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
	"github.com/bradenhc/kolob/internal/store/sqlite"
)

func TestMemberSqliteStore(t *testing.T) {
	// Setup the test
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "kolob-TestMemberSqliteStore.db")

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

	// Run the tests
	store := doTestMemberStoreSqliteCreate(t, db)
	id := doTestMemberStoreSqliteInsert(t, store, key)
	entity := doTestMemberStoreSqliteGet(t, store, key, id)
	doTestMemberStoreSqliteUpdate(t, store, key, entity)
	doTestMemberStoreSqliteList(t, store, key)
	doTestMemberStoreSqliteRemove(t, store, id)
}

func doTestMemberStoreSqliteCreate(t *testing.T, db *sql.DB) sqlite.MemberStore {
	store, err := sqlite.NewMemberStore(db)
	if err != nil {
		t.Fatalf("failed to create member store: %v", err)
	}

	return store
}

func doTestMemberStoreSqliteInsert(t *testing.T, s sqlite.MemberStore, key crypto.Key) model.Uuid {
	member, err := model.NewMember("TestUser", "Name")
	if err != nil {
		t.Fatalf("failed to create new group: %v", err)
	}

	pass, _ := crypto.NewPassword("Password123!")

	entity, err := store.NewMemberEntity(member, pass, key)
	if err != nil {
		t.Fatalf("failed to create new member entity: %v", err)
	}

	err = s.AddMemberEntity(context.Background(), entity)
	if err != nil {
		t.Fatalf("failed to add member entity: %v", err)
	}

	return model.Uuid(member.Id())
}

func doTestMemberStoreSqliteGet(
	t *testing.T, s sqlite.MemberStore, key crypto.Key, uid model.Uuid,
) store.MemberEntity {
	entity, err := s.GetMemberEntity(context.Background(), uid)
	if err != nil {
		t.Fatalf("failed to get member data: %v", err)
	}

	member, err := entity.Decrypt(key)
	if err != nil {
		t.Fatalf("failed to decrypt member: %v", err)
	}

	if !slices.Equal(member.Uname(), []byte("TestUser")) {
		t.Errorf("member usernames are not equal")
	}

	if !slices.Equal(member.Name(), []byte("Name")) {
		t.Errorf("member names are not equal")
	}

	if member.Created() != member.Updated() {
		t.Errorf("member created/updated times are different")
	}

	return entity
}

func doTestMemberStoreSqliteUpdate(
	t *testing.T, store sqlite.MemberStore, key crypto.Key, entity store.MemberEntity,
) {
	uname := []byte("UpdatedUname")
	name := []byte("New Name")

	_, err := entity.Update(key, uname, name)
	if err != nil {
		t.Fatalf("failed to update member entity: %v", err)
	}

	err = store.UpdateMemberEntity(context.Background(), entity)
	if err != nil {
		t.Fatalf("failed to store updated group entity: %v", err)
	}

	entity, err = store.GetMemberEntity(context.Background(), entity.Id)
	if err != nil {
		t.Fatalf("failed to get updated group entity from store: %v", err)
	}

	m, err := entity.Decrypt(key)
	if err != nil {
		t.Fatalf("failed to decrypt updated member entity: %v", err)
	}

	if !slices.Equal(m.Uname(), uname) {
		t.Errorf("updated member name does not match")
	}

	if !slices.Equal(m.Name(), name) {
		t.Errorf("updated member description does not match")
	}

	if m.Created() == m.Updated() {
		t.Errorf("updated member created/updated times are the same but should be different")
	}
}

func doTestMemberStoreSqliteList(t *testing.T, s store.MemberStore, key crypto.Key) {
	for i := range 3 {
		uname := fmt.Sprintf("TestUser%02d", i)
		member, err := model.NewMember(uname, "Name")
		if err != nil {
			t.Fatalf("failed to create new member: %v", err)
		}

		pass, _ := crypto.NewPassword("Password123!")

		entity, err := store.NewMemberEntity(member, pass, key)
		if err != nil {
			t.Fatalf("failed to create new member entity: %v", err)
		}

		err = s.AddMemberEntity(context.Background(), entity)
		if err != nil {
			t.Fatalf("failed to add member entity: %v", err)
		}
	}

	entities, err := s.ListMemberEntities(context.Background())
	if err != nil {
		t.Fatalf("failed to list members: %v", err)
	}

	if len(entities) != 4 {
		t.Errorf("not enough entities: expected 4")
	}
}

func doTestMemberStoreSqliteRemove(t *testing.T, s store.MemberStore, id model.Uuid) {
	err := s.RemoveMemberEntity(context.Background(), id)
	if err != nil {
		t.Fatalf("failed to remove member: %v", err)
	}

	entities, err := s.ListMemberEntities(context.Background())
	if err != nil {
		t.Fatalf("failed to list members after delete: %v", err)
	}

	if len(entities) != 3 {
		t.Errorf("too many entities after delete: expected 3")
	}
}

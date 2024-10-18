package sqlite_test

import (
	"context"
	"database/sql"
	"path"
	"slices"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
	"github.com/bradenhc/kolob/internal/store/sqlite"
)

func TestGroupSqliteStore(t *testing.T) {
	// Setup the test
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "kolob-TestGroupSqliteStore.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Run the tests
	store := doTestGroupStoreSqliteCreate(t, db)
	key := doTestGroupStoreSqliteInsert(t, store)
	entity, group := doTestGroupStoreSqliteGet(t, store, key)
	doTestGroupStoreSqliteUpdate(t, store, key, entity, group)
}

func doTestGroupStoreSqliteCreate(t *testing.T, db *sql.DB) sqlite.GroupStore {
	store, err := sqlite.NewGroupStore(db)
	if err != nil {
		t.Fatalf("failed to create group store: %v", err)
	}

	return store
}

func doTestGroupStoreSqliteInsert(t *testing.T, s sqlite.GroupStore) crypto.Key {
	group, err := model.NewGroup("Group123", "Name", "Description")
	if err != nil {
		t.Fatalf("failed to create new group: %v", err)
	}

	pass, _ := crypto.NewPassword("Password123!")

	key, err := crypto.NewRandomKey()
	if err != nil {
		t.Fatalf("failed to create random key: %v", err)
	}

	entity, err := store.NewGroupEntity(group, pass, key)
	if err != nil {
		t.Fatalf("failed to create new group entity: %v", err)
	}

	err = s.AddGroupEntity(context.Background(), entity)
	if err != nil {
		t.Fatalf("failed to add group data: %v", err)
	}

	return key
}

func doTestGroupStoreSqliteGet(t *testing.T, s sqlite.GroupStore, key crypto.Key) (store.GroupEntity, *model.Group) {
	entity, err := s.GetGroupEntity(context.Background())
	if err != nil {
		t.Fatalf("failed to get group data: %v", err)
	}

	data, err := crypto.Decrypt(key, entity.EncryptedData)
	if err != nil {
		t.Fatalf("failed to decrypt group: %v", err)
	}
	group := model.GetRootAsGroup(data, 0)

	if !slices.Equal(group.Gid(), []byte("Group123")) {
		t.Errorf("group ids are not equal")
	}

	if !slices.Equal(group.Name(), []byte("Name")) {
		t.Errorf("group names are not equal")
	}

	if !slices.Equal(group.Desc(), []byte("Description")) {
		t.Errorf("group descriptions are not equal")
	}

	if group.Created() != group.Updated() {
		t.Errorf("group created/updated times are different")
	}

	return entity, group
}

func doTestGroupStoreSqliteUpdate(
	t *testing.T, s sqlite.GroupStore, k crypto.Key, e store.GroupEntity, g *model.Group,
) {
	ngid := []byte("Group456")
	nname := []byte("New Name")
	ndesc := []byte("New Description")
	ngroup := model.GroupCloneWithUpdates(g, ngid, nname, ndesc)

	e.GroupHash = crypto.HashData(ngroup.Gid())
	e.UpdatedAt = ngroup.Updated()

	edata, err := crypto.Encrypt(k, ngroup.Table().Bytes)
	if err != nil {
		t.Fatalf("failed to encrypt updated group data: %v", err)
	}
	e.EncryptedData = edata

	err = s.UpdateGroupEntity(context.Background(), e)
	if err != nil {
		t.Fatalf("failed to store updated group entity: %v", err)
	}

	e, err = s.GetGroupEntity(context.Background())
	if err != nil {
		t.Fatalf("failed to get updated group entity from store: %v", err)
	}

	data, err := crypto.Decrypt(k, e.EncryptedData)
	if err != nil {
		t.Fatalf("failed to decrypt group data: %v", err)
	}

	g = model.GetRootAsGroup(data, 0)

	if !slices.Equal(g.Gid(), ngid) {
		t.Errorf("updated group id does not match")
	}
	if !slices.Equal(g.Name(), nname) {
		t.Errorf("updated group name does not match")
	}
	if !slices.Equal(g.Desc(), ndesc) {
		t.Errorf("updated group description does not match")
	}
	if g.Created() == g.Updated() {
		t.Errorf("updated group created/updated times are the same but should be different")
	}
}

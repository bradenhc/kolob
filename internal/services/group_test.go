// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package services_test

import (
	"context"
	"database/sql"
	"path"
	"slices"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/services"
	"github.com/bradenhc/kolob/internal/store/sqlite"
	flatbuffers "github.com/google/flatbuffers/go"
)

func TestGroupService(t *testing.T) {
	// Setup the test
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "kolob-TestGroupService.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Setup the store and create the service
	store := doTestGroupCreateStore(t, db)
	gs := services.NewGroupService(store)

	// Create a group and test
	ctx := context.Background()
	a := doTestGroupCreate(t, ctx, gs)

	// Authenticate to get the symmetric key
	dkey := doTestGroupAuth(t, ctx, gs)

	// Access encrypted group information
	b := doTestGroupGetInfo(t, ctx, gs, dkey, a)

	// Update group information
	c := doTestGroupUpdate(t, ctx, gs, dkey, b)

	// Change the group password and make sure we can still authenticate and access the group
	doTestGroupChangePassword(t, ctx, gs, dkey, c)
}

func doTestGroupCreateStore(t *testing.T, db *sql.DB) sqlite.GroupStore {
	// Create our group store
	//
	store, err := sqlite.NewGroupStore(db)
	if err != nil {
		t.Fatalf("failed to create group store: %v", err)
	}

	return store
}

func doTestGroupCreate(t *testing.T, ctx context.Context, gs services.GroupService) *model.Group {
	builder := flatbuffers.NewBuilder(64)
	offsetGid := builder.CreateString("TestGroup123")
	offsetName := builder.CreateString("Test Group")
	offsetDesc := builder.CreateString("A test group")
	offsetPass := builder.CreateString("Password12345678!")
	services.GroupInitRequestStart(builder)
	services.GroupInitRequestAddGroupId(builder, offsetGid)
	services.GroupInitRequestAddName(builder, offsetName)
	services.GroupInitRequestAddDescription(builder, offsetDesc)
	services.GroupInitRequestAddPassword(builder, offsetPass)

	offsetReq := services.GroupInitRequestEnd(builder)
	builder.Finish(offsetReq)

	reqCreate := services.GetRootAsGroupInitRequest(builder.FinishedBytes(), 0)

	group, err := gs.Create(ctx, reqCreate)
	if err != nil {
		t.Fatalf("failed to create group: %v", err)
	}

	if !slices.Equal(group.Gid(), reqCreate.GroupId()) {
		t.Errorf("%s != %s", string(group.Gid()), string(reqCreate.GroupId()))
	}
	if !slices.Equal(group.Name(), reqCreate.Name()) {
		t.Errorf("%s != %s", string(group.Name()), string(reqCreate.Name()))
	}
	if !slices.Equal(group.Desc(), reqCreate.Description()) {
		t.Errorf("%s != %s", string(group.Desc()), string(reqCreate.Description()))
	}

	return group
}

func doTestGroupAuth(t *testing.T, ctx context.Context, gs services.GroupService) crypto.Key {
	builder := flatbuffers.NewBuilder(64)
	offsetGid := builder.CreateString("TestGroup123")
	offsetPass := builder.CreateString("Password12345678!")
	services.GroupAuthenticateRequestStart(builder)
	services.GroupAuthenticateRequestAddGroupId(builder, offsetGid)
	services.GroupAuthenticateRequestAddPassword(builder, offsetPass)
	offsetReq := services.GroupAuthenticateRequestEnd(builder)
	builder.Finish(offsetReq)

	reqAuth := services.GetRootAsGroupAuthenticateRequest(builder.FinishedBytes(), 0)
	dkey, err := gs.Authenticate(ctx, reqAuth)
	if err != nil {
		t.Fatalf("failed to auth group: %v", err)
	}

	return dkey
}

func doTestGroupGetInfo(
	t *testing.T, ctx context.Context, gs services.GroupService, dkey crypto.Key, a *model.Group,
) *model.Group {
	b, err := gs.Get(ctx, dkey)
	if err != nil {
		t.Fatalf("failed to get group info: %v", err)
	}
	if !model.GroupEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}

	return b
}

func doTestGroupUpdate(
	t *testing.T, ctx context.Context, gs services.GroupService, dkey crypto.Key, b *model.Group,
) *model.Group {
	builder := flatbuffers.NewBuilder(64)
	ugid := "TestGroup456"
	offsetGid := builder.CreateString(ugid)
	uname := "Test Group Updated"
	offsetName := builder.CreateString(uname)
	udesc := "Updated test group description"
	offsetDesc := builder.CreateString(udesc)
	services.GroupUpdateRequestStart(builder)
	services.GroupUpdateRequestAddGroupId(builder, offsetGid)
	services.GroupUpdateRequestAddName(builder, offsetName)
	services.GroupUpdateRequestAddDescription(builder, offsetDesc)
	r := services.GroupUpdateRequestEnd(builder)
	builder.Finish(r)

	reqUpdate := services.GetRootAsGroupUpdateRequest(builder.FinishedBytes(), 0)
	c, err := gs.Update(ctx, reqUpdate, dkey)
	if err != nil {
		t.Fatalf("failed to update group information: %v", err)
	}

	if model.GroupEqual(b, c) {
		t.Fatalf("values should be different: %+v == %+v", b, c)
	}

	if string(c.Gid()) != ugid {
		t.Errorf("%s != %s", c.Gid(), ugid)
	}
	if string(c.Name()) != uname {
		t.Errorf("%s != %s", c.Name(), uname)
	}
	if string(c.Desc()) != udesc {
		t.Errorf("%s != %s", b.Desc(), udesc)
	}

	return c
}

func doTestGroupChangePassword(
	t *testing.T, ctx context.Context, gs services.GroupService, dkey crypto.Key, c *model.Group,
) {
	builder := flatbuffers.NewBuilder(64)
	offsetOldPass := builder.CreateString("Password12345678!")
	newp, _ := crypto.NewPassword("UpdatedPassword123456!")
	offsetNewPass := builder.CreateString(string(newp))
	services.GroupChangePasswordRequestStart(builder)
	services.GroupChangePasswordRequestAddOldPassword(builder, offsetOldPass)
	services.GroupChangePasswordRequestAddNewPassword(builder, offsetNewPass)
	r := services.GroupChangePasswordRequestEnd(builder)
	builder.Finish(r)

	reqChangePass := services.GetRootAsGroupChangePasswordRequest(builder.FinishedBytes(), 0)
	err := gs.ChangePassword(ctx, reqChangePass, dkey)
	if err != nil {
		t.Fatalf("failed to update password: %v", err)
	}

	builder = flatbuffers.NewBuilder(64)
	offsetGid := builder.CreateString("TestGroup456")
	offsetPass := builder.CreateString(string(newp))
	services.GroupAuthenticateRequestStart(builder)
	services.GroupAuthenticateRequestAddGroupId(builder, offsetGid)
	services.GroupAuthenticateRequestAddPassword(builder, offsetPass)
	r = services.GroupAuthenticateRequestEnd(builder)
	builder.Finish(r)

	reqAuth := services.GetRootAsGroupAuthenticateRequest(builder.FinishedBytes(), 0)
	dkey, err = gs.Authenticate(ctx, reqAuth)
	if err != nil {
		t.Fatalf("failed to auth group: %v", err)
	}

	d, err := gs.Get(ctx, dkey)
	if err != nil {
		t.Fatalf("failed to get group info after password update: %v", err)
	}

	if c.Updated() == d.Updated() {
		t.Errorf("updated times should be different after password upadte")
	}
}

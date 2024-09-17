// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package services_test

import (
	"context"
	"path"
	"slices"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/services"
	"github.com/bradenhc/kolob/internal/store/sqlite"
	flatbuffers "github.com/google/flatbuffers/go"
)

func TestGroupTable(t *testing.T) {
	// Setup the test
	//
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "group-table-test.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Create our group store
	//
	store, err := sqlite.NewGroupStore(db)
	if err != nil {
		t.Fatalf("failed to create group store: %v", err)
	}

	// Create the group service we will test
	//
	gs := services.NewGroupService(store)

	// Create and store a group
	//
	builder := flatbuffers.NewBuilder(64)
	gid := "TestGroup123"
	offsetGid := builder.CreateString("TestGroup123")
	offsetName := builder.CreateString("Test Group")
	offsetDesc := builder.CreateString("A test group")
	pw, _ := crypto.NewPassword("password")
	offsetPass := builder.CreateString(string(pw))
	services.GroupInitRequestStart(builder)
	services.GroupInitRequestAddGroupId(builder, offsetGid)
	services.GroupInitRequestAddName(builder, offsetName)
	services.GroupInitRequestAddDescription(builder, offsetDesc)
	services.GroupInitRequestAddPassword(builder, offsetPass)

	r := services.GroupInitRequestEnd(builder)
	builder.Finish(r)

	reqCreate := services.GetRootAsGroupInitRequest(builder.FinishedBytes(), 0)

	ctx := context.Background()
	a, err := gs.Create(ctx, reqCreate)
	if err != nil {
		t.Fatalf("failed to create group: %v", err)
	}

	if !slices.Equal(a.Gid(), reqCreate.GroupId()) {
		t.Errorf("%s != %s", string(a.Gid()), string(reqCreate.GroupId()))
	}
	if !slices.Equal(a.Name(), reqCreate.Name()) {
		t.Errorf("%s != %s", string(a.Name()), string(reqCreate.Name()))
	}
	if !slices.Equal(a.Desc(), reqCreate.Description()) {
		t.Errorf("%s != %s", string(a.Desc()), string(reqCreate.Description()))
	}

	// Authenticate with the group
	//
	builder = flatbuffers.NewBuilder(64)
	offsetGid = builder.CreateString(gid)
	offsetPass = builder.CreateString(string(pw))
	services.GroupAuthenticateRequestStart(builder)
	services.GroupAuthenticateRequestAddGroupId(builder, offsetGid)
	services.GroupAuthenticateRequestAddPassword(builder, offsetPass)
	r = services.GroupAuthenticateRequestEnd(builder)
	builder.Finish(r)

	reqAuth := services.GetRootAsGroupAuthenticateRequest(builder.FinishedBytes(), 0)
	dkey, err := gs.Authenticate(ctx, reqAuth)
	if err != nil {
		t.Fatalf("failed to auth group: %v", err)
	}

	// Access encrypted group information
	//
	b, err := gs.GetInfo(ctx, dkey)
	if err != nil {
		t.Fatalf("failed to get group info: %v", err)
	}
	if !model.GroupEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}

	// Update group information
	//
	builder = flatbuffers.NewBuilder(64)
	ugid := "TestGroup456"
	offsetGid = builder.CreateString(ugid)
	uname := "Test Group Updated"
	offsetName = builder.CreateString(uname)
	udesc := "Updated test group description"
	offsetDesc = builder.CreateString(udesc)
	services.GroupUpdateRequestStart(builder)
	services.GroupUpdateRequestAddGroupId(builder, offsetGid)
	services.GroupUpdateRequestAddName(builder, offsetName)
	services.GroupUpdateRequestAddDescription(builder, offsetDesc)
	r = services.GroupUpdateRequestEnd(builder)
	builder.Finish(r)

	reqUpdate := services.GetRootAsGroupUpdateRequest(builder.FinishedBytes(), 0)
	c, err := gs.Update(ctx, reqUpdate, dkey)
	if err != nil {
		t.Fatalf("failed to update group information: %v", err)
	}

	if !model.GroupEqual(b, c) {
		t.Fatalf("values should be different: %+v == %+v", b, c)
	}

	if string(c.Gid()) != ugid {
		t.Errorf("%v != %s", c.Gid(), ugid)
	}
	if string(c.Name()) != uname {
		t.Errorf("%v != %s", c.Name(), uname)
	}
	if string(c.Desc()) != udesc {
		t.Errorf("%v != %s", a.Desc(), udesc)
	}

	// Change the group password and make sure we can still authenticate and access the group
	//
	builder = flatbuffers.NewBuilder(64)
	offsetOldPass := builder.CreateString(string(pw))
	newp, _ := crypto.NewPassword("newpassword")
	offsetNewPass := builder.CreateString(string(newp))
	services.GroupChangePasswordRequestStart(builder)
	services.GroupChangePasswordRequestAddOldPassword(builder, offsetOldPass)
	services.GroupChangePasswordRequestAddNewPassword(builder, offsetNewPass)
	r = services.GroupChangePasswordRequestEnd(builder)
	builder.Finish(r)

	reqChangePass := services.GetRootAsGroupChangePasswordRequest(builder.FinishedBytes(), 0)
	err = gs.ChangePassword(ctx, reqChangePass, dkey)
	if err != nil {
		t.Fatalf("failed to update password: %v", err)
	}

	builder = flatbuffers.NewBuilder(64)
	offsetGid = builder.CreateString(ugid)
	offsetPass = builder.CreateString(string(newp))
	services.GroupAuthenticateRequestStart(builder)
	services.GroupAuthenticateRequestAddGroupId(builder, offsetGid)
	services.GroupAuthenticateRequestAddPassword(builder, offsetPass)
	r = services.GroupAuthenticateRequestEnd(builder)
	builder.Finish(r)

	reqAuth = services.GetRootAsGroupAuthenticateRequest(builder.FinishedBytes(), 0)
	dkey, err = gs.Authenticate(ctx, reqAuth)
	if err != nil {
		t.Fatalf("failed to auth group: %v", err)
	}

	d, err := gs.GetInfo(ctx, dkey)
	if err != nil {
		t.Fatalf("failed to get group info after password update: %v", err)
	}

	if !model.GroupEqual(c, d) {
		t.Fatalf("values should be the same after password update: %+v == %+v", c, d)
	}
}

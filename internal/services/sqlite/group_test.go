// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite_test

import (
	"context"
	"path"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/services/sqlite"
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

	// Create the group service we will test
	//
	gs := sqlite.NewGroupService(db)

	// Create and store a group
	//
	gid := "TestGroup123"
	name := "Test Group"
	desc := "A test group"
	pass, _ := crypto.NewPassword("password")

	ctx := context.Background()
	a, err := gs.InitGroup(ctx, model.InitGroupParams{
		GroupId:     gid,
		Name:        name,
		Description: desc,
		Password:    pass,
	})
	if err != nil {
		t.Fatalf("failed to create group: %v", err)
	}

	if a.GroupId != gid {
		t.Errorf("%s != %s", a.GroupId, gid)
	}
	if a.Name != name {
		t.Errorf("%s != %s", a.Name, name)
	}
	if a.Description != desc {
		t.Errorf("%s != %s", a.Description, desc)
	}

	// Authenticate with the group
	//
	err = gs.AuthenticateGroup(ctx, model.AuthenticateGroupParams{GroupId: gid, Password: pass})
	if err != nil {
		t.Fatalf("failed to auth group: %v", err)
	}

	// Access encrypted group information
	//
	pkey, err := gs.GetGroupPassKey(ctx, model.GetGroupPassKeyParams{Password: pass})
	if err != nil {
		t.Fatalf("failed to get group pass key: %v", err)
	}

	b, err := gs.GetGroupInfo(ctx, model.GetGroupInfoParams{Id: a.Id, PassKey: pkey})
	if err != nil {
		t.Fatalf("failed to get group info: %v", err)
	}

	if !a.Equal(&b) {
		t.Errorf("%+v != %+v", a, b)
	}

	// Update group information
	//
	ugid := "TestGroup456"
	uname := "Test Group Updated"
	udesc := "Updated test group description"
	err = gs.UpdateGroup(ctx, model.UpdateGroupParams{
		Id:          a.Id,
		GroupId:     &ugid,
		Name:        &uname,
		Description: &udesc,
		PassKey:     pkey,
	})
	if err != nil {
		t.Fatalf("failed to update group information: %v", err)
	}

	c, err := gs.GetGroupInfo(ctx, model.GetGroupInfoParams{Id: a.Id, PassKey: pkey})
	if err != nil {
		t.Fatalf("failed to get group info: %v", err)
	}

	if c.Equal(&b) {
		t.Fatalf("values should be different: %+v == %+v", b, c)
	}

	if c.GroupId != ugid {
		t.Errorf("%s != %s", c.GroupId, ugid)
	}
	if c.Name != uname {
		t.Errorf("%s != %s", c.Name, uname)
	}
	if c.Description != udesc {
		t.Errorf("%s != %s", a.Description, udesc)
	}

	// Change the group password and make sure we can still authenticate and access the group
	//
	newp, _ := crypto.NewPassword("newpassword")
	err = gs.ChangeGroupPassword(ctx, model.ChangeGroupPasswordParams{
		Id:      a.Id,
		OldPass: pass,
		NewPass: newp,
	})
	if err != nil {
		t.Fatalf("failed to update password: %v", err)
	}

	err = gs.AuthenticateGroup(ctx, model.AuthenticateGroupParams{GroupId: ugid, Password: newp})
	if err != nil {
		t.Fatalf("failed to auth group after password update: %v", err)
	}

	pkey, err = gs.GetGroupPassKey(ctx, model.GetGroupPassKeyParams{Password: pass})
	if err != nil {
		t.Fatalf("failed to get group pass key after password update: %v", err)
	}

	d, err := gs.GetGroupInfo(ctx, model.GetGroupInfoParams{Id: a.Id, PassKey: pkey})
	if err != nil {
		t.Fatalf("failed to get group info after password update: %v", err)
	}

	if !d.Equal(&c) {
		t.Fatalf("values should be the same after password update: %+v == %+v", c, d)
	}
}

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

func TestMemberService(t *testing.T) {
	// Setup the test
	//
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "member-service-test.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	ctx := context.Background()

	// Create the group this member will belong to
	//
	gpass, _ := crypto.NewPassword("Password123456!")
	gs := sqlite.NewGroupService(db)
	_, err = gs.InitGroup(ctx, model.InitGroupParams{
		GroupId:     "testgroup",
		Name:        "Test Group",
		Description: "A group for this test",
		Password:    gpass,
	})
	if err != nil {
		t.Fatalf("failed to create group: %v", err)
	}

	pkey, err := gs.GetGroupPassKey(ctx, model.GetGroupPassKeyParams{Password: gpass})
	if err != nil {
		t.Fatalf("failed to get group pass key: %v", err)
	}

	// Create the member service
	//
	ms := sqlite.NewMemberService(db)

	// Add a member
	//
	uname := "testuser"
	name := "Alice Ann"
	upass, _ := crypto.NewPassword("U$hal1~PAss!")
	a, err := ms.CreateMember(ctx, model.CreateMemberParams{
		Name:     name,
		Username: uname,
		Password: upass,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to create member: %v", err)
	}

	if a.Name != name {
		t.Errorf("member name incorrect: %s != %s", a.Name, name)
	}
	if a.Username != uname {
		t.Errorf("member username incorrect: %s != %s", a.Username, uname)
	}

	// Authenticate member
	//
	err = ms.AuthenticateMember(ctx, model.AuthenticateMemberParams{
		Username: uname,
		Password: upass,
	})
	if err != nil {
		t.Fatalf("failed to auth member: %v", err)
	}

	// Get member
	//
	b, err := ms.FindMemberByUsername(ctx, model.FindMemberByUsernameParams{
		Username: uname,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to find member by username: %v", err)
	}

	if !b.Equal(&a) {
		t.Errorf("member not the same: %+v != %+v", b, a)
	}

	// Update member
	//
	uname = "updatedtestuser"
	name = "Bob Bill"
	err = ms.UpdateMember(ctx, model.UpdateMemberParams{
		Id:       a.Id,
		PassKey:  pkey,
		Username: &uname,
		Name:     &name,
	})
	if err != nil {
		t.Fatalf("failed to update member: %v", err)
	}

	c, err := ms.FindMemberByUsername(ctx, model.FindMemberByUsernameParams{
		Username: uname,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to find member by username: %v", err)
	}

	if c.Equal(&b) {
		t.Errorf("member should be different after update: %+v == %+v", c, b)
	}

	// Add another member, then list all member
	//
	uname2 := "anotheruser"
	name2 := "Carol Chris"
	upass2, _ := crypto.NewPassword("G0b@k2th#shadows!")
	d, err := ms.CreateMember(ctx, model.CreateMemberParams{
		Name:     name2,
		Username: uname2,
		Password: upass2,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to create second member: %v", err)
	}

	l, err := ms.ListMembers(ctx, model.ListMembersParams{
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to list members: %v", err)
	}

	if len(l) != 2 {
		t.Fatalf("expected two members in list, got %d", len(l))
	}

	if !l[0].Equal(&c) {
		t.Errorf("first member is not what was expected: %+v != %+v", l[0], c)
	}
	if !l[1].Equal(&d) {
		t.Errorf("second member is not what was expected: %+v != %+v", l[1], d)
	}

	// Remove a member
	//
	err = ms.RemoveMember(ctx, model.RemoveMemberParams{
		Id: d.Id,
	})
	if err != nil {
		t.Errorf("failed to remove member: %v", err)
	}
	l, err = ms.ListMembers(ctx, model.ListMembersParams{
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to list members: %v", err)
	}
	if len(l) != 1 {
		t.Fatalf("expected only one member after delete, got %d", len(l))
	}
	if !l[0].Equal(&c) {
		t.Errorf("remaining member is not what was expected after delete: %+v != %+v", l[0], c)
	}
}

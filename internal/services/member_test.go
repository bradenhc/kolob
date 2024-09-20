// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package services_test

import (
	"context"
	"database/sql"
	"path"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/services"
	"github.com/bradenhc/kolob/internal/store/sqlite"
	flatbuffers "github.com/google/flatbuffers/go"
)

func TestMemberService(t *testing.T) {
	// Setup the test
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "member-service-test.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	ctx := context.Background()

	// Create our group store and service
	gstore := doTestGroupCreateStore(t, db)
	gs := services.NewGroupService(gstore)

	// Create and store a group to associate members with and get the key
	doTestGroupCreate(t, ctx, gs)
	key := doTestGroupAuth(t, ctx, gs)

	// Create the member store and service
	mstore := doTestMemberCreateStore(t, db)
	ms := services.NewMemberService(mstore)

	// Add a member
	a := doTestMemberAdd(t, ctx, ms, key)

	// Authenticate member
	doTestMemberAuth(t, ctx, ms, key)

	// Get member
	b := doTestMemberFindByUsername(t, ctx, ms, key, a)

	// Update member
	c := doTestMemberUpdate(t, ctx, ms, key, b)

	// Add another member, then list all member
	d := doTestMemberList(t, ctx, ms, key, c)

	// Remove a member
	doTestMemberRemove(t, ctx, ms, key, c, d)
}

func doTestMemberCreateStore(t *testing.T, db *sql.DB) sqlite.MemberStore {
	store, err := sqlite.NewMemberStore(db)
	if err != nil {
		t.Fatalf("failed to create member store: %v", err)
	}
	return store

}

func doTestMemberAdd(
	t *testing.T, ctx context.Context, ms services.MemberService, key crypto.Key,
) *model.Member {
	uname := "testuser"
	name := "Alice Ann"
	upass, err := crypto.NewPassword("U$hal1~PAss!")
	if err != nil {
		t.Fatalf("failed to create member password: %v", err)
	}

	builder := flatbuffers.NewBuilder(64)
	unameOffset := builder.CreateString(uname)
	nameOffset := builder.CreateString(name)
	upassOffset := builder.CreateString(string(upass))

	services.MemberCreateRequestStart(builder)
	services.MemberCreateRequestAddUsername(builder, unameOffset)
	services.MemberCreateRequestAddName(builder, nameOffset)
	services.MemberCreateRequestAddPassword(builder, upassOffset)

	r := services.MemberCreateRequestEnd(builder)
	builder.Finish(r)

	reqCreateMember := services.GetRootAsMemberCreateRequest(builder.FinishedBytes(), 0)
	a, err := ms.Create(ctx, reqCreateMember, key)
	if err != nil {
		t.Fatalf("failed to create member: %v", err)
	}

	if string(a.Name()) != name {
		t.Errorf("member name incorrect: %s != %s", a.Name, name)
	}
	if string(a.Uname()) != uname {
		t.Errorf("member username incorrect: %s != %s", a.Uname(), uname)
	}

	return a
}

func doTestMemberAuth(
	t *testing.T, ctx context.Context, ms services.MemberService, key crypto.Key,
) {
	uname := "testuser"
	upass := "U$hal1~PAss!"

	builder := flatbuffers.NewBuilder(64)
	unameOffset := builder.CreateString(uname)
	upassOffset := builder.CreateString(upass)
	services.MemberAuthenticateRequestStart(builder)
	services.MemberAuthenticateRequestAddUsername(builder, unameOffset)
	services.MemberAuthenticateRequestAddPassword(builder, upassOffset)
	r := services.MemberAuthenticateRequestEnd(builder)
	builder.Finish(r)

	req := services.GetRootAsMemberAuthenticateRequest(builder.FinishedBytes(), 0)
	_, err := ms.Authenticate(ctx, req, key)
	if err != nil {
		t.Fatalf("member authentication failed: %v", err)
	}
}

func doTestMemberChangePassword(t *testing.T, ctx context.Context, ms services.MemberService) {
	id := "testuser"
	oldPass := "U$hal1~PAss!"
	newPass := "And$0itBegins"

	builder := flatbuffers.NewBuilder(64)
	idOffset := builder.CreateString(id)
	oldPassOffset := builder.CreateString(oldPass)
	newPassOffset := builder.CreateString(newPass)
	services.MemberChangePasswordRequestStart(builder)
	services.MemberChangePasswordRequestAddId(builder, idOffset)
	services.MemberChangePasswordRequestAddOldPassword(builder, oldPassOffset)
	services.MemberChangePasswordRequestAddNewPassword(builder, newPassOffset)
	r := services.MemberChangePasswordRequestEnd(builder)
	builder.Finish(r)

	req := services.GetRootAsMemberChangePasswordRequest(builder.FinishedBytes(), 0)
	err := ms.ChangePassword(ctx, req)
	if err != nil {
		t.Fatalf("failed to change member password: %v", err)
	}
}

func doTestMemberFindByUsername(
	t *testing.T, ctx context.Context, ms services.MemberService, key crypto.Key, a *model.Member,
) *model.Member {
	uname := "testuser"

	builder := flatbuffers.NewBuilder(32)
	unameOffset := builder.CreateString(uname)
	services.MemberFindByUsernameRequestStart(builder)
	services.MemberFindByUsernameRequestAddUsername(builder, unameOffset)
	r := services.MemberFindByUsernameRequestEnd(builder)
	builder.Finish(r)

	req := services.GetRootAsMemberFindByUsernameRequest(builder.FinishedBytes(), 0)
	b, err := ms.FindMemberByUsername(ctx, req, key)
	if err != nil {
		t.Fatalf("find member by username failed: %v", err)
	}

	if !model.MemberEqual(a, b) {
		t.Errorf("member not the same: %+v != %+v", b, a)
	}

	return b
}

func doTestMemberUpdate(
	t *testing.T, ctx context.Context, ms services.MemberService, key crypto.Key, b *model.Member,
) *model.Member {
	uname := "updatedtestuser"
	name := "Bob Bill"

	builder := flatbuffers.NewBuilder(64)
	unameOffset := builder.CreateString(uname)
	nameOffset := builder.CreateString(name)
	services.MemberUpdateRequestStart(builder)
	services.MemberUpdateRequestAddUsername(builder, unameOffset)
	services.MemberUpdateRequestAddName(builder, nameOffset)
	r := services.MemberUpdateRequestEnd(builder)
	builder.Finish(r)

	req := services.GetRootAsMemberUpdateRequest(builder.FinishedBytes(), 0)
	c, err := ms.UpdateMember(ctx, req, key)
	if err != nil {
		t.Fatalf("failed to update member: %v", err)
	}

	if model.MemberEqual(b, c) {
		t.Errorf("member should be different after update: %+v == %+v", c, b)
	}

	return c
}

func doTestMemberList(
	t *testing.T, ctx context.Context, ms services.MemberService, key crypto.Key, c *model.Member,
) *model.Member {
	uname := "anotheruser"
	upass := "G0b@k2th#shadows!"
	name := "Carol Chris"

	builder := flatbuffers.NewBuilder(64)
	unameOffset := builder.CreateString(uname)
	upassOffset := builder.CreateString(upass)
	nameOffset := builder.CreateString(name)
	services.MemberCreateRequestStart(builder)
	services.MemberCreateRequestAddUsername(builder, unameOffset)
	services.MemberCreateRequestAddPassword(builder, upassOffset)
	services.MemberCreateRequestAddName(builder, nameOffset)
	r := services.MemberCreateRequestEnd(builder)
	builder.Finish(r)

	reqCreate := services.GetRootAsMemberCreateRequest(builder.FinishedBytes(), 0)
	d, err := ms.Create(ctx, reqCreate, key)
	if err != nil {
		t.Fatalf("failed to create second member: %v", err)
	}

	l, err := ms.ListMembers(ctx, key)
	if err != nil {
		t.Fatalf("failed to list members: %v", err)
	}

	if len(l) != 2 {
		t.Fatalf("expected two members in list, got %d", len(l))
	}

	if !model.MemberEqual(l[0], c) {
		t.Errorf("first member is not what was expected: %+v != %+v", l[0], c)
	}
	if !model.MemberEqual(l[1], d) {
		t.Errorf("second member is not what was expected: %+v != %+v", l[1], d)
	}

	return d
}

func doTestMemberRemove(
	t *testing.T,
	ctx context.Context,
	ms services.MemberService,
	key crypto.Key, c,
	d *model.Member,
) {
	builder := flatbuffers.NewBuilder(32)
	idOffset := builder.CreateByteString(d.Id())
	services.MemberRemoveRequestStart(builder)
	services.MemberRemoveRequestAddId(builder, idOffset)
	r := services.MemberRemoveRequestEnd(builder)
	builder.Finish(r)

	req := services.GetRootAsMemberRemoveRequest(builder.FinishedBytes(), 0)
	err := ms.RemoveMember(ctx, req)
	if err != nil {
		t.Errorf("failed to remove member: %v", err)
	}

	l, err := ms.ListMembers(ctx, key)
	if err != nil {
		t.Fatalf("failed to list members after removing one: %v", err)
	}
	if len(l) != 1 {
		t.Fatalf("expected only one member after delete, got %d", len(l))
	}
	if !model.MemberEqual(l[0], c) {
		t.Errorf("remaining member is not what was expected after delete: %+v != %+v", l[0], c)
	}
}

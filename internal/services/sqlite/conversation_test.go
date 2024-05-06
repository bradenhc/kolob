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

func TestConversationService(t *testing.T) {
	// Setup the test
	//
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "member-service-test.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	err = sqlite.CreateTables(db)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	ctx := context.Background()

	// Create the group this conversation and the mediator member will belong to
	//
	gpass, _ := crypto.NewPassword("Password123456!")
	gs := sqlite.NewGroupService(db)
	_, err = gs.CreateGroup(ctx, model.CreateGroupParams{
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

	// Add a member to use later as a mediator
	//
	ms := sqlite.NewMemberService(db)
	upass, _ := crypto.NewPassword("U$hal1~PAss!")
	m1, err := ms.CreateMember(ctx, model.CreateMemberParams{
		Name:     "alice",
		Username: "Alice Ann",
		Password: upass,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to create member: %v", err)
	}

	// Create the conversation service
	//
	cs := sqlite.NewConversationService(db)

	// Create a conversation
	//
	name := "First Conversation"
	c, err := cs.CreateConversation(ctx, model.CreateConversationParams{
		Name:      name,
		Moderator: m1.Id,
		PassKey:   pkey,
	})
	if err != nil {
		t.Fatalf("failed to create conversation: %v", err)
	}

	if c.Name != name {
		t.Errorf("conversation name does not match: %s != %s", c.Name, name)
	}
}

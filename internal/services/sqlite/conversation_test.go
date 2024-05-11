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

	ctx := context.Background()

	// Create the group this conversation and the mediator member will belong to
	//
	gpass, _ := crypto.NewPassword("Password123456!")
	gs := sqlite.NewGroupService(db)
	_, err = gs.InitGroup(ctx, model.CreateGroupParams{
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
	c1, err := cs.CreateConversation(ctx, model.CreateConversationParams{
		Name:      name,
		Moderator: m1.Id,
		PassKey:   pkey,
	})
	if err != nil {
		t.Fatalf("failed to create conversation: %v", err)
	}

	if c1.Name != name {
		t.Errorf("conversation name does not match: %s != %s", c1.Name, name)
	}

	// Update the conversation
	//
	name = "Updated First Conversation"
	err = cs.UpdateConversation(ctx, model.UpdateConversationParams{
		Id:      c1.Id,
		Name:    &name,
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to update conversation: %v", err)
	}

	c1, err = cs.FindConversationById(ctx, model.FindConversationByIdParams{
		Id:      c1.Id,
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to get conversation by ID after update: %v", err)
	}
	if c1.Name != name {
		t.Errorf("failed to update conversation name: %s != %s", c1.Name, name)
	}

	// Add another conversation and then list all available conversations
	//
	c2, err := cs.CreateConversation(ctx, model.CreateConversationParams{
		Name:      "Another Conversation",
		Moderator: m1.Id,
		PassKey:   pkey,
	})
	if err != nil {
		t.Fatalf("failed to create second conversation: %v", err)
	}

	clist, err := cs.ListConversations(ctx, model.ListConversationsParams{
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to list conversations: %v", err)
	}
	if len(clist) != 2 {
		t.Fatalf("not enough entries in conversation list")
	}

	// Remove the conversation we just added
	//
	err = cs.RemoveConversation(ctx, model.RemoveConversationParams{
		Id: c2.Id,
	})
	if err != nil {
		t.Fatalf("failed to remove second conversation: %v", err)
	}
	clist, err = cs.ListConversations(ctx, model.ListConversationsParams{
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to list conversations after removing one: %v", err)
	}
	if len(clist) != 1 {
		t.Fatalf("should only be one entry in conversation list after removing the second")
	}

	// Add additional moderators to the conversation
	//
	m2pass, _ := crypto.NewPassword("Password7654321!")
	m2, err := ms.CreateMember(ctx, model.CreateMemberParams{
		Name:     "Bob Bill",
		Username: "bob",
		Password: m2pass,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to create second member: %v", err)
	}
	m3pass, _ := crypto.NewPassword("Password098765!")
	m3, err := ms.CreateMember(ctx, model.CreateMemberParams{
		Name:     "Carol Chris",
		Username: "carol",
		Password: m3pass,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to create third member: %v", err)
	}

	err = cs.AddConversationMods(ctx, model.AddConversationModsParams{
		Id:         c1.Id,
		Moderators: []model.Uuid{m2.Id, m3.Id},
		PassKey:    pkey,
	})
	if err != nil {
		t.Fatalf("failed to add additional moderators: %v", err)
	}

	// List moderators
	//
	ids, err := cs.ListConversationMods(ctx, model.ListConversationModsParams{
		Id:      c1.Id,
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to list conversation mods: %v", err)
	}

	if len(ids) != 3 {
		t.Fatalf("failed to check moderator length: there should be three moderators")
	}
}

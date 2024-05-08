// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite_test

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/services/sqlite"
)

func TestMessageService(t *testing.T) {
	// Setup the test
	//
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "message-service-test.db")

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

	// Add a couple of members to use later as authors
	//
	mems := sqlite.NewMemberService(db)
	upass, _ := crypto.NewPassword("U$hal1~PAss!")
	m1, err := mems.CreateMember(ctx, model.CreateMemberParams{
		Name:     "Alice Ann",
		Username: "alice",
		Password: upass,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to create member alice: %v", err)
	}
	m2, err := mems.CreateMember(ctx, model.CreateMemberParams{
		Name:     "Bob Bill",
		Username: "bob",
		Password: upass,
		PassKey:  pkey,
	})
	if err != nil {
		t.Fatalf("failed to create member bob: %v", err)
	}

	// Create a couple of conversations to add the messages to
	//
	cs := sqlite.NewConversationService(db)
	c1, err := cs.CreateConversation(ctx, model.CreateConversationParams{
		Name:      "First Conversation",
		Moderator: m1.Id,
		PassKey:   pkey,
	})
	if err != nil {
		t.Fatalf("failed to create first conversation: %v", err)
	}
	c2, err := cs.CreateConversation(ctx, model.CreateConversationParams{
		Name:      "Second Conversation",
		Moderator: m2.Id,
		PassKey:   pkey,
	})
	if err != nil {
		t.Fatalf("failed to create second conversation: %v", err)
	}

	// Create the message service
	//
	ms := sqlite.NewMessageService(db)

	// Add several messages to each conversation
	//
	mlist1 := make([]model.Message, 0)
	for i := range 10 {
		m, err := ms.CreateMessage(ctx, model.CreateMessageParams{
			ConversationId: c1.Id,
			Author:         m1.Id,
			Content:        fmt.Sprintf("This is test message %d", i),
			PassKey:        pkey,
		})
		if err != nil {
			t.Fatalf("failed to add message %d to conversation: %v", i, err)
		}
		mlist1 = append(mlist1, m)
	}
	mlist2 := make([]model.Message, 0)
	for i := range 10 {
		m, err := ms.CreateMessage(ctx, model.CreateMessageParams{
			ConversationId: c2.Id,
			Author:         m2.Id,
			Content:        fmt.Sprintf("This is test message %d", i),
			PassKey:        pkey,
		})
		if err != nil {
			t.Fatalf("failed to add message %d to conversation: %v", i, err)
		}
		mlist2 = append(mlist2, m)
	}

	if len(mlist1) != len(mlist2) {
		t.Fatalf("failed to add same number of messages to both conversations")
	}

	// Get a specific message
	//
	_, err = ms.GetMessage(ctx, model.GetMessageParams{
		Id:      mlist1[0].Id,
		PassKey: pkey,
	})
	if err != nil {
		t.Fatalf("failed to get single message: %v", err)
	}

	// Get messages for only one conversation
	//
	mlist3, err := ms.ListMessages(ctx, model.ListMessagesParams{
		ConversationId: c1.Id,
		PassKey:        pkey,
	})
	if err != nil {
		t.Fatalf("failed to list messages for first conversation: %v", err)
	}
	if len(mlist1) != len(mlist3) {
		t.Fatalf("missing messages from first conversation list")
	}

	mlist4, err := ms.ListMessages(ctx, model.ListMessagesParams{
		ConversationId: c2.Id,
		PassKey:        pkey,
	})
	if err != nil {
		t.Fatalf("failed to list messages for second conversation: %v", err)
	}
	if len(mlist2) != len(mlist4) {
		t.Fatalf("missing messages from second conversation list")
	}

	// Remove a message
	//
	err = ms.RemoveMessage(ctx, model.RemoveMessageParams{
		Id: mlist1[0].Id,
	})
	if err != nil {
		t.Fatalf("failed to remove message from first conversation: %v", err)
	}
	mlist5, err := ms.ListMessages(ctx, model.ListMessagesParams{
		ConversationId: c1.Id,
		PassKey:        pkey,
	})
	if err != nil {
		t.Fatalf("failed to list messages for first conversation after remove: %v", err)
	}
	if len(mlist5) != len(mlist1)-1 {
		t.Fatalf("incorrect number of messages listed after removing from first conversation")
	}
}

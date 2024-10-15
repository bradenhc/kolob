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
	"github.com/bradenhc/kolob/internal/store"
	"github.com/bradenhc/kolob/internal/store/sqlite"
	flatbuffers "github.com/google/flatbuffers/go"
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

	ctx := context.Background()

	// Setup group
	groupStore := doTestGroupCreateStore(t, db)
	svcGroup := services.NewGroupService(groupStore)
	doTestGroupCreate(t, ctx, svcGroup)
	key := doTestGroupAuth(t, ctx, svcGroup)

	// Setup members
	memberStore := doTestMemberCreateStore(t, db)
	svcMember := services.NewMemberService(memberStore)
	member1 := doTestMemberAdd(t, ctx, svcMember, key)
	member2 := doTestMemberAdd(t, ctx, svcMember, key)

	// Setup conversations
	convoStore := doTestConversationCreateStore(t, db)
	svcConvo := services.NewConversationService(convoStore)
	convo1 := doTestConversationAdd(t, ctx, svcConvo, key, member1)
	convo2 := doTestConversationAdd(t, ctx, svcConvo, key, member2)

	// Create the message store and service
	messageStore := doTestMessageCreateStore(t, db)
	svcMessage := services.NewMessageService(messageStore)

	// Add messages to the first conversation
	message1 := doTestMessageAdd(t, ctx, svcMessage, key, convo1, member1, "Hey there!")
	message2 := doTestMessageAdd(t, ctx, svcMessage, key, convo1, member2, "Hey! How's it going?")
	message3 := doTestMessageAdd(t, ctx, svcMessage, key, convo1, member1, "Great! Wanna hack?")

	// Add messages to the second conversation
	message4 := doTestMessageAdd(t, ctx, svcMessage, key, convo2, member2, "Let's talk here.")
	message5 := doTestMessageAdd(t, ctx, svcMessage, key, convo2, member1, "Okay! Another convo!")
	message6 := doTestMessageAdd(t, ctx, svcMessage, key, convo2, member2, "Can never have enough!")

	doTestMessageList(t, ctx, svcMessage, key, convo1, message1, message2, message3)
	doTestMessageList(t, ctx, svcMessage, key, convo2, message4, message5, message6)
}

func doTestMessageCreateStore(t *testing.T, db *sql.DB) store.MessageStore {
	store, err := sqlite.NewMessageStore(db)
	if err != nil {
		t.Fatalf("failed to create message store: %v", err)
	}

	return store
}

func doTestMessageAdd(
	t *testing.T,
	ctx context.Context,
	ms services.MessageService,
	key crypto.Key,
	convo *model.Conversation,
	author *model.Member,
	content string) *model.Message {

	builder := flatbuffers.NewBuilder(256)
	offsetConvo := builder.CreateByteString(convo.Id())
	offsetAuthor := builder.CreateByteString(author.Id())
	offsetContent := builder.CreateString(content)
	services.MessageAddRequestStart(builder)
	services.MessageAddRequestAddConversation(builder, offsetConvo)
	services.MessageAddRequestAddAuthor(builder, offsetAuthor)
	services.MessageAddRequestAddContent(builder, offsetContent)

	offsetAddRequest := services.MessageAddRequestEnd(builder)
	builder.Finish(offsetAddRequest)

	addRequest := services.GetRootAsMessageAddRequest(builder.FinishedBytes(), 0)

	message, err := ms.Add(ctx, addRequest, key)
	if err != nil {
		t.Fatalf("failed to add message: %v", err)
	}

	if !slices.Equal(message.Content(), []byte(content)) {
		t.Errorf("incorrect message content: '%s' != '%s'", message.Content(), content)
	}
	if !slices.Equal(convo.Id(), message.Conversation()) {
		t.Errorf("incorrect message conversation: %s != %s", convo.Id(), message.Conversation())
	}
	if !slices.Equal(author.Id(), message.Author()) {
		t.Errorf("incorrect message author: %s != %s", author.Id(), message.Author())
	}

	return message
}

func doTestMessageList(
	t *testing.T,
	ctx context.Context,
	ms services.MessageService,
	key crypto.Key,
	convo *model.Conversation,
	expected ...*model.Message,
) {
	builder := flatbuffers.NewBuilder(64)
	offsetId := builder.CreateByteString(convo.Id())

	services.MessageListRequestStart(builder)
	services.MessageListRequestAddConversation(builder, offsetId)

	offsetRequest := services.MessageListRequestEnd(builder)
	builder.Finish(offsetRequest)

	request := services.GetRootAsMessageListRequest(builder.FinishedBytes(), 0)
	messages, err := ms.List(ctx, request, key)
	if err != nil {
		t.Fatalf("failed to list messages: %v", err)
	}

	if len(messages) != len(expected) {
		t.Fatalf("bad length: messages != expected: %d != %d", len(messages), len(expected))
	}
}

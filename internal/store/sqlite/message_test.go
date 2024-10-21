// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package sqlite_test

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"slices"
	"testing"
	"time"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
	"github.com/bradenhc/kolob/internal/store"
	"github.com/bradenhc/kolob/internal/store/sqlite"
)

func TestMessageStoreSqlite(t *testing.T) {
	// Setup the test
	t.Parallel()
	tempdir := t.TempDir()
	dbpath := path.Join(tempdir, "kolob.db")

	db, err := sqlite.Open(dbpath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Create a random key we will use for encryption/decryption
	key, err := crypto.NewRandomKey()
	if err != nil {
		t.Fatalf("failed to create encryption key: %v", err)
	}

	// Create a user id to use for a moderator
	moderator, err := model.NewUuid()
	if err != nil {
		t.Fatalf("failed to create member UUID: %v", err)
	}

	// Create a member
	memberStore := doTestMemberStoreSqliteCreate(t, db)
	memberId := doTestMemberStoreSqliteInsert(t, memberStore, key)

	// Create a conversation
	conversationStore := doTestConversationStoreSqliteCreate(t, db)
	conversationEntity := doTestConversationStoreSqliteInsert(t, conversationStore, key, moderator)

	// Run tests on the message store
	messageStore := doTestMessageStoreSqliteCreate(t, db)
	messageId := doTestMessageStoreSqliteInsert(t, messageStore, key, memberId, conversationEntity.Id)
	messageEntity := doTestMessageStoreSqliteGet(t, messageStore, key, memberId, conversationEntity.Id, messageId)
	doTestMessageStoreSqliteUpdate(t, messageStore, key, messageEntity)
	doTestMessageStoreSqliteList(t, messageStore, key, memberId, conversationEntity.Id)
	doTestMessageStoreSqliteRemove(t, messageStore, conversationEntity.Id, messageId)
}

func doTestMessageStoreSqliteCreate(t *testing.T, db *sql.DB) store.MessageStore {
	store, err := sqlite.NewMessageStore(db)
	if err != nil {
		t.Fatalf("failed to create message store: %v", err)
	}

	return store
}

func doTestMessageStoreSqliteInsert(
	t *testing.T, s store.MessageStore, k crypto.Key, memberId, conversationId model.Uuid,
) model.Uuid {
	m, err := model.NewMessage(memberId, conversationId, "Hello, world!")
	if err != nil {
		t.Fatalf("failed to create new message: %v", err)
	}

	entity, err := store.NewMessageEntity(m, k)
	if err != nil {
		t.Fatalf("failed to create message entity: %v", err)
	}

	err = s.AddMessageEntity(context.Background(), entity)
	if err != nil {
		t.Fatalf("failed to add message entity: %v", err)
	}
	return entity.Id
}

func doTestMessageStoreSqliteGet(
	t *testing.T, s store.MessageStore, k crypto.Key, memberId, conversationId, messageId model.Uuid,
) store.MessageEntity {
	entity, err := s.GetMessageEntity(context.Background(), messageId)
	if err != nil {
		t.Fatalf("failed to get message entity: %v", err)
	}

	m, err := entity.Decrypt(k)
	if err != nil {
		t.Fatalf("failed to decrypt message: %v", err)
	}

	if !slices.Equal(m.Author(), []byte(memberId)) {
		t.Errorf("author ids do not match")
	}

	if !slices.Equal(m.Conversation(), []byte(conversationId)) {
		t.Errorf("conversation ids do not match")
	}

	if !slices.Equal(m.Content(), []byte("Hello, world!")) {
		t.Errorf("message contents do not match")
	}

	if m.Created() != m.Updated() {
		t.Errorf("message create/update times should be the same but are not")
	}

	return entity
}

func doTestMessageStoreSqliteUpdate(
	t *testing.T, s store.MessageStore, k crypto.Key, e store.MessageEntity,
) {
	content := []byte("Hello, test!")
	next, err := e.Update(k, content)
	if err != nil {
		t.Fatalf("failed to update message entity: %v", err)
	}

	if !slices.Equal(next.Content(), content) {
		t.Errorf("content did not update in message entity")
	}

	if next.Created() == next.Updated() {
		t.Errorf("created/updated times should be different, but are the same")
	}

	err = s.UpdateMessageEntity(context.Background(), e)
	if err != nil {
		t.Fatalf("failed to store updated message entity: %v", err)
	}
}

func doTestMessageStoreSqliteList(
	t *testing.T, s store.MessageStore, k crypto.Key, memberId, conversationId model.Uuid,
) {
	now := time.Now().UnixMilli()
	time.Sleep(1 * time.Second)

	for i := range 3 {
		m, err := model.NewMessage(memberId, conversationId, fmt.Sprintf("TestMessage%d", i))
		if err != nil {
			t.Fatalf("failed to create message list number %d: %v", i, err)
		}

		e, err := store.NewMessageEntity(m, k)
		if err != nil {
			t.Fatalf("failed to create message list entity number %d: %v", i, err)
		}

		err = s.AddMessageEntity(context.Background(), e)
		if err != nil {
			t.Fatalf("failed to store message list entity number %d: %v", i, err)
		}
	}

	var query store.ListMessageDataQuery
	entities, err := s.ListMessageEntities(context.Background(), conversationId, query)
	if err != nil {
		t.Fatalf("failed to list all messages: %v", err)
	}

	if len(entities) != 4 {
		t.Errorf("expected a total of 4 messages: got %d", len(entities))
	}

	query.CreatedAfter = new(int64)
	*query.CreatedAfter = now

	entities, err = s.ListMessageEntities(context.Background(), conversationId, query)
	if err != nil {
		t.Fatalf("failed to list messages created after list test: %v", err)
	}

	if len(entities) != 3 {
		t.Errorf("expected a total of 3 messages: got %d", len(entities))
	}
}

func doTestMessageStoreSqliteRemove(
	t *testing.T, s store.MessageStore, conversationId, messageId model.Uuid,
) {
	err := s.RemoveMessageEntity(context.Background(), messageId)
	if err != nil {
		t.Fatalf("failed to remove message from store: %v", err)
	}

	var query store.ListMessageDataQuery
	entities, err := s.ListMessageEntities(context.Background(), conversationId, query)
	if err != nil {
		t.Fatalf("failed to list messages after remove: %v", err)
	}

	if len(entities) != 3 {
		t.Errorf("expected a total of 3 messages: got %d", len(entities))
	}
}

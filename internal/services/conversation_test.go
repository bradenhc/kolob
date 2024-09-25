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
	gstore := doTestGroupCreateStore(t, db)
	gs := services.NewGroupService(gstore)
	doTestGroupCreate(t, ctx, gs)
	key := doTestGroupAuth(t, ctx, gs)

	// Add a member to use later as a mediator
	mstore := doTestMemberCreateStore(t, db)
	ms := services.NewMemberService(mstore)
	m1 := doTestMemberAdd(t, ctx, ms, key)

	// Create the conversation store and service
	cstore := doTestConversationCreateStore(t, db)
	cs := services.NewConversationService(cstore)

	// Create a conversation
	a := doTestConversationAdd(t, ctx, cs, key, m1)

	// Update the conversation
	b := doTestConversationUpdate(t, ctx, cs, key, a)

	// List conversations
	c := doTestConversationListAll(t, ctx, cs, key, m1, b)

	// Remove conversations
	doTestConversationRemove(t, ctx, cs, key, b)

	// Add mods to conversation
	m2 := doTestMemberAdd(t, ctx, ms, key)
	m3 := doTestMemberAdd(t, ctx, ms, key)
	doTestConversationModsAdd(t, ctx, cs, key, c, m2, m3)

	// Remove mods from conversation
	doTestConversationModsRemove(t, ctx, cs, key, c, m2)
}

func doTestConversationCreateStore(t *testing.T, db *sql.DB) store.ConversationStore {
	cstore, err := sqlite.NewConversationStore(db)
	if err != nil {
		t.Fatalf("failed to setup conversation store: %v", err)
	}

	return cstore
}

func doTestConversationAdd(
	t *testing.T,
	ctx context.Context,
	cs services.ConversationService,
	key crypto.Key,
	m *model.Member,
) *model.Conversation {
	name := "Test Conversation"
	desc := "A description for the test conversation"
	builder := flatbuffers.NewBuilder(32)
	nameOffset := builder.CreateString(name)
	descOffset := builder.CreateString(desc)

	modsElOffset := builder.CreateByteString(m.Id())
	services.ConversationAddRequestStartModeratorsVector(builder, 1)
	builder.PrependUOffsetT(modsElOffset)
	modsOffset := builder.EndVector(1)

	services.ConversationAddRequestStart(builder)
	services.ConversationAddRequestAddName(builder, nameOffset)
	services.ConversationAddRequestAddDescription(builder, descOffset)
	services.ConversationAddRequestAddModerators(builder, modsOffset)
	reqOffset := services.ConversationAddRequestEnd(builder)
	builder.Finish(reqOffset)

	req := services.GetRootAsConversationAddRequest(builder.FinishedBytes(), 0)
	c, err := cs.Add(ctx, req, key)
	if err != nil {
		t.Fatalf("failed to add conversation: %v", err)
	}

	if name != string(c.Name()) {
		t.Errorf("conversation name incorrect: %s != %s", name, c.Name())
	}
	if desc != string(c.Desc()) {
		t.Errorf("conversation description incorrect: %s != %s", desc, c.Desc())
	}
	if c.ModsLength() != 1 {
		t.Errorf("conversation mods missing member with id %s", m.Id())
	}
	if !slices.Equal(c.Mods(0), m.Id()) {
		t.Errorf("conversation mod id incorrect: %s != %s", c.Mods(0), m.Id())
	}

	return c
}

func doTestConversationUpdate(
	t *testing.T,
	ctx context.Context,
	cs services.ConversationService,
	key crypto.Key,
	a *model.Conversation,
) *model.Conversation {
	name := "Updated Test Conversation"
	desc := "Update description for the test conversation"

	builder := flatbuffers.NewBuilder(64)
	idOffset := builder.CreateByteString(a.Id())
	nameOffset := builder.CreateString(name)
	descOffset := builder.CreateString(desc)

	services.ConversationUpdateRequestStart(builder)
	services.ConversationUpdateRequestAddId(builder, idOffset)
	services.ConversationUpdateRequestAddName(builder, nameOffset)
	services.ConversationUpdateRequestAddDescription(builder, descOffset)

	reqOffset := services.ConversationUpdateRequestEnd(builder)
	builder.Finish(reqOffset)

	reqUpdate := services.GetRootAsConversationUpdateRequest(builder.FinishedBytes(), 0)
	err := cs.Update(ctx, reqUpdate, key)
	if err != nil {
		t.Fatalf("failed to update conversation: %v", err)
	}

	builder.Reset()
	idOffset = builder.CreateByteString(a.Id())
	services.ConversationGetRequestStart(builder)
	services.ConversationGetRequestAddId(builder, idOffset)

	reqOffset = services.ConversationGetRequestEnd(builder)
	builder.Finish(reqOffset)

	reqGet := services.GetRootAsConversationGetRequest(builder.FinishedBytes(), 0)
	b, err := cs.Get(ctx, reqGet, key)
	if err != nil {
		t.Fatalf("failed to get updated conversation: %v", err)
	}

	if string(b.Name()) != name {
		t.Errorf("updated name is incorrect: %s != %s", b.Name(), name)
	}
	if string(b.Desc()) != desc {
		t.Errorf("updated description is incorrect: %s != %s", b.Desc(), desc)
	}
	if b.Updated() == a.Updated() {
		t.Errorf("updated conversation timestamp not modified")
	}

	return b
}

func doTestConversationListAll(
	t *testing.T,
	ctx context.Context,
	cs services.ConversationService,
	key crypto.Key,
	m *model.Member,
	b *model.Conversation,
) *model.Conversation {
	c := doTestConversationAdd(t, ctx, cs, key, m)
	d := doTestConversationAdd(t, ctx, cs, key, m)

	l, err := cs.ListAll(ctx, key)
	if err != nil {
		t.Fatalf("failed to list all conversations: %v", err)
	}

	if len(l) != 3 {
		t.Fatalf("expected 3 conversations: found %d", len(l))
	}

	if !model.ConversationEqual(l[0], b) {
		t.Errorf("first conversations do not match")
	}
	if !model.ConversationEqual(l[1], c) {
		t.Errorf("second conversations do not match")
	}
	if !model.ConversationEqual(l[2], d) {
		t.Errorf("third conversations do not match")
	}

	return c
}

func doTestConversationRemove(
	t *testing.T,
	ctx context.Context,
	cs services.ConversationService,
	key crypto.Key,
	b *model.Conversation,
) {
	builder := flatbuffers.NewBuilder(32)
	idOffset := builder.CreateByteString(b.Id())

	services.ConversationRemoveRequestStart(builder)
	services.ConversationRemoveRequestAddId(builder, idOffset)
	reqOffset := services.ConversationRemoveRequestEnd(builder)
	builder.Finish(reqOffset)

	req := services.GetRootAsConversationRemoveRequest(builder.FinishedBytes(), 0)
	err := cs.Remove(ctx, req)
	if err != nil {
		t.Fatalf("failed to remove conversation: %v", err)
	}

	l, err := cs.ListAll(ctx, key)
	if err != nil {
		t.Fatalf("failed to list conversations after removing one: %v", err)
	}

	if len(l) != 2 {
		t.Errorf("too many conversations in list after remove: %d != %d", len(l), 2)
	}
}

func doTestConversationModsAdd(
	t *testing.T,
	ctx context.Context,
	cs services.ConversationService,
	key crypto.Key,
	c *model.Conversation,
	m2, m3 *model.Member,
) {
	builder := flatbuffers.NewBuilder(64)

	cIdOffset := builder.CreateByteString(c.Id())
	m2IdOffset := builder.CreateByteString(m2.Id())
	m3IdOffset := builder.CreateByteString(m3.Id())

	services.ConversationModsAddRequestStartModeratorsVector(builder, 2)
	builder.PrependUOffsetT(m2IdOffset)
	builder.PrependUOffsetT(m3IdOffset)
	modsOffset := builder.EndVector(2)

	services.ConversationModsAddRequestStart(builder)
	services.ConversationModsAddRequestAddId(builder, cIdOffset)
	services.ConversationModsAddRequestAddModerators(builder, modsOffset)
	reqOffset := services.ConversationModsAddRequestEnd(builder)
	builder.Finish(reqOffset)

	reqAddMods := services.GetRootAsConversationModsAddRequest(builder.FinishedBytes(), 0)
	err := cs.AddMods(ctx, reqAddMods, key)
	if err != nil {
		t.Fatalf("failed to add moderators to conversation: %v", err)
	}

	builder.Reset()
	cIdOffset = builder.CreateByteString(c.Id())

	services.ConversationGetRequestStart(builder)
	services.ConversationGetRequestAddId(builder, cIdOffset)
	reqOffset = services.ConversationGetRequestEnd(builder)
	builder.Finish(reqOffset)

	reqGet := services.GetRootAsConversationGetRequest(builder.FinishedBytes(), 0)
	d, err := cs.Get(ctx, reqGet, key)
	if err != nil {
		t.Fatalf("failed to get conversation after adding moderators: %v", err)
	}

	if c.ModsLength() == d.ModsLength() {
		t.Errorf("expected conversation mods list to be different after add")
	}

	if d.ModsLength() != 3 {
		t.Errorf("incorrect number of mods after add: %d != %d", d.ModsLength(), 3)
	}
}

func doTestConversationModsRemove(
	t *testing.T,
	ctx context.Context,
	cs services.ConversationService,
	key crypto.Key,
	c *model.Conversation,
	m2 *model.Member,
) {
	builder := flatbuffers.NewBuilder(64)

	cIdOffset := builder.CreateByteString(c.Id())
	m2IdOffset := builder.CreateByteString(m2.Id())

	services.ConversationModsAddRequestStartModeratorsVector(builder, 1)
	builder.PrependUOffsetT(m2IdOffset)
	modsOffset := builder.EndVector(1)

	services.ConversationModsAddRequestStart(builder)
	services.ConversationModsAddRequestAddId(builder, cIdOffset)
	services.ConversationModsAddRequestAddModerators(builder, modsOffset)
	reqOffset := services.ConversationModsAddRequestEnd(builder)
	builder.Finish(reqOffset)

	reqAddMods := services.GetRootAsConversationModsRemoveRequest(builder.FinishedBytes(), 0)
	err := cs.RemoveMods(ctx, reqAddMods, key)
	if err != nil {
		t.Fatalf("failed to remove moderators from conversation: %v", err)
	}

	builder.Reset()
	cIdOffset = builder.CreateByteString(c.Id())

	services.ConversationGetRequestStart(builder)
	services.ConversationGetRequestAddId(builder, cIdOffset)
	reqOffset = services.ConversationGetRequestEnd(builder)
	builder.Finish(reqOffset)

	reqGet := services.GetRootAsConversationGetRequest(builder.FinishedBytes(), 0)
	d, err := cs.Get(ctx, reqGet, key)
	if err != nil {
		t.Fatalf("failed to get conversation after removing moderators: %v", err)
	}

	if c.ModsLength() == d.ModsLength() {
		t.Errorf("expected conversation mods list to be different after add")
	}

	if d.ModsLength() != 2 {
		t.Errorf("incorrect number of mods after add: %d != %d", d.ModsLength(), 2)
	}

	for i := range d.ModsLength() {
		if slices.Equal(d.Mods(i), m2.Id()) {
			t.Errorf("found moderator that should have been removed: %v", m2.Id())
		}
	}
}

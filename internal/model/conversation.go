// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package model

import (
	"fmt"
	"slices"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

func NewConversation(name, desc string, mods []Uuid) (*Conversation, error) {
	uuid, err := NewUuid()
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %v", err)
	}

	now := time.Now().UnixMilli()

	builder := flatbuffers.NewBuilder(1024)

	idOffsets := builder.CreateString(string(uuid))
	nameOffset := builder.CreateString(name)
	descOffset := builder.CreateString(desc)
	modsElsOffsets := make([]flatbuffers.UOffsetT, len(mods))
	for _, m := range mods {
		modsElsOffsets = append(modsElsOffsets, builder.CreateString(string(m)))
	}
	ConversationStartModsVector(builder, len(mods))
	for _, m := range modsElsOffsets {
		builder.PrependUOffsetT(m)
	}
	modsOffset := builder.EndVector(len(modsElsOffsets))

	ConversationStart(builder)
	ConversationAddId(builder, idOffsets)
	ConversationAddName(builder, nameOffset)
	ConversationAddDesc(builder, descOffset)
	ConversationAddMods(builder, modsOffset)
	ConversationAddCreated(builder, now)
	ConversationAddUpdated(builder, now)
	convOffset := ConversationEnd(builder)

	builder.Finish(convOffset)

	return GetRootAsConversation(builder.FinishedBytes(), 0), nil
}

func CloneConversationWithUpdates(prev *Conversation, name, desc []byte, mods [][]byte) *Conversation {
	now := time.Now().UnixMilli()

	builder := flatbuffers.NewBuilder(1024)

	idOffsets := builder.CreateByteString(prev.Id())

	var nameOffset flatbuffers.UOffsetT
	if name != nil {
		nameOffset = builder.CreateByteString(name)
	} else {
		nameOffset = builder.CreateByteString(prev.Name())
	}

	var descOffset flatbuffers.UOffsetT
	if desc != nil {
		descOffset = builder.CreateByteString(desc)
	} else {
		descOffset = builder.CreateByteString(prev.Desc())
	}

	modsElsOffsets := make([]flatbuffers.UOffsetT, 0)
	if mods != nil {
		for _, m := range mods {
			modsElsOffsets = append(modsElsOffsets, builder.CreateString(string(m)))
		}
	} else {
		for i := range prev.ModsLength() {
			modsElsOffsets = append(modsElsOffsets, builder.CreateByteString(prev.Mods(i)))
		}
	}

	ConversationStartModsVector(builder, len(modsElsOffsets))
	for _, m := range modsElsOffsets {
		builder.PrependUOffsetT(m)
	}
	modsOffset := builder.EndVector(len(modsElsOffsets))

	ConversationStart(builder)
	ConversationAddId(builder, idOffsets)
	ConversationAddName(builder, nameOffset)
	ConversationAddDesc(builder, descOffset)
	ConversationAddMods(builder, modsOffset)
	ConversationAddCreated(builder, prev.Created())
	ConversationAddUpdated(builder, now)
	convOffset := ConversationEnd(builder)

	builder.Finish(convOffset)

	return GetRootAsConversation(builder.FinishedBytes(), 0)
}

func ConversationEqual(a, b *Conversation) bool {
	if a == b {
		return true
	}

	if slices.Equal(a.Id(), b.Id()) &&
		slices.Equal(a.Name(), b.Name()) &&
		slices.Equal(a.Desc(), b.Desc()) &&
		a.Created() != b.Created() &&
		a.Updated() != b.Updated() &&
		a.ModsLength() == b.ModsLength() {
		// Make sure all the mods are equal. Order is not important.
		amods := make(map[string]bool, a.ModsLength())
		for i := range a.ModsLength() {
			amods[string(a.Mods(i))] = true
		}
		for i := range b.ModsLength() {
			_, ok := amods[string(b.Mods(i))]
			if !ok {
				return false
			}
		}

		return true
	}

	return false
}

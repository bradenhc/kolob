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

func NewMessage(author, convo Uuid, content string) (*Message, error) {
	uuid, err := NewUuid()
	if err != nil {
		return nil, fmt.Errorf("failed to create new message: %v", err)
	}

	now := time.Now().UnixMilli()

	builder := flatbuffers.NewBuilder(1024)
	idOffset := builder.CreateString(string(uuid))
	authorOffset := builder.CreateString(string(author))
	convoOffset := builder.CreateString(string(convo))
	contentOffset := builder.CreateString(content)

	MessageStart(builder)
	MessageAddId(builder, idOffset)
	MessageAddAuthor(builder, authorOffset)
	MessageAddConversation(builder, convoOffset)
	MessageAddContent(builder, contentOffset)
	MessageAddCreated(builder, now)
	MessageAddUpdated(builder, now)

	msgOffset := MessageEnd(builder)
	builder.Finish(msgOffset)

	return GetRootAsMessage(builder.FinishedBytes(), 0), nil

}

func CloneMessageWithUpdates(prev *Message, content []byte) *Message {
	now := time.Now().UnixMilli()

	builder := flatbuffers.NewBuilder(1024)
	idOffset := builder.CreateByteString(prev.Id())
	authorOffset := builder.CreateByteString(prev.Author())
	convoOffset := builder.CreateByteString(prev.Conversation())
	contentOffset := builder.CreateByteString(content)

	MessageStart(builder)
	MessageAddId(builder, idOffset)
	MessageAddAuthor(builder, authorOffset)
	MessageAddConversation(builder, convoOffset)
	MessageAddContent(builder, contentOffset)
	MessageAddCreated(builder, prev.Created())
	MessageAddUpdated(builder, now)

	msgOffset := MessageEnd(builder)
	builder.Finish(msgOffset)

	return GetRootAsMessage(builder.FinishedBytes(), 0)
}

func MessageEqual(a, b *Message) bool {
	if a == b {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if slices.Equal(a.Id(), b.Id()) &&
		slices.Equal(a.Author(), b.Author()) &&
		slices.Equal(a.Conversation(), b.Conversation()) &&
		slices.Equal(a.Content(), b.Content()) &&
		a.Created() == b.Created() &&
		a.Updated() == b.Updated() {
		return true
	}

	return false
}

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

func NewMember(username, name string) (*Member, error) {
	uuid, err := NewUuid()
	if err != nil {
		return nil, fmt.Errorf("failed to create new member: %v", err)
	}

	builder := flatbuffers.NewBuilder(256)

	mi := builder.CreateString(string(uuid))
	mu := builder.CreateString(username)
	mn := builder.CreateString(name)

	now := time.Now()

	MemberStart(builder)
	MemberAddId(builder, mi)
	MemberAddUname(builder, mu)
	MemberAddName(builder, mn)
	MemberAddCreated(builder, now.UnixMilli())
	MemberAddUpdated(builder, now.UnixMilli())

	m := MemberEnd(builder)

	builder.Finish(m)

	return GetRootAsMember(builder.FinishedBytes(), 0), nil
}

func CloneMemberWithUpdates(prev *Member, uname, name []byte) *Member {
	builder := flatbuffers.NewBuilder(64)
	mi := builder.CreateByteString(prev.Id())

	var mu flatbuffers.UOffsetT
	if uname != nil {
		mu = builder.CreateByteString(uname)
	} else {
		mu = builder.CreateByteString(prev.Uname())
	}

	var mn flatbuffers.UOffsetT
	if name != nil {
		mn = builder.CreateByteString(name)
	} else {
		mn = builder.CreateByteString(prev.Name())
	}

	updated := time.Now().UnixMilli()

	MemberStart(builder)
	MemberAddId(builder, mi)
	MemberAddUname(builder, mu)
	MemberAddName(builder, mn)
	MemberAddCreated(builder, prev.Created())
	MemberAddUpdated(builder, updated)

	m := MemberEnd(builder)
	builder.Finish(m)

	return GetRootAsMember(builder.FinishedBytes(), 0)
}

func MemberEqual(a, b *Member) bool {
	if a != b {
		if a == nil || b == nil ||
			!slices.Equal(a.Id(), b.Id()) ||
			!slices.Equal(a.Uname(), b.Uname()) ||
			!slices.Equal(a.Name(), b.Name()) ||
			a.Created() != b.Created() ||
			a.Updated() != b.Updated() {
			return false
		}
	}
	return true
}

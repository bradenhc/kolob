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

func NewGroup(gid, name, desc string) (*Group, error) {
	id, err := NewUuid()
	if err != nil {
		return nil, fmt.Errorf("failed to create UUID for new group: %v", err)
	}

	builder := flatbuffers.NewBuilder(64)

	gi := builder.CreateString(string(id))
	gg := builder.CreateString(gid)
	gn := builder.CreateString(name)
	gd := builder.CreateString(desc)

	now := time.Now()

	GroupStart(builder)
	GroupAddId(builder, gi)
	GroupAddGid(builder, gg)
	GroupAddName(builder, gn)
	GroupAddDesc(builder, gd)
	GroupAddCreated(builder, now.UnixMilli())
	GroupAddUpdated(builder, now.UnixMilli())

	g := GroupEnd(builder)

	builder.Finish(g)

	return GetRootAsGroup(builder.FinishedBytes(), 0), nil
}

func GroupCloneWithUpdates(prev *Group, gid, name, desc []byte) *Group {
	builder := flatbuffers.NewBuilder(64)
	gi := builder.CreateByteString(prev.Id())

	var gh flatbuffers.UOffsetT
	if gid != nil {
		gh = builder.CreateByteString(gid)
	} else {
		gh = builder.CreateByteString(prev.Gid())
	}

	var gn flatbuffers.UOffsetT
	if name != nil {
		gn = builder.CreateByteString(name)
	} else {
		gn = builder.CreateByteString(prev.Name())
	}

	var gd flatbuffers.UOffsetT
	if desc != nil {
		gd = builder.CreateByteString(desc)
	} else {
		gd = builder.CreateByteString(prev.Desc())
	}

	updated := time.Now().UnixMilli()

	GroupStart(builder)
	GroupAddId(builder, gi)
	GroupAddGid(builder, gh)
	GroupAddName(builder, gn)
	GroupAddDesc(builder, gd)
	GroupAddCreated(builder, prev.Created())
	GroupAddUpdated(builder, updated)

	g := GroupEnd(builder)

	builder.Finish(g)

	return GetRootAsGroup(builder.FinishedBytes(), 0)
}

func GroupEqual(a, b *Group) bool {
	if a != b {
		if a == nil || b == nil ||
			!slices.Equal(a.Id(), b.Id()) ||
			!slices.Equal(a.Gid(), b.Gid()) ||
			!slices.Equal(a.Name(), b.Name()) ||
			!slices.Equal(a.Desc(), b.Desc()) ||
			a.Created() != b.Created() ||
			a.Updated() != b.Updated() {
			return false
		}
	}
	return true
}

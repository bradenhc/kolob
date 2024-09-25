// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package model

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Group struct {
	_tab flatbuffers.Table
}

func GetRootAsGroup(buf []byte, offset flatbuffers.UOffsetT) *Group {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Group{}
	x.Init(buf, n+offset)
	return x
}

func FinishGroupBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsGroup(buf []byte, offset flatbuffers.UOffsetT) *Group {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Group{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedGroupBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Group) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Group) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Group) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Group) Gid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Group) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Group) Desc() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Group) Created() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Group) MutateCreated(n int64) bool {
	return rcv._tab.MutateInt64Slot(12, n)
}

func (rcv *Group) Updated() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Group) MutateUpdated(n int64) bool {
	return rcv._tab.MutateInt64Slot(14, n)
}

func GroupStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func GroupAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func GroupAddGid(builder *flatbuffers.Builder, gid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(gid), 0)
}
func GroupAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(name), 0)
}
func GroupAddDesc(builder *flatbuffers.Builder, desc flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(desc), 0)
}
func GroupAddCreated(builder *flatbuffers.Builder, created int64) {
	builder.PrependInt64Slot(4, created, 0)
}
func GroupAddUpdated(builder *flatbuffers.Builder, updated int64) {
	builder.PrependInt64Slot(5, updated, 0)
}
func GroupEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Member struct {
	_tab flatbuffers.Table
}

func GetRootAsMember(buf []byte, offset flatbuffers.UOffsetT) *Member {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Member{}
	x.Init(buf, n+offset)
	return x
}

func FinishMemberBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMember(buf []byte, offset flatbuffers.UOffsetT) *Member {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Member{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMemberBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Member) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Member) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Member) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Member) Uname() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Member) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Member) Created() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Member) MutateCreated(n int64) bool {
	return rcv._tab.MutateInt64Slot(10, n)
}

func (rcv *Member) Updated() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Member) MutateUpdated(n int64) bool {
	return rcv._tab.MutateInt64Slot(12, n)
}

func MemberStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func MemberAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MemberAddUname(builder *flatbuffers.Builder, uname flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(uname), 0)
}
func MemberAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(name), 0)
}
func MemberAddCreated(builder *flatbuffers.Builder, created int64) {
	builder.PrependInt64Slot(3, created, 0)
}
func MemberAddUpdated(builder *flatbuffers.Builder, updated int64) {
	builder.PrependInt64Slot(4, updated, 0)
}
func MemberEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Conversation struct {
	_tab flatbuffers.Table
}

func GetRootAsConversation(buf []byte, offset flatbuffers.UOffsetT) *Conversation {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Conversation{}
	x.Init(buf, n+offset)
	return x
}

func FinishConversationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConversation(buf []byte, offset flatbuffers.UOffsetT) *Conversation {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Conversation{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConversationBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Conversation) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Conversation) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Conversation) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Conversation) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Conversation) Desc() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Conversation) Mods(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *Conversation) ModsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Conversation) Created() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Conversation) MutateCreated(n int64) bool {
	return rcv._tab.MutateInt64Slot(12, n)
}

func (rcv *Conversation) Updated() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Conversation) MutateUpdated(n int64) bool {
	return rcv._tab.MutateInt64Slot(14, n)
}

func ConversationStart(builder *flatbuffers.Builder) {
	builder.StartObject(6)
}
func ConversationAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func ConversationAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(name), 0)
}
func ConversationAddDesc(builder *flatbuffers.Builder, desc flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(desc), 0)
}
func ConversationAddMods(builder *flatbuffers.Builder, mods flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(mods), 0)
}
func ConversationStartModsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ConversationAddCreated(builder *flatbuffers.Builder, created int64) {
	builder.PrependInt64Slot(4, created, 0)
}
func ConversationAddUpdated(builder *flatbuffers.Builder, updated int64) {
	builder.PrependInt64Slot(5, updated, 0)
}
func ConversationEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type Message struct {
	_tab flatbuffers.Table
}

func GetRootAsMessage(buf []byte, offset flatbuffers.UOffsetT) *Message {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Message{}
	x.Init(buf, n+offset)
	return x
}

func FinishMessageBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMessage(buf []byte, offset flatbuffers.UOffsetT) *Message {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Message{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMessageBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Message) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Message) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Message) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Message) Author() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Message) Content() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Message) Created() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Message) MutateCreated(n int64) bool {
	return rcv._tab.MutateInt64Slot(10, n)
}

func (rcv *Message) Updated() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Message) MutateUpdated(n int64) bool {
	return rcv._tab.MutateInt64Slot(12, n)
}

func MessageStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func MessageAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MessageAddAuthor(builder *flatbuffers.Builder, author flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(author), 0)
}
func MessageAddContent(builder *flatbuffers.Builder, content flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(content), 0)
}
func MessageAddCreated(builder *flatbuffers.Builder, created int64) {
	builder.PrependInt64Slot(3, created, 0)
}
func MessageAddUpdated(builder *flatbuffers.Builder, updated int64) {
	builder.PrependInt64Slot(4, updated, 0)
}
func MessageEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

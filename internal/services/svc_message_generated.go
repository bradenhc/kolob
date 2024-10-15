// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package services

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type MessageAddRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMessageAddRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageAddRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MessageAddRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMessageAddRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMessageAddRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageAddRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MessageAddRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMessageAddRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MessageAddRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MessageAddRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MessageAddRequest) Conversation() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MessageAddRequest) Author() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MessageAddRequest) Content() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MessageAddRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func MessageAddRequestAddConversation(builder *flatbuffers.Builder, conversation flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(conversation), 0)
}
func MessageAddRequestAddAuthor(builder *flatbuffers.Builder, author flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(author), 0)
}
func MessageAddRequestAddContent(builder *flatbuffers.Builder, content flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(content), 0)
}
func MessageAddRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MessageGetRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMessageGetRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageGetRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MessageGetRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMessageGetRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMessageGetRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageGetRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MessageGetRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMessageGetRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MessageGetRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MessageGetRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MessageGetRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MessageGetRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func MessageGetRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MessageGetRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MessageUpdateRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMessageUpdateRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageUpdateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MessageUpdateRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMessageUpdateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMessageUpdateRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageUpdateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MessageUpdateRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMessageUpdateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MessageUpdateRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MessageUpdateRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MessageUpdateRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MessageUpdateRequest) Content() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MessageUpdateRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func MessageUpdateRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MessageUpdateRequestAddContent(builder *flatbuffers.Builder, content flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(content), 0)
}
func MessageUpdateRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MessageRemoveRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMessageRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MessageRemoveRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMessageRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMessageRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MessageRemoveRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMessageRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MessageRemoveRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MessageRemoveRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MessageRemoveRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MessageRemoveRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func MessageRemoveRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MessageRemoveRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MessageListRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMessageListRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageListRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MessageListRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMessageListRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMessageListRequest(buf []byte, offset flatbuffers.UOffsetT) *MessageListRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MessageListRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMessageListRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MessageListRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MessageListRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MessageListRequest) Conversation() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MessageListRequest) Author() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MessageListRequest) CreatedAfter() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MessageListRequest) MutateCreatedAfter(n int64) bool {
	return rcv._tab.MutateInt64Slot(8, n)
}

func (rcv *MessageListRequest) CreatedBefore() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MessageListRequest) MutateCreatedBefore(n int64) bool {
	return rcv._tab.MutateInt64Slot(10, n)
}

func (rcv *MessageListRequest) Pattern() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MessageListRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func MessageListRequestAddConversation(builder *flatbuffers.Builder, conversation flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(conversation), 0)
}
func MessageListRequestAddAuthor(builder *flatbuffers.Builder, author flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(author), 0)
}
func MessageListRequestAddCreatedAfter(builder *flatbuffers.Builder, createdAfter int64) {
	builder.PrependInt64Slot(2, createdAfter, 0)
}
func MessageListRequestAddCreatedBefore(builder *flatbuffers.Builder, createdBefore int64) {
	builder.PrependInt64Slot(3, createdBefore, 0)
}
func MessageListRequestAddPattern(builder *flatbuffers.Builder, pattern flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(pattern), 0)
}
func MessageListRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package services

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ConversationAddRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsConversationAddRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationAddRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConversationAddRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishConversationAddRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConversationAddRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationAddRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConversationAddRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConversationAddRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConversationAddRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConversationAddRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConversationAddRequest) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ConversationAddRequest) Description() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ConversationAddRequest) Moderators(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *ConversationAddRequest) ModeratorsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ConversationAddRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ConversationAddRequestAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(name), 0)
}
func ConversationAddRequestAddDescription(builder *flatbuffers.Builder, description flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(description), 0)
}
func ConversationAddRequestAddModerators(builder *flatbuffers.Builder, moderators flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(moderators), 0)
}
func ConversationAddRequestStartModeratorsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ConversationAddRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type ConversationGetRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsConversationGetRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationGetRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConversationGetRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishConversationGetRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConversationGetRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationGetRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConversationGetRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConversationGetRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConversationGetRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConversationGetRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConversationGetRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ConversationGetRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ConversationGetRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func ConversationGetRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type ConversationUpdateRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsConversationUpdateRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationUpdateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConversationUpdateRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishConversationUpdateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConversationUpdateRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationUpdateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConversationUpdateRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConversationUpdateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConversationUpdateRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConversationUpdateRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConversationUpdateRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ConversationUpdateRequest) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ConversationUpdateRequest) Description() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ConversationUpdateRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ConversationUpdateRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func ConversationUpdateRequestAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(name), 0)
}
func ConversationUpdateRequestAddDescription(builder *flatbuffers.Builder, description flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(description), 0)
}
func ConversationUpdateRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type ConversationModsAddRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsConversationModsAddRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationModsAddRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConversationModsAddRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishConversationModsAddRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConversationModsAddRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationModsAddRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConversationModsAddRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConversationModsAddRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConversationModsAddRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConversationModsAddRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConversationModsAddRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ConversationModsAddRequest) Moderators(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *ConversationModsAddRequest) ModeratorsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ConversationModsAddRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func ConversationModsAddRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func ConversationModsAddRequestAddModerators(builder *flatbuffers.Builder, moderators flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(moderators), 0)
}
func ConversationModsAddRequestStartModeratorsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ConversationModsAddRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type ConversationModsRemoveRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsConversationModsRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationModsRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConversationModsRemoveRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishConversationModsRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConversationModsRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationModsRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConversationModsRemoveRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConversationModsRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConversationModsRemoveRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConversationModsRemoveRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConversationModsRemoveRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ConversationModsRemoveRequest) Moderators(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *ConversationModsRemoveRequest) ModeratorsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ConversationModsRemoveRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func ConversationModsRemoveRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func ConversationModsRemoveRequestAddModerators(builder *flatbuffers.Builder, moderators flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(moderators), 0)
}
func ConversationModsRemoveRequestStartModeratorsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func ConversationModsRemoveRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type ConversationRemoveRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsConversationRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ConversationRemoveRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishConversationRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsConversationRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *ConversationRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &ConversationRemoveRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedConversationRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *ConversationRemoveRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ConversationRemoveRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ConversationRemoveRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ConversationRemoveRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ConversationRemoveRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func ConversationRemoveRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
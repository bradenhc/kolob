// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package services

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type MemberCreateRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMemberCreateRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberCreateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MemberCreateRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMemberCreateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMemberCreateRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberCreateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MemberCreateRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMemberCreateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MemberCreateRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MemberCreateRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MemberCreateRequest) Username() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MemberCreateRequest) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MemberCreateRequest) Password() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MemberCreateRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func MemberCreateRequestAddUsername(builder *flatbuffers.Builder, username flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(username), 0)
}
func MemberCreateRequestAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(name), 0)
}
func MemberCreateRequestAddPassword(builder *flatbuffers.Builder, password flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(password), 0)
}
func MemberCreateRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MemberAuthenticateRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMemberAuthenticateRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberAuthenticateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MemberAuthenticateRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMemberAuthenticateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMemberAuthenticateRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberAuthenticateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MemberAuthenticateRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMemberAuthenticateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MemberAuthenticateRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MemberAuthenticateRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MemberAuthenticateRequest) Username() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MemberAuthenticateRequest) Password() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MemberAuthenticateRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func MemberAuthenticateRequestAddUsername(builder *flatbuffers.Builder, username flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(username), 0)
}
func MemberAuthenticateRequestAddPassword(builder *flatbuffers.Builder, password flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(password), 0)
}
func MemberAuthenticateRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MemberChangePasswordRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMemberChangePasswordRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberChangePasswordRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MemberChangePasswordRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMemberChangePasswordRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMemberChangePasswordRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberChangePasswordRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MemberChangePasswordRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMemberChangePasswordRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MemberChangePasswordRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MemberChangePasswordRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MemberChangePasswordRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MemberChangePasswordRequest) OldPassword() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MemberChangePasswordRequest) NewPassword() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MemberChangePasswordRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func MemberChangePasswordRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MemberChangePasswordRequestAddOldPassword(builder *flatbuffers.Builder, oldPassword flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(oldPassword), 0)
}
func MemberChangePasswordRequestAddNewPassword(builder *flatbuffers.Builder, newPassword flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(newPassword), 0)
}
func MemberChangePasswordRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MemberUpdateRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMemberUpdateRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberUpdateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MemberUpdateRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMemberUpdateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMemberUpdateRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberUpdateRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MemberUpdateRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMemberUpdateRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MemberUpdateRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MemberUpdateRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MemberUpdateRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MemberUpdateRequest) Username() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MemberUpdateRequest) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MemberUpdateRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func MemberUpdateRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MemberUpdateRequestAddUsername(builder *flatbuffers.Builder, username flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(username), 0)
}
func MemberUpdateRequestAddName(builder *flatbuffers.Builder, name flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(name), 0)
}
func MemberUpdateRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MemberRemoveRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMemberRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MemberRemoveRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMemberRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMemberRemoveRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberRemoveRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MemberRemoveRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMemberRemoveRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MemberRemoveRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MemberRemoveRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MemberRemoveRequest) Id() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MemberRemoveRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func MemberRemoveRequestAddId(builder *flatbuffers.Builder, id flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(id), 0)
}
func MemberRemoveRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
type MemberFindByUsernameRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsMemberFindByUsernameRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberFindByUsernameRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MemberFindByUsernameRequest{}
	x.Init(buf, n+offset)
	return x
}

func FinishMemberFindByUsernameRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsMemberFindByUsernameRequest(buf []byte, offset flatbuffers.UOffsetT) *MemberFindByUsernameRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MemberFindByUsernameRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedMemberFindByUsernameRequestBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *MemberFindByUsernameRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MemberFindByUsernameRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MemberFindByUsernameRequest) Username() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func MemberFindByUsernameRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func MemberFindByUsernameRequestAddUsername(builder *flatbuffers.Builder, username flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(username), 0)
}
func MemberFindByUsernameRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}

// Code generated by protoc-gen-go.
// source: game.proto
// DO NOT EDIT!

/*
Package game_msg is a generated protocol buffer package.

It is generated from these files:
	game.proto

It has these top-level messages:
	LoginMsg
	LoginMsgRes
	CreateRoleMsg
	CreateRoleMsgRes
	LoadUserMsg
	LoadUserMsgRes
	LoginMsgEnd
	ChatMsg
	ChatMsgRes
	NotifyChatMsgRes
	UpdateEventMsg
	UpdateEventMsgRes
	NotifyItemMsg
	NotifyItemMsgRes
	SysInformationMsg
	SysInformationMsgRes
	SysNewDayMsg
	SysNewDayMsgRes
*/
package game_msg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import "common_msg"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type LoginMsg struct {
	Platform         *uint32 `protobuf:"varint,1,opt,name=platform" json:"platform,omitempty"`
	Account          *string `protobuf:"bytes,2,opt,name=account" json:"account,omitempty"`
	Password         *string `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *LoginMsg) Reset()         { *m = LoginMsg{} }
func (m *LoginMsg) String() string { return proto.CompactTextString(m) }
func (*LoginMsg) ProtoMessage()    {}

func (m *LoginMsg) GetPlatform() uint32 {
	if m != nil && m.Platform != nil {
		return *m.Platform
	}
	return 0
}

func (m *LoginMsg) GetAccount() string {
	if m != nil && m.Account != nil {
		return *m.Account
	}
	return ""
}

func (m *LoginMsg) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

type LoginMsgRes struct {
	HasRole          *uint32 `protobuf:"varint,1,opt,name=has_role" json:"has_role,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *LoginMsgRes) Reset()         { *m = LoginMsgRes{} }
func (m *LoginMsgRes) String() string { return proto.CompactTextString(m) }
func (*LoginMsgRes) ProtoMessage()    {}

func (m *LoginMsgRes) GetHasRole() uint32 {
	if m != nil && m.HasRole != nil {
		return *m.HasRole
	}
	return 0
}

type CreateRoleMsg struct {
	Nick             *string `protobuf:"bytes,1,opt,name=nick" json:"nick,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CreateRoleMsg) Reset()         { *m = CreateRoleMsg{} }
func (m *CreateRoleMsg) String() string { return proto.CompactTextString(m) }
func (*CreateRoleMsg) ProtoMessage()    {}

func (m *CreateRoleMsg) GetNick() string {
	if m != nil && m.Nick != nil {
		return *m.Nick
	}
	return ""
}

type CreateRoleMsgRes struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *CreateRoleMsgRes) Reset()         { *m = CreateRoleMsgRes{} }
func (m *CreateRoleMsgRes) String() string { return proto.CompactTextString(m) }
func (*CreateRoleMsgRes) ProtoMessage()    {}

type LoadUserMsg struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *LoadUserMsg) Reset()         { *m = LoadUserMsg{} }
func (m *LoadUserMsg) String() string { return proto.CompactTextString(m) }
func (*LoadUserMsg) ProtoMessage()    {}

type LoadUserMsgRes struct {
	User             *common_msg.UserT    `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	Events           []*common_msg.EventT `protobuf:"bytes,2,rep,name=events" json:"events,omitempty"`
	Items            []*common_msg.ItemT  `protobuf:"bytes,6,rep,name=items" json:"items,omitempty"`
	Head             *common_msg.HeadT    `protobuf:"bytes,24,opt,name=head" json:"head,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *LoadUserMsgRes) Reset()         { *m = LoadUserMsgRes{} }
func (m *LoadUserMsgRes) String() string { return proto.CompactTextString(m) }
func (*LoadUserMsgRes) ProtoMessage()    {}

func (m *LoadUserMsgRes) GetUser() *common_msg.UserT {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *LoadUserMsgRes) GetEvents() []*common_msg.EventT {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *LoadUserMsgRes) GetItems() []*common_msg.ItemT {
	if m != nil {
		return m.Items
	}
	return nil
}

func (m *LoadUserMsgRes) GetHead() *common_msg.HeadT {
	if m != nil {
		return m.Head
	}
	return nil
}

type LoginMsgEnd struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *LoginMsgEnd) Reset()         { *m = LoginMsgEnd{} }
func (m *LoginMsgEnd) String() string { return proto.CompactTextString(m) }
func (*LoginMsgEnd) ProtoMessage()    {}

// //////////////////////////////////////////////
// 聊天
// //////////////////////////////////////////////
type ChatMsg struct {
	Uid              *uint32 `protobuf:"varint,2,opt,name=uid" json:"uid,omitempty"`
	Msg              *string `protobuf:"bytes,3,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ChatMsg) Reset()         { *m = ChatMsg{} }
func (m *ChatMsg) String() string { return proto.CompactTextString(m) }
func (*ChatMsg) ProtoMessage()    {}

func (m *ChatMsg) GetUid() uint32 {
	if m != nil && m.Uid != nil {
		return *m.Uid
	}
	return 0
}

func (m *ChatMsg) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

type ChatMsgRes struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *ChatMsgRes) Reset()         { *m = ChatMsgRes{} }
func (m *ChatMsgRes) String() string { return proto.CompactTextString(m) }
func (*ChatMsgRes) ProtoMessage()    {}

type NotifyChatMsgRes struct {
	Uid              *uint32 `protobuf:"varint,2,opt,name=uid" json:"uid,omitempty"`
	Msg              *string `protobuf:"bytes,4,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *NotifyChatMsgRes) Reset()         { *m = NotifyChatMsgRes{} }
func (m *NotifyChatMsgRes) String() string { return proto.CompactTextString(m) }
func (*NotifyChatMsgRes) ProtoMessage()    {}

func (m *NotifyChatMsgRes) GetUid() uint32 {
	if m != nil && m.Uid != nil {
		return *m.Uid
	}
	return 0
}

func (m *NotifyChatMsgRes) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

// //////////////////////////////////////////////
// 事件
// //////////////////////////////////////////////
type UpdateEventMsg struct {
	Event            *common_msg.EventT `protobuf:"bytes,1,opt,name=event" json:"event,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *UpdateEventMsg) Reset()         { *m = UpdateEventMsg{} }
func (m *UpdateEventMsg) String() string { return proto.CompactTextString(m) }
func (*UpdateEventMsg) ProtoMessage()    {}

func (m *UpdateEventMsg) GetEvent() *common_msg.EventT {
	if m != nil {
		return m.Event
	}
	return nil
}

type UpdateEventMsgRes struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *UpdateEventMsgRes) Reset()         { *m = UpdateEventMsgRes{} }
func (m *UpdateEventMsgRes) String() string { return proto.CompactTextString(m) }
func (*UpdateEventMsgRes) ProtoMessage()    {}

// //////////////////////////////////////////////
// 道具
// //////////////////////////////////////////////
type NotifyItemMsg struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *NotifyItemMsg) Reset()         { *m = NotifyItemMsg{} }
func (m *NotifyItemMsg) String() string { return proto.CompactTextString(m) }
func (*NotifyItemMsg) ProtoMessage()    {}

type NotifyItemMsgRes struct {
	Items            *common_msg.ItemT `protobuf:"bytes,1,opt,name=items" json:"items,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *NotifyItemMsgRes) Reset()         { *m = NotifyItemMsgRes{} }
func (m *NotifyItemMsgRes) String() string { return proto.CompactTextString(m) }
func (*NotifyItemMsgRes) ProtoMessage()    {}

func (m *NotifyItemMsgRes) GetItems() *common_msg.ItemT {
	if m != nil {
		return m.Items
	}
	return nil
}

// //////////////////////////////////////////////
// 系统
// //////////////////////////////////////////////
type SysInformationMsg struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *SysInformationMsg) Reset()         { *m = SysInformationMsg{} }
func (m *SysInformationMsg) String() string { return proto.CompactTextString(m) }
func (*SysInformationMsg) ProtoMessage()    {}

type SysInformationMsgRes struct {
	TimeSecond       *uint32 `protobuf:"varint,1,opt,name=time_second" json:"time_second,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SysInformationMsgRes) Reset()         { *m = SysInformationMsgRes{} }
func (m *SysInformationMsgRes) String() string { return proto.CompactTextString(m) }
func (*SysInformationMsgRes) ProtoMessage()    {}

func (m *SysInformationMsgRes) GetTimeSecond() uint32 {
	if m != nil && m.TimeSecond != nil {
		return *m.TimeSecond
	}
	return 0
}

type SysNewDayMsg struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *SysNewDayMsg) Reset()         { *m = SysNewDayMsg{} }
func (m *SysNewDayMsg) String() string { return proto.CompactTextString(m) }
func (*SysNewDayMsg) ProtoMessage()    {}

type SysNewDayMsgRes struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *SysNewDayMsgRes) Reset()         { *m = SysNewDayMsgRes{} }
func (m *SysNewDayMsgRes) String() string { return proto.CompactTextString(m) }
func (*SysNewDayMsgRes) ProtoMessage()    {}

func init() {
	proto.RegisterType((*LoginMsg)(nil), "game_msg.login_msg")
	proto.RegisterType((*LoginMsgRes)(nil), "game_msg.login_msg_res")
	proto.RegisterType((*CreateRoleMsg)(nil), "game_msg.create_role_msg")
	proto.RegisterType((*CreateRoleMsgRes)(nil), "game_msg.create_role_msg_res")
	proto.RegisterType((*LoadUserMsg)(nil), "game_msg.load_user_msg")
	proto.RegisterType((*LoadUserMsgRes)(nil), "game_msg.load_user_msg_res")
	proto.RegisterType((*LoginMsgEnd)(nil), "game_msg.login_msg_end")
	proto.RegisterType((*ChatMsg)(nil), "game_msg.chat_msg")
	proto.RegisterType((*ChatMsgRes)(nil), "game_msg.chat_msg_res")
	proto.RegisterType((*NotifyChatMsgRes)(nil), "game_msg.notify_chat_msg_res")
	proto.RegisterType((*UpdateEventMsg)(nil), "game_msg.update_event_msg")
	proto.RegisterType((*UpdateEventMsgRes)(nil), "game_msg.update_event_msg_res")
	proto.RegisterType((*NotifyItemMsg)(nil), "game_msg.notify_item_msg")
	proto.RegisterType((*NotifyItemMsgRes)(nil), "game_msg.notify_item_msg_res")
	proto.RegisterType((*SysInformationMsg)(nil), "game_msg.sys_information_msg")
	proto.RegisterType((*SysInformationMsgRes)(nil), "game_msg.sys_information_msg_res")
	proto.RegisterType((*SysNewDayMsg)(nil), "game_msg.sys_new_day_msg")
	proto.RegisterType((*SysNewDayMsgRes)(nil), "game_msg.sys_new_day_msg_res")
}

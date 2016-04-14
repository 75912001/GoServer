// icartravel project main.go
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"ict_user"
	"pb_square"
	"zzcommon"
)

////////////////////////////////////////////////////////////////////////////////
//protobuf func
type PbFunHandle struct {
	pbFun        func(user *ict_user.User_t, protoMessage proto.Message) (ret int)
	protoMessage proto.Message
}

var pbFunMap map[zzcommon.MESSAGE_ID]PbFunHandle

func initPbFun() (ret int) {
	///////////////////////////////////////////////////////////////////
	//pb message
	pbFunMap = make(map[zzcommon.MESSAGE_ID]PbFunHandle)
	{
		var cmd_id zzcommon.MESSAGE_ID = 0x100101
		var pbFunHandle PbFunHandle
		pbFunHandle.pbFun = OnLoginMsg
		pbFunHandle.protoMessage = new(pb_square.LoginMsg)
		pbFunMap[cmd_id] = pbFunHandle
	}

	//注册新的消息
	return 0
}

func onRecv(peerConn *zzcommon.PeerConn_t) (ret int) {
	PacketLength := peerConn.RecvProtoHead.PacketLength //总包长度
	MessageId := peerConn.RecvProtoHead.MessageId       //消息号
	SessionId := peerConn.RecvProtoHead.SessionId       //会话id
	UserId := peerConn.RecvProtoHead.UserId             //用户id
	ResultId := peerConn.RecvProtoHead.ResultId         //结果id
	fmt.Println(PacketLength, MessageId, SessionId, UserId, ResultId)

	pbFunHandle, ok := pbFunMap[MessageId]
	if !ok {
		fmt.Println("######pbFunMap[MessageId]", MessageId)
		return zzcommon.ERROR_DISCONNECT_PEER
	}
	user, ok := ict_user.GuserMgr.UserMap[peerConn]
	if !ok {
		fmt.Println("######ict_user.GuserMgr.UserMap[peerConn]", peerConn)
		return zzcommon.ERROR_DISCONNECT_PEER
	}

	err := proto.Unmarshal(peerConn.RecvBuf[20:PacketLength], pbFunHandle.protoMessage)
	if nil != err {
		fmt.Println("######proto.Unmarshal", MessageId)
		return zzcommon.ERROR_DISCONNECT_PEER
	}
	ret = pbFunHandle.pbFun(user, pbFunHandle.protoMessage)

	return ret
}

package zzcommon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"runtime"
	"strconv"
)

//////////////////////////////////////////////////////////////////////////////
//消息包头
type PACKET_LENGTH uint32
type MESSAGE_ID uint32
type SESSION_ID uint32
type USER_ID uint32
type RESULT_ID uint32

//消息包头
type ProtoHead struct {
	PacketLength PACKET_LENGTH //总包长度
	MessageId    MESSAGE_ID    //消息号
	SessionId    SESSION_ID    //会话id
	UserId       USER_ID       //用户id
	ResultId     RESULT_ID     //结果id
}

const (
	//消息包头长度
	ProtoHeadLength uint32 = 20
)

//////////////////////////////////////////////////////////////////////////////
//对端连接信息
type PeerConn struct {
	Conn    *net.TCPConn //连接
	RecvBuf []byte
}

func (p *PeerConn) Send(messageId MESSAGE_ID, req proto.Message, sessionId SESSION_ID, userId USER_ID, resultId RESULT_ID) (err error) {
	reqBuf, err := proto.Marshal(req)
	if nil != err {
		fmt.Printf("######proto.Marshal err:", err)
		return err
	}
	var reqBufLen = uint32(len(reqBuf))
	var sendBufAllLength uint32
	sendBufAllLength = reqBufLen + ProtoHeadLength

	headBuf := new(bytes.Buffer)
	var data = []interface{}{
		sendBufAllLength,
		messageId,
		sessionId,
		userId,
		resultId,
	}
	for _, v := range data {
		err := binary.Write(headBuf, binary.LittleEndian, v)
		if nil != err {
			fmt.Println("######binary.Write failed:", err)
		}
	}

	//[优化]使用一个发送
	_, err = p.Conn.Write(headBuf.Bytes())
	_, err = p.Conn.Write(reqBuf)
	if nil != err {
		fmt.Printf("######PeerConn.Conn.Write err:", err)
		return err
	}
	fmt.Println("Send body len:", reqBufLen, sendBufAllLength)
	headBuf = nil
	return err
}

//////////////////////////////////////////////////////////////////////////////
const (
	//成功
	SUCC int = 0
	//错误
	ERROR int = -1
	//断开对方的连接
	ERROR_DISCONNECT_PEER int = -2
)

//////////////////////////////////////////////////////////////////////////////
//字符串转
func StringToUint32(s string) (value uint32) {
	vaule, err := strconv.ParseUint(s, 10, 32)
	if nil != err {
		return 0
	}
	return uint32(vaule)
}

func StringToInt(s string) (value int) {
	vaule, err := strconv.ParseInt(s, 10, 0)
	if nil != err {
		return 0
	}
	return int(vaule)
}

func StringToUint16(s string) (value uint16) {
	vaule, err := strconv.ParseUint(s, 10, 16)
	if nil != err {
		return 0
	}
	return uint16(vaule)
}

func StringSubstr(str string, length int) string {
	rs := []rune(str)
	return string(rs[0:length])
}

//////////////////////////////////////////////////////////////////////////////
func IsWindows() bool {
	return `windows` == runtime.GOOS
}

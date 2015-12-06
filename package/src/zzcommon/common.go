package zzcommon

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"runtime"
	"strconv"
)

//////////////////////////////////////////////////////////////////////////////
const (
	//成功
	SUCC int = 0
	//错误
	ERROR int = -1
	//断开对方的连接
	ERROR_DISCONNECT_PEER int = -2

	ERROR_SMS_SENDING       int = 10000 //短信已发出,请收到后重试
	ERROR_SMS_REGISTER_CODE int = 10001 // 短信注册码失败,请重新请求短信注册
	ERROR_USER_EXIST        int = 10002 //用户已存在
	ERROR_PHONE_NUM_BIND    int = 10003 //手机号已绑定
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
			return err
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

func StringSubstr(s string, length int) (value string) {
	r := []rune(s)
	return string(r[0:length])
}

////////////////////////////////////////////////////////////////////////////////
//md5
func GenMd5(s string) (value string) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	cipherStr := md5Ctx.Sum(nil)
	value = hex.EncodeToString(cipherStr)
	return value
}

//////////////////////////////////////////////////////////////////////////////
func IsWindows() bool {
	return `windows` == runtime.GOOS
}

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

	ERROR_SYS               int = 1     //系统错误
	ERROR_PARAM             int = 2     //参数错误
	ERROR_SMS_SENDING       int = 10000 //短信已发出,请收到后重试
	ERROR_SMS_REGISTER_CODE int = 10001 // 短信注册码失败,请重新请求短信注册
	ERROR_USER_EXIST        int = 10002 //用户已存在
	ERROR_PHONE_NUM_BIND    int = 10003 //手机号已绑定
	ERROR_PHONE_NUM_NO_BIND int = 10004 //手机号未绑定
)

//////////////////////////////////////////////////////////////////////////////
//消息包头
type PACKET_LENGTH uint32
type MESSAGE_ID uint32
type SESSION_ID uint32
type USER_ID uint32
type RESULT_ID uint32

//消息包头
type ProtoHead_t struct {
	PacketLength PACKET_LENGTH //总包长度
	MessageId    MESSAGE_ID    //消息号
	SessionId    SESSION_ID    //会话id
	UserId       USER_ID       //用户id
	ResultId     RESULT_ID     //结果id
}

const (
	//消息包头长度
	//	ProtoHeadPacketLength uint32 = 4
	ProtoHeadLength uint32 = 20
)

//////////////////////////////////////////////////////////////////////////////
//对端连接信息
type PeerConn_t struct {
	Conn          *net.TCPConn //连接
	RecvBuf       []byte
	RecvProtoHead ProtoHead_t
}

//解析协议包头长度
func (p *PeerConn_t) ParseProtoHeadPacketLength() {
	buf1 := bytes.NewBuffer(p.RecvBuf[0:4])

	binary.Read(buf1, binary.LittleEndian, &p.RecvProtoHead.PacketLength)
}

//解析协议包头
func (p *PeerConn_t) ParseProtoHead() {
	buf1 := bytes.NewBuffer(p.RecvBuf[0:4])
	buf2 := bytes.NewBuffer(p.RecvBuf[4:8])
	buf3 := bytes.NewBuffer(p.RecvBuf[8:12])
	buf4 := bytes.NewBuffer(p.RecvBuf[12:16])
	buf5 := bytes.NewBuffer(p.RecvBuf[16:ProtoHeadLength])

	binary.Read(buf1, binary.LittleEndian, &p.RecvProtoHead.PacketLength)
	binary.Read(buf2, binary.LittleEndian, &p.RecvProtoHead.MessageId)
	binary.Read(buf3, binary.LittleEndian, &p.RecvProtoHead.SessionId)
	binary.Read(buf4, binary.LittleEndian, &p.RecvProtoHead.UserId)
	binary.Read(buf5, binary.LittleEndian, &p.RecvProtoHead.ResultId)
}

func (p *PeerConn_t) Send(req proto.Message, messageId MESSAGE_ID, sessionId SESSION_ID, userId USER_ID, resultId RESULT_ID) (err error) {
	reqBuf, err := proto.Marshal(req)
	if nil != err {
		fmt.Printf("######proto.Marshal err:", err)
		return err
	}
	var reqBufLen = uint32(len(reqBuf))
	var sendBufAllLength uint32 = reqBufLen + ProtoHeadLength

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

////////////////////////////////////////////////////////////////////////////////
func IsWindows() bool {
	return `windows` == runtime.GOOS
}

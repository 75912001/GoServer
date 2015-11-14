package zzcli

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"zzcommon"
	"zzser"
)

//服务端连接建立
type ON_SER_CONN func(peerConn *zzser.PeerConn) int

//服务端连接关闭
type ON_SER_CONN_CLOSED func(peerConn *zzser.PeerConn) int

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
type ON_SER_GET_PACKET_LEN func(peerConn *zzser.PeerConn, packetLength int) int

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
type ON_SER_PACKET func(peerConn *zzser.PeerConn, packetLength int) int

//己方作为客户端
type Client struct {
	OnSerConn         ON_SER_CONN
	OnSerConnClosed   ON_SER_CONN_CLOSED
	OnSerGetPacketLen ON_SER_GET_PACKET_LEN
	OnSerPacket       ON_SER_PACKET

	Conn *net.TCPConn
}

//连接
func (p *Client) Connect(ip string, port uint16, recvBufMax int) (err error) {

	var addr = ip + ":" + strconv.Itoa(int(port))
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if nil != err {
		fmt.Println("######net.ResolveTCPAddr err:", err, addr)
		return err
	}
	p.Conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		fmt.Println("######net.Dial err:", err, addr)
		return err
	}
	go p.recv(recvBufMax)
	return err
}

//todo
func (p *Client) Send(messageId uint32, req proto.Message) (ret int, err error) {
	reqBuf, err := proto.Marshal(req)
	if nil != err {
		fmt.Printf("######proto.Marshal err:", err)
		return ret, err
	}
	var reqBufLen = len(reqBuf)
	var sendBufAllLength uint32
	sendBufAllLength = uint32(reqBufLen + 20)

	head_buf := new(bytes.Buffer)
	var data = []interface{}{
		sendBufAllLength,
		messageId,
	}
	for _, v := range data {
		err := binary.Write(head_buf, binary.LittleEndian, v)
		if nil != err {
			fmt.Println("binary.Write failed:", err)
		}
	}

	//todo
	ret, err = p.Conn.Write(head_buf.Bytes())
	ret, err = p.Conn.Write(reqBuf)
	if nil != err {
		fmt.Printf("######user.Conn.Write err:", err)
		return ret, err
	}
	fmt.Println("Send body len:", reqBufLen)
	return ret, err
}

func (p *Client) recv(recvBufMax int) {
	var peerConn zzser.PeerConn
	peerConn.Conn = p.Conn
	p.OnSerConn(&peerConn)
	defer p.Conn.Close()

	defer p.OnSerConnClosed(&peerConn)

	//优化[消耗内存过大]
	peerConn.RecvBuf = make([]byte, recvBufMax)

	var readIndex int

	for {
		readNum, err := p.Conn.Read(peerConn.RecvBuf[readIndex:])
		if nil != err {
			fmt.Println("######Conn.Read err:", readNum, err)
			break
		}

		readIndex += readNum
		packetLength := p.OnSerGetPacketLen(&peerConn, readIndex)
		if zzcommon.ERROR_DISCONNECT_PEER == packetLength {
			fmt.Println("######OnSerGetPacketLen:", zzcommon.ERROR_DISCONNECT_PEER)
			break
		}
		if packetLength > 0 { //有完整的包
			ret := p.OnSerPacket(&peerConn, packetLength)
			if zzcommon.ERROR_DISCONNECT_PEER == ret {
				fmt.Println("######OnSerPacket:", zzcommon.ERROR_DISCONNECT_PEER)
				break
			}
			peerConn.RecvBuf = peerConn.RecvBuf[readIndex-packetLength : readIndex]
			readIndex = readIndex - packetLength
		}
	}
	peerConn.RecvBuf = nil
}

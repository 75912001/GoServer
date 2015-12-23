package zztcp

import (
	"fmt"
	"net"
	"strconv"
	"zzcommon"
)

//服务端连接建立
type ON_SER_CONN func(peerConn *zzcommon.PeerConn_t) int

//服务端连接关闭
type ON_SER_CONN_CLOSED func(peerConn *zzcommon.PeerConn_t) int

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
type ON_SER_GET_PACKET_LEN func(peerConn *zzcommon.PeerConn_t, packetLength int) int

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
type ON_SER_PACKET func(peerConn *zzcommon.PeerConn_t, packetLength int) int

//己方作为客户端
type Client_t struct {
	OnSerConn         ON_SER_CONN
	OnSerConnClosed   ON_SER_CONN_CLOSED
	OnSerGetPacketLen ON_SER_GET_PACKET_LEN
	OnSerPacket       ON_SER_PACKET
	PeerConn          zzcommon.PeerConn_t
}

//连接
func (p *Client_t) Connect(ip string, port uint16, recvBufMax int) (err error) {
	var addr = ip + ":" + strconv.Itoa(int(port))
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if nil != err {
		fmt.Println("######net.ResolveTCPAddr err:", err, addr)
		return err
	}
	p.PeerConn.Conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		fmt.Println("######net.Dial err:", err, addr)
		return err
	}
	go p.recv(recvBufMax)
	return err
}

func (p *Client_t) recv(recvBufMax int) {
	var peerConn zzcommon.PeerConn_t
	peerConn.Conn = p.PeerConn.Conn
	p.OnSerConn(&peerConn)
	defer p.PeerConn.Conn.Close()

	defer p.OnSerConnClosed(&peerConn)

	//优化[消耗内存过大]
	peerConn.RecvBuf = make([]byte, recvBufMax)

	var readIndex int

	for {
		readNum, err := p.PeerConn.Conn.Read(peerConn.RecvBuf[readIndex:])
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
	peerConn.Conn = nil
}

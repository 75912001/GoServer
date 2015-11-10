package zzcli

import (
	"fmt"
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
}

//连接
func (p *Client) Connect(ip string, port uint16) (conn *net.TCPConn, err error) {

	var addr = ip + ":" + strconv.Itoa(int(port))
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if nil != err {
		fmt.Println("######net.ResolveTCPAddr err:", err, addr)
		return conn, err
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		fmt.Println("######net.Dial err:", err, addr)
		return conn, err
	}
	return conn, err
}

func (p *Client) ClientRecv(conn *net.TCPConn, recvBufMax int) {
	var peerConn zzser.PeerConn
	peerConn.Conn = conn
	p.OnSerConn(&peerConn)
	defer conn.Close()

	defer p.OnSerConnClosed(&peerConn)

	//优化[消耗内存过大]
	recvBuf := make([]byte, recvBufMax)

	var readIndex int

	for {
		readNum, err := conn.Read(recvBuf[readIndex:])
		if nil != err {
			fmt.Println("######conn.Read err:", readNum, err)
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
			recvBuf = recvBuf[readIndex-packetLength : readIndex]
			readIndex = readIndex - packetLength
		}
	}
	recvBuf = nil
}

package zztcp

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"zzcommon"
)

//初始化服务器
type ON_INIT func() int

//服务器结束
type ON_FINI func() int

//客户端连上
type ON_CLI_CONN func(peerConn *zzcommon.PeerConn_t) int

//客户端连接关闭
type ON_CLI_CONN_CLOSED func(peerConn *zzcommon.PeerConn_t) int

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,客户端断开
type ON_CLI_GET_PACKET_LEN func(peerConn *zzcommon.PeerConn_t, packetLength int) int

//客户端消息
//返回ERROR_DISCONNECT_PEER断开客户端
type ON_CLI_PACKET func(peerConn *zzcommon.PeerConn_t, packetLength int) int

//己方作为服务
type Server_t struct {
	IsRun                bool //是否运行
	OnInit               ON_INIT
	OnFini               ON_FINI
	OnCliConnClosed      ON_CLI_CONN_CLOSED
	OnCliConn            ON_CLI_CONN
	OnCliGetPacketLength ON_CLI_GET_PACKET_LEN
	OnCliPacket          ON_CLI_PACKET
	PacketLengthMax      uint32
}

//运行
func (p *Server_t) Run(ip string, port uint16, delay bool) (err error) {
	p.IsRun = true
	var addr = ip + ":" + strconv.Itoa(int(port))
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if nil != err {
		fmt.Println("######net.ResolveTCPAddr err:", err)
		return err
	}
	//优化[设置地址复用]
	//优化[设置监听的缓冲数量]
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		fmt.Println("######net.Listen err:", err)
		return err
	}
	defer listen.Close()

	p.OnInit()
	defer p.OnFini()

	go p.handleAccept(listen, delay)

	for p.IsRun {
		time.Sleep(60 * time.Second)
		fmt.Println("server runing...", time.Now())
	}

	fmt.Println("######server done...", time.Now())

	return err
}

func (p *Server_t) handleAccept(listen *net.TCPListener, delay bool) {
	for {
		conn, err := listen.AcceptTCP()
		if nil != err {
			fmt.Println("######listen.Accept err:", err)
			p.IsRun = false
			return
		}

		conn.SetNoDelay(delay)
		conn.SetReadBuffer((int)(p.PacketLengthMax))
		conn.SetWriteBuffer((int)(p.PacketLengthMax))
		go p.handleConnection(conn)
	}
}

func (p *Server_t) handleConnection(conn *net.TCPConn) {
	var peerIp = conn.RemoteAddr().String()
	fmt.Println("connection from:", peerIp)

	var peerConn zzcommon.PeerConn_t
	peerConn.Conn = conn

	defer peerConn.Conn.Close()

	p.OnCliConn(&peerConn)
	defer p.OnCliConnClosed(&peerConn)

	//优化[消耗内存过大]
	peerConn.RecvBuf = make([]byte, p.PacketLengthMax)

	var readIndex int
	for {
		readNum, err := peerConn.Conn.Read(peerConn.RecvBuf[readIndex:])
		if nil != err {
			fmt.Println("######peerConn.Conn.Read err:", readNum, err)
			break
		}

		readIndex = readIndex + readNum
		packetLength := p.OnCliGetPacketLength(&peerConn, readIndex)
		if packetLength > 0 { //有完整的包
			ret := p.OnCliPacket(&peerConn, packetLength)
			if zzcommon.ERROR_DISCONNECT_PEER == ret {
				break
			}
			copy(peerConn.RecvBuf, peerConn.RecvBuf[packetLength:readIndex])
			readIndex -= packetLength
		} else if zzcommon.ERROR_DISCONNECT_PEER == packetLength {
			break
		}
	}
}

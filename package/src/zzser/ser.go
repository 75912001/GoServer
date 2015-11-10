package zzser

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"
	"zzcommon"
	"zzini"
)

//初始化服务器
type ON_INIT func() int

//服务器结束
type ON_FINI func() int

//客户端连上
type ON_CLI_CONN func(peerConn *PeerConn) int

//客户端连接关闭
type ON_CLI_CONN_CLOSED func(peerConn *PeerConn) int

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,客户端断开
type ON_CLI_GET_PACKET_LEN func(peerConn *PeerConn, packetLength int) int

//客户端消息
//返回ERROR_DISCONNECT_PEER断开客户端
type ON_CLI_PACKET func(peerConn *PeerConn, packetLength int) int

//己方作为服务
type Server struct {
	FileIni              zzini.ZZIni //ini配置文件
	Ip                   string
	Port                 uint16
	PacketLengthMax      int  //设置包最大
	Delay                bool //tcp中是否延迟,默认为true
	GoProcessMax         int  //并行执行的数量
	IsRun                bool //是否运行
	OnInit               ON_INIT
	OnFini               ON_FINI
	OnCliConnClosed      ON_CLI_CONN_CLOSED
	OnCliConn            ON_CLI_CONN
	OnCliGetPacketLength ON_CLI_GET_PACKET_LEN
	OnCliPacket          ON_CLI_PACKET
}

//对端连接信息
type PeerConn struct {
	Conn    *net.TCPConn //连接
	recvBuf []byte
}

//加载配置文件
func (p *Server) LoadConfig() (err error) {
	err = p.FileIni.Load()
	if nil != err {
		return err
	}

	p.Ip = p.FileIni.Get("server", "ip", "")
	p.Port = zzcommon.StringToUint16(p.FileIni.Get("server", "port", "0"))
	p.PacketLengthMax = zzcommon.StringToInt(p.FileIni.Get("common", "packet_length_max", "81920"))
	str_num_cpu := strconv.Itoa(runtime.NumCPU())
	p.GoProcessMax = zzcommon.StringToInt(p.FileIni.Get("common", "go_process_max", str_num_cpu))
	runtime.GOMAXPROCS(p.GoProcessMax)
	return err
}

//运行
func (p *Server) Run() (err error) {
	p.IsRun = true
	var addr = p.Ip + ":" + strconv.Itoa(int(p.Port))
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

	go handleAccept(listen, p)

	for p.IsRun {
		time.Sleep(1000000000)
		fmt.Println("runing...")
	}
	return err
}

func handleAccept(listen *net.TCPListener, server *Server) {
	for {
		conn, err := listen.AcceptTCP()
		if nil != err {
			fmt.Println("######listen.Accept err:", err)
			return
		}

		conn.SetNoDelay(!server.Delay)
		conn.SetReadBuffer(server.PacketLengthMax)
		conn.SetWriteBuffer(server.PacketLengthMax)
		go handleConnection(server, conn)
	}
}

func handleConnection(server *Server, conn *net.TCPConn) {
	var peerIp = conn.RemoteAddr().String()
	fmt.Println("connection from:", peerIp)

	var peerConn PeerConn
	peerConn.Conn = conn

	defer conn.Close()

	server.OnCliConn(&peerConn)
	defer server.OnCliConnClosed(&peerConn)

	//优化[消耗内存过大]
	peerConn.recvBuf = make([]byte, server.PacketLengthMax)

	var readIndex int
	for {
		readNum, err := conn.Read(peerConn.recvBuf[readIndex:])
		if nil != err {
			fmt.Println("######conn.Read err:", readNum, err)
			break
		}

		readIndex = readIndex + readNum
		packetLength := server.OnCliGetPacketLength(&peerConn, readIndex)
		if packetLength > 0 { //有完整的包
			ret := server.OnCliPacket(&peerConn, packetLength)
			if zzcommon.ERROR_DISCONNECT_PEER == ret {
				break
			}
			copy(peerConn.recvBuf, peerConn.recvBuf[packetLength:readIndex])
			readIndex -= packetLength
		} else if zzcommon.ERROR_DISCONNECT_PEER == packetLength {
			break
		}
	}
}

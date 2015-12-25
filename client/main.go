// icartravel project main.go
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"ict_cfg"
	"pb_square"
	"time"
	"zzcommon"
	"zztcp"
	"zztimer"
)

////////////////////////////////////////////////////////////////////////////////
//服务器相关回调函数

//服务器连接成功
func onSerConn(peerConn *zzcommon.PeerConn_t) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerConn")
	return 0
}

//服务端连接关闭
func onSerConnClosed(peerConn *zzcommon.PeerConn_t) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerConnClosed")
	return 0
}

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
func onSerGetPacketLen(peerConn *zzcommon.PeerConn_t, packetLength int) (ret int) {
	if (uint32)(packetLength) < zzcommon.ProtoHeadLength { //长度不足一个包头中的长度大小
		return 0
	}
	peerConn.ParseProtoHead()
	if (uint32)(peerConn.RecvProtoHead.PacketLength) < zzcommon.ProtoHeadLength {
		return zzcommon.ERROR_DISCONNECT_PEER
	}
	fmt.Print(peerConn.RecvProtoHead)
	if gTcpServer.PacketLengthMax <= (uint32)(peerConn.RecvProtoHead.PacketLength) {
		return zzcommon.ERROR_DISCONNECT_PEER
	}
	if packetLength < int(peerConn.RecvProtoHead.PacketLength) {
		return 0
	}
	fmt.Println(peerConn)
	return int(peerConn.RecvProtoHead.PacketLength)
}

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
func onSerPacket(peerConn *zzcommon.PeerConn_t, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerPacket")
	return 0
}

func main() {
	fmt.Println("client runing...", time.Now())

	///////////////////////////////////////////////////////////////////
	//测试
	///////////////////////////////////////////////////////////////////

	///////////////////////////////////////////////////////////////////
	//加载配置文件bench.ini
	{
		if zzcommon.IsWindows() {
			ict_cfg.Gbench.Load("./bench.ini.bak")
		} else {
			ict_cfg.Gbench.Load("/Users/mlc/Desktop/GoServer/icartravel/bench.ini.bak")
		}
	}

	//////////////////////////////////////////////////////////////////
	//做为客户端
	{
		var gzztcpClient zztcp.Client_t
		gzztcpClient.OnSerConn = onSerConn
		gzztcpClient.OnSerConnClosed = onSerConnClosed
		gzztcpClient.OnSerGetPacketLen = onSerGetPacketLen
		gzztcpClient.OnSerPacket = onSerPacket

		ip := ict_cfg.Gbench.FileIni.Get("square_server", "ip", "999")
		port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("square_server", "port", "0"))
		packetLengthMax := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get("square_server", "packet_length_max", "81920"))
		err := gzztcpClient.Connect(ip, port, packetLengthMax)
		if nil != err {
			fmt.Println("######zzcliClient.Connect err:", err)
			return
		}

		//发送pb测试包
		//登录
		{
			req := &pb_square.LoginMsg{
				Account:  proto.String("17721027200"),
				Password: proto.String("7883df2788b1098886f99b0a7563a5a8"),
			}
			fmt.Println(111)
			gzztcpClient.PeerConn.Send(req, 0x100101, 0, 0, 0)
		}
		//登录
		{
			req := &pb_square.LoginMsg{
				Account:  proto.String("17721027200"),
				Password: proto.String("7883df2788b1098886f99b0a7563a5a8"),
			}
			fmt.Println(222)
			gzztcpClient.PeerConn.Send(req, 0x100101, 0, 0, 0)
		}
	}

	//////////////////////////////////////////////////////////////////
	//定时器
	zztimer.Second(1, timerSecondTest)
	fmt.Println("OK")
	for {
		time.Sleep(10 * time.Second)
		gLock.Lock()
		gLock.Unlock()
	}

	///////////////////////////////////////////////////////////////////
	//测试chan
	/*
		ch := make(chan int, 0)
		end := make(chan int, 0)
		go func(ch chan int) {
			//		time.Sleep(1000000000 * 100)
			ch <- 1
			//		time.Sleep(1000000000 * 10)
			ch <- 2
		}(ch)
		var ii int
		var jj int
		//L:
		for {
			select {
			case <-ch:
				ii++
				fmt.Println("iiiiii", ii)
				if ii >= 2 {
					//				break L
				}
			case <-end:
				jj = 100
				fmt.Println("jjjjjj", jj)
			}
		}
		var ti uint32
		for {
			ti = ti + 1
			time.Sleep(1000000000)
			//		fmt.Println("===", i)
			//		fmt.Println("end")
		}
	*/
	//////////////////////////////////////////////////////////////////

	fmt.Println("!!!!!!server done!")
}

//定时器,秒,测试
func timerSecondTest() {
	gLock.Lock()
	defer gLock.Unlock()

	//	fmt.Println("timerSecondTest...")
	zztimer.Second(1, timerSecondTest)
}

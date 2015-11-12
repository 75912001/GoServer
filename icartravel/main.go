// icartravel project main.go
package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"zzcli"
	"zzcommon"
	"zzhttp"
	"zzser"
	"zztimer"
)

var gLock = &sync.Mutex{}

func onInit() (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onInit")
	return 0
}

func onFini() (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onFini")
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//客户端相关的回调函数
func onCliConn(peerConn *zzser.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onCliConn")
	return 0
}

func onCliConnClosed(peerConn *zzser.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onCliConnClosed")
	return 0
}

func onCliGetPacketLen(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	//	fmt.Println("onCliGetPacketLen")
	return 5
	return packetLength
	return 0
}

func onCliPacket(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	//	fmt.Println("on_cli_conn")
	//	peer_conn.conn.Write("123")
	//	peer_conn.Conn.Write([]byte("123"))
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//服务器相关回调函数

//服务器连接成功
func onSerConn(peerConn *zzser.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	//	fmt.Println("onSerConn")
	return 0
}

//服务端连接关闭
func onSerConnClosed(peerConn *zzser.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerConnClosed")
	return 0
}

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
func onSerGetPacketLen(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerGetPacketLen")
	return 0
}

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
func onSerPacket(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerPacket")
	return 0
}

func main() {
	///////////////////////////////////////////////////////////////////
	//加载配置文件bench.ini
	var benchFile BenchFile
	if zzcommon.IsWindows() {
		benchFile.FileIni.Path = "./bench.ini"
	} else {
		benchFile.FileIni.Path = "/Users/mlc/Desktop/GoServer/icartravel/bench.ini"
	}
	benchFile.Load()
	//////////////////////////////////////////////////////////////////
	//作为HTTP CLIENT Weather
	httpClientWeather.Url = benchFile.FileIni.Get("weather", "url", " ")
	httpClientWeather.Get()
	//////////////////////////////////////////////////////////////////
	//作为HTTP SERVER
	var httpServer zzhttp.HttpServer
	httpServer.Ip = benchFile.FileIni.Get("http_server", "ip", "999")
	httpServer.Port = zzcommon.StringToUint16(benchFile.FileIni.Get("http_server", "port", "0"))
	var weather Weather
	weather.Register(&httpServer)
	go httpServer.Run()

	//////////////////////////////////////////////////////////////////
	//定时器
	//zztimer.Second(10, timerSecondTest)

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

	//////////////////////////////////////////////////////////////////
	//做为服务端
	var zzserServer zzser.Server
	//设置回调函数
	zzserServer.OnInit = onInit
	zzserServer.OnFini = onFini
	zzserServer.OnCliConnClosed = onCliConnClosed
	zzserServer.OnCliConn = onCliConn
	zzserServer.OnCliGetPacketLength = onCliGetPacketLen
	zzserServer.OnCliPacket = onCliPacket
	//运行
	go zzserServer.Run(benchFile.Ip, benchFile.Port, benchFile.PacketLengthMax, benchFile.Delay)

	for {
		time.Sleep(10 * time.Second)
		gLock.Lock()
		gLock.Unlock()
	}
	//////////////////////////////////////////////////////////////////
	//做为客户端
	var zzcliClient zzcli.Client
	zzcliClient.OnSerConn = onSerConn
	zzcliClient.OnSerConnClosed = onSerConnClosed
	zzcliClient.OnSerGetPacketLen = onSerGetPacketLen
	zzcliClient.OnSerPacket = onSerPacket

	gameServerIp := benchFile.FileIni.Get("game_server", "ip", "999")
	gameServerPort := zzcommon.StringToUint16(benchFile.FileIni.Get("game_server", "port", "0"))

	GUserMgr.Init()
	var conn_time uint32
	for {
		conn_time++
		fmt.Println("conn time", conn_time)
		for i := 1; i <= 10000; i++ {
			var user User
			user.Account = "mm" + strconv.Itoa(i)
			conn, err := zzcliClient.Connect(gameServerIp, gameServerPort)
			if nil != err {
				fmt.Println("######zzcli_client.Connect err:", err)
			} else {
				user.Conn = conn

				{ //登录
					/*
						req := &game_msg.LoginMsg{
							Platform: proto.Uint32(0),
							Account:  proto.String(user.Account),
							Password: proto.String(user.Account),
						}
						user.Send(0x00010101, req)
					*/
				}

				GUserMgr.UserMap[conn] = user
				go zzcliClient.ClientRecv(conn, benchFile.PacketLengthMax)
			}
			if 0 == i%1000 {
				fmt.Println(i)
			}
		}
		fmt.Println("conn done")
		time.Sleep(1000000000 * 10)
		fmt.Println("close")
		var close_idx uint32
		for k, v := range GUserMgr.UserMap {
			close_idx++
			if 0 == close_idx%1000 {
				fmt.Println(close_idx)
			}
			v.Conn.Close()
			delete(GUserMgr.UserMap, k)
		}
		fmt.Println("close done")
		time.Sleep(1000000000)
	}

	fmt.Println("!!!!!!server done!")
}

//定时器,秒,测试
func timerSecondTest() {
	fmt.Println("timerSecondTest...")
	zztimer.Second(1, timerSecondTest)
}

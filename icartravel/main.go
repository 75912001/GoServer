// icartravel project main.go
package main

import (
	//	"common_msg"
	"fmt"
	"game_msg"
	"strconv"
	"sync"
	"time"
	"zzcommon"
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
	if zzcommon.IsWindows() {
		gBenchFile.FileIni.Path = "./bench.ini"
	} else {
		gBenchFile.FileIni.Path = "/Users/mlc/Desktop/GoServer/icartravel/bench.ini"
	}
	gBenchFile.Load()
	//////////////////////////////////////////////////////////////////
	//做为客户端
	gzzcliClient.OnSerConn = onSerConn
	gzzcliClient.OnSerConnClosed = onSerConnClosed
	gzzcliClient.OnSerGetPacketLen = onSerGetPacketLen
	gzzcliClient.OnSerPacket = onSerPacket

	gameServerIp := gBenchFile.FileIni.Get("game_server", "ip", "999")
	gameServerPort := zzcommon.StringToUint16(gBenchFile.FileIni.Get("game_server", "port", "0"))

	var userCount = 10000
	GUserMgr.Init()
	var connCount uint32
	for {
		connCount++
		fmt.Println("conn time", connCount)
		for i := 1; i <= userCount; i++ {
			var user User
			user.Account = "mm" + strconv.Itoa(i)
			conn, err := gzzcliClient.Connect(gameServerIp, gameServerPort)
			if nil != err {
				fmt.Println("######zzcli_client.Connect err:", err)
			} else {
				user.Conn = conn

				{ //登录

					req := &game_msg.LoginMsg{
						Platform: proto.Uint32(0),
						Account:  proto.String(user.Account),
						Password: proto.String(user.Account),
					}
					user.Send(0x00010101, req)
				}

				GUserMgr.UserMap[conn] = user
				go gzzcliClient.ClientRecv(conn, gBenchFile.PacketLengthMax)
			}
			if 0 == i%100 {
				fmt.Println(i)
			}
		}
		fmt.Println("conn done")
		time.Sleep(10 * time.Second)
		fmt.Println("close")
		var closeCount uint32
		for k, v := range GUserMgr.UserMap {
			closeCount++
			if 0 == closeCount%100 {
				fmt.Println(closeCount)
			}
			v.Conn.Close()
			delete(GUserMgr.UserMap, k)
		}
		fmt.Println("close done")
		time.Sleep(10 * time.Second)
	}
	//////////////////////////////////////////////////////////////////
	//作为HTTP CLIENT Weather
	gHttpClientWeather.Url = gBenchFile.FileIni.Get("weather", "url", " ")
	gHttpClientWeather.Get()
	//////////////////////////////////////////////////////////////////
	//作为HTTP SERVER
	gHttpServer.Ip = gBenchFile.FileIni.Get("http_server", "ip", "999")
	gHttpServer.Port = zzcommon.StringToUint16(gBenchFile.FileIni.Get("http_server", "port", "0"))
	gHttpServer.AddHandler(pattern, WeatherHttpHandler)
	go gHttpServer.Run()

	//////////////////////////////////////////////////////////////////
	//定时器
	//zztimer.Second(10, timerSecondTest)

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
	go zzserServer.Run(gBenchFile.Ip, gBenchFile.Port, gBenchFile.PacketLengthMax, gBenchFile.Delay)

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
	fmt.Println("timerSecondTest...")
	zztimer.Second(1, timerSecondTest)
}

// icartravel project main.go
package main

import (
	//"common_msg"
	"fmt"
	//	"game_msg"
	//	"github.com/golang/protobuf/proto"
	//	"strconv"
	"time"
	"zzcli"
	"zzcommon"
	//	"zzser"

	"math/rand"
	//	"strconv"
	"zztimer"
)

func onInit() (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	gUserMgr.Init()
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
func onCliConn(peerConn *zzcommon.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	var user User
	user.PeerConn = peerConn
	gUserMgr.UserMap[user.PeerConn] = user

	fmt.Println("onCliConn")
	return 0
}

func onCliConnClosed(peerConn *zzcommon.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	delete(gUserMgr.UserMap, peerConn)
	fmt.Println("onCliConnClosed")
	return 0
}

func onCliGetPacketLen(peerConn *zzcommon.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	//	fmt.Println("onCliGetPacketLen")
	return packetLength
	return 0
}

func onCliPacket(peerConn *zzcommon.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	//fmt.Println("on_cli_conn")
	//	peer_conn.conn.Write("123")
	//	peer_conn.Conn.Write([]byte("123"))
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//服务器相关回调函数

//服务器连接成功
func onSerConn(peerConn *zzcommon.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	//	fmt.Println("onSerConn")
	return 0
}

//服务端连接关闭
func onSerConnClosed(peerConn *zzcommon.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerConnClosed")
	return 0
}

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
func onSerGetPacketLen(peerConn *zzcommon.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerGetPacketLen packetLength:", packetLength)
	return 0
}

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
func onSerPacket(peerConn *zzcommon.PeerConn, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerPacket")
	return 0
}

func main() {
	rand.Seed(time.Now().Unix())

	fmt.Println("server runing...", time.Now())
	///////////////////////////////////////////////////////////////////
	//加载配置文件bench.ini
	if zzcommon.IsWindows() {
		gBenchFile.FileIni.Path = "./bench.ini"
	} else {
		gBenchFile.FileIni.Path = "/Users/mlc/Desktop/GoServer/icartravel/bench.ini"
	}
	gBenchFile.Load()
	//////////////////////////////////////////////////////////////////
	//做为服务端
	//设置回调函数
	gzzserServer.OnInit = onInit
	gzzserServer.OnFini = onFini
	gzzserServer.OnCliConnClosed = onCliConnClosed
	gzzserServer.OnCliConn = onCliConn
	gzzserServer.OnCliGetPacketLength = onCliGetPacketLen
	gzzserServer.OnCliPacket = onCliPacket
	//运行
	go gzzserServer.Run(gBenchFile.Ip, gBenchFile.Port, gBenchFile.PacketLengthMax, gBenchFile.Delay)

	//////////////////////////////////////////////////////////////////
	//作为HTTP CLIENT Weather
	//	gHttpClientWeather.Url = gBenchFile.FileIni.Get("weather", "url", " ")
	//	gHttpClientWeather.Get()
	//////////////////////////////////////////////////////////////////
	//作为HTTP SERVER
	gHttpServer.Ip = gBenchFile.FileIni.Get("http_server", "ip", "999")
	gHttpServer.Port = zzcommon.StringToUint16(gBenchFile.FileIni.Get("http_server", "port", "0"))
	gHttpServer.AddHandler(weatherPattern, WeatherHttpHandler)
	gHttpServer.AddHandler(loginPattern, LoginHttpHandler)

	{ //启动手机注册功能
		err := gSmsPhoneRegister.Init()
		if nil != err {
			return
		}
		gHttpServer.AddHandler(gSmsPhoneRegister.Pattern, SmsPhoneRegisterHttpHandler)
		err = gPhoneRegister.Init()
		if nil != err {
			return
		}
		gHttpServer.AddHandler(gPhoneRegister.Pattern, PhoneRegisterHttpHandler)
	}

	go gHttpServer.Run()

	//////////////////////////////////////////////////////////////////

	//////////////////////////////////////////////////////////////////
	//定时器
	zztimer.Second(1, timerSecondTest)
	fmt.Println("OK")
	for {
		time.Sleep(10 * time.Second)
		gLock.Lock()
		gLock.Unlock()
	}

	//////////////////////////////////////////////////////////////////
	//做为客户端
	var gzzcliClient zzcli.Client
	gzzcliClient.OnSerConn = onSerConn
	gzzcliClient.OnSerConnClosed = onSerConnClosed
	gzzcliClient.OnSerGetPacketLen = onSerGetPacketLen
	gzzcliClient.OnSerPacket = onSerPacket

	gameServerIp := gBenchFile.FileIni.Get("game_server", "ip", "999")
	gameServerPort := zzcommon.StringToUint16(gBenchFile.FileIni.Get("game_server", "port", "0"))
	err := gzzcliClient.Connect(gameServerIp, gameServerPort, gBenchFile.PacketLengthMax)
	if nil != err {
		fmt.Println("######zzcliClient.Connect err:", err)
	} else {
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

// icartravel project main.go
package main

import (
	//"common_msg"
	"fmt"
	//	"game_msg"
	//	"github.com/golang/protobuf/proto"
	//	"strconv"
	"ict_cfg"
	"ict_login"
	"ict_register"
	"ict_user"
	"math/rand"
	"runtime"
	"strconv"
	"time"
	"zzcommon"
	"zztcp"
	"zztimer"
)

func onInit() (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	ict_user.GuserMgr.Init()
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

	var user ict_user.User
	user.PeerConn = peerConn
	ict_user.GuserMgr.UserMap[user.PeerConn] = user

	fmt.Println("onCliConn")
	return 0
}

func onCliConnClosed(peerConn *zzcommon.PeerConn) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	delete(ict_user.GuserMgr.UserMap, peerConn)
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
	{
		if zzcommon.IsWindows() {
			ict_cfg.Gbench.Load("./bench.ini")
		} else {
			ict_cfg.Gbench.Load("/Users/mlc/Desktop/GoServer/icartravel/bench.ini.bak")
		}
	}

	//////////////////////////////////////////////////////////////////
	//做为服务端
	//设置回调函数
	{
		gTcpServer.OnInit = onInit
		gTcpServer.OnFini = onFini
		gTcpServer.OnCliConnClosed = onCliConnClosed
		gTcpServer.OnCliConn = onCliConn
		gTcpServer.OnCliGetPacketLength = onCliGetPacketLen
		gTcpServer.OnCliPacket = onCliPacket

		//运行
		delay := true
		ip := ict_cfg.Gbench.FileIni.Get("server", "ip", "")
		port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("server", "port", "0"))
		packetLengthMax := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get("common", "packet_length_max", "81920"))
		str_num_cpu := strconv.Itoa(runtime.NumCPU())
		goProcessMax := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get("common", "go_process_max", str_num_cpu))
		runtime.GOMAXPROCS(goProcessMax)
		go gTcpServer.Run(ip, port, packetLengthMax, delay)
	}

	//////////////////////////////////////////////////////////////////
	//作为HTTP CLIENT Weather
	//	gHttpClientWeather.Url = ict_bench_file.GbenchFile.FileIni.Get("weather", "url", " ")
	//	gHttpClientWeather.Get()
	//////////////////////////////////////////////////////////////////
	//作为HTTP SERVER
	{
		ip := ict_cfg.Gbench.FileIni.Get("http_server", "ip", "999")
		port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("http_server", "port", "0"))
		gHttpServer.AddHandler(weatherPattern, WeatherHttpHandler)
		gHttpServer.AddHandler(ict_login.LoginPattern, ict_login.LoginHttpHandler)

		{ //启动手机注册功能
			err := ict_register.GphoneSms.Init()
			if nil != err {
				return
			}
			gHttpServer.AddHandler(ict_register.GphoneSms.Pattern, ict_register.PhoneSmsHttpHandler)

			err = ict_register.Gphone.Init()
			if nil != err {
				return
			}
			gHttpServer.AddHandler(ict_register.Gphone.Pattern, ict_register.PhoneHttpHandler)

			//			err = ict_register.Gphone.Init()
			//			if nil != err {
			//				return
			//			}

			err = ict_user.GuidMgr.Init()
			if nil != err {
				return
			}
		}

		go gHttpServer.Run(ip, port)
	}
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
	{
		var gzztcpClient zztcp.Client
		gzztcpClient.OnSerConn = onSerConn
		gzztcpClient.OnSerConnClosed = onSerConnClosed
		gzztcpClient.OnSerGetPacketLen = onSerGetPacketLen
		gzztcpClient.OnSerPacket = onSerPacket

		ip := ict_cfg.Gbench.FileIni.Get("game_server", "ip", "999")
		port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("game_server", "port", "0"))
		packetLengthMax := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get("game_server", "packet_length_max", "81920"))
		err := gzztcpClient.Connect(ip, port, packetLengthMax)
		if nil != err {
			fmt.Println("######zzcliClient.Connect err:", err)
		} else {
		}
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

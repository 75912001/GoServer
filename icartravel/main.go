// icartravel project main.go
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"ict_account"
	"ict_cfg"
	"ict_common"
	"ict_login"
	"ict_phone_sms"
	"ict_user"
	"math/rand"
	"pb_square"
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
func onCliConn(peerConn *zzcommon.PeerConn_t) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	var user ict_user.User_t
	user.PeerConn = peerConn
	ict_user.GuserMgr.UserMap[user.PeerConn] = &user

	fmt.Println("onCliConn")
	return 0
}

func onCliConnClosed(peerConn *zzcommon.PeerConn_t) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	delete(ict_user.GuserMgr.UserMap, peerConn)
	fmt.Println("onCliConnClosed")
	return 0
}

func onCliGetPacketLen(peerConn *zzcommon.PeerConn_t, packetLength int) (ret int) {
	//fmt.Println("onCliGetPacketLen")
	if uint32(packetLength) < zzcommon.ProtoHeadLength { //长度不足一个包头中的长度大小
		return 0
	}
	peerConn.ParseProtoHeadPacketLength()
	if uint32(peerConn.RecvProtoHead.PacketLength) < zzcommon.ProtoHeadLength {
		return zzcommon.ERROR_DISCONNECT_PEER
	}
	if gTcpServer.PacketLengthMax <= uint32(peerConn.RecvProtoHead.PacketLength) {
		return zzcommon.ERROR_DISCONNECT_PEER
	}
	if packetLength < int(peerConn.RecvProtoHead.PacketLength) {
		return 0
	}
	fmt.Println("onCliGetPacketLen:", peerConn.RecvProtoHead.PacketLength)
	return int(peerConn.RecvProtoHead.PacketLength)
}

func onCliPacket(peerConn *zzcommon.PeerConn_t) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onCliPacket")
	peerConn.ParseProtoHead()

	ret = onRecv(peerConn)
	return ret
}

////////////////////////////////////////////////////////////////////////////////
//服务器相关回调函数

//服务器连接成功
func onSerConn(peerConn *zzcommon.PeerConn_t) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	//	fmt.Println("onSerConn")
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
	//	gLock.Lock()
	//	defer gLock.Unlock()

	fmt.Println("onSerGetPacketLen packetLength:", packetLength)
	return 0
}

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
func onSerPacket(peerConn *zzcommon.PeerConn_t, packetLength int) (ret int) {
	gLock.Lock()
	defer gLock.Unlock()

	fmt.Println("onSerPacket")
	return 0
}

//////////////////////////////////////////////////////////////////////
//测试

func OnLoginMsg(user *ict_user.User_t, protoMessage proto.Message) (ret int) {
	fmt.Println("OnLoginMsg")
	fmt.Println(user)
	fmt.Println(protoMessage)
	var ss *pb_square.LoginMsg = protoMessage.(*pb_square.LoginMsg)
	fmt.Println(ss.GetAccount())
	return ret
}

func fn1(ss *[]int) {
	*ss = append(*ss, 1, 2, 3, 4, 5)

	(*ss)[2] = 555
	fmt.Println(*ss)
	//	fmt.Println(st)
}

func init() {
	fmt.Println("init")
}
func main() {
	fmt.Println("main")
	var s1 = []int{1, 2, 3, 4, 5}
	//var s2 = s1

	fn1(&s1)
	fmt.Println(s1)
	//fmt.Println(s2)
	fmt.Println("OK")
	//	var ss []byte
	//	for {
	//		time.Sleep(1 * time.Second)
	//		fmt.Println("OK")
	//		ss = make([]byte, 1024*1024*100)
	//		fmt.Println(ss[0:1])
	//	}

	////////////////////////////////////////////////////////////////////
	rand.Seed(time.Now().Unix())

	ret := initPbFun()
	if zzcommon.SUCC != ret {
		fmt.Println("######initPbFun")
		return
	}

	fmt.Println("server runing...", time.Now())

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
	//redis
	{
		const benchFileSection string = "redis_server"
		ip := ict_cfg.Gbench.FileIni.Get(benchFileSection, "ip", " ")
		port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get(benchFileSection, "port", " "))
		redisDatabases := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get(benchFileSection, "databases", " "))

		//链接redis
		err := ict_common.GRedisClient.Connect(ip, port, redisDatabases)
		if nil != err {
			fmt.Println("######ict_common.GRedisClient.Connect(ip, port, redisDatabases) err:", err)
			return
		}
	}

	//////////////////////////////////////////////////////////////////
	//作为HTTP SERVER
	{
		//phome sms
		{
			err := ict_phone_sms.GphoneSms.Init()
			if nil != err {
				fmt.Println("######ict_phone_sms.GphoneSms.Init()")
				return
			}
		}

		gHttpServer.AddHandler(weatherPattern, WeatherHttpHandler)
		gHttpServer.AddHandler(ict_login.LoginPattern, ict_login.LoginHttpHandler)

		//启动手机注册功能
		{
			err := ict_account.GphoneSmsRegister.Init()
			if nil != err {
				fmt.Println("错误ict_account.GphoneSmsRegister.Init()")
				return
			}
			gHttpServer.AddHandler(ict_account.GphoneSmsRegister.Pattern, ict_account.PhoneSmsRegisterHttpHandler)

			err = ict_account.GphoneSmsChangePwd.Init()
			if nil != err {
				fmt.Println("错误ict_account.GphoneSmsChangePwd.Init()")
				return
			}
			gHttpServer.AddHandler(ict_account.GphoneSmsChangePwd.Pattern, ict_account.PhoneSmsChangePwdHttpHandler)

			err = ict_account.GphoneChangePwd.Init()
			if nil != err {
				fmt.Println("错误ict_account.GphoneChangePwd.Init()")
				return
			}
			gHttpServer.AddHandler(ict_account.GphoneChangePwd.Pattern, ict_account.PhoneChangePwdHttpHandler)

			err = ict_account.GphoneRegister.Init()
			if nil != err {
				fmt.Println("错误ict_account.GphoneRegister.Init()")
				return
			}
			gHttpServer.AddHandler(ict_account.GphoneRegister.Pattern, ict_account.PhoneRegisterHttpHandler)

			err = ict_user.Gbase.Init()
			if nil != err {
				fmt.Println("错误ict_user.Gbase.Init()")
				return
			}

			err = ict_user.GuidMgr.Init()
			if nil != err {
				fmt.Println("错误ict_user.GuidMgr.Init()")
				return
			}

			err = ict_login.Glogin.Init()
			if nil != err {
				fmt.Println("错误ict_login.Glogin.Init()")
				return
			}
		}

		ip := ict_cfg.Gbench.FileIni.Get("http_server", "ip", "999")
		port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("http_server", "port", "0"))

		go gHttpServer.Run(ip, port)
	}

	//////////////////////////////////////////////////////////////////
	//做为服务端
	{ //设置回调函数
		gTcpServer.OnInit = onInit
		gTcpServer.OnFini = onFini
		gTcpServer.OnCliConnClosed = onCliConnClosed
		gTcpServer.OnCliConn = onCliConn
		gTcpServer.OnCliGetPacketLen = onCliGetPacketLen
		gTcpServer.OnCliPacket = onCliPacket

		//运行
		noDelay := true
		ip := ict_cfg.Gbench.FileIni.Get("server", "ip", "")
		port := zzcommon.StringToUint16(ict_cfg.Gbench.FileIni.Get("server", "port", "0"))
		gTcpServer.PacketLengthMax = zzcommon.StringToUint32(ict_cfg.Gbench.FileIni.Get("common", "packet_length_max", "81920"))
		str_num_cpu := strconv.Itoa(runtime.NumCPU())
		goProcessMax := zzcommon.StringToInt(ict_cfg.Gbench.FileIni.Get("common", "go_process_max", str_num_cpu))
		runtime.GOMAXPROCS(goProcessMax)
		go gTcpServer.Run(ip, port, noDelay)
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

	//////////////////////////////////////////////////////////////////
	//作为HTTP CLIENT Weather
	//	gHttpClientWeather.Url = ict_bench_file.GbenchFile.FileIni.Get("weather", "url", " ")
	//	gHttpClientWeather.Get()

	//////////////////////////////////////////////////////////////////
	//做为客户端
	{
		var gzztcpClient zztcp.Client_t
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

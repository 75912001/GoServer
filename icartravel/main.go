package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"zzcli"
	"zzser"
	"zztimer"
)

var globalLock = &sync.Mutex{}

func onInit() (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fmt.Println("onInit")
	return 0
}

func onFini() (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fmt.Println("onFini")
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//客户端相关的回调函数
func onCliConn(peerConn *zzser.PeerConn) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fmt.Println("onCliConn")
	return 0
}

func onCliConnClosed(peerConn *zzser.PeerConn) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fmt.Println("onCliConnClosed")
	return 0
}

func onCliGetPacketLen(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	//	fmt.Println("onCliGetPacketLen")
	return 5
	return packetLength
	return 0
}

func onCliPacket(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	//	fmt.Println("on_cli_conn")
	//	peer_conn.conn.Write("123")
	//	peer_conn.Conn.Write([]byte("123"))
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//服务器相关回调函数

//服务器连接成功
func onSerConn(peerConn *zzser.PeerConn) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	//	fmt.Println("onSerConn")
	return 0
}

//服务端连接关闭
func onSerConnClosed(peerConn *zzser.PeerConn) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fmt.Println("onSerConnClosed")
	return 0
}

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
func onSerGetPacketLen(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fmt.Println("onSerGetPacketLen")
	return 0
}

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
func onSerPacket(peerConn *zzser.PeerConn, packetLength int) (ret int) {
	globalLock.Lock()
	defer globalLock.Unlock()

	fmt.Println("onSerPacket")
	return 0
}

func main() {
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
	//测试HTTP SERVER

	//////////////////////////////////////////////////////////////////
	//定时器
	zztimer.Second(1, timerSecondTest)

	//////////////////////////////////////////////////////////////////
	//做为服务端
	var zzserServer zzser.Server
	zzserServer.FileIni.Path = "./bench.ini"

	zzserServer.Delay = true
	//设置回调函数
	zzserServer.OnInit = onInit
	zzserServer.OnFini = onFini
	zzserServer.OnCliConnClosed = onCliConnClosed
	zzserServer.OnCliConn = onCliConn
	zzserServer.OnCliGetPacketLength = onCliGetPacketLen
	zzserServer.OnCliPacket = onCliPacket

	zzserServer.LoadConfig()

	go zzserServer.Run()
	for {
		time.Sleep(10 * time.Second)
		globalLock.Lock()
		globalLock.Unlock()
	}

	//////////////////////////////////////////////////////////////////
	//做为客户端
	var zzcliClient zzcli.Client
	zzcliClient.OnSerConn = onSerConn
	zzcliClient.OnSerConnClosed = onSerConnClosed
	zzcliClient.OnSerGetPacketLen = onSerGetPacketLen
	zzcliClient.OnSerPacket = onSerPacket

	game_server_ip := zzserServer.FileIni.Get("game_server", "ip", "999")
	_port, _ := (strconv.ParseUint(zzserServer.FileIni.Get("game_server", "port", "0"), 10, 16))
	game_server_port := uint16(_port)

	G_user_mgr.Init()
	var conn_time uint32
	for {
		conn_time++
		fmt.Println("conn time", conn_time)
		for i := 1; i <= 10000; i++ {
			var user User_t
			user.Account = "mm" + strconv.Itoa(i)
			conn, err := zzcliClient.Connect(game_server_ip, game_server_port)
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

				G_user_mgr.User_map[conn] = user
				go zzcliClient.ClientRecv(conn, zzserServer.PacketLengthMax)
			}
			if 0 == i%1000 {
				fmt.Println(i)
			}
		}
		fmt.Println("conn done")
		time.Sleep(1000000000 * 10)
		fmt.Println("close")
		var close_idx uint32
		for k, v := range G_user_mgr.User_map {
			close_idx++
			if 0 == close_idx%1000 {
				fmt.Println(close_idx)
			}
			v.Conn.Close()
			delete(G_user_mgr.User_map, k)
		}
		fmt.Println("close done")
		time.Sleep(1000000000)
	}

	fmt.Println("done!")
	var i uint32
	for {
		i = i + 1
		time.Sleep(1000000000)
		//		fmt.Println("===", i)
		//		fmt.Println("end")
	}
}

//定时器,秒,测试
func timerSecondTest() {
	fmt.Println("timerSecondTest...")
	zztimer.Second(1, timerSecondTest)
}

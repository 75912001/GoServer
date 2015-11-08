// mygo project main.go
package main

import (
	"fmt"
	//	"game_msg"
	//	"proto"
	"strconv"
	"sync"
	"time"
	"zzcli"
	"zzser"
)

var G_lock = &sync.Mutex{}

func on_init() (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()

	fmt.Println("on_init")
	return 0
}

func on_fini() (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()

	fmt.Println("on_fini")
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//客户端相关的回调函数
func on_cli_conn(peer_conn *zzser.PeerConn) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()

	fmt.Println("on_cli_conn")
	return 0
}

func on_cli_conn_closed(peer_conn *zzser.PeerConn) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()

	fmt.Println("on_cli_conn_closed")
	return 0
}

func on_cli_get_pkg_len(peer_conn *zzser.PeerConn, pkg_len int) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()

	//	fmt.Println("on_get_pkg_len")
	return 5
	return pkg_len
	return 0
}

func on_cli_pkg(peer_conn *zzser.PeerConn, pkg_len int) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()

	//	fmt.Println("on_cli_conn")
	//	peer_conn.conn.Write("123")
	//	peer_conn.Conn.Write([]byte("123"))
	return 0
}

////////////////////////////////////////////////////////////////////////////////
//服务器相关回调函数

//服务器连接成功
func on_ser_conn(peer_conn *zzser.PeerConn) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()

	//	fmt.Println("on_ser_conn")
	return 0
}

//服务端连接关闭
func on_ser_conn_closed(peer_conn *zzser.PeerConn) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()
	//	fmt.Println("on_ser_conn_closed")
	return 0
}

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
func on_ser_get_pkg_len(peer_conn *zzser.PeerConn, pkg_len int) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()
	fmt.Println("on_ser_get_pkg_len")
	return 0
}

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
func on_ser_pkg(peer_conn *zzser.PeerConn, pkg_len int) (ret int) {
	G_lock.Lock()
	defer G_lock.Unlock()
	fmt.Println("on_ser_pkg")
	return 0
}

func f1() {
	//	fmt.Println("f1 done !")
	time.AfterFunc(1*time.Second, f1)
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
	time.AfterFunc(1*time.Second, f1)

	//做为服务端
	var zzser_server zzser.Server
	zzser_server.FileIni.Path = "./bench.ini"

	zzser_server.Delay = true
	//设置回调函数
	zzser_server.OnInit = on_init
	zzser_server.OnFini = on_fini
	zzser_server.OnCliConnClosed = on_cli_conn_closed
	zzser_server.OnCliConn = on_cli_conn
	zzser_server.OnCliGetPacketLength = on_cli_get_pkg_len
	zzser_server.OnCliPacket = on_cli_pkg

	zzser_server.LoadConfig()
	go zzser_server.Run()
	var jj uint32
	for {
		time.Sleep(1000000000 * 10)
		G_lock.Lock()
		G_lock.Unlock()
		jj++
		if 0 == jj%100000000 {
			//			fmt.Println(jj)
		}
	}

	//做为客户端
	zzcli_client := new(zzcli.Client_t)
	zzcli_client.On_ser_conn = on_ser_conn
	zzcli_client.On_ser_conn_closed = on_ser_conn_closed
	zzcli_client.On_ser_get_pkg_len = on_ser_get_pkg_len
	zzcli_client.On_ser_pkg = on_ser_pkg

	game_server_ip := zzser_server.FileIni.Get("game_server", "ip", "999")
	_port, _ := (strconv.ParseUint(zzser_server.FileIni.Get("game_server", "port", "0"), 10, 16))
	game_server_port := uint16(_port)

	G_user_mgr.Init()
	var conn_time uint32
	for {
		conn_time++
		fmt.Println("conn time", conn_time)
		for i := 1; i <= 10000; i++ {
			var user User_t
			user.Account = "mm" + strconv.Itoa(i)
			conn, _, err := zzcli_client.Connect(game_server_ip, game_server_port)
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
				go zzcli_client.Client_recv(conn, zzser_server.PacketLengthMax)
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

package zzcli

import (
	"fmt"
	"net"
	"strconv"
	"zzcommon"
	"zzser"
)

//服务端连接建立
type ON_SER_CONN func(peer_conn *zzser.PeerConn) int

//服务端连接关闭
type ON_SER_CONN_CLOSED func(peer_conn *zzser.PeerConn) int

//获取消息的长度,0表示消息还未接受完成,
//ERROR_DISCONNECT_PEER表示长度有误,服务端断开
type ON_SER_GET_PKG_LEN func(peer_conn *zzser.PeerConn, pkg_len int) int

//服务端消息
//返回ERROR_DISCONNECT_PEER断开服务端
type ON_SER_PKG func(peer_conn *zzser.PeerConn, pkg_len int) int

//客户端
type Client_t struct {
	On_ser_conn        ON_SER_CONN
	On_ser_conn_closed ON_SER_CONN_CLOSED
	On_ser_get_pkg_len ON_SER_GET_PKG_LEN
	On_ser_pkg         ON_SER_PKG
}

//连接
func (client *Client_t) Connect(ip string, port uint16) (conn *net.TCPConn, ret int32, err error) {

	var str_addr = ip + ":" + strconv.Itoa(int(port))
	tcpAddr, err := net.ResolveTCPAddr("tcp4", str_addr)
	if nil != err {
		fmt.Println("######net.ResolveTCPAddr err:", err)
		return conn, ret, err
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		fmt.Println("######net.Dial err:", err)
		return conn, ret, err
	}
	//	fmt.Println("Connect conn:", conn)
	return conn, ret, err
}

func (client *Client_t) Client_recv(conn *net.TCPConn, recv_buf_max int) {
	//	var peer_ip = conn.RemoteAddr().String()
	//	fmt.Println("connection server:", peer_ip)

	var peer_conn zzser.PeerConn
	peer_conn.Conn = conn
	client.On_ser_conn(&peer_conn)
	defer conn.Close()

	defer client.On_ser_conn_closed(&peer_conn)

	//todo 优化[消耗内存过大]
	recv_buf := make([]byte, recv_buf_max)

	var read_pos int

	for {
		read_num, err := conn.Read(recv_buf[read_pos:])
		//		fmt.Println("######conn.Read read_num:", read_num)
		if err != nil {
			//			fmt.Println("######conn.Read err:", read_num, err)
			break
		}

		read_pos = read_pos + read_num
		ok_len := client.On_ser_get_pkg_len(&peer_conn, read_pos)
		if ok_len == zzcommon.ERROR_DISCONNECT_PEER {
			fmt.Println("######On_ser_get_pkg_len:", zzcommon.ERROR_DISCONNECT_PEER)
			break
		}
		if ok_len > 0 { //有完整的包
			ret := client.On_ser_pkg(&peer_conn, ok_len)
			if ret == zzcommon.ERROR_DISCONNECT_PEER {
				fmt.Println("######On_ser_pkg:", zzcommon.ERROR_DISCONNECT_PEER)
				break
			}
			recv_buf = recv_buf[read_pos-ok_len : read_pos]
			read_pos = read_pos - ok_len
		}
	}
	recv_buf = nil
}

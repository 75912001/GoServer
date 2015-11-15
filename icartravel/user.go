package main

import (
	"zzcommon"
)

var gUserMgr UserMgr

type User struct {
	PeerConn *zzcommon.PeerConn
	Account  string
	Uid      zzcommon.USER_ID
}

type USER_MAP map[*zzcommon.PeerConn]User

type UserMgr struct {
	UserMap USER_MAP
}

func (p *UserMgr) Init() {
	p.UserMap = make(USER_MAP)
}

/*
	user.Account = "mm" + strconv.Itoa(i)
	//登录
	req := &game_msg.LoginMsg{
		Platform: proto.Uint32(0),
		Account:  proto.String(user.Account),
		Password: proto.String(user.Account),
	}
	user.Send(0x00010101, req)
}
*/

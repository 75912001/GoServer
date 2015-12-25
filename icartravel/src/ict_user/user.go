package ict_user

import (
	"zzcommon"
)

var GuserMgr userMgr_t

type User_t struct {
	PeerConn *zzcommon.PeerConn_t
	Account  string
	Uid      zzcommon.USER_ID
}

type USER_MAP map[*zzcommon.PeerConn_t]*User_t

type userMgr_t struct {
	UserMap USER_MAP
}

func (p *userMgr_t) Init() {
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

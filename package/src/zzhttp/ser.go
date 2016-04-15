/*
////////////////////////////////////////////////////////////////////////////////
//使用方法
import (
	"zzhttp"
)

func main() {
	var gHttpServer zzhttp.Server_t
	gHttpServer.AddHandler("/PhoneRegister", PhoneRegisterHttpHandler)
	go gHttpServer.Run(ip, port)
}

func PhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
}
*/

package zzhttp

import (
	"fmt"
	"net/http"
	"strconv"
)

type Server_t struct {
}

func (p *Server_t) AddHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

func (p *Server_t) Run(ip string, port uint16) {
	httpAddr := ip + ":" + strconv.Itoa(int(port))
	err := http.ListenAndServe(httpAddr, nil)
	if nil != err {
		fmt.Println("######http.ListenAndServe: ", err, ip, port, httpAddr)
	}
}

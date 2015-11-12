package zzhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	Url    string
	Result []byte
}

func (p *HttpClient) Get() {
	resp, err := http.Get(p.Url)
	if nil != err {
		fmt.Println("######HttpClient.Get err:", err, p.Url)
		return
	}
	defer resp.Body.Close()

	p.Result = make([]byte, resp.ContentLength)

	p.Result, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println("######HttpClient.Get err:", err, resp.Body)
		return
	}
}

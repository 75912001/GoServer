package zzhttp

import (
	//	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	Url    string
	Result []byte
}
type Book struct {
	Title       string
	Publisher   string
	IsPublished bool
	Price       float32
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
	//	var strResult string = string(p.Result)
	//	fmt.Println("######")
	//	fmt.Println(strResult)
	//	fmt.Println("######")
	//测试json
	//var gobook Book
	//gobook.Title = "Go语言编程"
	//gobook.Publisher = "ituring.com.cn"
	//gobook.IsPublished = true
	//gobook.Price = 9.99
	//p.Result, err = json.Marshal(gobook)
	//	fmt.Println(p.Result)
}

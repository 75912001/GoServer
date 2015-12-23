package zzhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client_t struct {
	Result []byte
}

//type Book struct {
//	Title       string
//	Publisher   string
//	IsPublished bool
//	Price       float32
//}

func (p *Client_t) Get(url string) (err error) {
	resp, err := http.Get(url)
	if nil != err {
		fmt.Println("######HttpClient.Get err:", err, url)
		return
	}
	defer resp.Body.Close()

	p.Result = make([]byte, resp.ContentLength)

	p.Result, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println("######HttpClient.Get err:", err, resp.Body)
		return err
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
	return err
}

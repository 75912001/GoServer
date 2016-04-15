/*
////////////////////////////////////////////////////////////////////////////////
//使用方法
import (
	"zzhttp"
)
func main() {
	var gHttpClient zzhttp.Client_t
	err := gHttpClient.Get(url)
	if nil != err {
		fmt.Println("######HttpClient.Get err:", err)
		return err
	}
	fmt.Println(gHttpClient.Result)
}
*/

package zzhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client_t struct {
	Result []byte
}

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
	var strResult string = string(p.Result)
	fmt.Println("~~~~~~")
	fmt.Println(url)
	fmt.Println(resp)
	fmt.Println(strResult)
	fmt.Println(p.Result)
	fmt.Println("~~~~~~")
	return err
}

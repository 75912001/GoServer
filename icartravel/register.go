package main

import (
	"fmt"
	"net/http"
)

const phoneRegisterPattern string = "/phoneRegister"

//?number=17721027200

//手机号码长度
const phoneNumberLen int = 11

func PhoneRegisterHttpHandler(w http.ResponseWriter, req *http.Request) {
	var phoneNumber string
	req.ParseForm()
	if len(req.Form["number"]) > 0 {
		phoneNumber = req.Form["number"][0]
	}
	if phoneNumberLen != len(phoneNumber) {
		return
	}
	fmt.Println(phoneNumber)
}

type PhoneRegister struct {
}

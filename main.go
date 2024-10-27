package main

import (
	"Idcsmart_login/model"
	"fmt"
)

func main() {
	//邮箱登录测试
	domain := "idc_smart_domain"
	ck, err := model.GetVerifyFrontInfo(domain)
	if err != nil {
		fmt.Println(err)
	}
	ck.Email = "123@123.com"
	ck.Password = "123456"
	result, err := model.VerifyUser(domain, ck, "email")
	if err != nil {
		fmt.Println(err)
	}
	if result {
		fmt.Println("登录成功")
	}
}

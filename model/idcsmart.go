package model

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type IdcsmartVerify struct {
	Cookie    string
	Token     string
	Email     string
	Phone     string
	PhoneCode int
	Password  string
}

var c = colly.NewCollector()

func GetVerifyFrontInfo(domain string) (IdcsmartVerify, error) {
	var verify_info IdcsmartVerify
	var err error
	// 注册回调函数来处理HTML响应
	c.OnResponse(func(response *colly.Response) {
		// 获取并打印Cookies
		verify_info.Cookie = strings.TrimPrefix(strings.Split(response.Headers.Get("Set-Cookie"), ";")[0], "PHPSESSID=")
	})
	// 注册回调函数来处理HTML响应
	c.OnHTML("#phone > form > input[type=hidden]", func(e *colly.HTMLElement) {
		// 提取元素的value属性
		verify_info.Token = e.Attr("value")
	})
	// 处理错误
	c.OnError(func(_ *colly.Response, e error) {
		err = e
	})
	url := "http://" + domain + "/login"
	c.Visit(url)
	c.Wait()
	return verify_info, err
}

func VerifyUser(domain string, verify_info IdcsmartVerify, login_type string) (bool, error) {
	var url string
	var result bool
	var err error
	switch login_type {
	case "email":
		url = "http://" + domain + "/login?action=email"
		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Cookie", "PHPSESSID="+verify_info.Cookie)
		})
		// 魔方通过跳转的网页title判断是否登录成功
		c.OnHTML("head > title", func(e *colly.HTMLElement) {
			if strings.Contains(e.Text, "用户中心") {
				result = true
			} else {
				result = false
			}
		})
		// 处理错误
		c.OnError(func(_ *colly.Response, e error) {
			err = e
		})
		aes := AesBase64{key: []byte("idcsmart.finance"), iv: []byte("9311019310287172")}
		aes_password, err := aes.Encrypt([]byte(verify_info.Password))
		if err != nil {
			return false, err
		}
		c.Post(url, map[string]string{"token": verify_info.Token, "email": verify_info.Email, "password": aes_password})
		c.Wait()
	case "phone":
		url = "http://" + domain + "/login?action=phone"
		// 魔方通过跳转的网页title判断是否登录成功
		c.OnHTML("head > title", func(e *colly.HTMLElement) {
			if strings.Contains(e.Text, "用户中心") {
				result = true
			} else {
				result = false
			}
		})
		// 处理错误
		c.OnError(func(_ *colly.Response, e error) {
			err = e
		})
		aes := AesBase64{key: []byte("idcsmart.finance"), iv: []byte("9311019310287172")}
		aes_password, err := aes.Encrypt([]byte(verify_info.Password))
		if err != nil {
			return false, err
		}
		c.Post(url, map[string]string{"token": verify_info.Token, "phone": verify_info.Phone, "phone_code": "+" + strconv.Itoa(verify_info.PhoneCode), "password": aes_password})
		c.Wait()
	}
	return result, err
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var loginUrl = "https://v2ex.com/signin"
var missionUrl = "https://v2ex.com/mission/daily"

var username = "123"
var password = "123"

type myCookieJar struct {
	cookies []*http.Cookie
}

func (c *myCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if c.cookies == nil {
		c.cookies = make([]*http.Cookie, 0)
	}

	for _, it := range cookies {
		c.cookies = append(c.cookies, it)
	}
}

func (c *myCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return c.cookies
}

func main() {
	cookieJar := &myCookieJar{}
	client := http.Client{Jar: cookieJar}
	// login
	loginreq, err := http.NewRequest("GET", loginUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	loginresp, err := client.Do(loginreq)
	if err != nil {
		log.Fatal(err)
	}
	defer loginresp.Body.Close()
	loginHtml, err := ioutil.ReadAll(loginresp.Body)
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile(`<input type="hidden" value="(\d+)" name="once" />`)
	once := re.FindStringSubmatch(string(loginHtml))
	//log.Println(once)
	re = regexp.MustCompile(`<input type="text" class="sl" name="(\w+)"`)
	user := re.FindStringSubmatch(string(loginHtml))
	//log.Println(user)
	re = regexp.MustCompile(`<input type="password" class="sl" name="(\w+)"`)
	passwd := re.FindStringSubmatch(string(loginHtml))
	//log.Println(passwd)

	params := url.Values{
		"next":    {"/", "/"},
		user[1]:   {username},
		passwd[1]: {password},
		"once":    {once[1]}}
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Host", "v2ex.com")
	req.Header.Add("Origin", "http://v2ex.com")
	req.Header.Add("Referer", loginUrl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookieJar.cookies {
		req.AddCookie(c)
	}
	loginResult, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer loginResult.Body.Close()
	//loginResultS, err := ioutil.ReadAll(loginResult.Body)
	//if err != nil {
	//log.Fatal(err)
	//}
	//log.Printf("login result:\n%s\n", string(loginResultS))

	// try finish mission
	missionReq, err := http.NewRequest("GET", missionUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	missionResp, err := client.Do(missionReq)
	if err != nil {
		log.Fatal(err)
	}
	defer missionResp.Body.Close()
	missionHtml, err := ioutil.ReadAll(missionResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	missionHtmlStr := string(missionHtml)
	if strings.Index(missionHtmlStr, "每日登录奖励已领取") > 0 {
		fmt.Println("每日登录奖励已领取")
	} else {
		re = regexp.MustCompile(`location.href = '(.*)'`)
		matched := re.FindStringSubmatch(missionHtmlStr)
		req, _ = http.NewRequest("GET", "https://www.v2ex.com"+matched[1], nil)
		//log.Printf("%v\n", req)
		result, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer result.Body.Close()
		//resultS, err := ioutil.ReadAll(result.Body)
		//if err != nil {
		//log.Fatalln(err)
		//}
		//log.Printf("result:\n%s\n", string(resultS))
	}
}

package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type v2ex struct{}

func (v v2ex) checkin() error {
	const (
		loginUrl   = "https://v2ex.com/signin"
		missionUrl = "https://v2ex.com/mission/daily"
		username   = "123"
		password   = "123"
	)

	cookieJar := &myCookieJar{}
	client := http.Client{Jar: cookieJar}
	// login
	loginreq, err := http.NewRequest("GET", loginUrl, nil)
	if err != nil {
		return err
	}
	loginresp, err := client.Do(loginreq)
	if err != nil {
		return err
	}
	defer loginresp.Body.Close()
	loginHtml, err := ioutil.ReadAll(loginresp.Body)
	if err != nil {
		return err
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
		return err
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
		return err
	}
	defer loginResult.Body.Close()
	//loginResultS, err := ioutil.ReadAll(loginResult.Body)
	//if err != nil {
	//return err
	//}
	//log.Printf("login result:\n%s\n", string(loginResultS))

	// try finish mission
	missionReq, err := http.NewRequest("GET", missionUrl, nil)
	if err != nil {
		return err
	}
	missionResp, err := client.Do(missionReq)
	if err != nil {
		return err
	}
	defer missionResp.Body.Close()
	missionHtml, err := ioutil.ReadAll(missionResp.Body)
	if err != nil {
		return err
	}
	missionHtmlStr := string(missionHtml)
	if strings.Index(missionHtmlStr, "每日登录奖励已领取") > 0 {
		log.Println("v2ex: already checkin")
	} else {
		re = regexp.MustCompile(`location.href = '(.*)'`)
		matched := re.FindStringSubmatch(missionHtmlStr)
		req, _ = http.NewRequest("GET", "https://www.v2ex.com"+matched[1], nil)
		//log.Printf("%v\n", req)
		result, err := client.Do(req)
		if err != nil {
			return err
		}
		defer result.Body.Close()
		resultS, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return err
		}
		if strings.Index(string(resultS), "每日登录奖励已领取") > 0 {
			log.Println("v2ex: checkin success\n")
		} else {
			return errors.New("v2ex: checkin failed")
		}
	}
	return nil
}

func init() {
	registerMission(v2ex{})
}

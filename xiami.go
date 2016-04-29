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

type xiami struct{}

func (xm xiami) checkin() error {
	const (
		loginUrl   = "https://login.xiami.com/web/login"
		missionUrl = "http://www.xiami.com/web"
		email      = "123"
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

	params := url.Values{
		"email":       {email},
		"password":    {password},
		"LoginButton": {"登录"},
	}
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Host", "login.xiami.com")
	req.Header.Add("Origin", "https://login.xiami.com")
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
	//log.Printf("\n%s\n", string(loginResultS))

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
	if strings.Index(missionHtmlStr, "已连续签到") > 0 {
		log.Println("xiami: already checkin")
	} else {
		re := regexp.MustCompile(`<a class="check_in" href="/web/checkin/id/(\d+)">每日签到</a>`)
		matched := re.FindStringSubmatch(missionHtmlStr)
		req, _ = http.NewRequest("GET", "http://www.xiami.com/web/checkin/id/"+matched[1], nil)
		req.Header.Add("Referer", missionUrl)
		result, err := client.Do(req)
		if err != nil {
			return err
		}
		defer result.Body.Close()
		resultS, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return err
		}
		if strings.Index(string(resultS), "已连续签到") > 0 {
			log.Println("xiami: checkin success")
		} else {
			return errors.New("xiami: checkin failed")
		}
	}
	return nil
}

func init() {
	registerMission(xiami{})
}

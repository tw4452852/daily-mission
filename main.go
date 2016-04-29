package main

import (
	"log"
	"net/http"
	"net/url"
	"sync"
)

type missioner interface {
	checkin() error
}

var missions []missioner

func registerMission(m missioner) {
	missions = append(missions, m)
}

func unregisterMission(m missioner) {
	for i, mission := range missions {
		if mission == m {
			missions = append(missions[:i], missions[i+1:]...)
			return
		}
	}
}

func main() {
	var waiter sync.WaitGroup

	waiter.Add(len(missions))

	for _, mission := range missions {
		go func(mission missioner) {
			defer waiter.Done()

			if err := mission.checkin(); err != nil {
				log.Printf("mission[%#v] checkin error: %s\n", mission, err)
			}
		}(mission)
	}
	waiter.Wait()
}

type myCookieJar struct {
	cookies []*http.Cookie
}

func (c *myCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, it := range cookies {
		c.cookies = append(c.cookies, it)
	}
}

func (c *myCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return c.cookies
}

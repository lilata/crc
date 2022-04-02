package main

import (
	"io"
	"log"
	"net/http"
)

type httpSession struct {
	client *http.Client
	cookies []*http.Cookie
}
func (s *httpSession) setCookie(cookie *http.Cookie) {
	updated := false
	for idx, c := range s.cookies {
		if c.Name == cookie.Name {
			s.cookies[idx] = cookie
			updated = true
			break
		}
	}
	if !updated {
		s.cookies = append(s.cookies, cookie)
	}
}
func (s *httpSession) GetBody(url string) string {
	resp, err := s.MakeRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	return RespBody(resp)
}
func (s *httpSession) MakeRequest(method string, url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("User-Agent", DefaultUA)
	for _, cookie := range s.cookies {
		req.AddCookie(cookie)
	}
	resp, err = s.client.Do(req)
	for _, cookie := range resp.Cookies() {
		s.setCookie(cookie)
	}
	return
}
func NewHttpSession() *httpSession {
	return &httpSession{
		client:  &http.Client{},
		cookies: nil,
	}
}
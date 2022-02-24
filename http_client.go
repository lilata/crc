package main

import (
	"io"
	"net/http"
)

type httpSession struct {
	client *http.Client
	cookies []*http.Cookie
}

func (s *httpSession) MakeRequest(method string, url string, body io.Reader) (resp *http.Response, err error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("User-Agent", DefaultUA)
	for _, cookie := range s.cookies {
		req.AddCookie(cookie)
	}
	resp, err = s.client.Do(req)
	for _, cookie := range resp.Cookies() {
		s.cookies = append(s.cookies, cookie)
	}
	return
}
func NewHttpSession() *httpSession {
	return &httpSession{
		client:  &http.Client{},
		cookies: nil,
	}
}
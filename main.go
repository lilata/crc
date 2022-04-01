package main

import "time"

func main() {
	// searchKeywords("sm", "捆", "瘦", "SM", "虐")
	session := NewHttpSession()
	session.MakeRequest("GET", "https://cn.bing.com", nil)
	time.Sleep(time.Second)
	session.MakeRequest("GET", "https://cn.bing.com", nil)
	time.Sleep(time.Second)
	session.MakeRequest("GET", "https://cn.bing.com", nil)
	time.Sleep(time.Second)
	session.MakeRequest("GET", "https://cn.bing.com", nil)
}

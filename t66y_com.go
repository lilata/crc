package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)
func fetchRmdownMagnetURLInSession(hash string, session *httpSession) (string, error) {
	pageUrl := fmt.Sprintf("http://www.rmdown.com/link.php?hash=%s", hash)
	url := fmt.Sprintf("http://www.rmdown.com/download.php?action=magnet&ref=%s", hash)
	session.MakeRequest("GET", pageUrl, nil)
	resp, err := session.MakeRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body), nil
}
func fetchRmdownMagnetURL(hash string) (string, error) {
	session := NewHttpSession()
	return fetchRmdownMagnetURLInSession(hash, session)
}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func searchKeywords(keywords ...string) (urls []string) {
	baseUrl := "https://" + t66y_com_Hostname + "/thread0806.php?fid=25&search=&page=%s"
	c := colly.NewCollector(
		colly.Async(false),
		colly.MaxDepth(1),
		)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay: 13 * time.Second,
	})
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", DefaultUA)
		r.Headers.Set("Accept", "*/*")
	})
	c.OnHTML("tr.tr3.t_one.tac td.tal h3 a", func(e *colly.HTMLElement) {
		for _, kw := range keywords {
			if strings.Contains(e.Text, kw) {
				url := e.Attr("href")
				if !strings.HasPrefix(url, "http") {
					url = fmt.Sprintf("https://%s/%s", t66y_com_Hostname, url)
				}
				if !contains(urls, url) {
					urls = append(urls, url)
				}
			}
		}
	})
	for i := 1; i <= t66y_com_MaxPage; i++ {
		c.Visit(fmt.Sprintf(baseUrl, strconv.Itoa(i)))
	}
	return
}
package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
func getRmdownHash(url string) string {
	splitted := strings.Split(url, "=")
	return splitted[len(splitted) - 1]
}
func fetchRmdownMagnetURLInSession(hash string, session *httpSession) (string, error) {
	pageUrl := fmt.Sprintf("http://www.rmdown.com/link.php?hash=%s", hash)
	url := fmt.Sprintf("http://www.rmdown.com/download.php?action=magnet&ref=%s", hash)
	session.MakeRequest("GET", pageUrl, nil)
	resp, err := session.MakeRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	return RespBody(resp), nil
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

func t66ySearchKeywords(keywords ...string) {
	session := NewHttpSession()
	baseUrl := "https://" + t66y_com_Hostname + "/thread0806.php?fid=25&search=&page=%s"
	for i := 1; i <= t66y_com_MaxPage; i++ {
		url := fmt.Sprintf(baseUrl, strconv.Itoa(i))
		body := session.GetBody(url)
		parsedBody := soup.HTMLParse(body)
		tbody := parsedBody.Find("tbody", "id", "tbody")
		posts := tbody.FindAll("tr")
		for _, p := range posts {
			entry := NewDataEntry()
			tds := p.FindAll("td")
			titleTag := tds[1].Find("h3").Find("a")
			title := titleTag.Text()
			postUrl := fmt.Sprintf("https://%s/%s", t66y_com_Hostname, titleTag.Attrs()["href"])
			postHTML := session.GetBody(postUrl)
			post := soup.HTMLParse(postHTML)
			imgs := post.FindAll("img")
			if len(imgs) > 0 {
				mediaDir := fmt.Sprintf("media/%s", entry.uuid)
				_, err := os.Stat(mediaDir)
				if os.IsNotExist(err) {
					os.Mkdir(mediaDir, 0755)
				}
			}
			for idx, i := range imgs {
				src, ok := i.Attrs()["ess-data"]
				if !ok {
					continue
				}
				resp, err := http.Get(src)
				if err != nil {
					log.Println(err)
				}
				s := strings.Split(src, ".")
				ext := s[len(s) - 1]
				filename := fmt.Sprintf("%s.%s", strconv.Itoa(idx), ext)
				f, err := os.OpenFile(
					fmt.Sprintf("media/%s/%s", entry.uuid, filename),
					os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
				io.Copy(f, resp.Body)
				resp.Body.Close()
				entry.images = append(entry.images, filename)
				time.Sleep(time.Second)
			}
			downloadTags := post.FindAll("a")
			for _, d := range downloadTags {
				href, ok := d.Attrs()["href"]
				if !ok {
					continue
				}
				if strings.HasPrefix(href, "http://www.rmdown.com/link.php?hash=") {
					magnetUrl, err := fetchRmdownMagnetURLInSession(getRmdownHash(href), session)
					if err != nil {
						log.Println(err)
					}
					entry.downloadLink = magnetUrl
					break
				}
			}

			//if !contains(keywords, title) {
			//	continue
			//}
			dateTag := tds[2].Find("div").Find("span")
			dateStr := dateTag.Attrs()["title"]
			dateStr = strings.Split(dateStr, " - ")[1]
			date, _ := time.Parse("2006-01-02", dateStr)
			t := date.Unix()
			entry.title = title
			entry.timestamp = t
			entry.description = ""
			entry.saveToDatabase()
			time.Sleep(t66y_com_Interval * time.Second)
		}
		time.Sleep(t66y_com_Interval * time.Second)
	}
}
package main

import "net/http"

func main() {
	migrateDatabase()
	go t66ySearchKeywords("sm", "捆", "瘦", "SM", "虐")
	http.ListenAndServe(":1443", getRouter())

}

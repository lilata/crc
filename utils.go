package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func RespBody(resp *http.Response) string {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return string(body)
}

func migrateDatabase() {
	db := getDatabaseConn()
	defer db.Close()
	db.Exec(`CREATE TABLE IF NOT EXISTS data_entries 
                   (id SERIAL PRIMARY KEY, Uuid TEXT, Title TEXT, ts BIGINT, Description TEXT,
                    Images TEXT [], download_link TEXT)`)
}
package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"strconv"
	"strings"
)
type DataEntry struct {
	Uuid      string `json:"uuid"`
	Title       string `json:"title"`
	Timestamp   int64 `json:"timestamp"`
	Description  string `json:"description"`
	Images       []string `json:"images"`
	DownloadLink string `json:"download_link"`
}

func NewDataEntry() *DataEntry {
	return &DataEntry{
		Uuid: uuid.NewString(),
	}
}
func (e *DataEntry) PQImageArray () string {
	var imgs []string
	for _, img := range e.Images {
		imgs = append(imgs, fmt.Sprintf("\"%s\"", img))
	}
	joined := strings.Join(imgs, ",")
	return fmt.Sprintf("{%s}", joined)
}
func (e *DataEntry) saveToDatabase() {
	db := getDatabaseConn()
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO data_entries " +
		"(Uuid, Title, ts, Description, Images, download_link)" +
		" VALUES('%s', '%s', %s, '%s', '%s', '%s')",
		e.Uuid, e.Title,  strconv.FormatInt(e.Timestamp, 10),
		e.Description, e.PQImageArray(), e.DownloadLink)
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func hasTitle(title string) bool {
	r := false
	db := getDatabaseConn()
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT id FROM data_entries WHERE Title='%s'", title))
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		r = true
		break
	}
	return r
}

func getPage(page int) []byte {
	var entries []DataEntry
	var largestId int64
	db := getDatabaseConn()
	defer db.Close()
	idq := "SELECT id FROM data_entries ORDER BY id DESC LIMIT 1"
	rows, err := db.Query(idq)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		rows.Scan(&largestId)
		break
	}
	query := fmt.Sprintf("SELECT Uuid, Title, ts, Description, Images, download_link" +
		" FROM data_entries WHERE id <= %s ORDER BY id DESC",
		strconv.FormatInt(largestId - int64((page - 1) * 10), 10))
	rows, err = db.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var e DataEntry
		rows.Scan(&e.Uuid, &e.Title, &e.Timestamp, &e.Description, pq.Array(&e.Images), &e.DownloadLink)
		entries = append(entries, e)
	}
	j, err := json.Marshal(entries)
	if err != nil {
		log.Println(err)
	}
	return j
}
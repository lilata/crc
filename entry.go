package main

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)
type DataEntry struct {
	uuid string
	title string
	timestamp int64
	description string
	images []string
	downloadLink string
}

func NewDataEntry() *DataEntry {
	return &DataEntry{
		uuid: uuid.NewString(),
	}
}
func (e *DataEntry) PQImageArray () string {
	var imgs []string
	for _, img := range e.images {
		imgs = append(imgs, fmt.Sprintf("\"%s\"", img))
	}
	joined := strings.Join(imgs, ",")
	return fmt.Sprintf("'{%s}'", joined)
}
func (e *DataEntry) saveToDatabase() {
	db := getDatabaseConn()
	defer db.Close()
	b64downloadLink := base64.StdEncoding.EncodeToString([]byte(e.downloadLink))
	query := fmt.Sprintf("INSERT INTO data_entries VALUES('%s', '%s', %s, '%s', %s, '%s')",
		e.uuid, e.title,  strconv.FormatInt(e.timestamp, 10), e.description, e.PQImageArray(), b64downloadLink)
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}
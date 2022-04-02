package main
import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)
func getDatabaseConn() *sql.DB {
	connStr := "host=127.0.0.1 port=5432 user=crc password=crc dbname=crc sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}
	return db
}
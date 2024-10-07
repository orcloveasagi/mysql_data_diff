package util

import (
	"database/sql"
	"log"
)

func Query(db *sql.DB, sql string) (*sql.Rows, error) {
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("do query error: %v_%s", db, sql)
		return nil, err
	}
	return rows, nil
}

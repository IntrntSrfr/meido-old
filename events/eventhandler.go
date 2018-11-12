package events

import (
	"database/sql"
)

var (
	db *sql.DB
)

func Initialize(DB *sql.DB) {
	db = DB
}

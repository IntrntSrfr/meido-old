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

const (
	dColorRed   = 13107200
	dColorGreen = 51200
	dColorWhite = 16777215
)

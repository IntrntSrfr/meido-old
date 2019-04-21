package events

import (
	"database/sql"

	"go.uber.org/zap"
)

var (
	db     *sql.DB
	logger *zap.Logger
)

func Initialize(DB *sql.DB, Logger *zap.Logger) {
	db = DB
	logger = Logger
}

const (
	dColorRed    = 13107200
	dColorOrange = 15761746
	dColorLBlue  = 6410733
	dColorGreen  = 51200
	dColorWhite  = 16777215
)

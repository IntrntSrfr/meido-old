package models

import "time"

type Strikes struct {
	Uid        int
	Guildid    string
	Userid     string
	Reason     string
	Executorid string
	Tstamp     time.Time
}

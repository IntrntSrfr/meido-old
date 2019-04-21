package models

import "time"

type DiscordUser struct {
	Uid            int
	Userid         string
	Username       string
	Discriminator  string
	Xp             int
	Nextxpgaintime time.Time
	Xpexcluded     bool
	Reputation     int
	Cangivereptime time.Time
}

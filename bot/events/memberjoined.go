package events

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/intrntsrfr/meido/bot/models"

	"github.com/bwmarrin/discordgo"
)

func MemberJoinedHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	totalUsers++
	sqlstr := "INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

	stmt, err := db.Prepare(sqlstr)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if m.User.Bot {
		return
	}

	row := db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", m.User.ID)

	user := models.DiscordUser{}

	err = row.Scan(
		&user.Uid,
		&user.Userid,
		&user.Username,
		&user.Discriminator,
		&user.Xp,
		&user.Nextxpgaintime,
		&user.Xpexcluded,
		&user.Reputation,
		&user.Cangivereptime)

	currentTime := time.Now()

	if err != nil {
		logger.Error(err.Error())
		if err == sql.ErrNoRows {
			//var lastInsertID int
			_, err := stmt.Exec(m.User.ID, m.User.Username, m.User.Discriminator, 0, currentTime, false, 0, currentTime)
			if err != nil {
				logger.Error(err.Error())
				fmt.Println(err)
				return
			}
		}
	}
}

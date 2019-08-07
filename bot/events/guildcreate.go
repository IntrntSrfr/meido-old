package events

import (
	"database/sql"
	"fmt"

	"github.com/intrntsrfr/meido/bot/models"
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

func (eh *EventHandler) guildAvailableHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	/*
		sqlstr := "INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

		stmt, err := db.Prepare(sqlstr)
		if err != nil {
			return
		}

		loadTimeStart := time.Now()

		fmt.Println(g.Name)
		for i := range g.Members {
			m := g.Members[i]

			if m.User.Bot {
				continue
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
				if err == sql.ErrNoRows {
					//var lastInsertID int
					_, err := stmt.Exec(m.User.ID, m.User.Username, m.User.Discriminator, 0, currentTime, false, 0, currentTime)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
		}

		loadTimeEnd := time.Now()
		totalLoadTime := loadTimeEnd.Sub(loadTimeStart)
		fmt.Println(fmt.Sprintf("Loaded %v in %v", g.Name, totalLoadTime.String()))
	*/

	if len(g.Members) != g.MemberCount {
		s.RequestGuildMembers(g.ID, "", 0)
	}

	dbg := &models.DiscordGuild{}

	err := eh.db.Get(dbg, "SELECT guildid FROM discordguilds WHERE guildid = $1;", g.ID)
	if err != nil && err != sql.ErrNoRows {
		eh.logger.Error("error", zap.Error(err))
	} else if err == sql.ErrNoRows {
		eh.db.Exec("INSERT INTO discordguilds(guildid, usestrikes, maxstrikes) VALUES($1, $2, $3)", g.ID, false, 3)
		eh.logger.Info(fmt.Sprintf("Inserted new guild: %v [%v]", g.Name, g.ID))
	}

	eh.logger.Info(fmt.Sprintf("Loaded %v", g.Name))
	fmt.Println(fmt.Sprintf("Loaded %v", g.Name))
}

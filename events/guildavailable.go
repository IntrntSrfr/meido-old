package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GuildAvailableHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
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
	fmt.Println(fmt.Sprintf("Loaded %v", g.Name))
}

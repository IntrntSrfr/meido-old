package events

import (
	"github.com/bwmarrin/discordgo"
)

func (eh *EventHandler) guildMembersChunkHandler(s *discordgo.Session, m *discordgo.GuildMembersChunk) {
	/*
		for _, mem := range m.Members {

			if mem.User.Bot {
				return
			}

			row := eh.db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", mem.User.ID)

			user := models.DiscordUser{}

			err := row.Scan(
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
				eh.logger.Error(err.Error())
				if err == sql.ErrNoRows {
					//var lastInsertID int
					_, err := eh.db.Exec("INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)", mem.User.ID, mem.User.Username, mem.User.Discriminator, 0, currentTime, false, 0, currentTime)
					if err != nil {
						eh.logger.Error(err.Error())
						fmt.Println(err)
						return
					}
				}
			}
		}
	*/
}

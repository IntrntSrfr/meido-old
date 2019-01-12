package events

import (
	"fmt"
	"meido/models"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func MessageUpdateHandler(s *discordgo.Session, m *discordgo.MessageUpdate) {

	if m.Author == nil {
		return
	}

	if m.Author.Bot {
		return
	}

	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return
	}

	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	perms, err := s.State.UserChannelPermissions(m.Author.ID, ch.ID)
	if err != nil {
		return
	}

	isIllegal := false
	trigger := ""

	if perms&discordgo.PermissionManageMessages == 0 {

		var count int

		row := db.QueryRow("SELECT COUNT(*) FROM filterignorechannels WHERE channelid = $1;", m.ChannelID)
		err := row.Scan(&count)
		if err != nil {
			return
		}

		if count > 0 {
			return
		}

		rows, _ := db.Query("SELECT phrase FROM filters WHERE guildid = $1", m.GuildID)

		for rows.Next() {
			filter := models.Filter{}
			err := rows.Scan(&filter.Filter)
			if err != nil {
				continue
			}

			if strings.Contains(m.Content, filter.Filter) {
				trigger = filter.Filter
				isIllegal = true
				break
			}
		}

		if isIllegal {
			row := db.QueryRow("SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid = $1;", m.GuildID)

			dbg := models.DiscordGuild{}

			err := row.Scan(&dbg.UseStrikes, &dbg.MaxStrikes)
			if err != nil {
				return
			}

			if dbg.UseStrikes {

				reason := fmt.Sprintf("Triggering filter: %v", trigger)

				strikeCount := 0

				row := db.QueryRow("SELECT COUNT(*) FROM strikes WHERE guildid = $1 AND userid = $2;", m.Message.GuildID, m.Author.ID)
				err := row.Scan(&strikeCount)
				if err != nil {
					return
				}

				if strikeCount+1 >= dbg.MaxStrikes {
					//ban
					userch, _ := s.UserChannelCreate(m.Author.ID)
					s.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for acquiring %v strikes.\nLast warning was: %v", g.Name, dbg.MaxStrikes, reason))
					err = s.GuildBanCreateWithReason(m.Message.GuildID, m.Author.ID, fmt.Sprintf("Acquired %v strikes.", dbg.MaxStrikes), 0)
					if err != nil {
						return
					}

					s.ChannelMessageSend(ch.ID, fmt.Sprintf("%v has been banned after acquiring too many strikes. Miss them.", m.Author.Mention()))
					_, err := db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", m.Author.ID, g.ID)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					//insert warn
					userch, _ := s.UserChannelCreate(m.Author.ID)
					s.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been warned in %v.\nWarned for: %v\nYou are currently at strike %v/%v", g.Name, reason, strikeCount+1, dbg.MaxStrikes))
					s.ChannelMessageSend(ch.ID, fmt.Sprintf("%v has been warned\nThey are currently at strike %v/%v", m.Author.Mention(), strikeCount+1, dbg.MaxStrikes))
					db.Exec("INSERT INTO strikes(guildid, userid, reason, executorid, tstamp) VALUES ($1, $2, $3, $4, $5);", g.ID, m.Author.ID, reason, s.State.User.ID, time.Now())
				}
			} else {
				s.ChannelMessageDelete(ch.ID, m.ID)
				s.ChannelMessageSend(ch.ID, fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!", m.Author.Mention()))
			}
		}
	}

}

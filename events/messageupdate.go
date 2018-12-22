package events

import (
	"database/sql"
	"fmt"
	"meido-test/models"
	"strings"

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
				dbs := models.Strikes{}

				row := db.QueryRow("SELECT * FROM strikes WHERE guildid = $1 AND userid = $2;", m.GuildID, m.Author.ID)
				err := row.Scan(&dbs.Uid, &dbs.Guildid, &dbs.Userid, &dbs.Strikes)
				if err != nil {
					if err == sql.ErrNoRows {
						if dbg.MaxStrikes < 2 {
							s.ChannelMessageDelete(m.ChannelID, m.ID)
							userch, _ := s.UserChannelCreate(m.Author.ID)
							s.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for triggering the filter.\n- %v", g.Name, trigger))
							err = s.GuildBanCreateWithReason(m.GuildID, m.Author.ID, fmt.Sprintf("Triggering filter: %v", trigger), 0)
							if err != nil {
								s.ChannelMessageSend(ch.ID, err.Error())
								return
							}

							embed := &discordgo.MessageEmbed{
								Title:       "User banned",
								Description: "Filter triggered",
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:   "Username",
										Value:  fmt.Sprintf("%v", m.Author.Mention()),
										Inline: true,
									},
									{
										Name:   "ID",
										Value:  fmt.Sprintf("%v", m.Author.ID),
										Inline: true,
									},
								},
								Color: dColorRed,
							}

							s.ChannelMessageSendEmbed(ch.ID, embed)

						} else {
							s.ChannelMessageDelete(ch.ID, m.ID)
							s.ChannelMessageSend(ch.ID, fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!\nYou are currently at strike %v/%v", m.Author.Mention(), dbs.Strikes+1, dbg.MaxStrikes))
							db.Exec("INSERT INTO strikes(guildid, userid, strikes) VALUES ($1, $2, $3);", g.ID, m.Author.ID, 1)
						}
					}
				} else {
					if dbs.Strikes+1 >= dbg.MaxStrikes {
						s.ChannelMessageDelete(m.ChannelID, m.ID)
						userch, _ := s.UserChannelCreate(m.Author.ID)
						s.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for triggering the filter.\n- %v", g.Name, trigger))
						err = s.GuildBanCreateWithReason(g.ID, m.Author.ID, fmt.Sprintf("Triggering filter: %v", trigger), 0)
						if err != nil {
							s.ChannelMessageSend(ch.ID, err.Error())
							return
						}

						embed := &discordgo.MessageEmbed{
							Title:       "User banned",
							Description: "Filter triggered",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Username",
									Value:  fmt.Sprintf("%v", m.Author.Mention()),
									Inline: true,
								},
								{
									Name:   "ID",
									Value:  fmt.Sprintf("%v", m.Author.ID),
									Inline: true,
								},
							},
							Color: dColorRed,
						}

						s.ChannelMessageSendEmbed(ch.ID, embed)

						_, err := db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", m.Author.ID, g.ID)
						if err != nil {
							fmt.Println(err)
						}

					} else {
						s.ChannelMessageDelete(ch.ID, m.ID)
						s.ChannelMessageSend(ch.ID, fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!\nYou are currently at strike %v/%v", m.Author.Mention(), dbs.Strikes+1, dbg.MaxStrikes))
						db.Exec("UPDATE strikes SET strikes = $1 WHERE userid = $2 AND guildid = $3;", dbs.Strikes+1, m.Author.ID, g.ID)
					}
				}

			} else {
				s.ChannelMessageDelete(ch.ID, m.ID)
				s.ChannelMessageSend(ch.ID, fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!", m.Author.Mention()))
			}
		}
	}

}

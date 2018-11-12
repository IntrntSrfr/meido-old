package events

import (
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

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return
	}

	perms, err := s.UserChannelPermissions(m.Author.ID, ch.ID)
	if err != nil {
		return
	}

	if perms&discordgo.PermissionManageMessages == 0 {

		rows, _ := db.Query("SELECT phrase FROM filters WHERE guildid = $1", ch.GuildID)

		isIllegal := false

		for rows.Next() {
			filter := models.Filter{}
			err := rows.Scan(&filter.Filter)
			if err != nil {
				continue
			}

			if strings.Contains(m.Content, filter.Filter) {
				isIllegal = true
				break
			}
		}

		if isIllegal {
			s.ChannelMessageDelete(ch.ID, m.ID)
			s.ChannelMessageSend(ch.ID, fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!", m.Author.Mention()))
		}
	}
}

package events

import (
	"time"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (eh *EventHandler) messageUpdateHandler(s *discordgo.Session, m *discordgo.MessageUpdate) {

	if m.Author == nil {
		return
	}

	if m.Author.Bot {
		return
	}

	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		eh.logger.Error(err.Error())
		return
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return
	}

	startTime := time.Now()

	context, err := service.NewContext(s, m.Message, startTime)
	if err != nil {
		return
	}

	perms, err := s.State.UserChannelPermissions(m.Author.ID, ch.ID)
	if err != nil {
		eh.logger.Error(err.Error())
		return
	}

	isIllegal := eh.checkFilter(&context, &perms, m.Message)
	if isIllegal {
		return
	}
}

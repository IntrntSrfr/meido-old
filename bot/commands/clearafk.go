package commands

import (
	"fmt"

	"github.com/intrntsrfr/meido/bot/service"
)

func (ch *CommandHandler) clearAFK(args []string, ctx *service.Context) {
	if ctx.Guild.AfkChannelID == "" {
		ctx.Send("There is no AFK channel")
		return
	}

	memberList := []string{}

	for _, val := range ctx.Guild.VoiceStates {
		if val.SelfMute && val.SelfDeaf && val.ChannelID != ctx.Guild.AfkChannelID {
			memberList = append(memberList, val.UserID)
		}
	}

	if len(memberList) < 1 {
		ctx.Send("There is no one to move.")
		return
	}

	for _, val := range memberList {
		ctx.Session.GuildMemberMove(ctx.Guild.ID, val, ctx.Guild.AfkChannelID)
	}

	ctx.Send(fmt.Sprintf("Moved %v users.", len(memberList)))
}

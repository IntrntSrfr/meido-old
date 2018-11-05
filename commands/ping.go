package commands

import (
	"meido-test/service"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Ping = Command{
	Name:          "ping",
	Description:   "Shows delay.",
	Aliases:       []string{},
	Usage:         "m?ping",
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {
		sendTime := time.Now()

		msg, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Pong")
		if err != nil {
			return
		}

		receiveTime := time.Now()

		delay := receiveTime.Sub(sendTime)

		ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, "Pong - "+delay.String())
	},
}

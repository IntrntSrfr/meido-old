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
	Function: func(context *service.Context) {
		sendTime := time.Now()

		msg, err := context.Session.ChannelMessageSend(context.Message.ChannelID, "Pong")
		if err != nil {
			return
		}

		receiveTime := time.Now()

		delay := receiveTime.Sub(sendTime)

		context.Session.ChannelMessageEdit(context.Message.ChannelID, msg.ID, "Pong - "+delay.String())
	},
}

package commands

import (
	"fmt"
	"time"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var Ping = Command{
	Name:          "Ping",
	Description:   "Displays bot latency.",
	Triggers:      []string{"m?ping"},
	Usage:         "m?ping",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {
		sendTime := time.Now()

		msg, err := ctx.Send("Pong")
		if err != nil {
			return
		}

		receiveTime := time.Now()

		botLatency := receiveTime.Sub(ctx.StartTime)
		latency := receiveTime.Sub(sendTime)

		ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong!\nDiscord delay: %v\nBot delay: %v", latency.String(), botLatency.String()))
	},
}

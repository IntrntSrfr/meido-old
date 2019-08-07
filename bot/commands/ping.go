package commands

import (
	"fmt"
	"time"

	"github.com/intrntsrfr/meido/bot/service"
)

func (ch *CommandHandler) ping(args []string, ctx *service.Context) {
	sendTime := time.Now()

	msg, err := ctx.Send("Pong")
	if err != nil {
		return
	}

	receiveTime := time.Now()

	botLatency := receiveTime.Sub(ctx.StartTime)
	latency := receiveTime.Sub(sendTime)

	ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, msg.ID, fmt.Sprintf("Pong!\nDiscord delay: %v\nBot delay: %v", latency.String(), botLatency.String()))
}

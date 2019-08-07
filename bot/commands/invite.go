package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"
)

func (ch *CommandHandler) invite(args []string, ctx *service.Context) {
	botLink := "https://discordapp.com/oauth2/authorize?client_id=394162399348785152&scope=bot"
	serverLink := "https://discord.gg/KgMEGK3"
	ctx.Send(fmt.Sprintf("Invite me to your server: %v\nSupport server: %v", botLink, serverLink))
}

func (ch *CommandHandler) feedback(args []string, ctx *service.Context) {
	if len(args) < 2 {
		return
	}

	text := fmt.Sprintf("Message from %v - %v (%v) from channel %v (%v) in server %v (%v)\n", ctx.User.Mention(), ctx.User.String(), ctx.User.ID, ctx.Channel.Name, ctx.Channel.ID, ctx.Guild.Name, ctx.Guild.ID)
	text += fmt.Sprintf("`%v`", (strings.Join(args[1:], " ")))
	_, err := ctx.Session.ChannelMessageSend("533009623188242443", text)
	if err != nil {
		ch.logger.Error("error", zap.Error(err))
		ctx.Send(fmt.Sprintf("error: %v", err.Error()))
		return
	}
	ctx.Send("Feedback left")
}

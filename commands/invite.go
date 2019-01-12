package commands

import (
	"fmt"
	"meido/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Invite = Command{
	Name:          "invite",
	Description:   "Sends bot invite link and support server invite.",
	Triggers:      []string{"m?invite"},
	Usage:         "m?invite",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: false,
	Execute: func(args []string, ctx *service.Context) {
		botLink := "https://discordapp.com/oauth2/authorize?client_id=394162399348785152&scope=bot"
		serverLink := "https://discord.gg/KgMEGK3"
		ctx.Send(fmt.Sprintf("Invite me to your server: %v\nSupport server: %v", botLink, serverLink))
	},
}

var Feedback = Command{
	Name:          "feedback",
	Description:   "sends your very nice and helpful feedback to the Meido Café.",
	Triggers:      []string{"m?feedback"},
	Usage:         "m?feedback wow what a really COOL and NICE bot that works flawlessly",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: false,
	Execute: func(args []string, ctx *service.Context) {
		if len(args) > 1 {
			text := fmt.Sprintf("Message from %v - %v (%v) from channel %v (%v) in server %v (%v)\n", ctx.User.Mention(), ctx.User.String(), ctx.User.ID, ctx.Channel.Name, ctx.Channel.ID, ctx.Guild.Name, ctx.Guild.ID)
			text += fmt.Sprintf("`%v`", (strings.Join(args[1:], " ")))
			_, err := ctx.Session.ChannelMessageSend("533009623188242443", text)
			if err != nil {
				fmt.Println(err)
			}
			ctx.Send("Feedback left")
		}
	},
}

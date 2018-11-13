package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Test = Command{
	Name:          "test",
	Description:   "does epic testing.",
	Triggers:      []string{"m?test"},
	Usage:         "m?test",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		ctx.Send("epic")

	},
}

var Dm = Command{
	Name:          "dm",
	Description:   "Sends a direct message. Owner only.",
	Triggers:      []string{"m?dm"},
	Usage:         "m?dm 163454407999094786 jeff",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		userch, err := ctx.Session.UserChannelCreate(args[1])

		if err != nil {
			return
		}

		ctx.Session.ChannelMessageSend(userch.ID, strings.Join(args[2:], " "))
		ctx.Send(fmt.Sprintf("Message sent to %v", userch.Recipients[0]))

	},
}

var Msg = Command{
	Name:          "msg",
	Description:   "Sends a message to a channel. Owner only.",
	Triggers:      []string{"m?msg"},
	Usage:         "m?msg 497106582144942101 jeff",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		var ch string

		if strings.HasPrefix(args[1], "<#") && strings.HasSuffix(args[1], ">") {
			ch = args[1]
			ch = ch[2 : len(ch)-1]
		} else {
			ch = args[1]
		}

		chn, err := ctx.Session.Channel(ch)
		if err != nil {
			return
		}

		ctx.Session.ChannelMessageSend(chn.ID, strings.Join(args[2:], " "))
		ctx.Send(fmt.Sprintf("Message sent to %v", chn.Name))

	},
}

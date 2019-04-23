package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var Test = Command{
	Name:          "Test",
	Description:   "Does epic testing.",
	Triggers:      []string{"m?test"},
	Usage:         "m?test",
	RequiredPerms: discordgo.PermissionSendMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {
		ctx.Send(fmt.Sprintf("Top role position: %v", HighestRole(ctx.Guild, ctx.User.ID)))
		ctx.Send(fmt.Sprintf("Top color: #" + FullHex(fmt.Sprintf("%X", UserColor(ctx.Guild, ctx.User.ID)))))
	},
}

var Dm = Command{
	Name:          "DM",
	Description:   "Sends a direct message. Owner only.",
	Triggers:      []string{"m?dm"},
	Usage:         "m?dm 163454407999094786 jeff",
	Category:      Owner,
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) < 3 {
			ctx.Send("no")
			return
		}

		userch, err := ctx.Session.UserChannelCreate(args[1])

		if err != nil {
			return
		}

		if userch.Type != discordgo.ChannelTypeDM {
			return
		}

		ctx.Session.ChannelMessageSend(userch.ID, strings.Join(args[2:], " "))
		ctx.Send(fmt.Sprintf("Message sent to %v", userch.Recipients[0]))
	},
}

var Msg = Command{
	Name:          "Msg",
	Description:   "Sends a message to a channel. Owner only.",
	Triggers:      []string{"m?msg"},
	Usage:         "m?msg 497106582144942101 jeff",
	Category:      Owner,
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) < 3 {
			ctx.Send("no")
			return
		}

		var ch string

		if strings.HasPrefix(args[1], "<#") && strings.HasSuffix(args[1], ">") {
			ch = args[1]
			ch = ch[2 : len(ch)-1]
		} else {
			ch = args[1]
		}

		chn, err := ctx.Session.State.Channel(ch)
		if err != nil {
			return
		}

		ctx.Session.ChannelMessageSend(chn.ID, strings.Join(args[2:], " "))
		ctx.Send(fmt.Sprintf("Message sent to %v [<#%v>]", chn.Name, chn.ID))
	},
}

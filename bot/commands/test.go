package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

/*
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
*/
func (ch *CommandHandler) dm(args []string, ctx *service.Context) {

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
}

func (ch *CommandHandler) msg(args []string, ctx *service.Context) {

	if len(args) < 3 {
		ctx.Send("no")
		return
	}

	var chID string

	if strings.HasPrefix(args[1], "<#") && strings.HasSuffix(args[1], ">") {
		chID = args[1]
		chID = chID[2 : len(chID)-1]
	} else {
		chID = args[1]
	}

	chn, err := ctx.Session.State.Channel(chID)
	if err != nil {
		return
	}

	ctx.Session.ChannelMessageSend(chn.ID, strings.Join(args[2:], " "))
	ctx.Send(fmt.Sprintf("Message sent to %v [<#%v>]", chn.Name, chn.ID))
}

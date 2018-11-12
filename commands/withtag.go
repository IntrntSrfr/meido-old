package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var WithTag = Command{
	Name:          "withtag",
	Description:   "Shows how many has an input discriminator.",
	Triggers:      []string{"m?withtag"},
	Usage:         "m?withtag <0001/#0001>",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		var (
			tag           string
			matchingUsers int
		)

		if len(args) != 2 {
			return
		}

		tag = args[1]

		if strings.HasPrefix(tag, "#") {
			tag = strings.TrimPrefix(tag, "#")
		}

		if len(tag) != 4 {
			ctx.Send("Invalid tag")
			return
		}

		for i := 0; i < len(ctx.Guild.Members); i++ {
			member := ctx.Guild.Members[i]

			if strings.ToLower(member.User.Discriminator) == strings.ToLower(tag) {
				matchingUsers++
			}
		}

		ctx.Send(fmt.Sprintf("Users with tag %v: %v", tag, matchingUsers))
	},
}

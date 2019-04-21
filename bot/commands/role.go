package commands

import (
	"fmt"
	"sort"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var ListRoles = Command{
	Name:          "List roles",
	Description:   "Lists all roles in the server with the respective amount of members in each role.",
	Triggers:      []string{"m?listroles"},
	Usage:         "m?listroles",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		roles := ctx.Guild.Roles

		sort.Sort(discordgo.Roles(roles))

		text := fmt.Sprintf("Roles in %v [%v]\n\n\n", ctx.Guild.Name, len(roles))
		var count int
		for _, role := range ctx.Guild.Roles {
			for _, mem := range ctx.Guild.Members {
				count = 0
				for _, r := range mem.Roles {
					if r == role.ID {
						count++
					}
				}
			}
			text += fmt.Sprintf("%v [%v]\n", role.Name, count)
		}

		link, err := OWOApi.Upload(text)
		if err != nil {
			ctx.Send("Error getting role list")
			return
		}
		ctx.Send(fmt.Sprintf("Rolecount in %v - [%v]\n%v", ctx.Guild.Name, len(roles), link))
	},
}

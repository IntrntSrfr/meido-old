package commands

import (
	"fmt"
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Hackban = Command{
	Name:          "hackban",
	Description:   "Hackbans one or several users. Prunes 14 days.",
	Triggers:      []string{"m?hackban", "m?hb"},
	Usage:         "m?hb 123 123 12 31 23 123 ",
	RequiredPerms: discordgo.PermissionBanMembers,
	Execute: func(args []string, ctx *service.Context) {
		if len(args) < 2 {
			return
		}

		userList := []string{}

		for _, mention := range ctx.Message.Mentions {
			userList = append(userList, mention.ID)
		}

		for _, userID := range args[1:] {
			userList = append(userList, userID)
		}

		badbans := 0

		for _, userID := range userList {
			err := ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, userID, fmt.Sprintf("[%v] - Hackban", ctx.User.String()), 14)
			if err != nil {
				badbans++
			}
		}

		ctx.Send(fmt.Sprintf("Banned %v out of %v users provided.", len(userList)-badbans, len(userList)))
	},
}

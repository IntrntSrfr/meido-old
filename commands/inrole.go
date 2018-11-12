package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Inrole = Command{
	Name:          "inrole",
	Description:   "does epic testing.",
	Triggers:      []string{"m?inrole"},
	Usage:         "m?inrole gamers",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {
		if len(args) < 2 {
			return
		}

		var selectedRole *discordgo.Role

		for _, val := range ctx.Guild.Roles {
			if args[1] == val.ID {
				selectedRole = val
				break
			}
			if strings.ToLower(strings.Join(args[1:], " ")) == strings.ToLower(val.Name) {
				selectedRole = val
				break
			}
		}
		if selectedRole == nil {
			ctx.Send("Could not find that role")
			return
		}

		var memberList []*discordgo.Member

		for _, u := range ctx.Guild.Members {
			if len(memberList) > 20 {
				break
			}
			for _, r := range u.Roles {
				if selectedRole.ID == r {
					memberList = append(memberList, u)
				}
			}
		}

		lenlist := 0

		if len(memberList) > 20 {
			lenlist = 20
		} else {
			lenlist = len(memberList)
		}

		board := fmt.Sprintf("```\nFirst %v users in **%v**\n", lenlist, selectedRole.Name)
		for i := 0; i < lenlist; i++ {
			m := memberList[i]
			board += fmt.Sprintf("- %v", m.User.String())
		}
		board += "```"
		ctx.Send(board)
	},
}

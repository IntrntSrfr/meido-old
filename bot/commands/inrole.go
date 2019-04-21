package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var Inrole = Command{
	Name:          "Inrole",
	Description:   "Shows a list of who and how many users who are in a specified role.",
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
			for _, r := range u.Roles {
				if selectedRole.ID == r {
					memberList = append(memberList, u)
				}
			}
		}

		if len(memberList) <= 0 {
			ctx.Send("No users in that role.")
			return
		}

		board := fmt.Sprintf("Total users in role %v: %v\n", selectedRole.Name, len(memberList))
		if len(memberList) <= 20 {
			for i := 0; i < len(memberList); i++ {
				mem := memberList[i]

				board += fmt.Sprintf("%v\t[%v]\n", mem.User.String(), mem.User.ID)

			}

		} else if len(memberList) > 20 && len(memberList) < 1000 {
			text := fmt.Sprintf("Total users in role %v: %v\n\n\n", selectedRole.Name, len(memberList))
			for _, mem := range memberList {
				text += fmt.Sprintf("%v\t-%v\n", mem.User.String(), mem.User.ID)
			}

			res, err := OWOApi.Upload(text)
			if err != nil {
				board += "Error getting list."
			}
			board += res
		}
		ctx.Send(board)
	},
}

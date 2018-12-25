package commands

import (
	"fmt"
	"meido-test/service"
	"strings"
	"time"

	"github.com/ninedraft/simplepaste"

	"github.com/bwmarrin/discordgo"
)

var Inrole = Command{
	Name:          "inrole",
	Description:   "does epic testing.",
	Triggers:      []string{"m?inrole"},
	Usage:         "m?inrole gamers",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
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
		}

		board := fmt.Sprintf("Total users in role %v: %v\n", selectedRole.Name, len(memberList))
		if len(memberList) > 20 && len(memberList) < 5000 {
			text := ""
			for _, mem := range memberList {
				text += fmt.Sprintf("%v\t-%v\n", mem.User.String(), mem.User.ID)
			}
			paste := simplepaste.NewPaste(fmt.Sprintf("Users in role %v: %v", selectedRole.Name, time.Now().Format(time.RFC1123)), text)
			res, err := pbAPI.SendPaste(paste)
			if err != nil {
				board += "Error getting list."
			}
			board += res
		} else {
			board += "```\n"
			for _, mem := range memberList {
				board += fmt.Sprintf("%v\t(%v)\n", mem.User.String(), mem.User.ID)
			}
			board += "```"
		}
		ctx.Send(board)
	},
}

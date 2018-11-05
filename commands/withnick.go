package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var WithNick = Command{
	Name:          "withnick",
	Description:   "Shows how many has an input user- or nickname.",
	Aliases:       []string{},
	Usage:         "m?withnick meido",
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {

		var (
			newName       string
			matchingUsers int
			userList      []*discordgo.User
		)

		if len(args) <= 1 {
			return
		}

		newName = strings.Join(args[1:], " ")

		for i := 0; i < len(ctx.Guild.Members); i++ {
			member := ctx.Guild.Members[i]

			if strings.ToLower(member.Nick) == strings.ToLower(newName) {
				userList = append(userList, member.User)
				matchingUsers++
			} else if strings.ToLower(member.User.Username) == strings.ToLower(newName) {
				userList = append(userList, member.User)
				matchingUsers++
			}
		}

		if matchingUsers < 1 {
			ctx.Send("No users with that name")
			return
		}

		userBoard := "```\n"

		for i := 0; i < len(userList); i++ {
			u := userList[i]
			userBoard += fmt.Sprintf("%v\n", u.String())
		}
		userBoard += "```"

		ctx.Send(fmt.Sprintf("Total users with name %v: %v\nPreview:\n%v", newName, matchingUsers, userBoard))
	},
}

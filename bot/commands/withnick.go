package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var WithNick = Command{
	Name:          "With nick",
	Description:   "Shows how many has an input user- or nickname.",
	Triggers:      []string{"m?withnick"},
	Usage:         "m?withnick meido",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		var (
			newName    string
			memberList []*discordgo.Member
		)

		if len(args) <= 1 {
			return
		}

		newName = strings.Join(args[1:], " ")

		for _, mem := range ctx.Guild.Members {
			if strings.ToLower(mem.Nick) == strings.ToLower(newName) {
				memberList = append(memberList, mem)
			} else if strings.ToLower(mem.User.Username) == strings.ToLower(newName) {
				memberList = append(memberList, mem)
			}
		}

		if len(memberList) < 1 {
			ctx.Send("No users with that name")
			return
		}

		board := strings.Builder{}
		board.WriteString(fmt.Sprintf("Total users with name %v: %v\n", newName, len(memberList)))

		if len(memberList) <= 20 {
			for i := 0; i < len(memberList); i++ {
				mem := memberList[i]
				board.WriteString(fmt.Sprintf("%v\t[%v]\n", mem.User.String(), mem.User.ID))
			}

		} else if len(memberList) > 20 && len(memberList) < 1000 {
			text := strings.Builder{}
			text.WriteString(fmt.Sprintf("Total users with name %v: %v\n\n\n", newName, len(memberList)))
			for _, mem := range memberList {
				text.WriteString(fmt.Sprintf("%v\t-%v\n", mem.User.String(), mem.User.ID))
			}

			res, err := OWOApi.Upload(text.String())
			if err != nil {
				board.WriteString("Error getting list.")
			}
			board.WriteString(res)
		}
		ctx.Send(board.String())
	},
}

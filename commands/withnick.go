package commands

import (
	"fmt"
	"meido-test/service"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ninedraft/simplepaste"
)

var WithNick = Command{
	Name:          "withnick",
	Description:   "Shows how many has an input user- or nickname.",
	Triggers:      []string{"m?withnick"},
	Usage:         "m?withnick meido",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner:true,
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

		board := fmt.Sprintf("Total users with name %v: %v\n", newName, len(memberList))
		if len(memberList) > 20 && len(memberList) < 5000 {
			text := ""
			for _, mem := range memberList {
				text += fmt.Sprintf("%v\t%v\n", mem.User.String(), mem.User.ID)
			}
			paste := simplepaste.NewPaste(fmt.Sprintf("Users with name %v: %v", newName, time.Now().Format(time.RFC1123)), text)
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

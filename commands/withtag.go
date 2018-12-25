package commands

import (
	"fmt"
	"meido-test/service"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ninedraft/simplepaste"
)

var WithTag = Command{
	Name:          "withtag",
	Description:   "Shows how many has an input discriminator.",
	Triggers:      []string{"m?withtag"},
	Usage:         "m?withtag <0001/#0001>",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		var (
			tag        string
			memberList []*discordgo.Member
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

		for _, mem := range ctx.Guild.Members {
			if strings.ToLower(mem.User.Discriminator) == strings.ToLower(tag) {
				memberList = append(memberList, mem)
			}
		}

		if len(memberList) < 1 {
			ctx.Send("No users with that discriminator")
			return
		}

		board := fmt.Sprintf("Total users with discriminator %v: %v\n", tag, len(memberList))
		if len(memberList) > 20 && len(memberList) < 5000 {
			text := ""
			for _, mem := range memberList {
				text += fmt.Sprintf("%v\t%v\n", mem.User.String(), mem.User.ID)
			}
			paste := simplepaste.NewPaste(fmt.Sprintf("Users with discriminator %v: %v", tag, time.Now().Format(time.RFC1123)), text)
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

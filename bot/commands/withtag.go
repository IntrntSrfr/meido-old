package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) withTag(args []string, ctx *service.Context) {

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

	board := fmt.Sprintf("Total users with discriminator #%v: %v\n", tag, len(memberList))
	if len(memberList) <= 20 {
		for i := 0; i < len(memberList); i++ {
			mem := memberList[i]

			board += fmt.Sprintf("%v\t[%v]\n", mem.User.String(), mem.User.ID)

		}

	} else if len(memberList) > 20 && len(memberList) < 1000 {
		text := fmt.Sprintf("Total users with discriminator #%v: %v\n\n\n", tag, len(memberList))
		for _, mem := range memberList {
			text += fmt.Sprintf("%v\t-%v\n", mem.User.String(), mem.User.ID)
		}

		res, err := ch.owo.Upload(text)
		if err != nil {
			board += "Error getting list."
		}
		board += res
	}
	ctx.Send(board)
}

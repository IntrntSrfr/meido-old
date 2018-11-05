package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var CoolNameBro = Command{
	Name:          "coolnamebro",
	Description:   "Renames attentionseeking nick- or usernames.",
	Aliases:       []string{"cnb"},
	Usage:         "m?coolnamebro my name is shit",
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) < 2 {
			ctx.Send("Please choose a proper name.")
			return
		}

		newName := strings.Join(args[1:], " ")

		memberList := []string{}

		for _, val := range ctx.Guild.Members {
			if badName(val) {
				memberList = append(memberList, val.User.ID)
			}
		}

		if len(memberList) < 1 {
			ctx.Send("There is no one rename.")
			return
		} else {
			ctx.Send(fmt.Sprintf("Starting rename of %v user(s).", len(memberList)))
		}

		var successfulRenames, failedRenames int

		for _, val := range memberList {
			err := ctx.Session.GuildMemberNickname(ctx.Guild.ID, val, newName)
			if err != nil {
				failedRenames++
			} else {
				successfulRenames++
			}
		}

		ctx.Send(fmt.Sprintf("Rename finished. Successful: %v. Failed: %v.", successfulRenames, failedRenames))
	},
}

func badName(u *discordgo.Member) bool {

	if u.Nick != "" {
		if u.Nick[0] < 46 {
			return true
		}
	}
	if u.User.Username[0] < 46 {
		return true
	}
	return false

}

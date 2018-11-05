package commands

import (
	"fmt"
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Test = Command{
	Name:          "test",
	Description:   "does epic testing.",
	Aliases:       []string{},
	Usage:         "m?test",
	RequiredPerms: discordgo.PermissionVoiceMoveMembers,
	Execute: func(args []string, ctx *service.Context) {
		ctx.Send(fmt.Sprintf("https://cdn.discordapp.com/icons/%v/%v.png", ctx.Guild.ID, ctx.Guild.Icon))
	},
}

package commands

import (
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Test = Command{
	Name:          "test",
	Description:   "does epic testing.",
	Triggers:      []string{"m?test"},
	Usage:         "m?test",
	RequiredPerms: discordgo.PermissionSendMessages,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		ctx.Send("epic")

	},
}

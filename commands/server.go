package commands

import (
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Server = Command{
	Name:          "server",
	Description:   "Shows information about the current server.",
	Aliases:       []string{"serverinfo", "sa"},
	Usage:         "m?withtag <0001/#0001>",
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {

	},
}

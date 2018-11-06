package commands

import (
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Inrole = Command{
	Name:          "test",
	Description:   "does epic testing.",
	Aliases:       []string{},
	Usage:         "m?test",
	RequiredPerms: discordgo.PermissionVoiceMoveMembers,
	Execute: func(args []string, ctx *service.Context) {

	},
}

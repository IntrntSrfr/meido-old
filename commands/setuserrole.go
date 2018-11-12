package commands

import (
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var SetUserRole = Command{
	Name:          "setuserrole",
	Description:   "Sets a users custom role. First provide the user, followed by the role.",
	Triggers:      []string{"m?setuserrole"},
	Usage:         "m?setuserrole 163454407999094786 kumiko",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

	},
}

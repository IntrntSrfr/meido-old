package commands

import (
	"fmt"

	"github.com/intrntsrfr/meido/service"

	"github.com/bwmarrin/discordgo"
)

var Unban = Command{
	Name:          "Unban",
	Description:   "Unbans a user.",
	Triggers:      []string{"m?unban", "m?ub", ".ub", ".unban"},
	Usage:         ".unban 163454407999094786",
	RequiredPerms: discordgo.PermissionBanMembers,
	Execute: func(args []string, ctx *service.Context) {
		if len(args) <= 1 {
			return
		}

		userID := args[1]

		err := ctx.Session.GuildBanDelete(ctx.Guild.ID, userID)
		if err != nil {
			return
		}

		targetUser, err := ctx.Session.User(userID)
		if err != nil {
			return
		}

		embed := &discordgo.MessageEmbed{
			Description: fmt.Sprintf("**Unbanned** %v - %v#%v (%v)", targetUser.Mention(), targetUser.Username, targetUser.Discriminator, targetUser.ID),
			Color:       dColorGreen,
		}

		ctx.SendEmbed(embed)
	},
}

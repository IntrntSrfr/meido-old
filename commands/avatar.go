package commands

import (
	"fmt"
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Avatar = Command{
	Name:          "avatar",
	Description:   "Displays a users profile picture.",
	Triggers:      []string{"m?avatar", ">av", "m?av"},
	Usage:         ">av\n>av @internet surfer#0001\n>av 163454407999094786",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		var targetUser *discordgo.User
		var err error

		if len(args) > 1 {

			if len(ctx.Message.Mentions) >= 1 {
				targetUser = ctx.Message.Mentions[0]
			} else {
				targetUser, err = ctx.Session.User(args[1])
				if err != nil {
					return
				}
			}
		}

		if targetUser == nil {
			targetUser = ctx.Message.Author
		}

		if targetUser.Avatar == "" {
			ctx.SendEmbed(&discordgo.MessageEmbed{
				Color:       dColorRed,
				Description: fmt.Sprintf("%v has no avatar set.", targetUser.String()),
			})
		} else {
			ctx.SendEmbed(&discordgo.MessageEmbed{
				Color: dColorGreen,
				Title: targetUser.String(),
				Image: &discordgo.MessageEmbedImage{URL: targetUser.AvatarURL("1024")},
			})
		}
	},
}

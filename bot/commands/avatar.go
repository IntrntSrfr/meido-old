package commands

import (
	"fmt"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) avatar(args []string, ctx *service.Context) {

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
		return
	}
	/*
		mem, err := ctx.Session.GuildMember(ctx.Guild.ID, targetUser.ID)
		if err != nil {
			return
		} */

	if targetUser.Avatar == "" {
		ctx.SendEmbed(&discordgo.MessageEmbed{
			Color:       dColorRed,
			Description: fmt.Sprintf("%v has no avatar set.", targetUser.String()),
		})
	} else {
		ctx.SendEmbed(&discordgo.MessageEmbed{
			Color: ctx.Session.State.UserColor(targetUser.ID, ctx.Channel.ID),
			Title: targetUser.String(),
			Image: &discordgo.MessageEmbedImage{URL: targetUser.AvatarURL("1024")},
		})
	}
}

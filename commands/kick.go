package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Kick = Command{
	Name:          "kick",
	Description:   "kick a user. Reason and prune days is optional.",
	Triggers:      []string{"m?kick", "m?k", ".k"},
	Usage:         "m?k @internet surfer#0001\n.k 163454407999094786",
	RequiredPerms: discordgo.PermissionKickMembers,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) <= 1 {
			return
		}

		var targetUser *discordgo.User
		var err error

		reason := strings.Join(args[2:], " ")

		if len(ctx.Message.Mentions) >= 1 {
			targetUser = ctx.Message.Mentions[0]
		} else {
			targetUser, err = ctx.Session.User(args[1])
			if err != nil {
				return
			}
		}

		if targetUser.ID == ctx.Message.Author.ID {
			ctx.Send("no")
			return
		}

		if HighestRole(ctx.Guild, ctx.User.ID) <= HighestRole(ctx.Guild, targetUser.ID) {
			ctx.Send("no")
			return
		}

		okCh := true

		userchannel, err := ctx.Session.UserChannelCreate(targetUser.ID)
		if err != nil {
			okCh = false
		}

		if okCh {
			if reason == "" {
				ctx.Session.ChannelMessageSend(userchannel.ID, fmt.Sprintf("You have been kicked from %v.", ctx.Guild.Name))

			} else {
				ctx.Session.ChannelMessageSend(userchannel.ID, fmt.Sprintf("You have been kicked from %v for the following reason: %v", ctx.Guild.Name, reason))
			}
		}

		err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, targetUser.ID, fmt.Sprintf("%v#%v - %v", ctx.Message.Author.Username, ctx.Message.Author.Discriminator, reason))
		if err != nil {
			ctx.Send(err.Error())
			return
		}

		embed := &discordgo.MessageEmbed{
			Title: "User kicked",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Username",
					Value:  fmt.Sprintf("%v", targetUser.Mention()),
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  fmt.Sprintf("%v", targetUser.ID),
					Inline: true,
				},
			},
			Color: dColorRed,
		}

		ctx.SendEmbed(embed)

	},
}

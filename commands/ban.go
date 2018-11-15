package commands

import (
	"fmt"
	"meido-test/service"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Ban = Command{
	Name:          "ban",
	Description:   "Bans a user. Reason and prune days is optional.",
	Triggers:      []string{"m?ban", "m?b", ".b"},
	Usage:         ".b @internet surfer#0001\n.b 163454407999094786\n.b 163454407999094786 being very mean\n.b 163454407999094786 1 being very mean\n.b 163454407999094786 1",
	RequiredPerms: discordgo.PermissionBanMembers,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) <= 1 {
			return
		}

		var targetUser *discordgo.User
		var reason string
		var pruneDays int
		var err error

		if len(args) >= 3 {
			pruneDays, err = strconv.Atoi(args[2])
			if err != nil {
				pruneDays = 0
				reason = strings.Join(args[2:], " ")
			} else {
				reason = strings.Join(args[3:], " ")
			}
			if pruneDays > 7 {
				pruneDays = 7
			}
		}

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

		currentMem, err := ctx.Session.GuildMember(ctx.Guild.ID, ctx.User.ID)
		if err != nil {
			return
		}

		targetMem, err := ctx.Session.GuildMember(ctx.Guild.ID, targetUser.ID)
		if err != nil {
			return
		}

		if HighestRole(ctx.Guild, currentMem) <= HighestRole(ctx.Guild, targetMem) {
			ctx.Send("no")
			return
		}

		userchannel, err := ctx.Session.UserChannelCreate(targetUser.ID)
		if err != nil {
			return
		}

		if reason == "" {
			ctx.Session.ChannelMessageSend(userchannel.ID, fmt.Sprintf("You have been banned from %v.", ctx.Guild.Name))

		} else {
			ctx.Session.ChannelMessageSend(userchannel.ID, fmt.Sprintf("You have been banned from %v for the following reason: %v", ctx.Guild.Name, reason))
		}

		err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, targetUser.ID, fmt.Sprintf("%v#%v - %v", ctx.Message.Author.Username, ctx.Message.Author.Discriminator, reason), pruneDays)
		if err != nil {
			ctx.Send(err.Error())
			return
		}

		embed := &discordgo.MessageEmbed{
			Title: "User banned",
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

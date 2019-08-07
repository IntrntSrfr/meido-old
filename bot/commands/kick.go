package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) kick(args []string, ctx *service.Context) {

	if len(args) <= 1 {
		return
	}

	var targetUser *discordgo.Member
	var err error

	reason := ""
	if len(args) > 2 {
		reason = strings.Join(args[2:], " ")
	}

	if len(ctx.Message.Mentions) >= 1 {
		targetUser, err = ctx.Session.State.Member(ctx.Guild.ID, ctx.Message.Mentions[0].ID)
		if err != nil {
			ctx.Send("that person isnt even here wtf :(", err)
			return
		}
	} else {
		targetUser, err = ctx.Session.State.Member(ctx.Guild.ID, args[1])
		if err != nil {
			ctx.Send("that person isnt even here wtf :(", err)
			return
		}
	}

	if targetUser.User.ID == ctx.Session.State.User.ID {
		ctx.Send("no")
		return
	}
	/*
		_, err = ctx.Session.State.Member(ctx.Guild.ID, args[1])
		if err != nil {
			ctx.Send("didnt work: ", err.Error())
			return
		} */

	if targetUser.User.ID == ctx.Message.Author.ID {
		ctx.Send("no")
		return
	}

	topUserrole := ch.HighestRole(ctx.Guild, ctx.User.ID)
	topTargetrole := ch.HighestRole(ctx.Guild, targetUser.User.ID)
	topBotrole := ch.HighestRole(ctx.Guild, ctx.Session.State.User.ID)

	if topUserrole <= topTargetrole || topBotrole <= topTargetrole {
		ctx.Send("no")
		return
	}

	if topTargetrole > 0 {

		okCh := true

		userchannel, err := ctx.Session.UserChannelCreate(targetUser.User.ID)
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
	}
	err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, targetUser.User.ID, fmt.Sprintf("%v#%v - %v", ctx.Message.Author.Username, ctx.Message.Author.Discriminator, reason))
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
				Value:  fmt.Sprintf("%v", targetUser.User.ID),
				Inline: true,
			},
		},
		Color: dColorRed,
	}

	ctx.SendEmbed(embed)

}

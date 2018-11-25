package commands

import (
	"fmt"
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Server = Command{
	Name:          "server",
	Description:   "Shows information about the current server.",
	Triggers:      []string{"m?server", "m?serverinfo", "m?sa"},
	Usage:         "m?server",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		tc := 0
		vc := 0

		for _, val := range ctx.Guild.Channels {
			switch val.Type {
			case discordgo.ChannelTypeGuildText:
				tc++
			case discordgo.ChannelTypeGuildVoice:
				vc++
			}
		}

		owner, err := ctx.Session.User(ctx.Guild.OwnerID)
		if err != nil {
			ctx.Send("error occured")
			return
		}

		t, err := ctx.Guild.JoinedAt.Parse()
		if err != nil {
			ctx.Send("error occured")
			return
		}

		embed := discordgo.MessageEmbed{
			Color: dColorWhite,
			Author: &discordgo.MessageEmbedAuthor{
				Name: ctx.Guild.Name,
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Members",
					Value:  fmt.Sprintf("%v", ctx.Guild.MemberCount),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Channels",
					Value:  fmt.Sprintf("Total: %v\nText: %v\nVoice: %v", len(ctx.Guild.Channels), tc, vc),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Roles",
					Value:  fmt.Sprintf("%v roles", len(ctx.Guild.Roles)),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Owner",
					Value:  fmt.Sprintf("%v\n(%v)", owner.Mention(), ctx.Guild.OwnerID),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Created",
					Value:  fmt.Sprintf("%v", t.Format("15-04-2021 18:00pm")),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Verification level",
					Value:  fmt.Sprintf("%v", verificationMap[int(ctx.Guild.VerificationLevel)]),
					Inline: true,
				},
			},
		}

		if ctx.Guild.Icon != "" {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://cdn.discordapp.com/icons/%v/%v.png", ctx.Guild.ID, ctx.Guild.Icon),
			}
			embed.Author.IconURL = fmt.Sprintf("https://cdn.discordapp.com/icons/%v/%v.png", ctx.Guild.ID, ctx.Guild.Icon)
		}

		ctx.SendEmbed(&embed)
	},
}

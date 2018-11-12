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
		embed := discordgo.MessageEmbed{
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://cdn.discordapp.com/icons/%v/%v.png", ctx.Guild.ID, ctx.Guild.Icon),
			},
			Author: &discordgo.MessageEmbedAuthor{
				IconURL: fmt.Sprintf("https://cdn.discordapp.com/icons/%v/%v.png", ctx.Guild.ID, ctx.Guild.Icon),
				Name:    ctx.Guild.Name,
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Members",
					Value:  fmt.Sprintf("%v members", ctx.Guild.MemberCount),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Channels",
					Value:  fmt.Sprintf("%v channels", len(ctx.Guild.Channels)),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Verification level",
					Value:  fmt.Sprintf("%v", ctx.Guild.VerificationLevel),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Roles",
					Value:  fmt.Sprintf("%v roles", len(ctx.Guild.Roles)),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Owner",
					Value:  fmt.Sprintf("%v (%v)", ctx.Guild.OwnerID, ctx.Guild.OwnerID),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Created",
					Value:  fmt.Sprintf("%v", ctx.Guild.JoinedAt),
					Inline: true,
				},
			},
		}

		ctx.SendEmbed(&embed)
	},
}

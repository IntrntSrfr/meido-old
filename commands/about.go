package commands

import (
	"fmt"
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var About = Command{
	Name:          "about",
	Description:   "Shows info about Meido.",
	Aliases:       []string{},
	Usage:         "m?help",
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {
		var (
			totalUsers    int
			botUsers      int
			humanUsers    int
			totalChannels int
			textChannels  int
			voiceChannels int
		)
		for _, g := range ctx.Session.State.Guilds {
			totalUsers += g.MemberCount
			for _, m := range g.Members {
				if m.User.Bot {
					botUsers++
				} else {
					humanUsers++
				}
			}
			totalChannels += len(g.Channels)

			for _, ch := range g.Channels {
				if ch.Type == discordgo.ChannelTypeGuildText {
					textChannels++
				} else {
					voiceChannels++
				}
			}
		}

		embed := discordgo.MessageEmbed{
			Title: "About",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Users",
					Value:  fmt.Sprintf("Total: %v\nHuman: %v\nBot: %v", totalUsers, humanUsers, botUsers),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Channels",
					Value:  fmt.Sprintf("Total: %v\nText: %v\nVoice: %v", totalChannels, textChannels, voiceChannels),
					Inline: true,
				},
			},
		}

		ctx.SendEmbed(&embed)
	},
}

package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

/*
var About = Command{
	Name:          "About",
	Description:   "Shows info about Meido.",
	Triggers:      []string{"m?about"},
	Usage:         "m?about",
	Category:      Utility,
	RequiredPerms: discordgo.PermissionSendMessages,
	//Execute:       func(a []string, c *service.Context) {},
} */

func (ch *CommandHandler) about(args []string, ctx *service.Context) {
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

	var totalCommands int
	ctx.Db.Get(&totalCommands, "SELECT COUNT(*) FROM commandlog")

	thisTime := time.Now()

	timespan := thisTime.Sub(ch.botStartTime)

	embed := discordgo.MessageEmbed{
		Title: "About",
		Color: dColorWhite,
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
			&discordgo.MessageEmbedField{
				Name:   "Servers",
				Value:  fmt.Sprintf("%v servers", len(ctx.Session.State.Guilds)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Uptime",
				Value:  fmt.Sprintf("Uptime: %v", timespan.String()),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Total commands ran",
				Value:  strconv.Itoa(totalCommands),
				Inline: true,
			},
		},
	}

	ctx.SendEmbed(&embed)
}

package commands

import (
	"fmt"
	"meido-test/service"
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var About = Command{
	Name:          "about",
	Description:   "Shows info about Meido.",
	Triggers:      []string{"m?about"},
	Usage:         "m?about",
	RequiredPerms: discordgo.PermissionSendMessages,
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

		var totalCommands int
		row := db.QueryRow("SELECT COUNT(*) FROM commandlog;")
		err := row.Scan(&totalCommands)
		if err != nil {
			totalCommands = 0
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		thisTime := time.Now()

		timespan := thisTime.Sub(startTime)

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
					Name:   "Memory usage",
					Value:  fmt.Sprintf("(%vmb/%vmb)", m.TotalAlloc/1024/1024, m.Sys/1024/1024),
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
	},
}

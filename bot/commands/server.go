package commands

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) server(args []string, ctx *service.Context) {

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

	owner, err := ctx.Session.State.Member(ctx.Guild.ID, ctx.Guild.OwnerID)
	if err != nil {
		ctx.Send("error occured ", err)
		return
	}

	id, err := strconv.ParseInt(ctx.Guild.ID, 0, 63)
	if err != nil {
		return
	}

	id = ((id >> 22) + 1420070400000) / 1000

	dur := time.Since(time.Unix(int64(id), 0))

	ts := time.Unix(id, 0)

	embed := discordgo.MessageEmbed{
		Color: dColorWhite,
		Author: &discordgo.MessageEmbedAuthor{
			Name: ctx.Guild.Name,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Owner",
				Value:  fmt.Sprintf("%v\n(%v)", owner.Mention(), ctx.Guild.OwnerID),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:  "Creation date",
				Value: fmt.Sprintf("%v\n%v days ago", ts.Format(time.RFC1123), math.Floor(dur.Hours()/float64(24))),
			},
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
}

func (ch *CommandHandler) guildbanner(args []string, ctx *service.Context) {

	if ctx.Guild.Banner == "" {
		ctx.Send("This server has no banner")
		return
	}

	ctx.SendEmbed(&discordgo.MessageEmbed{
		Color: dColorWhite,
		Title: ctx.Guild.Name,
		Image: &discordgo.MessageEmbedImage{URL: discordgo.EndpointGuildBanner(ctx.Guild.ID, ctx.Guild.Banner)},
	})
}

func (ch *CommandHandler) guildsplash(args []string, ctx *service.Context) {

	if ctx.Guild.Splash == "" {
		ctx.Send("This server has no splash")
		return
	}

	ctx.SendEmbed(&discordgo.MessageEmbed{
		Color: dColorWhite,
		Title: ctx.Guild.Name,
		Image: &discordgo.MessageEmbedImage{URL: discordgo.EndpointGuildSplash(ctx.Guild.ID, ctx.Guild.Splash)},
	})
}

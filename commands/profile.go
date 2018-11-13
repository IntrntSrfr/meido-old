package commands

import (
	"fmt"
	"meido-test/models"
	"meido-test/service"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var Profile = Command{
	Name:          "profile",
	Description:   "Shows a user profile.",
	Triggers:      []string{"m?profile"},
	Usage:         "m?profile\nm?profile @internet surfer#0001\nm?profile 163454407999094786",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

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
		} else {
			targetUser = ctx.User
		}

		if targetUser == nil {
			ctx.Send("Could not find that user.")
		}

		if targetUser.Bot {
			ctx.Send("Bots dont get to join the fun")
			return
		}

		row := db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", targetUser.ID)

		dbu := models.Discorduser{}

		err = row.Scan(
			&dbu.Uid,
			&dbu.Userid,
			&dbu.Username,
			&dbu.Discriminator,
			&dbu.Xp,
			&dbu.Nextxpgaintime,
			&dbu.Xpexcluded,
			&dbu.Reputation,
			&dbu.Cangivereptime)

		if err != nil {
			ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: "User not available"})
			return
		}

		embed := discordgo.MessageEmbed{
			Color:     dColorGreen,
			Title:     fmt.Sprintf("Profile for %v", targetUser.String()),
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: targetUser.AvatarURL("1024")},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Experience",
					Value:  strconv.Itoa(dbu.Xp),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Reputation",
					Value:  strconv.Itoa(dbu.Reputation),
					Inline: true,
				},
			},
		}

		ctx.SendEmbed(&embed)
	},
}

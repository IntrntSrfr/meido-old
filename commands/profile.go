package commands

import (
	"fmt"
	"meido-test/models"
	"meido-test/service"
	"strconv"
	"strings"
	"time"

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

		dbu := models.DiscordUser{}

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

var Rep = Command{
	Name:          "rep",
	Description:   "Gives a user a reputation point or checks whether you can give it or not.",
	Triggers:      []string{"m?rep"},
	Usage:         "m?rep\nm?rep @internet surfer#0001\nm?rep 163454407999094786",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		u := ctx.User

		currentTime := time.Now()

		row := db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", u.ID)

		dbu := models.DiscordUser{}

		err := row.Scan(
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
			return
		}

		diff := dbu.Cangivereptime.Sub(currentTime.Add(time.Hour * time.Duration(2)))

		if len(args) < 2 {
			if diff > 0 {
				ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: strings.TrimSuffix(fmt.Sprintf("You can award a reputation point in %v", diff.Round(time.Minute).String()), "0s")})
			} else {
				ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorGreen, Description: "You can award a reputation point."})
			}
			return
		}

		var targetUser *discordgo.User
		if len(ctx.Message.Mentions) >= 1 {
			targetUser = ctx.Message.Mentions[0]
		} else {
			targetUser, err = ctx.Session.User(args[1])
			if err != nil {
				//s.ChannelMessageSend(ch.ID, err.Error())
				return
			}
		}

		if targetUser.Bot {
			ctx.Send("Bots dont get to join the fun")
			return
		}

		if u.ID == targetUser.ID {
			ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: "You cannot award yourself a reputation point."})
			return
		}

		if diff > 0 {
			ctx.SendEmbed(&discordgo.MessageEmbed{
				Color:       dColorRed,
				Description: strings.TrimSuffix(fmt.Sprintf("You can award a reputation point in %v", diff.Round(time.Minute).String()), "0s")})
			return
		}

		row = db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", targetUser.ID)

		dbtu := models.DiscordUser{}

		err = row.Scan(
			&dbtu.Uid,
			&dbtu.Userid,
			&dbtu.Username,
			&dbtu.Discriminator,
			&dbtu.Xp,
			&dbtu.Nextxpgaintime,
			&dbtu.Xpexcluded,
			&dbtu.Reputation,
			&dbtu.Cangivereptime)

		if err != nil {
			return
		}

		_, err = db.Exec("UPDATE discordusers SET reputation = $1, cangivereptime = $2 WHERE userid = $3", dbtu.Reputation+1, currentTime.Add(time.Hour*time.Duration(24)), dbtu.Userid)
		if err != nil {
			return
		}

		ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorGreen, Description: fmt.Sprintf("%v awarded %v a reputation point!", u.Mention(), targetUser.Mention())})
	},
}

var Repleaderboard = Command{
	Name:          "repleaderboard",
	Description:   "Gives a user a reputation point or checks whether you can give it or not.",
	Triggers:      []string{"m?rplb"},
	Usage:         "m?rplb",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		rows, err := db.Query("SELECT * FROM discordusers WHERE reputation > 0 ORDER BY reputation DESC LIMIT 10 ")
		if err != nil {
			fmt.Println(err)
			return
		}

		if rows.Err() != nil {
			fmt.Println(rows.Err())
		}

		leaderboard := "```\n"

		place := 1

		for rows.Next() {
			dbu := models.DiscordUser{}

			err = rows.Scan(
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
				fmt.Println(err)
				return
			}

			leaderboard += fmt.Sprintf("#%v - %v#%v - %v reputation points\n", place, dbu.Username, dbu.Discriminator, dbu.Reputation)
			place++
		}
		leaderboard += "```"

		ctx.Send(leaderboard)

	},
}

var Xpleaderboard = Command{
	Name:          "xpleaderboard",
	Description:   "Gives a user a reputation point or checks whether you can give it or not.",
	Triggers:      []string{"m?xplb"},
	Usage:         "m?xplb",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		rows, err := db.Query("SELECT * FROM discordusers WHERE xp > 0 ORDER BY xp DESC LIMIT 10 ")
		if err != nil {
			fmt.Println(err)
			return
		}

		if rows.Err() != nil {
			fmt.Println(rows.Err())
		}

		leaderboard := "```\n"

		place := 1

		for rows.Next() {
			dbu := models.DiscordUser{}

			err = rows.Scan(
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
				fmt.Println(err)
				return
			}

			leaderboard += fmt.Sprintf("#%v - %v#%v - %vxp\n", place, dbu.Username, dbu.Discriminator, dbu.Xp)
			place++
		}
		leaderboard += "```"

		ctx.Send(leaderboard)

	},
}

package commands

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var ShowProfile = Command{
	Name:          "Profile",
	Description:   "Shows a user profile.",
	Triggers:      []string{"m?profile", "m?p"},
	Usage:         "m?profile\nm?profile @internet surfer#0001\nm?profile 163454407999094786",
	Category:      Profile,
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

		dbu := models.DiscordUser{}
		lcxp := models.Localxp{}
		gbxp := models.Globalxp{}

		row := db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", targetUser.ID)
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
			SetupProfile(targetUser, ctx, 0)
			//ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: "User not available"})
		}
		row = db.QueryRow("SELECT xp FROM globalxp WHERE userid = $1", targetUser.ID)
		err = row.Scan(
			&gbxp.Xp,
		)
		if err != nil {
			//ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: "User not available"})
			return
		}
		row = db.QueryRow("SELECT xp FROM localxp WHERE userid = $1 AND guildid = $2", targetUser.ID, ctx.Guild.ID)
		err = row.Scan(
			&lcxp.Xp,
		)
		if err != nil {
			//ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: "User not available"})
			return
		}

		embed := discordgo.MessageEmbed{
			Color:     UserColor(ctx.Guild, targetUser.ID),
			Title:     fmt.Sprintf("Profile for %v", targetUser.String()),
			Thumbnail: &discordgo.MessageEmbedThumbnail{URL: targetUser.AvatarURL("1024")},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "Local xp",
					Value:  strconv.Itoa(lcxp.Xp),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "Global xp",
					Value:  strconv.Itoa(gbxp.Xp),
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
	Name:          "Rep",
	Description:   "Gives a user a reputation point or checks whether you can give it or not.",
	Triggers:      []string{"m?rep"},
	Usage:         "m?rep\nm?rep @internet surfer#0001\nm?rep 163454407999094786",
	Category:      Profile,
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

		diff := dbu.Cangivereptime.Sub(currentTime)

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

		row = db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", targetUser.ID)
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
			SetupProfile(targetUser, ctx, 1)
			//ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: "User not available"})
		}

		db.Exec("UPDATE discordusers SET reputation = $1 WHERE userid = $2", dbtu.Reputation+1, dbtu.Userid)
		db.Exec("UPDATE discordusers SET cangivereptime = $1 WHERE userid = $2", currentTime.Add(time.Hour*time.Duration(24)), dbu.Userid)

		ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorGreen, Description: fmt.Sprintf("%v awarded %v a reputation point!", u.Mention(), targetUser.Mention())})
	},
}

var Repleaderboard = Command{
	Name:          "Rep leaderboard",
	Description:   "Checks the reputation leaderboard.",
	Triggers:      []string{"m?rplb"},
	Usage:         "m?rplb",
	Category:      Profile,
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		rows, err := db.Query("SELECT userid, reputation FROM discordusers WHERE reputation > 0 ORDER BY reputation DESC LIMIT 10;")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		if rows.Err() != nil {
			fmt.Println(rows.Err())
		}

		leaderboard := "```\n"

		place := 1

		for rows.Next() {
			dbu := models.DiscordUser{}

			err = rows.Scan(
				&dbu.Userid,
				&dbu.Reputation,
			)

			if err != nil {
				fmt.Println(err)
				return
			}

			u, err := ctx.Session.User(dbu.Userid)
			if err != nil {
				continue
			}

			leaderboard += fmt.Sprintf("#%v - %v#%v - %v reputation points\n", place, u.Username, u.Discriminator, dbu.Reputation)
			place++
		}

		leaderboard += "\n"

		row := db.QueryRow("SELECT reputation FROM discordusers WHERE userid=$1;", ctx.User.ID)

		dbu := models.DiscordUser{}
		err = row.Scan(
			&dbu.Reputation,
		)
		if err != nil {
			ctx.Send("Error getting leaderboard.")
			return
		}

		leaderboard += fmt.Sprintf("Your stats\nReputation: %v\n", dbu.Reputation)

		leaderboard += "```"

		ctx.Send(leaderboard)

	},
}

var XpLeaderboard = Command{
	Name:          "XP leaderboard",
	Description:   "Checks local leaderboard.",
	Triggers:      []string{"m?xplb"},
	Usage:         "m?xplb",
	Category:      Profile,
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		rows, err := db.Query("SELECT userid, xp FROM localxp WHERE xp > 0 AND guildid = $1 ORDER BY xp DESC LIMIT 10;", ctx.Guild.ID)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		if rows.Err() != nil {
			fmt.Println(rows.Err())
		}

		leaderboard := "```\n"

		place := 1

		for rows.Next() {
			dbxp := models.Localxp{}

			err = rows.Scan(
				&dbxp.Userid,
				&dbxp.Xp,
			)

			if err != nil {
				fmt.Println(err)
				return
			}

			mem, err := ctx.Session.State.Member(ctx.Guild.ID, dbxp.Userid)
			if err != nil {
				leaderboard += fmt.Sprintf("#%v - User not here (%v) - %vxp\n", place, dbxp.Userid, dbxp.Xp)
			} else {
				leaderboard += fmt.Sprintf("#%v - %v#%v - %vxp\n", place, mem.User.Username, mem.User.Discriminator, dbxp.Xp)
			}
			place++
		}
		leaderboard += "```"

		ctx.Send(leaderboard)

	},
}

var GlobalXpLeaderboard = Command{
	Name:          "Global XP Leaderboard",
	Description:   "Checks the global xp leaderboard.",
	Triggers:      []string{"m?gxplb"},
	Usage:         "m?gxplb",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		rows, err := db.Query("SELECT userid, xp FROM globalxp WHERE xp > 0 ORDER BY xp DESC LIMIT 10;")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		if rows.Err() != nil {
			fmt.Println(rows.Err())
		}

		leaderboard := "```\n"

		place := 1

		for rows.Next() {
			dbgxp := models.Globalxp{}

			err = rows.Scan(
				&dbgxp.Userid,
				&dbgxp.Xp,
			)

			if err != nil {
				fmt.Println(err)
				return
			}

			user, err := ctx.Session.User(dbgxp.Userid)
			if err != nil {
				return
			}

			leaderboard += fmt.Sprintf("#%v - %v#%v - %vxp\n", place, user.Username, user.Discriminator, dbgxp.Xp)
			place++
		}
		leaderboard += "```"

		ctx.Send(leaderboard)

	},
}

var XpIgnoreChannel = Command{
	Name:          "XP ignore channel",
	Description:   "Adds or removes a channel to or from the xp ignored list.",
	Triggers:      []string{"m?xpignorechannel", "m?xpigch"},
	Usage:         "m?xpigch\nm?xpigch 123123123123",
	Category:      Profile,
	RequiredPerms: discordgo.PermissionManageChannels,
	Execute: func(args []string, ctx *service.Context) {

		var (
			err     error
			channel *discordgo.Channel
			ch      string
		)

		if len(args) > 1 {
			if strings.HasPrefix(args[1], "<#") && strings.HasSuffix(args[1], ">") {
				ch = args[1]
				ch = ch[2 : len(ch)-1]
			} else {
				ch = args[1]
			}
			channel, err = ctx.Session.State.Channel(ch)
			if err != nil {
				ctx.Send("Invalid channel.")
				return
			}

			if channel.GuildID != ctx.Guild.ID {
				ctx.Send("Channel not found.")
				return
			}
		} else {
			channel = ctx.Channel
		}

		dbigch := models.Xpignoredchannel{}

		row := db.QueryRow("SELECT channelid FROM xpignoredchannels WHERE channelid = $1;", channel.ID)
		err = row.Scan(&dbigch.Channelid)
		if err != nil {
			if err == sql.ErrNoRows {
				db.Exec("INSERT INTO xpignoredchannels (channelid) VALUES ($1);", channel.ID)
				ctx.Send(fmt.Sprintf("Added %v to the xp ignore list.", channel.Mention()))
			}
		} else {
			db.Exec("DELETE FROM xpignoredchannels WHERE channelid=$1;", channel.ID)
			ctx.Send(fmt.Sprintf("Removed %v from the xp ignore list.", channel.Mention()))
		}
	},
}

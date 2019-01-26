package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/intrntsrfr/meido/models"
	"github.com/intrntsrfr/meido/service"

	"github.com/bwmarrin/discordgo"
)

var Warn = Command{
	Name:          "warn",
	Description:   "Warns a user, adding a strike. Does not work if strike system is disabled.",
	Triggers:      []string{"m?warn", ".warn"},
	Usage:         "m?warn 163454407999094786\n.warn @internet surfer#0001",
	RequiredPerms: discordgo.PermissionBanMembers,
	Execute: func(args []string, ctx *service.Context) {

		row := db.QueryRow("SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid = $1;", ctx.Guild.ID)

		dbg := models.DiscordGuild{}

		err := row.Scan(&dbg.UseStrikes, &dbg.MaxStrikes)
		if err != nil {
			return
		}

		if dbg.UseStrikes {

			if len(args) < 2 {
				ctx.Send("no")
				return
			}

			var targetUser *discordgo.Member

			if len(ctx.Message.Mentions) >= 1 {
				targetUser, err = ctx.Session.State.Member(ctx.Guild.ID, ctx.Message.Mentions[0].ID)
				if err != nil {
					ctx.Send("error occured:", err)
					return
				}
			} else {
				targetUser, err = ctx.Session.State.Member(ctx.Guild.ID, args[1])
				if err != nil {
					ctx.Send("error occured:", err)
					return
				}
			}

			if targetUser.User.ID == ctx.Session.State.User.ID {
				ctx.Send("no")
				return
			}

			_, err = ctx.Session.State.Member(ctx.Guild.ID, targetUser.User.ID)
			if err != nil {
				ctx.Send("that person isnt even here wtf :(")
				return
			}

			if targetUser.User.ID == ctx.Message.Author.ID {
				ctx.Send("no")
				return
			}

			topUserrole := HighestRole(ctx.Guild, ctx.User.ID)
			topTargetrole := HighestRole(ctx.Guild, targetUser.User.ID)

			if topUserrole <= topTargetrole {
				ctx.Send("no")
				return
			}

			reason := "no reason"
			if len(args) > 2 {
				reason = strings.Join(args[2:], " ")
			}

			strikeCount := 0

			row := db.QueryRow("SELECT COUNT(*) FROM strikes WHERE guildid = $1 AND userid = $2;", ctx.Guild.ID, targetUser.User.ID)
			err := row.Scan(&strikeCount)
			if err != nil {
				ctx.Send(err.Error())
				return
			}

			if strikeCount+1 >= dbg.MaxStrikes {
				//ban
				userch, _ := ctx.Session.UserChannelCreate(targetUser.User.ID)
				ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for acquiring %v strikes.\nLast warning was: %v", ctx.Guild.Name, dbg.MaxStrikes, reason))
				err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, targetUser.User.ID, fmt.Sprintf("Acquired %v strikes.", dbg.MaxStrikes), 0)
				if err != nil {
					ctx.Send(err.Error())
					return
				}

				ctx.Send(fmt.Sprintf("%v has been banned after acquiring too many strikes. miss them.", targetUser.Mention()))
				_, err := db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", targetUser.User.ID, ctx.Guild.ID)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				//insert warn
				_, err = db.Exec("INSERT INTO strikes(guildid, userid, reason, executorid, tstamp) VALUES ($1, $2, $3, $4, $5);", ctx.Guild.ID, targetUser.User.ID, reason, ctx.User.ID, time.Now())
				if err != nil {
					ctx.Send(err)
					return
				}
				userch, _ := ctx.Session.UserChannelCreate(targetUser.User.ID)
				ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been warned in %v.\nWarned for: %v\nYou are currently at strike %v/%v", ctx.Guild.Name, reason, strikeCount+1, dbg.MaxStrikes))
				ctx.Send(fmt.Sprintf("%v has been warned\nThey are currently at strike %v/%v", targetUser.Mention(), strikeCount+1, dbg.MaxStrikes))
			}
		} else {
			ctx.Send("Strike system is not enabled.")
		}
	},
}

var StrikeLog = Command{
	Name:          "strikelog",
	Description:   "Shows a users strikes.",
	Triggers:      []string{"m?strikelog"},
	Usage:         "m?strikelog 163454407999094786\nm?strikelog @internet surfer#0001",
	RequiredPerms: discordgo.PermissionBanMembers,
	RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		var err error

		if len(args) < 2 {
			ctx.Send("no")
			return
		}

		var targetUser *discordgo.User

		if len(ctx.Message.Mentions) >= 1 {
			targetUser = ctx.Message.Mentions[0]
		} else {
			targetUser, err = ctx.Session.User(args[1])
			if err != nil {
				ctx.Send("error occured:", err)
				return
			}
		}
		if targetUser == nil {
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:  fmt.Sprintf("Strikes for %v", targetUser.String()),
			Fields: []*discordgo.MessageEmbedField{},
		}

		rows, err := db.Query("SELECT * FROM strikes WHERE userid=$1 AND guildid=$2;", targetUser.ID, ctx.Guild.ID)
		if err != nil {
			ctx.Send("No strikes.")
			return
		}
		count := 0
		for rows.Next() {
			count++
		}

		if count <= 0 {
			embed.Description = "No strikes."
		} else {

			for rows.Next() {
				dbs := models.Strikes{}
				err = rows.Scan(&dbs.Uid, &dbs.Guildid, &dbs.Userid, &dbs.Reason, &dbs.Executorid, &dbs.Tstamp)
				if err != nil {
					fmt.Println(err)
					return
				}
				exec := ""
				mem, err := ctx.Session.State.Member(ctx.Guild.ID, dbs.Executorid)
				if err != nil {
					exec = dbs.Executorid
				} else {
					exec = mem.User.String()
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:  fmt.Sprintf("Strike ID: %v - At %v by %v", dbs.Uid, dbs.Tstamp.Format(time.RFC1123), exec),
					Value: fmt.Sprintf("- %v", dbs.Reason),
				})
			}
		}
		ctx.SendEmbed(embed)
	},
}

var StrikeLogAll = Command{
	Name:          "strikelogall",
	Description:   "Shows all strikes in the guild.",
	Triggers:      []string{"m?strikelogall"},
	Usage:         "m?strikelog",
	RequiredPerms: discordgo.PermissionBanMembers,
	Execute: func(args []string, ctx *service.Context) {
		/*
			var err error
			embed := &discordgo.MessageEmbed{
				Title:  fmt.Sprintf("Strikes in %v", targetUser.String()),
				Fields: []*discordgo.MessageEmbedField{},
			}

			rows, err := db.Query("SELECT * FROM strikes WHERE userid=$1 AND guildid = $2;", targetUser.ID, ctx.Guild.ID)
			if err != nil {
				ctx.Send("No strikes.")
				return
			}
			for rows.Next() {
				dbs := models.Strikes{}
				err = rows.Scan(&dbs)
				if err != nil {
					return
				}
				exec := ""
				mem, err := ctx.Session.State.Member(ctx.Guild.ID, dbs.Executorid)
				if err != nil {
					exec = dbs.Executorid
				} else {
					exec = mem.User.String()
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:  fmt.Sprintf("Strike ID: %v - At %v by %v", dbs.Uid, dbs.Tstamp.Format(time.RFC1123), exec),
					Value: fmt.Sprintf("- %v", dbs.Reason),
				})
			}
			ctx.SendEmbed(embed) */
	},
}

var RemoveStrike = Command{
	Name:          "removestrike",
	Description:   "Removes a strike from a user.",
	Triggers:      []string{"m?removestrike", "m?rmstrike"},
	Usage:         "m?removestrike 163454407999094786\nm?rmstrike @internet surfer#0001",
	RequiredPerms: discordgo.PermissionBanMembers,
	Execute: func(args []string, ctx *service.Context) {
		if len(args) < 2 {
			return
		}

		uid, err := strconv.Atoi(args[1])
		if err != nil {
			ctx.Send("no")
			return
		}

		dbs := models.Strikes{}
		row := db.QueryRow("SELECT guildid FROM strikes WHERE uid = $1;", uid)

		err = row.Scan(&dbs.Guildid)
		if err != nil {
			ctx.Send("Strike does not exist.")
			return
		}
		if dbs.Guildid != ctx.Guild.ID {
			ctx.Send("nice try")
			return
		}

		_, err = db.Exec("DELETE FROM strikes WHERE uid=$1;", uid)
		if err != nil {
			ctx.Send("Error deleting strike.")
		} else {
			ctx.Send(fmt.Sprintf("Removed strike with ID: %v", args[1]))
		}
	},
}

var ClearStrikes = Command{
	Name:          "clearstrikes",
	Description:   "Clears the strikes on a user.",
	Triggers:      []string{"m?clearstrikes", "m?cs"},
	Usage:         "m?clearstrikes @internet surfer#0001\nm?cs 163454407999094786",
	RequiredPerms: discordgo.PermissionManageMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) < 2 {
			return
		}

		var (
			targetUser *discordgo.User
			err        error
		)
		if len(ctx.Message.Mentions) >= 1 {
			targetUser = ctx.Message.Mentions[0]
		} else {
			targetUser, err = ctx.Session.User(args[1])
			if err != nil {
				return
			}
		}

		res, err := db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", targetUser.ID, ctx.Guild.ID)
		if err != nil {
			ctx.Send("error occured :" + err.Error())
			return
		}
		ar, _ := res.RowsAffected()
		if ar < 1 {
			ctx.Send("User has no strikes.")
		} else {
			ctx.Send(fmt.Sprintf("Removed strikes from user: %v", targetUser.Mention()))
		}
	},
}

package commands

import (
	"database/sql"
	"fmt"
	"meido-test/models"
	"meido-test/service"
	"strings"

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

			_, err = ctx.Session.State.Member(ctx.Guild.ID, targetUser.ID)
			if err != nil {
				ctx.Send("that person isnt even here wtf :(")
				return
			}

			if targetUser.ID == ctx.Message.Author.ID {
				ctx.Send("no")
				return
			}

			topUserrole := HighestRole(ctx.Guild, ctx.User.ID)
			topTargetrole := HighestRole(ctx.Guild, targetUser.ID)

			if topUserrole <= topTargetrole {
				ctx.Send("no")
				return
			}

			dbs := models.Strikes{}

			reason := "no reason"
			if len(args) > 2 {
				reason = strings.Join(args[2:], " ")
			}

			row := db.QueryRow("SELECT * FROM strikes WHERE guildid = $1 AND userid = $2;", ctx.Guild.ID, targetUser.ID)
			err := row.Scan(&dbs.Uid, &dbs.Guildid, &dbs.Userid, &dbs.Strikes)
			if err != nil {
				if err == sql.ErrNoRows {
					if dbg.MaxStrikes < 2 {
						userch, _ := ctx.Session.UserChannelCreate(targetUser.ID)
						ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for acquiring %v strikes.\nLast warning was: %v", ctx.Guild.Name, dbg.MaxStrikes, reason))
						err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, targetUser.ID, fmt.Sprintf("Acquired %v strikes.", dbg.MaxStrikes), 0)
						if err != nil {
							ctx.Send(err.Error())
							return
						}

						ctx.Send(fmt.Sprintf("%v has been banned after acquiring too many strikes. miss them.", targetUser.Mention()))

					} else {
						userch, _ := ctx.Session.UserChannelCreate(targetUser.ID)
						ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been warned in %v.\nWarned for: %v", ctx.Guild.Name, reason))
						ctx.Send(fmt.Sprintf("%v has been warned\nThey are currently at strike %v/%v", targetUser.Mention(), dbs.Strikes+1, dbg.MaxStrikes))
						db.Exec("INSERT INTO strikes(guildid, userid, strikes) VALUES ($1, $2, $3);", ctx.Guild.ID, targetUser.ID, 1)
					}
				}
			} else {
				if dbs.Strikes+1 >= dbg.MaxStrikes {
					userch, _ := ctx.Session.UserChannelCreate(targetUser.ID)
					ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for acquiring %v strikes.\nLast warning was: %v", ctx.Guild.Name, dbg.MaxStrikes, reason))
					err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, targetUser.ID, fmt.Sprintf("Acquired %v strikes.", dbg.MaxStrikes), 0)
					if err != nil {
						ctx.Send(err.Error())
						return
					}

					ctx.Send(fmt.Sprintf("%v has been banned after acquiring too many strikes. miss them.", targetUser.Mention()))

					_, err := db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", targetUser.ID, ctx.Guild.ID)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					userch, _ := ctx.Session.UserChannelCreate(targetUser.ID)
					ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been warned in %v.\nWarned for: %v", ctx.Guild.Name, reason))
					ctx.Send(fmt.Sprintf("%v has been warned\nThey are currently at strike %v/%v", targetUser.Mention(), dbs.Strikes+1, dbg.MaxStrikes))
					db.Exec("UPDATE strikes SET strikes = $1 WHERE userid = $2 AND guildid = $3;", dbs.Strikes+1, targetUser.ID, ctx.Guild.ID)
				}
			}
		} else {
			ctx.Send("Strike system is not enabled.")
		}
	},
}

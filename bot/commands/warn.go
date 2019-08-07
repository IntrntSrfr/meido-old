package commands

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) warn(args []string, ctx *service.Context) {

	dbg := &models.DiscordGuild{}
	err := ch.db.Get(dbg, "SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid = $1;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
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
				ctx.Send("that person isnt even here wtf :(")
				return
			}
		} else {
			targetUser, err = ctx.Session.State.Member(ctx.Guild.ID, args[1])
			if err != nil {
				ctx.Send("that person isnt even here wtf :(")
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

		topUserrole := ch.HighestRole(ctx.Guild, ctx.User.ID)
		topTargetrole := ch.HighestRole(ctx.Guild, targetUser.User.ID)
		topBotrole := ch.HighestRole(ctx.Guild, ctx.Session.State.User.ID)

		if topUserrole <= topTargetrole || topBotrole <= topTargetrole {
			ctx.Send("no")
			return
		}

		reason := "no reason"
		if len(args) > 2 {
			reason = strings.Join(args[2:], " ")
		}

		strikeCount := 0

		err = ch.db.Get(&strikeCount, "SELECT COUNT(*) FROM strikes WHERE guildid = $1 AND userid = $2;", ctx.Guild.ID, targetUser.User.ID)
		if err != nil {
			ctx.Send("there was an error, please try again")
			ch.logger.Error("error", zap.Error(err))
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
			_, err := ch.db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", targetUser.User.ID, ctx.Guild.ID)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			//insert warn
			_, err = ch.db.Exec("INSERT INTO strikes(guildid, userid, reason, executorid, tstamp) VALUES ($1, $2, $3, $4, $5);", ctx.Guild.ID, targetUser.User.ID, reason, ctx.User.ID, time.Now())
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
}

func (ch *CommandHandler) strikeLog(args []string, ctx *service.Context) {

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

	dbs := []models.Strikes{}
	err = ch.db.Select(&dbs, "SELECT * FROM strikes WHERE userid=$1 AND guildid=$2;", targetUser.ID, ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	if len(dbs) < 1 {
		embed.Description = "No strikes."
	} else {
		for _, strk := range dbs {
			exec := ""
			mem, err := ctx.Session.State.Member(ctx.Guild.ID, strk.Executorid)
			if err != nil {
				exec = strk.Executorid
			} else {
				exec = mem.User.String()
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("Strike ID: %v - At %v by %v", strk.Uid, strk.Tstamp.Format(time.RFC822), exec),
				Value: fmt.Sprintf("- %v", strk.Reason),
			})
		}
	}
	ctx.SendEmbed(embed)
}

var StrikeLogAll = Command{
	Name:          "Full strike log",
	Description:   "Shows all strikes in the guild.",
	Triggers:      []string{"m?strikelogall"},
	Usage:         "m?strikelog",
	Category:      Moderation,
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

func (ch *CommandHandler) removeStrike(args []string, ctx *service.Context) {
	if len(args) < 2 {
		return
	}

	uid, err := strconv.Atoi(args[1])
	if err != nil {
		ctx.Send("no")
		return
	}

	dbs := &models.Strikes{}
	err = ch.db.Get(dbs, "SELECT guildid FROM strikes WHERE uid = $1;", uid)
	if err != nil && err != sql.ErrNoRows {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		ctx.Send("Strike does not exist.")
		return
	}

	if dbs.Guildid != ctx.Guild.ID {
		ctx.Send("nice try")
		return
	}

	_, err = ch.db.Exec("DELETE FROM strikes WHERE uid=$1;", uid)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	ctx.Send(fmt.Sprintf("Removed strike with ID: %v", args[1]))
}

func (ch *CommandHandler) clearStrikes(args []string, ctx *service.Context) {

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

	_, err = ch.db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", targetUser.ID, ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	ctx.Send(fmt.Sprintf("Removed strikes from user: %v", targetUser.Mention()))

}

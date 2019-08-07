package events

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"
)

func (eh *EventHandler) messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	//fmt.Println(fmt.Sprintf("[%v] - context in %v", id, time.Now().Sub(startTime)))

	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if ch.Type == discordgo.ChannelTypeDM {
		dmembed := discordgo.MessageEmbed{
			Color:       dColorWhite,
			Title:       fmt.Sprintf("Message from %v", m.Author.String()),
			Description: m.Content,
			Footer:      &discordgo.MessageEmbedFooter{Text: m.Author.ID},
			Timestamp:   string(m.Timestamp),
		}

		if len(m.Attachments) > 0 {
			dmembed.Image = &discordgo.MessageEmbedImage{URL: m.Attachments[0].URL}
		}

		for i := range eh.dmLogChannels {
			dmch := eh.dmLogChannels[i]

			_, err := s.ChannelMessageSendEmbed(dmch, &dmembed)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		return
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return
	}

	startTime := time.Now()

	context, err := service.NewContext(s, m.Message, startTime)
	if err != nil {
		return
	}

	fmt.Println(fmt.Sprintf("%v - %v - %v: %v", context.Guild.Name, ch.Name, m.Author.String(), m.Content))

	perms, err := s.State.UserChannelPermissions(m.Author.ID, ch.ID)
	if err != nil {
		return
	}
	botPerms, err := s.State.UserChannelPermissions(s.State.User.ID, ch.ID)
	if err != nil {
		return
	}

	args := strings.Split(m.Content, " ")

	isIllegal := eh.checkFilter(&context, &perms, m.Message)
	if isIllegal {
		return
	}
	//fmt.Println(fmt.Sprintf("[%v] - filter in %v", id, time.Now().Sub(startTime)))

	eh.doXp(&context)
	//fmt.Println(fmt.Sprintf("[%v] - xp in %v", id, time.Now().Sub(startTime)))

	triggerCommand := ""
	for _, val := range eh.ch.GetCommandMap() {
		for _, com := range val.Triggers {
			if strings.ToLower(args[0]) == strings.ToLower(com) {
				triggerCommand = val.Name
			}
		}
	}
	//fmt.Println(fmt.Sprintf("[%v] - checked command in %v", id, time.Now().Sub(startTime)))

	if triggerCommand != "" {

		if cmd, ok := eh.ch.GetCommandMap()[triggerCommand]; ok {

			isOwner := false

			for _, val := range eh.ownerIds {
				if m.Author.ID == val {
					isOwner = true
				}
			}
			if cmd.RequiresOwner {
				if !isOwner {
					context.Send("Owner only.")
					return
				}
			}
			/*
				if !cmd.RequiresOwner {
				}
			*/
			if perms&cmd.RequiredPerms == 0 && perms&discordgo.PermissionAdministrator == 0 {
				//fmt.Println(perms, cmd.RequiredPerms, permMap[cmd.RequiredPerms], perms&cmd.RequiredPerms)
				return
			}

			if botPerms&cmd.RequiredPerms == 0 && perms&discordgo.PermissionAdministrator == 0 {
				context.Send(fmt.Sprintf("I am missing permissions: %v", permMap[cmd.RequiredPerms]))
				return
			}

			go cmd.Execute(args, &context)
			//fmt.Println(fmt.Sprintf("[%v] - executed command in %v\n", id, time.Now().Sub(startTime)))
			eh.db.Exec("INSERT INTO commandlog(command, args, userid, guildid, channelid, messageid, tstamp) VALUES($1, $2, $3, $4, $5, $6, $7)", cmd.Name, strings.Join(args, " "), m.Author.ID, context.Guild.ID, ch.ID, m.ID, time.Now())
			fmt.Println(fmt.Sprintf("\nCommand executed\nCommand: %v\nUser: %v [%v]\nSource: %v [%v] - #%v [%v]", args, m.Author.String(), m.Author.ID, context.Guild.Name, context.Guild.ID, ch.Name, ch.ID))
		}
	}
}

func (eh *EventHandler) checkFilter(ctx *service.Context, perms *int, msg *discordgo.Message) bool {

	isIllegal := false
	trigger := ""

	if *perms&discordgo.PermissionManageMessages == 0 {

		var count int

		err := eh.db.Get(&count, "SELECT COUNT(*) FROM filterignorechannels WHERE channelid = $1;", ctx.Channel.ID)
		if err != nil {
			eh.logger.Error("error", zap.Error(err))
			return false
		}

		if count > 0 {
			return false
		}

		guildFilters := []models.Filter{}

		err = eh.db.Select(&guildFilters, "SELECT phrase FROM filters WHERE guildid = $1", ctx.Guild.ID)
		if err != nil {
			eh.logger.Error("error", zap.Error(err))
			return false
		}

		for _, filter := range guildFilters {

			if strings.Contains(strings.ToLower(msg.Content), strings.ToLower(filter.Phrase)) {
				trigger = filter.Phrase
				isIllegal = true
				break
			}
		}

		if isIllegal {
			ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)

			dbg := &models.DiscordGuild{}
			err = eh.db.Get(dbg, "SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid = $1;", ctx.Guild.ID)
			if err != nil {
				eh.logger.Error("error", zap.Error(err))
				return false
			}

			if dbg.UseStrikes {

				reason := fmt.Sprintf("Triggering filter: %v", trigger)

				strikeCount := 0

				err = eh.db.Get(&strikeCount, "SELECT COUNT(*) FROM strikes WHERE guildid = $1 AND userid = $2;", ctx.Guild.ID, ctx.User.ID)
				if err != nil {
					eh.logger.Error("error", zap.Error(err))
					return false
				}

				if strikeCount+1 >= dbg.MaxStrikes {
					//ban
					userch, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
					ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for acquiring %v strikes.\nLast warning was: %v", ctx.Guild.Name, dbg.MaxStrikes, reason))
					err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, ctx.User.ID, fmt.Sprintf("Acquired %v strikes.", dbg.MaxStrikes), 0)
					if err != nil {
						return false
					}

					ctx.Send(fmt.Sprintf("%v has been banned after acquiring too many strikes. Miss them.", ctx.User.Mention()))
					_, err := eh.db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", ctx.User.ID, ctx.Guild.ID)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					//insert warn
					userch, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
					ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been warned in %v.\nWarned for: %v", ctx.Guild.Name, reason))
					ctx.Send(fmt.Sprintf("%v has been warned\nThey are currently at strike %v/%v", ctx.User.Mention(), strikeCount+1, dbg.MaxStrikes))
					eh.db.Exec("INSERT INTO strikes(guildid, userid, reason, executorid, tstamp) VALUES ($1, $2, $3, $4, $5);", ctx.Guild.ID, ctx.User.ID, reason, ctx.Session.State.User.ID, time.Now())
				}

			} else {
				ctx.Send(fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!", msg.Author.Mention()))
			}
		}
	}

	return isIllegal
}

func (eh *EventHandler) doXp(ctx *service.Context) {

	dbu := &models.DiscordUser{}

	currentTime := time.Now()
	xpTime := time.Now()
	isIgnored := false

	err := eh.db.Get(dbu, "SELECT * FROM discordusers WHERE userid = $1", ctx.User.ID)
	if err != nil {
		eh.logger.Error("error", zap.Error(err))
		return
	}

	if err != nil {
		if err == sql.ErrNoRows {
			eh.db.Exec("INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
				ctx.User.ID,
				ctx.User.Username,
				ctx.User.Discriminator,
				0,
				currentTime,
				false,
				0,
				currentTime,
			)
		}
	} else {
		isIgnored = dbu.Xpexcluded
		xpTime = dbu.Nextxpgaintime
	}

	diff := xpTime.Sub(currentTime)

	if diff <= 0 {
		//igu := models.Xpignoreduser{}
		lcxp := &models.Localxp{}
		gbxp := &models.Globalxp{}

		newXp := Random(15, 26)

		count := 0
		err = eh.db.Get(&count, "SELECT COUNT(*) FROM xpignoredchannels WHERE channelid = $1;", ctx.Channel.ID)
		if err != nil {
			eh.logger.Error("error", zap.Error(err))
			return
		}

		if count > 0 {
			newXp = 0
		}
		/*
			row = db.QueryRow("SELECT * FROM xpignoreduser WHERE userid = $1;", ctx.User.ID)
			err = row.Scan(
				&igu.Uid,
				&igu.Userid,
			)

			if err != nil {
				if igu.Userid == ctx.User.ID {
					newXp = 0
				}
			} */

		if isIgnored {
			newXp = 0
		}

		err = eh.db.Get(lcxp, "SELECT * FROM localxp WHERE userid = $1 AND guildid = $2;", ctx.User.ID, ctx.Guild.ID)
		if err != nil && err != sql.ErrNoRows {
			eh.logger.Error("error", zap.Error(err))
			return
		} else if err == sql.ErrNoRows {
			eh.db.Exec("INSERT INTO localxp(guildid, userid, xp) VALUES($1, $2, $3);", ctx.Guild.ID, ctx.User.ID, newXp)
		} else {
			if !isIgnored {
				eh.db.Exec("UPDATE localxp SET xp = $1 WHERE guildid = $2 AND userid = $3;", lcxp.Xp+newXp, ctx.Guild.ID, ctx.User.ID)
			}
		}

		err = eh.db.Get(gbxp, "SELECT * FROM globalxp WHERE userid = $1;", ctx.User.ID)
		if err != nil && err != sql.ErrNoRows {
			eh.logger.Error("error", zap.Error(err))
			return
		} else if err == sql.ErrNoRows {
			eh.db.Exec("INSERT INTO globalxp(userid, xp) VALUES($1, $2);", ctx.User.ID, newXp)
		} else {
			if !isIgnored {
				eh.db.Exec("UPDATE globalxp SET xp = $1 WHERE userid = $2;", gbxp.Xp+newXp, ctx.User.ID)
			}
		}

		eh.db.Exec("UPDATE discordusers SET nextxpgaintime = $1 WHERE userid = $2;", currentTime.Add(time.Minute*time.Duration(2)), ctx.User.ID)
	}
}

func (eh *EventHandler) SetupProfile(target *discordgo.User, ctx *service.Context, rep int) {

	currentTime := time.Now()
	dbu := &models.DiscordUser{}

	err := eh.db.Get(dbu, "SELECT * FROM discordusers WHERE userid = $1", target.ID)
	if err != nil && err != sql.ErrNoRows {
		eh.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		eh.db.Exec("INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
			target.ID,
			target.Username,
			target.Discriminator,
			0,
			currentTime,
			false,
			rep,
			currentTime,
		)
	}

	lcxp := &models.Localxp{}
	gbxp := &models.Globalxp{}

	newXp := 0

	err = eh.db.Get(lcxp, "SELECT * FROM localxp WHERE userid = $1 AND guildid = $2;", target.ID, ctx.Guild.ID)
	if err != nil && err != sql.ErrNoRows {
		eh.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		eh.db.Exec("INSERT INTO localxp(guildid, userid, xp) VALUES($1, $2, $3);", ctx.Guild.ID, target.ID, newXp)
	}

	err = eh.db.Get(gbxp, "SELECT * FROM globalxp WHERE userid = $1;", target.ID)
	if err != nil && err != sql.ErrNoRows {
		eh.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		eh.db.Exec("INSERT INTO globalxp(userid, xp) VALUES($1, $2);", target.ID, newXp)
	}

}

func (eh *EventHandler) HighestRole(g *discordgo.Guild, userID string) int {

	user, err := eh.client.State.Member(g.ID, userID)
	if err != nil {
		return -1
	}

	if user.User.ID == g.OwnerID {
		return math.MaxInt64
	}

	topRole := 0

	for _, val := range user.Roles {
		for _, role := range g.Roles {
			if val == role.ID {
				if role.Position > topRole {
					topRole = role.Position
				}
			}
		}
	}

	return topRole
}

func (eh *EventHandler) UserColor(g *discordgo.Guild, userID string) int {

	member, err := eh.client.State.Member(g.ID, userID)
	if err != nil {
		return 0
	}

	roles := discordgo.Roles(g.Roles)
	sort.Sort(roles)

	for _, role := range roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID {
				if role.Color != 0 {
					return role.Color
				}
			}
		}
	}

	return 0
}

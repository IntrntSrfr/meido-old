package commands

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) showProfile(args []string, ctx *service.Context) {

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

	dbu := &models.DiscordUser{}
	lcxp := &models.Localxp{}
	gbxp := &models.Globalxp{}

	err = ch.db.Get(dbu, "SELECT * FROM discordusers WHERE userid = $1", targetUser.ID)
	if err != nil && err != sql.ErrNoRows {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		ch.SetupProfile(targetUser, ctx, 0)
	}

	err = ch.db.Get(gbxp, "SELECT xp FROM globalxp WHERE userid = $1", targetUser.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	err = ch.db.Get(lcxp, "SELECT xp FROM localxp WHERE userid = $1 AND guildid = $2", targetUser.ID, ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	embed := &discordgo.MessageEmbed{
		Color:     ch.UserColor(ctx.Guild, targetUser.ID),
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
	ctx.SendEmbed(embed)
}

func (ch *CommandHandler) rep(args []string, ctx *service.Context) {

	u := ctx.User

	currentTime := time.Now()

	dbu := &models.DiscordUser{}

	err := ch.db.Get(dbu, "SELECT * FROM discordusers WHERE userid = $1", u.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
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

	dbtu := &models.DiscordUser{}

	err = ch.db.Get(dbtu, "SELECT * FROM discordusers WHERE userid = $1", targetUser.ID)
	if err != nil && err != sql.ErrNoRows {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		ch.SetupProfile(targetUser, ctx, 1)
	}

	ch.db.Exec("UPDATE discordusers SET reputation = $1 WHERE userid = $2", dbtu.Reputation+1, dbtu.Userid)
	ch.db.Exec("UPDATE discordusers SET cangivereptime = $1 WHERE userid = $2", currentTime.Add(time.Hour*time.Duration(24)), dbu.Userid)

	ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorGreen, Description: fmt.Sprintf("%v awarded %v a reputation point!", u.Mention(), targetUser.Mention())})
}

func (ch *CommandHandler) repleaderboard(args []string, ctx *service.Context) {

	users := []models.DiscordUser{}

	err := ch.db.Select(&users, "SELECT userid, reputation FROM discordusers WHERE reputation > 0 ORDER BY reputation DESC LIMIT 10;")
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	leaderboard := "```\n"

	place := 1

	for _, user := range users {

		u, err := ctx.Session.User(user.Userid)
		if err != nil {
			continue
		}

		leaderboard += fmt.Sprintf("#%v - %v#%v - %v reputation points\n", place, u.Username, u.Discriminator, user.Reputation)
		place++
	}

	leaderboard += "\n"

	user := &models.DiscordUser{}
	err = ch.db.Get(user, "SELECT reputation FROM discordusers WHERE userid=$1;", ctx.User.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	leaderboard += fmt.Sprintf("Your stats\nReputation: %v\n", user.Reputation)

	leaderboard += "```"

	ctx.Send(leaderboard)
}

func (ch *CommandHandler) xpLeaderboard(args []string, ctx *service.Context) {

	xplist := []models.Localxp{}

	err := ch.db.Select(&xplist, "SELECT userid, xp FROM localxp WHERE xp > 0 AND guildid = $1 ORDER BY xp DESC LIMIT 10;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	leaderboard := "```\n"

	place := 1

	for _, xp := range xplist {
		mem, err := ctx.Session.State.Member(ctx.Guild.ID, xp.Userid)
		if err != nil {
			leaderboard += fmt.Sprintf("#%v - User not here (%v) - %vxp\n", place, xp.Userid, xp.Xp)
		} else {
			leaderboard += fmt.Sprintf("#%v - %v#%v - %vxp\n", place, mem.User.Username, mem.User.Discriminator, xp.Xp)
		}
		place++
	}

	leaderboard += "```"

	ctx.Send(leaderboard)
}

func (ch *CommandHandler) globalXpLeaderboard(args []string, ctx *service.Context) {

	xplist := []models.Globalxp{}

	err := ch.db.Select(&xplist, "SELECT userid, xp FROM globalxp WHERE xp > 0 ORDER BY xp DESC LIMIT 10;")
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	leaderboard := "```\n"

	place := 1

	for _, xp := range xplist {
		user, err := ctx.Session.User(xp.Userid)
		if err != nil {
			return
		}

		leaderboard += fmt.Sprintf("#%v - %v#%v - %vxp\n", place, user.Username, user.Discriminator, xp.Xp)
		place++
	}

	leaderboard += "```"

	ctx.Send(leaderboard)
}

func (ch *CommandHandler) xpIgnoreChannel(args []string, ctx *service.Context) {

	var (
		err     error
		channel *discordgo.Channel
		chn     string
	)

	if len(args) > 1 {
		if strings.HasPrefix(args[1], "<#") && strings.HasSuffix(args[1], ">") {
			chn = args[1]
			chn = chn[2 : len(chn)-1]
		} else {
			chn = args[1]
		}
		channel, err = ctx.Session.State.Channel(chn)
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

	dbigch := &models.Xpignoredchannel{}

	err = ch.db.Get(dbigch, "SELECT channelid FROM xpignoredchannels WHERE channelid = $1;", channel.ID)
	switch err {
	case nil:
		ch.db.Exec("DELETE FROM xpignoredchannels WHERE channelid=$1;", channel.ID)
		ctx.Send(fmt.Sprintf("Removed %v from the xp ignore list.", channel.Mention()))
	case sql.ErrNoRows:
		ch.db.Exec("INSERT INTO xpignoredchannels (channelid) VALUES ($1);", channel.ID)
		ctx.Send(fmt.Sprintf("Added %v to the xp ignore list.", channel.Mention()))
	default:
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
	}
}

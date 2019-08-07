package commands

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) filterWord(args []string, ctx *service.Context) {

	if len(args) > 1 {

		phrase := strings.Join(args[1:], " ")

		if len(phrase) < 1 {
			return
		}

		dbf := &models.Filter{}

		err := ch.db.Get(dbf, "SELECT phrase FROM filters WHERE phrase = $1 AND guildid = $2;", phrase, ctx.Guild.ID)
		switch err {
		case nil:
			ch.db.Exec("DELETE FROM filters WHERE guildid=$1 AND phrase=$2;", ctx.Guild.ID, phrase)
			ctx.Send(fmt.Sprintf("Removed `%v` from the filter.", phrase))
		case sql.ErrNoRows:
			ch.db.Exec("INSERT INTO filters (guildid, phrase) VALUES ($1,$2);", ctx.Guild.ID, phrase)
			ctx.Send(fmt.Sprintf("Added `%v` to the filter.", phrase))
		default:
			ctx.Send("there was an error, please try again")
			ch.logger.Error("error", zap.Error(err))
		}
	}
}
func (ch *CommandHandler) filterInfo(args []string, ctx *service.Context) {

	dbch := []models.FilterIgnoreChannel{}
	err := ch.db.Select(&dbch, "SELECT channelid FROM filterignorechannels WHERE guildid=$1;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	filterlist := ""

	if len(dbch) < 1 {
		filterlist = "None"
	} else {
		for _, ch := range dbch {
			filterlist += fmt.Sprintf("<#%v>\n", ch.Channelid)
		}
	}

	dbg := &models.DiscordGuild{}
	err = ch.db.Get(dbg, "SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid=$1;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	embed := discordgo.MessageEmbed{
		Title: "Filter info",
		Color: dColorWhite,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Strikes currently enabled",
				Value:  fmt.Sprint(dbg.UseStrikes),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Max strikes",
				Value:  fmt.Sprint(dbg.MaxStrikes),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Ignored channels",
				Value:  filterlist,
				Inline: false,
			},
		},
	}

	ctx.SendEmbed(&embed)

}

func (ch *CommandHandler) filterWordList(args []string, ctx *service.Context) {

	fwl := []models.Filter{}
	err := ch.db.Select(&fwl, "SELECT * FROM filters WHERE guildid=$1;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	if len(fwl) < 1 {
		ctx.Send("The filter is empty.")
		return
	}

	filterlist := "```\nList of currently filtered phrases\n"

	for _, fw := range fwl {
		filterlist += fmt.Sprintf("- %v\n", fw.Phrase)
	}

	filterlist += "```"

	ctx.Send(filterlist)

}

func (ch *CommandHandler) clearFilter(args []string, ctx *service.Context) {

	_, err := ch.db.Exec("DELETE FROM filters WHERE guildid = $1;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	ctx.Send("filter was cleared")
}

func (ch *CommandHandler) useStrikes(args []string, ctx *service.Context) {
	dbg := &models.DiscordGuild{}
	err := ch.db.Get(dbg, "SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid=$1;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	if dbg.UseStrikes {
		ch.db.Exec("UPDATE discordguilds SET usestrikes = $1 WHERE guildid = $2;", false, ctx.Guild.ID)
		ctx.Send("Disabled strike system.")
	} else {
		ch.db.Exec("UPDATE discordguilds SET usestrikes = $1 WHERE guildid = $2;", true, ctx.Guild.ID)
		ctx.Send(fmt.Sprintf("Enabled strike system.\nMax strikes currently set to %v", dbg.MaxStrikes))
	}
}

func (ch *CommandHandler) setMaxStrikes(args []string, ctx *service.Context) {

	if len(args) != 2 {
		return
	}

	num, err := strconv.ParseInt(args[1], 0, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	if num <= 1 {
		num = 1
	}
	if num >= 10 {
		num = 10
	}

	ch.db.Exec("UPDATE discordguilds SET maxstrikes = $1 WHERE guildid = $2;", num, ctx.Guild.ID)

	ctx.Send(fmt.Sprintf("Set max strikes to %v.", num))
}

func (ch *CommandHandler) filterIgnoreChannel(args []string, ctx *service.Context) {

	var err error
	gch := ctx.Channel

	if len(args) > 1 {

		gchid := ""

		if strings.HasPrefix(args[1], "<#") && strings.HasSuffix(args[1], ">") {
			gchid = args[1]
			gchid = gchid[2 : len(gchid)-1]
		} else {
			gchid = args[1]
		}

		gch, err = ctx.Session.Channel(gchid)
		if err != nil {
			ctx.Send("Channel not found.")
			return
		}

		if gch.GuildID != ctx.Guild.ID {
			ctx.Send("Channel not found.")
			return
		}
	}

	dbf := &models.FilterIgnoreChannel{}

	err = ch.db.Get(dbf, "SELECT channelid FROM filterignorechannels WHERE guildid = $1 AND channelid = $2;", ctx.Guild.ID, gch.ID)
	if err != nil && err != sql.ErrNoRows {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		ch.db.Exec("INSERT INTO filterignorechannels (guildid, channelid) VALUES ($1,$2);", ctx.Guild.ID, gch.ID)
		ctx.Send(fmt.Sprintf("Added <#%v> to the list of filter ignored channels.", gch.ID))
	} else {
		ch.db.Exec("DELETE FROM filterignorechannels WHERE guildid = $1 AND channelid = $2;", ctx.Guild.ID, gch.ID)
		ctx.Send(fmt.Sprintf("Removed <#%v> from the list of filter ignored channels.", gch.ID))
	}

}

package commands

import (
	"database/sql"
	"fmt"
	"meido/models"
	"meido/service"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var FilterWord = Command{
	Name:          "filterword",
	Description:   "filters stuff.",
	Triggers:      []string{"m?filterword", "m?fw"},
	Usage:         "m?fw jeff\nm?filterword jeff",
	RequiredPerms: discordgo.PermissionManageMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) > 1 {

			phrase := strings.Join(args[1:], " ")

			if len(phrase) < 1 {
				return
			}

			dbf := models.Filter{}

			row := db.QueryRow("SELECT phrase FROM filters WHERE phrase = $1 AND guildid = $2;", phrase, ctx.Guild.ID)
			err := row.Scan(&dbf.Filter)
			if err != nil {
				if err == sql.ErrNoRows {
					db.Exec("INSERT INTO filters (guildid, phrase) VALUES ($1,$2);", ctx.Guild.ID, phrase)
					ctx.Send(fmt.Sprintf("Added `%v` to the filter.", phrase))
				}
			} else {
				db.Exec("DELETE FROM filters WHERE guildid=$1 AND phrase=$2;", ctx.Guild.ID, phrase)
				ctx.Send(fmt.Sprintf("Removed `%v` from the filter.", phrase))
			}
		}
	},
}

var FilterInfo = Command{
	Name:          "filterinfo",
	Description:   "Shows filter info.",
	Triggers:      []string{"m?filterinfo", "m?fi"},
	Usage:         "m?filterinfo\nm?fi",
	RequiredPerms: discordgo.PermissionManageMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		rows, err := db.Query("SELECT channelid FROM filterignorechannels WHERE guildid=$1;", ctx.Guild.ID)
		if err != nil {
			if err == sql.ErrNoRows {
			} else {
				ctx.Send("error occured: " + err.Error())
				return
			}
		}

		filterlist := ""

		for rows.Next() {
			dbch := models.FilterIgnoreChannel{}

			err = rows.Scan(&dbch.Channelid)
			if err != nil {
				return
			}

			filterlist += fmt.Sprintf("<#%v>\n", dbch.Channelid)
		}

		if filterlist == "" {
			filterlist = "None"
		}

		row := db.QueryRow("SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid=$1;", ctx.Guild.ID)
		dbg := models.DiscordGuild{}
		err = row.Scan(&dbg.UseStrikes, &dbg.MaxStrikes)
		if err != nil {
			if err == sql.ErrNoRows {
			} else {
				ctx.Send("error occured: " + err.Error())
				return
			}
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

	},
}

var FilterWordList = Command{
	Name:          "filterwordlist",
	Description:   "Shows filtered words.",
	Triggers:      []string{"m?filterwordlist", "m?fwl"},
	Usage:         "m?filterwordlist\nm?fwl",
	RequiredPerms: discordgo.PermissionSendMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		rows, err := db.Query("SELECT * FROM filters WHERE guildid=$1;", ctx.Guild.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.Send("The filter is empty.")
				return
			}
		}

		filterlist := "```\nList of currently filtered phrases\n"

		for rows.Next() {
			f := models.Filter{}

			err = rows.Scan(&f.Uid, &f.Guildid, &f.Filter)
			if err != nil {
				return
			}

			filterlist += fmt.Sprintf("- %v\n", f.Filter)
		}

		filterlist += "```"

		ctx.Send(filterlist)

	},
}

var ClearFilter = Command{
	Name:          "clearfilter",
	Description:   "Clears filtered words.",
	Triggers:      []string{"m?clearfilter"},
	Usage:         "m?clearfilter",
	RequiredPerms: discordgo.PermissionManageMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		res, err := db.Exec("DELETE FROM filters WHERE guildid = $1;", ctx.Guild.ID)
		if err != nil {
			return
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return
		}

		if affected == 0 {
			ctx.Send("The filter is empty.")
		} else {
			ctx.Send("Cleared the filter.")
		}
	},
}

var UseStrikes = Command{
	Name:          "UseStrikes",
	Description:   "Toggles strike system.",
	Triggers:      []string{"m?usestrikes"},
	Usage:         "m?usestrikes",
	RequiredPerms: discordgo.PermissionManageMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {
		row := db.QueryRow("SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid=$1;", ctx.Guild.ID)
		dbg := models.DiscordGuild{}
		err := row.Scan(&dbg.UseStrikes, &dbg.MaxStrikes)
		if err != nil {
			ctx.Send("error occured", err)
			return
		}
		if dbg.UseStrikes {
			db.Exec("UPDATE discordguilds SET usestrikes = $1 WHERE guildid = $2;", false, ctx.Guild.ID)
			ctx.Send("Disabled strike system.")
		} else {
			db.Exec("UPDATE discordguilds SET usestrikes = $1 WHERE guildid = $2;", true, ctx.Guild.ID)
			ctx.Send(fmt.Sprintf("Enabled strike system.\nMax strikes currently set to %v", dbg.MaxStrikes))
		}
	},
}

var SetMaxStrikes = Command{
	Name:          "SetMaxStrikes",
	Description:   "Sets max strikes. Max 10.",
	Triggers:      []string{"m?maxstrikes"},
	Usage:         "m?maxstrikes 5",
	RequiredPerms: discordgo.PermissionManageMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) != 2 {
			return
		}

		num, err := strconv.ParseInt(args[1], 0, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		if num <= 0 {
			num = 0
		}
		if num >= 10 {
			num = 10
		}
		/*

			row := db.QueryRow("SELECT maxstrikes FROM discordguilds WHERE guildid=$1;", ctx.Guild.ID)
			dbg := models.DiscordGuild{}
			err = row.Scan(&dbg.MaxStrikes)
			if err != nil {
				ctx.Send("error occured", err)
				return
			}
		*/
		db.Exec("UPDATE discordguilds SET maxstrikes = $1 WHERE guildid = $2;", num, ctx.Guild.ID)

		ctx.Send(fmt.Sprintf("Set max strikes to %v.", num))
	},
}

var FilterIgnoreChannel = Command{
	Name:          "filterignorechannel",
	Description:   "sets a channel to be ignored by filter.",
	Triggers:      []string{"m?filterignorechannel", "m?figch"},
	Usage:         "m?figch\nm?figch 393558442977263619\nm?filterignorechannel #gamers",
	RequiredPerms: discordgo.PermissionManageMessages,
	//RequiresOwner: true,
	Execute: func(args []string, ctx *service.Context) {

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

		dbf := models.FilterIgnoreChannel{}

		row := db.QueryRow("SELECT channelid FROM filterignorechannels WHERE guildid = $1 AND channelid = $2;", ctx.Guild.ID, gch.ID)
		err = row.Scan(&dbf.Channelid)
		if err != nil {
			if err == sql.ErrNoRows {
				db.Exec("INSERT INTO filterignorechannels (guildid, channelid) VALUES ($1,$2);", ctx.Guild.ID, gch.ID)
				ctx.Send(fmt.Sprintf("Added <#%v> to the list of filter ignored channels.", gch.ID))
			}
		} else {
			db.Exec("DELETE FROM filterignorechannels WHERE guildid = $1 AND channelid = $2;", ctx.Guild.ID, gch.ID)
			ctx.Send(fmt.Sprintf("Removed <#%v> from the list of filter ignored channels.", gch.ID))
		}
	},
}

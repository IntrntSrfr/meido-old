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

func (ch *CommandHandler) myRole(args []string, ctx *service.Context) {

	var err error

	// m?myrole color #123123
	if len(args) > 2 {

		ur := &models.Userrole{}
		err = ch.db.Get(ur, "SELECT * FROM userroles WHERE guildid=$1 AND userid=$2", ctx.Guild.ID, ctx.User.ID)
		if err != nil && err != sql.ErrNoRows {
			ctx.Send("there was an error, please try again")
			ch.logger.Error("error", zap.Error(err))
			return
		} else if err == sql.ErrNoRows {
			ctx.Send("No custom role set.")
			return
		}

		var oldRole *discordgo.Role

		for _, role := range ctx.Guild.Roles {
			if role.ID == ur.Roleid {
				oldRole = role
			}
		}

		switch args[1] {
		case "name":

			newName := strings.Join(args[2:], " ")

			_, err = ctx.Session.GuildRoleEdit(ctx.Guild.ID, oldRole.ID, newName, oldRole.Color, oldRole.Hoist, oldRole.Permissions, oldRole.Mentionable)
			if err != nil {
				if strings.Contains(err.Error(), strconv.Itoa(discordgo.ErrCodeMissingPermissions)) {
					ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Missing permissions.", Color: dColorRed})
					return
				}
				ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Some error occured: `" + err.Error() + "`.", Color: dColorRed})
				return
			}

			embed := discordgo.MessageEmbed{
				Color:       int(oldRole.Color),
				Description: fmt.Sprintf("Role name changed from %v to %v", oldRole.Name, newName),
			}
			ctx.SendEmbed(&embed)

		case "color":

			if strings.HasPrefix(args[2], "#") {
				args[2] = args[2][1:]
			}

			color, err := strconv.ParseInt("0x"+args[2], 0, 64)
			if err != nil {
				ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Invalid color code.", Color: dColorRed})
				return
			}

			_, err = ctx.Session.GuildRoleEdit(ctx.Guild.ID, oldRole.ID, oldRole.Name, int(color), oldRole.Hoist, oldRole.Permissions, oldRole.Mentionable)
			if err != nil {
				if strings.Contains(err.Error(), strconv.Itoa(discordgo.ErrCodeMissingPermissions)) {
					ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Missing permissions.", Color: dColorRed})
					return
				}
				ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Invalid color code.", Color: dColorRed})
				return
			}

			embed := discordgo.MessageEmbed{
				Color:       int(color),
				Description: fmt.Sprintf("Color changed from #%v to #%v", FullHex(fmt.Sprintf("%X", oldRole.Color)), FullHex(fmt.Sprintf("%X", color))),
			}
			ctx.SendEmbed(&embed)
		default:
		}

		return
	}

	// m?myrole or m?myrole 123123123123
	if len(args) > 0 {

		var target *discordgo.Member

		if len(args) > 1 {

			if len(ctx.Message.Mentions) >= 1 {
				target, err = ctx.Session.State.Member(ctx.Guild.ID, ctx.Message.Mentions[0].ID)
				if err != nil {
					//s.ChannelMessageSend(ch.ID, err.Error())
					return
				}
			} else {
				target, err = ctx.Session.State.Member(ctx.Guild.ID, args[1])
				if err != nil {
					//s.ChannelMessageSend(ch.ID, err.Error())
					return
				}
			}
		}

		if target == nil {
			target, err = ctx.Session.State.Member(ctx.Guild.ID, ctx.User.ID)
			if err != nil {
				//s.ChannelMessageSend(ch.ID, err.Error())
				return
			}
		}

		ur := &models.Userrole{}
		err = ch.db.Get(ur, "SELECT * FROM userroles WHERE guildid=$1 AND userid=$2", ctx.Guild.ID, target.User.ID)
		if err != nil && err != sql.ErrNoRows {
			ctx.Send("there was an error, please try again")
			ch.logger.Error("error", zap.Error(err))
			return
		} else if err == sql.ErrNoRows {
			ctx.Send("No custom role set.")
			return
		}

		var customRole *discordgo.Role

		for i := range ctx.Guild.Roles {
			role := ctx.Guild.Roles[i]

			if role.ID == ur.Roleid {
				customRole = role
			}
		}

		if customRole == nil {
			ctx.Send("the custom role is broken, wait for someone to fix it or try setting a new userrole")
			return
		}

		embed := discordgo.MessageEmbed{
			Color: int(customRole.Color),
			Title: fmt.Sprintf("Custom role for %v", target.User.String()),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Name",
					Value:  customRole.Name,
					Inline: true,
				},
				{
					Name:   "Color",
					Value:  fmt.Sprintf("#" + FullHex(fmt.Sprintf("%X", customRole.Color))),
					Inline: true,
				},
			},
		}
		ctx.SendEmbed(&embed)

		return
	}
}

package commands

import (
	"database/sql"
	"fmt"
	"meido-test/models"
	"meido-test/service"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var MyRole = Command{
	Name:          "myrole",
	Description:   "Gets information about a custom role, or lets the owner of the role edit its name or color.",
	Triggers:      []string{"m?myrole"},
	Usage:         "m?myrole\nm?myrole 163454407999094786\nm?myrole color c0ff33\nm?myrole name kumiko",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		var err error

		if len(args) >= 2 {

			if args[1] == "color" {
				if len(args) != 3 {
					return
				}

				u := ctx.Message.Author

				g, err := ctx.Session.Guild(ctx.Guild.ID)
				if err != nil {
					ctx.Send(err.Error())
					return
				}

				row := db.QueryRow("SELECT * FROM userroles WHERE guildid=$1 AND userid=$2", g.ID, u.ID)

				ur := models.Userrole{}

				err = row.Scan(&ur.Uid,
					&ur.Guildid,
					&ur.Userid,
					&ur.Roleid)
				if err != nil {
					if err == sql.ErrNoRows {
						ctx.Send("You dont have a custom role set.")
					}
					return
				}

				if strings.HasPrefix(args[2], "#") {
					args[2] = args[2][1:]
				}

				color, err := strconv.ParseInt("0x"+args[2], 0, 64)
				if err != nil {
					ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Invalid color code.", Color: dColorRed})
					return
				}

				var oldRole *discordgo.Role

				for i := range g.Roles {
					role := g.Roles[i]

					if role.ID == ur.Roleid {
						oldRole = role
						_, err = ctx.Session.GuildRoleEdit(g.ID, role.ID, role.Name, int(color), role.Hoist, role.Permissions, role.Mentionable)
						if err != nil {
							if strings.Contains(err.Error(), strconv.Itoa(discordgo.ErrCodeMissingPermissions)) {
								ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Missing permissions.", Color: dColorRed})
								return
							}
							ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Invalid color code.", Color: dColorRed})
							return
						}
					}
				}

				embed := discordgo.MessageEmbed{
					Color:       int(color),
					Description: fmt.Sprintf("Color changed from #%v to #%v", fullHex(fmt.Sprintf("%X", oldRole.Color)), fullHex(fmt.Sprintf("%X", color))),
				}
				ctx.SendEmbed(&embed)
			}

			if args[1] == "name" {

				if len(args) < 3 {
					return
				}

				newName := strings.Join(args[2:], " ")

				u := ctx.Message.Author

				g, err := ctx.Session.Guild(ctx.Guild.ID)
				if err != nil {
					ctx.Send(err.Error())
					return
				}

				row := db.QueryRow("SELECT * FROM userroles WHERE guildid=$1 AND userid=$2", g.ID, u.ID)

				ur := models.Userrole{}

				err = row.Scan(&ur.Uid,
					&ur.Guildid,
					&ur.Userid,
					&ur.Roleid)
				if err != nil {
					if err == sql.ErrNoRows {
						ctx.Send("You dont have a custom role set.")
					}
					return
				}

				var oldRole *discordgo.Role

				for i := range g.Roles {
					role := g.Roles[i]

					if role.ID == ur.Roleid {
						oldRole = role
						_, err = ctx.Session.GuildRoleEdit(g.ID, role.ID, newName, role.Color, role.Hoist, role.Permissions, role.Mentionable)
						if err != nil {
							if strings.Contains(err.Error(), strconv.Itoa(discordgo.ErrCodeMissingPermissions)) {
								ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Missing permissions.", Color: dColorRed})
								return
							}
							ctx.SendEmbed(&discordgo.MessageEmbed{Description: "Some error occured: `" + err.Error() + "`.", Color: dColorRed})
							return
						}
					}
				}

				embed := discordgo.MessageEmbed{
					Color:       int(oldRole.Color),
					Description: fmt.Sprintf("Role name changed from %v to %v", oldRole.Name, newName),
				}
				ctx.SendEmbed(&embed)
			}
		}
		var targetUser *discordgo.User

		if len(args) > 1 {

			if len(ctx.Message.Mentions) >= 1 {
				targetUser = ctx.Message.Mentions[0]
			} else {
				targetUser, err = ctx.Session.User(args[1])
				if err != nil {
					//s.ChannelMessageSend(ch.ID, err.Error())
					return
				}
			}
		}

		if targetUser == nil {
			targetUser = ctx.Message.Author
		}

		if targetUser.Bot {
			ctx.Send("Bots dont get to join the fun")
			return
		}

		u := targetUser

		g, err := ctx.Session.Guild(ctx.Guild.ID)
		if err != nil {
			ctx.Send(err.Error())
			return
		}

		row := db.QueryRow("SELECT * FROM userroles WHERE guildid=$1 AND userid=$2", g.ID, u.ID)

		ur := models.Userrole{}

		err = row.Scan(&ur.Uid,
			&ur.Guildid,
			&ur.Userid,
			&ur.Roleid)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.Send("No custom role set.")
			}
			return
		}

		var customRole *discordgo.Role

		for i := range g.Roles {
			role := g.Roles[i]

			if role.ID == ur.Roleid {
				customRole = role
			}
		}

		embed := discordgo.MessageEmbed{
			Color: int(customRole.Color),
			Title: fmt.Sprintf("Custom role for %v", u.String()),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Name",
					Value:  customRole.Name,
					Inline: true,
				},
				{
					Name:   "Color",
					Value:  fmt.Sprintf("#" + fullHex(fmt.Sprintf("%X", customRole.Color))),
					Inline: true,
				},
			},
		}
		ctx.SendEmbed(&embed)
	},
}

func fullHex(hex string) string {
	i := len(hex)

	for i < 6 {
		hex = "0" + hex
		i++
	}

	return hex
}

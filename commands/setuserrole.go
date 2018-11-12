package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var SetUserRole = Command{
	Name:          "setuserrole",
	Description:   "Sets a users custom role. First provide the user, followed by the role.",
	Triggers:      []string{"m?setuserrole"},
	Usage:         "m?setuserrole 163454407999094786 kumiko",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) < 3 {
			return
		}

		perms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Channel.ID)
		if err != nil {
			return
		}

		if perms&discordgo.PermissionManageRoles == 0 {
			ctx.SendEmbed(&discordgo.MessageEmbed{Color: dColorRed, Description: "You do not have the required permissions."})
			return
		}

		var targetUser *discordgo.User

		if len(ctx.Message.Mentions) >= 1 {
			targetUser = ctx.Message.Mentions[0]
		} else {
			targetUser, err = ctx.Session.User(args[1])
			if err != nil {
				ctx.Send(err.Error())
				return
			}
		}

		if targetUser.Bot {
			ctx.Send("Bots dont get to join the fun")
			return
		}

		g, err := ctx.Session.Guild(ctx.Guild.ID)
		if err != nil {
			ctx.Send(err.Error())
			return
		}

		var selectedRole *discordgo.Role

		for i := range g.Roles {
			role := g.Roles[i]

			if role.ID == args[2] {
				selectedRole = role
			} else if strings.ToLower(role.Name) == strings.ToLower(strings.Join(args[2:], " ")) {
				selectedRole = role
			}
		}

		if selectedRole == nil {
			ctx.Send("Role not found")
			return
		}

		var lastinsertid int
		err = db.QueryRow("INSERT INTO userroles(guildid, userid, roleid) VALUES($1, $2, $3) returning uid", g.ID, targetUser.ID, selectedRole.ID).Scan(&lastinsertid)
		if err != nil {
			ctx.Send(err.Error())
			return
		}

		ctx.Send(fmt.Sprintf("Bound role **%v** to user **%v#%v**", selectedRole.Name, targetUser.Username, targetUser.Discriminator))

	},
}

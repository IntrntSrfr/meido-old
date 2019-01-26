package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/models"
	"github.com/intrntsrfr/meido/service"

	"github.com/bwmarrin/discordgo"
)

var SetUserRole = Command{
	Name:          "setuserrole",
	Description:   "Sets a users custom role. First provide the user, followed by the role.",
	Triggers:      []string{"m?setuserrole"},
	Usage:         "m?setuserrole 163454407999094786 kumiko",
	RequiredPerms: discordgo.PermissionManageRoles,
	Execute: func(args []string, ctx *service.Context) {

		var err error

		if len(args) < 3 {
			return
		}

		var targetUser *discordgo.Member

		if len(ctx.Message.Mentions) >= 1 {
			targetUser, err = ctx.Session.State.Member(ctx.Guild.ID, ctx.Message.Mentions[0].ID)
			if err != nil {
				//s.ChannelMessageSend(ch.ID, err.Error())
				return
			}
		} else {
			targetUser, err = ctx.Session.State.Member(ctx.Guild.ID, args[1])
			if err != nil {
				//s.ChannelMessageSend(ch.ID, err.Error())
				return
			}
		}
		if targetUser.User.Bot {
			ctx.Send("Bots dont get to join the fun")
			return
		}

		g, err := ctx.Session.State.Guild(ctx.Guild.ID)
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

		row := db.QueryRow("SELECT COUNT(*) FROM userroles WHERE guildid=$1 AND userid=$2 AND roleid=$3;", g.ID, targetUser.User.ID, selectedRole.ID)
		userrolecount := 0
		err = row.Scan(&userrolecount)
		if err != nil {
			ctx.Send("error occured", err)
			return
		}

		if userrolecount <= 0 {
			_, err = db.Exec("INSERT INTO userroles(guildid, userid, roleid) VALUES($1, $2, $3);", g.ID, targetUser.User.ID, selectedRole.ID)
			if err != nil {
				ctx.Send(err.Error())
				return
			}
			ctx.Send(fmt.Sprintf("Bound role **%v** to user **%v#%v**", selectedRole.Name, targetUser.User.Username, targetUser.User.Discriminator))
		} else {
			_, err = db.Exec("DELETE FROM userroles WHERE guildid=$1 AND userid=$2 AND roleid=$3;", g.ID, targetUser.User.ID, selectedRole.ID)
			if err != nil {
				ctx.Send(err.Error())
				return
			}
			ctx.Send(fmt.Sprintf("Unbound role **%v** from user **%v#%v**", selectedRole.Name, targetUser.User.Username, targetUser.User.Discriminator))
		}
	},
}

var ListUserRoles = Command{
	Name:          "listuserroles",
	Description:   "Sets a users custom role. First provide the user, followed by the role.",
	Triggers:      []string{"m?listuserroles"},
	Usage:         "m?listuserroles",
	RequiredPerms: discordgo.PermissionManageRoles,
	Execute: func(args []string, ctx *service.Context) {
		rows, err := db.Query("SELECT roleid, userid FROM userroles WHERE guildid=$1;", ctx.Guild.ID)
		if err != nil {
			return
		}
		text := fmt.Sprintf("Userroles in %v\n\n\n", ctx.Guild.Name)
		for rows.Next() {
			count := 0
			dbur := models.Userrole{}
			err := rows.Scan(&dbur.Roleid, &dbur.Userid)
			if err != nil {
				fmt.Println(err)
				continue
			}
			role, err := ctx.Session.State.Role(ctx.Guild.ID, dbur.Roleid)
			if err != nil {
				fmt.Println(err)
				continue
			}
			mem, err := ctx.Session.State.Member(ctx.Guild.ID, dbur.Userid)
			if err != nil {
				fmt.Println(err)
				text += fmt.Sprintf("Role #%v: %v (%v) | Bound user: %v - User no longer in guild.", count, role.Name, role.ID, dbur.Userid)
			} else {
				text += fmt.Sprintf("Role #%v: %v (%v) | Bound user: %v (%v)", count, role.Name, role.ID, mem.User.String(), mem.User.ID)
			}
			count++
		}
		fmt.Println(text)
		link, err := OWOApi.Upload(text)
		if err != nil {
			ctx.Send("Error getting user roles.")
		}
		ctx.Send(fmt.Sprintf("User roles in %v\n%v", ctx.Guild.Name, link))
	},
}

package commands

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) setUserRole(args []string, ctx *service.Context) {

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

	userrole := &models.Userrole{}

	err = ch.db.Get(userrole, "SELECT * FROM userroles WHERE guildid=$1 AND userid=$2", g.ID, targetUser.User.ID)
	switch err {
	case nil:
		if selectedRole.ID == userrole.Roleid {
			ch.db.Exec("DELETE FROM userroles WHERE guildid=$1 AND userid=$2 AND roleid=$3;", g.ID, targetUser.User.ID, selectedRole.ID)
			ctx.Send(fmt.Sprintf("Unbound role **%v** from user **%v**", selectedRole.Name, targetUser.User.String()))
		} else {
			ch.db.Exec("UPDATE userroles SET roleid=$1 WHERE guildid=$2 AND userid=$3", selectedRole.ID, g.ID, targetUser.User.ID)
			ctx.Send(fmt.Sprintf("Updated userrole for **%v** to **%v**", targetUser.User.String(), selectedRole.Name))
		}
	case sql.ErrNoRows:
		ch.db.Exec("INSERT INTO userroles(guildid, userid, roleid) VALUES($1, $2, $3);", g.ID, targetUser.User.ID, selectedRole.ID)
		ctx.Send(fmt.Sprintf("Bound role **%v** to user **%v#%v**", selectedRole.Name, targetUser.User.Username, targetUser.User.Discriminator))
	default:
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
	}
}

func (ch *CommandHandler) listUserRoles(args []string, ctx *service.Context) {
	userroles := []models.Userrole{}

	err := ch.db.Select(&userroles, "SELECT roleid, userid FROM userroles WHERE guildid=$1;", ctx.Guild.ID)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	text := fmt.Sprintf("Userroles in %v\n\n", ctx.Guild.Name)
	count := 0
	for _, ur := range userroles {

		role, err := ctx.Session.State.Role(ctx.Guild.ID, ur.Roleid)
		if err != nil {
			fmt.Println(err)
			continue
		}
		mem, err := ctx.Session.State.Member(ctx.Guild.ID, ur.Userid)
		if err != nil {
			text += fmt.Sprintf("Role #%v: %v (%v) | Bound user: %v - User no longer in guild.\n", count, role.Name, role.ID, ur.Userid)
		} else {
			text += fmt.Sprintf("Role #%v: %v (%v) | Bound user: %v (%v)\n", count, role.Name, role.ID, mem.User.String(), mem.User.ID)
		}
		count++
	}

	link, err := ch.owo.Upload(text)
	if err != nil {
		ctx.Send("Error getting user roles.")
		return
	}
	ctx.Send(fmt.Sprintf("User roles in %v\n%v", ctx.Guild.Name, link))
}

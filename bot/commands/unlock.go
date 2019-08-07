package commands

import (
	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) unlock(args []string, ctx *service.Context) {
	var erole *discordgo.Role

	for _, val := range ctx.Guild.Roles {
		if val.ID == ctx.Guild.ID {
			erole = val
		}
	}

	var eperms *discordgo.PermissionOverwrite

	for _, val := range ctx.Channel.PermissionOverwrites {
		if val.ID == erole.ID {
			eperms = val
		}
	}

	if erole == nil || eperms == nil {
		return
	}

	if eperms.Allow&discordgo.PermissionSendMessages == 0 && eperms.Deny&discordgo.PermissionSendMessages == 0 {
		// DEFAULT
		ctx.Send("Channel is already unlocked.")
	} else if eperms.Allow&discordgo.PermissionSendMessages != 0 && eperms.Deny&discordgo.PermissionSendMessages == 0 {
		// IS ALLOWED
		ctx.Send("Channel is already unlocked.")
	} else if eperms.Allow&discordgo.PermissionSendMessages == 0 && eperms.Deny&discordgo.PermissionSendMessages != 0 {
		// IS DENIED
		err := ctx.Session.ChannelPermissionSet(
			ctx.Channel.ID,
			erole.ID,
			"role",
			eperms.Allow,
			eperms.Deny-discordgo.PermissionSendMessages,
		)
		if err != nil {
			ctx.Send("Could not unlock channel")
			return
		}
		ctx.Send("Channel unlocked")
	}
}

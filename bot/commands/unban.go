package commands

import (
	"fmt"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

func (ch *CommandHandler) unban(args []string, ctx *service.Context) {
	if len(args) <= 1 {
		return
	}

	userID := args[1]

	err := ctx.Session.GuildBanDelete(ctx.Guild.ID, userID)
	if err != nil {
		return
	}

	targetUser, err := ctx.Session.User(userID)
	if err != nil {
		return
	}

	embed := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**Unbanned** %v - %v#%v (%v)", targetUser.Mention(), targetUser.Username, targetUser.Discriminator, targetUser.ID),
		Color:       dColorGreen,
	}

	ctx.SendEmbed(embed)
}

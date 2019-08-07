package commands

import (
	"fmt"
	"strconv"

	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"
)

func (ch *CommandHandler) hackban(args []string, ctx *service.Context) {
	if len(args) < 2 {
		return
	}

	userList := []string{}

	for _, mention := range ctx.Message.Mentions {
		userList = append(userList, mention.ID)
	}

	for _, userID := range args[1:] {
		userList = append(userList, userID)
	}

	badbans := 0
	badIDs := 0

	for _, userID := range userList {
		_, err := strconv.Atoi(userID)
		if err != nil {
			badIDs++
			continue
		}
		err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, userID, fmt.Sprintf("[%v] - Hackban", ctx.User.String()), 7)
		if err != nil {
			badbans++
			continue
		}

		_, err = ch.db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", userID, ctx.Guild.ID)
		if err != nil {
			ch.logger.Error("error", zap.Error(err))
			continue
		}
	}

	ctx.Send(fmt.Sprintf("Banned %v out of %v users provided.", len(userList)-badbans-badIDs, len(userList)-badIDs))
}

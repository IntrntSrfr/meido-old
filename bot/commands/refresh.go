package commands

import (
	"github.com/intrntsrfr/meido/bot/database"
	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"
)

func (ch *CommandHandler) refresh(args []string, ctx *service.Context) {

	ctx.Send("Refreshing. This might take a second")

	err := database.Refresh(ch.db, ch.logger, ctx.Session.State.Guilds)
	if err != nil {
		ctx.Send("there was an error, please try again")
		ch.logger.Error("error", zap.Error(err))
		return
	}

	ctx.Send("Refresh finished, things should be fine now")
}

package commands

import (
	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"
	"go.uber.org/zap"
)

func (ch *CommandHandler) refresh(args []string, ctx *service.Context) {

	ctx.Send("Refreshing. This might take a second")

	for _, g := range ctx.Session.State.Guilds {

		userroles := []models.Userrole{}

		ch.db.Select(&userroles, "SELECT * FROM userroles WHERE guildid=$1", g.ID)

		for _, ur := range userroles {

			hasRole := false

			for _, gr := range g.Roles {
				if ur.Roleid == gr.ID {
					hasRole = true
				}
			}

			if !hasRole {
				_, err := ch.db.Exec("DELETE FROM userroles WHERE guildid=$1 AND roleid=$2", g.ID, ur.Roleid)
				if err != nil {
					ch.logger.Error("error", zap.Error(err))
				}
			}
		}
	}

	ctx.Send("Refresh finished, things should be fine now")
}

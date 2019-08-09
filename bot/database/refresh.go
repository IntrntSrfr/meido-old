package database

import (
	"time"

	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/jmoiron/sqlx"
)

func Refresh(db *sqlx.DB, z *zap.Logger, guilds []*discordgo.Guild) error {

	z.Info("running refresh")

	for _, g := range guilds {
		userroles := []models.Userrole{}

		err := db.Select(&userroles, "SELECT * FROM userroles WHERE guildid=$1", g.ID)
		if err != nil {
			return err
		}

		for _, ur := range userroles {

			hasRole := false

			for _, gr := range g.Roles {
				if ur.Roleid == gr.ID {
					hasRole = true
				}
			}

			if !hasRole {
				_, err := db.Exec("DELETE FROM userroles WHERE uid=$1", ur.Uid)
				if err != nil {
					z.Error("error", zap.Error(err))
				}
			}
		}

		strikes := []models.Strikes{}

		err = db.Select(&strikes, "SELECT * FROM strikes WHERE guildid=$1", g.ID)
		if err != nil {
			return err
		}

		for _, strike := range strikes {
			if strike.Tstamp.Unix() < time.Now().Add(time.Hour*24*30*-1).Unix() {
				_, err := db.Exec("DELETE FROM strikes WHERE uid=$1", strike.Uid)
				if err != nil {
					z.Error("error", zap.Error(err))
				}
			}
		}
	}
	return nil
}

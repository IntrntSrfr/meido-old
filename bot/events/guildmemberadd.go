package events

import (
	"database/sql"
	"time"

	"github.com/intrntsrfr/meido/bot/models"
	"go.uber.org/zap"

	"github.com/bwmarrin/discordgo"
)

func (eh *EventHandler) guildMemberAddHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	if m.User.Bot {
		return
	}

	user := &models.DiscordUser{}
	err := eh.db.Get(user, "SELECT * FROM discordusers WHERE userid = $1", m.User.ID)
	if err != nil && err != sql.ErrNoRows {
		eh.logger.Error("error", zap.Error(err))
	} else if err == sql.ErrNoRows {

		currentTime := time.Now()

		_, err := eh.db.Exec("INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)", m.User.ID, m.User.Username, m.User.Discriminator, 0, currentTime, false, 0, currentTime)
		if err != nil {
			eh.logger.Error(err.Error())
		}
	}
}

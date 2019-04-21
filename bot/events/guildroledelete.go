package events

import (
	"github.com/intrntsrfr/meido/bot/models"

	"github.com/bwmarrin/discordgo"
)

func GuildRoleDeleteHandler(s *discordgo.Session, m *discordgo.GuildRoleDelete) {
	row := db.QueryRow("SELECT * FROM userroles WHERE guildid=$1 AND roleid=$2", m.GuildID, m.RoleID)

	ur := models.Userrole{}

	err := row.Scan(&ur.Uid,
		&ur.Guildid,
		&ur.Roleid,
		&ur.Userid)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	stmt, err := db.Prepare("DELETE FROM userroles WHERE guildid=$1 AND roleid=$2")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	_, err = stmt.Exec(m.GuildID, m.RoleID)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

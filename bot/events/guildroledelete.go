package events

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func (eh *EventHandler) guildRoleDeleteHandler(s *discordgo.Session, m *discordgo.GuildRoleDelete) {
	_, err := eh.db.Exec("DELETE FROM userroles WHERE guildid=$1 AND roleid=$2", m.GuildID, m.RoleID)
	if err != nil {
		eh.logger.Error("error", zap.Error(err))
	}
}

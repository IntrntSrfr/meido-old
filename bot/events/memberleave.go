package events

import "github.com/bwmarrin/discordgo"

func MemberLeaveHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	totalUsers--
}

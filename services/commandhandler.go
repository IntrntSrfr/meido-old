package services

import (
	"github.com/bwmarrin/discordgo"
)

//Command is the base for all commands
type Command struct {
	name        string
	aliases     []string
	description string
	usage       string
	function    func(s *discordgo.Session, m *discordgo.MessageCreate)
}

//Commands is the map for all commands
type Commands map[string]Command

func init() {
	commands := make(Commands)

	commands["jeff"] = Command{name: "jeff",
		aliases:     []string{},
		description: "does epic",
		usage:       "command",
		function: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			s.ChannelMessageSend(m.ChannelID, "jeff")
		}}
}

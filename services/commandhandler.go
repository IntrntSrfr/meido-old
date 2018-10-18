package services

import (
	"github.com/bwmarrin/discordgo"
)

//Command is the base for all commands
type Command struct {
	name     string
	aliases  []string
	function func(s *discordgo.Session, m *discordgo.MessageCreate)
}

//Commands is the map for all commands
type Commands map[string]Command

func init() {
	commands := make(Commands)
}

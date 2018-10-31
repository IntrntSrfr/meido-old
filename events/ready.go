package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, m *discordgo.Ready) {
	fmt.Println("Logged in.")
}

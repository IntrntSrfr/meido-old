package services

import (
	"meido-test/events"

	"github.com/bwmarrin/discordgo"
)

// AddHandlers does the job
func AddHandlers(s *discordgo.Session) {
	go s.AddHandler(events.ReadyHandler)
	go s.AddHandler(events.GuildCreateHandler)
}

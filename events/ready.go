package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, m *discordgo.Ready) {
	/*
		data := discordgo.UpdateStatusData{
			Game: &discordgo.Game{
				Type: discordgo.GameTypeWatching,
				Name: "22 jump street",
			},
		}

		s.UpdateStatusComplex(data)
	*/
	fmt.Println("Logged in.")
}

package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var totalUsers = 0

func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
/* 
	data := discordgo.UpdateStatusData{
		Game: &discordgo.Game{
			Type: discordgo.GameTypeWatching,
			Name: fmt.Sprintf("over all %v of you"),
		},
	}

	s.UpdateStatusComplex(data)
 */
	fmt.Println(fmt.Sprintf("Logged in as %v.", r.User.String()))
}

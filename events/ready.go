package events

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	totalUsers = 0
	timer      *time.Ticker
)

func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {

	oldTotalUsers := 0
	timer := time.NewTicker(15 * time.Second)
	go func() {
		for range timer.C {
			if totalUsers != oldTotalUsers {

				data := discordgo.UpdateStatusData{
					Game: &discordgo.Game{
						Type: discordgo.GameTypeWatching,
						Name: fmt.Sprintf("over all %v of you", totalUsers),
					},
				}
				s.UpdateStatusComplex(data)
				//fmt.Println(fmt.Sprintf("Status update - [%v users]", totalUsers))
				oldTotalUsers = totalUsers
			}
		}
	}()

	totalUsers = 0
	fmt.Println(fmt.Sprintf("Logged in as %v.", r.User.String()))
}

func DisconnectHandler(s *discordgo.Session, d *discordgo.Disconnect) {

	fmt.Println("Disconnected at: " + time.Now().String())
	timer = nil
}

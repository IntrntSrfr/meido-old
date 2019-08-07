package events

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	timer *time.Ticker
)

func (eh *EventHandler) readyHandler(s *discordgo.Session, r *discordgo.Ready) {

	timer := time.NewTicker(15 * time.Second)
	go func() {
		for range timer.C {
			memCount := 0
			oldMemCount := 0
			for _, g := range eh.client.State.Guilds {
				memCount += g.MemberCount
			}

			if memCount != oldMemCount {
				data := discordgo.UpdateStatusData{
					Game: &discordgo.Game{
						Type: discordgo.GameTypeWatching,
						Name: fmt.Sprintf("over all %v of you", memCount),
					},
				}
				s.UpdateStatusComplex(data)
				oldMemCount = memCount
				//fmt.Println(fmt.Sprintf("Status update - [%v users]", totalUsers))
			}
		}
	}()

	fmt.Println(fmt.Sprintf("Logged in as %v.", r.User.String()))
}

func (eh *EventHandler) disconnectHandler(s *discordgo.Session, d *discordgo.Disconnect) {

	fmt.Println("Disconnected at: " + time.Now().String())
	timer = nil
}

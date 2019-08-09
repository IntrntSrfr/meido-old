package events

import (
	"fmt"
	"time"

	"github.com/intrntsrfr/meido/bot/database"

	"github.com/bwmarrin/discordgo"
)

var (
	timer *time.Ticker
)

func (eh *EventHandler) readyHandler(s *discordgo.Session, r *discordgo.Ready) {

	statusTimer := time.NewTicker(time.Second * 15)
	refreshTimer := time.NewTicker(time.Minute * 10)

	go func() {
		for range statusTimer.C {
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

	go func() {
		for range refreshTimer.C {
			database.Refresh(eh.db, eh.logger, eh.client.State.Guilds)
		}
	}()

	fmt.Println(fmt.Sprintf("Logged in as %v.", r.User.String()))
}

func (eh *EventHandler) disconnectHandler(s *discordgo.Session, d *discordgo.Disconnect) {

	fmt.Println("Disconnected at: " + time.Now().String())
	timer = nil
}

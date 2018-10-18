package events

import (
	"fmt"

	"test/services"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println(r.User.Username)

	fmt.Println(services.Commands)
}

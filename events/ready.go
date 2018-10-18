package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Println(r.User.Username)
}

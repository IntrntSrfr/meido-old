package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

//GuildCreateHandler for when a guild is available
func GuildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	fmt.Println(g.Name)
}

package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name          string
	Aliases       []string
	Description   string
	Usage         string
	RequiredPerms int
	Function      func(context *service.Context)
}

type Commandmap map[string]Command

var comms = Commandmap{}

func LoadCommands(cmap *Commandmap) {

	cmap.RegisterCommand(Help)
	cmap.RegisterCommand(Ping)

	comms = *cmap

}

func GetCommandMap() Commandmap {

	return comms

}

func (cmap *Commandmap) RegisterCommand(cmd Command) {

	(*cmap)[cmd.Name] = cmd
}

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	context := service.NewContext(s, m.Message)
	//fmt.Println(context)
	//context.Load(s, m.Message)

	if m.Author.Bot {
		return
	}
	/*
		context.Send("jeff")
	*/
	//service.Send("jeff")

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	g, err := s.Guild(ch.GuildID)
	if err != nil {
		return
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return
	}

	perms, err := s.UserChannelPermissions(m.Author.ID, ch.ID)
	if err != nil {
		return
	}

	args := strings.Split(m.Content, " ")

	triggerCommand := ""
	for _, val := range comms {
		name := "m?" + val.Name

		if args[0] == name {
			triggerCommand = val.Name
		}

		for _, com := range val.Aliases {
			if args[0] == "m?"+com {
				triggerCommand = val.Name
			}
		}
	}
	if triggerCommand != "" {

		if cmd, ok := comms[triggerCommand]; ok {
			if perms&cmd.RequiredPerms == 0 {
				return
			}
			cmd.Function(&context)
			fmt.Println(fmt.Sprintf("Command executed\nCommand: %v\nUser: %v [%v]\nSource: %v [%v] - #%v [%v]\n", cmd.Name, m.Author.String(), m.Author.ID, g.Name, g.ID, ch.Name, ch.ID))
		}
	}

}

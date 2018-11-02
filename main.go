package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"meido-test/commands"
	"meido-test/events"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

//type commandmap map[string]commands.Command

type Config struct {
	Token string `json:"Token"`
}

var (
	comms  commands.Commandmap
	config Config
)

func main() {

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("Config file not found.")
		return
	}

	json.Unmarshal(file, &config)

	token := config.Token

	client, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println(err)
		return
	}

	comms = make(commands.Commandmap)

	commands.LoadCommands(&comms)

	AddHandlers(client)

	err = client.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	client.Close()
}

// AddHandlers does the job
func AddHandlers(s *discordgo.Session) {
	s.AddHandler(events.ReadyHandler)
	s.AddHandler(commands.MessageCreateHandler)
}

/*
func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	ch, err := s.Channel(m.ChannelID)
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
			if args[0] == com {
				triggerCommand = val.Name
			}
		}
	}
	if triggerCommand != "" {

		if cmd, ok := comms[triggerCommand]; ok {
			if perms&cmd.RequiredPerms == 0 {
				return
			}
			cmd.Function(s, m)
		}
	}
}
*/
/*
Command{
	Name:        "jeff",
	Description: "jeffe",
	Aliases:     []string{"jeffer", "jeffette"},
	Usage:       "m?jeff",
	Function: func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, "jeff")
	},
} */

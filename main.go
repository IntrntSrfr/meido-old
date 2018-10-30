package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type commandmap map[string]Command

var commands commandmap

func main() {

	token := "NDg1NzIwNzI1MDkzODc1NzI0.DqnVwg.zbBZIxVSHVQjnX0Aqt2ws4XucXE"
	client, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println(err)
		return
	}

	commands = make(commandmap)

	commandList := []Command{
		{
			Name:        "jeff",
			Description: "jeffe",
			Aliases:     []string{"jeffer", "jeffette"},
			Usage:       "m?jeff",
			Function: func(s *discordgo.Session, m *discordgo.MessageCreate) {
				list := "```\n"
				for _, val := range commands {
					list += val.Name + "\n"
				}
				list += "```"
				s.ChannelMessageSend(m.ChannelID, list)
			},
		},
	}

	commands.LoadCommands(commandList)

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
	s.AddHandler(messageCreateHandler)
}

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if cmd, ok := commands[strings.Split(m.Content, " ")[0]]; ok {
		cmd.Function(s, m)
	}
}

func (cmap *commandmap) LoadCommands(cmds []Command) {

	for i := range cmds {
		cmd := cmds[i]

		(*cmap)["m?"+cmd.Name] = cmd
	}
}

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Function    func(s *discordgo.Session, m *discordgo.MessageCreate)
}

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

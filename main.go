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
var commandList []Command

func main() {

	token := "NDg1NzIwNzI1MDkzODc1NzI0.DqnVwg.zbBZIxVSHVQjnX0Aqt2ws4XucXE"
	client, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println(err)
		return
	}

	commands = make(commandmap)

	//commandList = []Command{}

	commandList = append(commandList, Command{
		Name:          "jeff",
		Description:   "jeffe",
		Aliases:       []string{"jeffer", "jeffette"},
		Usage:         "m?jeff",
		RequiredPerms: discordgo.PermissionManageMessages,
		Function: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			list := "```\n"
			for _, val := range commands {
				list += val.Name + "\n"
			}
			list += "```"
			s.ChannelMessageSend(m.ChannelID, list)
		}})

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
	for _, val := range commands {
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

		if cmd, ok := commands[triggerCommand]; ok {
			if perms&cmd.RequiredPerms == 0 {
				return
			}
			cmd.Function(s, m)
		}
	}
	/*
	 */
}

func (cmap *commandmap) LoadCommands(cmds []Command) {

	for i := range cmds {

		cmd := cmds[i]

		(*cmap)[cmd.Name] = cmd
	}
}

type Command struct {
	Name          string
	Aliases       []string
	Description   string
	Usage         string
	RequiredPerms int
	Function      func(s *discordgo.Session, m *discordgo.MessageCreate)
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

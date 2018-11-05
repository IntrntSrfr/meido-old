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
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

//type commandmap map[string]commands.Command

type Config struct {
	Token            string `json:"Token"`
	ConnectionString string `json:"ConnectionString"`
}

type Bot struct {
	StartTime time.Time
}

var (
	comms  commands.Commandmap
	config Config
)

func main() {
	bot := Bot{}
	bot.Run()
}

func (b *Bot) Run() {
	b.StartTime = time.Now()

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
	/*
		db, err = sql.Open("postgres", config.Connectionstring)
		if err != nil {
			panic("could not connect to db " + err.Error())
		} */

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

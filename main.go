package main

import (
	"database/sql"
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

type Config struct {
	Token            string   `json:"Token"`
	ConnectionString string   `json:"Connectionstring"`
	DmLogChannels    []string `json:"DmLogChannels"`
	OwnerIds         []string `json:"OwnerIds"`
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

	db, err := sql.Open("postgres", config.ConnectionString)
	if err != nil {
		panic("could not connect to db " + err.Error())
	}

	commands.Initialize(&config.OwnerIds, &config.DmLogChannels, db)
	events.Initialize(db)

	addHandlers(client)

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

func addHandlers(s *discordgo.Session) {
	go s.AddHandler(events.GuildAvailableHandler)
	go s.AddHandler(events.GuildRoleDeleteHandler)
	go s.AddHandler(events.MemberJoinedHandler)
	go s.AddHandler(events.MessageUpdateHandler)
	go s.AddHandler(events.ReadyHandler)
	go s.AddHandler(commands.MessageCreateHandler)
}

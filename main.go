package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/intrntsrfr/meido/commands"
	"github.com/intrntsrfr/meido/events"
	"github.com/intrntsrfr/meido/owo"

	"net/http"
	_ "net/http/pprof"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

type Config struct {
	Token            string   `json:"Token"`
	OWOToken         string   `json:"OWOToken"`
	ConnectionString string   `json:"Connectionstring"`
	DmLogChannels    []string `json:"DmLogChannels"`
	OwnerIds         []string `json:"OwnerIds"`
}

type Bot struct {
}

var config Config

func main() {
	bot := Bot{}
	bot.Run()
}

func (b *Bot) Run() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

	defer db.Close()

	OWOApi := owo.NewOWOClient(config.OWOToken)

	commands.Initialize(client, &config.OwnerIds, &config.DmLogChannels, db, OWOApi)
	events.Initialize(db)

	addHandlers(client)

	err = client.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func addHandlers(s *discordgo.Session) {
	s.AddHandler(events.GuildAvailableHandler)
	s.AddHandler(events.GuildRoleDeleteHandler)
	s.AddHandler(events.MemberJoinedHandler)
	s.AddHandler(events.MemberLeaveHandler)
	s.AddHandler(events.MessageUpdateHandler)
	s.AddHandler(events.ReadyHandler)
	s.AddHandler(events.DisconnectHandler)
	s.AddHandler(commands.MessageCreateHandler)
}

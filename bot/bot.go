package bot

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/intrntsrfr/owo"

	"github.com/bwmarrin/discordgo"
	"github.com/intrntsrfr/meido/bot/commands"
	"github.com/intrntsrfr/meido/bot/events"
	"go.uber.org/zap"
)

type Bot struct {
	logger    *zap.Logger
	db        *sql.DB
	client    *discordgo.Session
	config    *Config
	starttime time.Time
	owoAPI    *owo.OWOClient
}

func NewBot(Config *Config, Log *zap.Logger) (*Bot, error) {

	client, err := discordgo.New("Bot " + Config.Token)
	if err != nil {
		fmt.Println(err)
		Log.Error(err.Error())
		return nil, err
	}
	Log.Info("created discord client")

	psql, err := sql.Open("postgres", Config.ConnectionString)
	if err != nil {
		fmt.Println("could not connect to db " + err.Error())
		Log.Error(err.Error())
		return nil, err
	}
	Log.Info("Established postgres connection")

	OWOApi := owo.NewOWOClient(Config.OwoAPIKey)

	return &Bot{
		client:    client,
		db:        psql,
		config:    Config,
		logger:    Log,
		starttime: time.Now(),
		owoAPI:    OWOApi,
	}, nil
}

func (b *Bot) Close() {
	b.logger.Info("Shutting down bot.")
	b.db.Close()
	b.client.Close()
}

func (b *Bot) Run() error {
	commands.Initialize(b.client, &b.config.OwnerIds, &b.config.DmLogChannels, b.db, b.owoAPI, b.logger)
	events.Initialize(b.db, b.logger)
	b.addHandlers()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	return b.client.Open()
}

func (b *Bot) addHandlers() {
	b.client.AddHandler(events.GuildAvailableHandler)
	b.client.AddHandler(events.GuildRoleDeleteHandler)
	b.client.AddHandler(events.MemberJoinedHandler)
	b.client.AddHandler(events.MemberLeaveHandler)
	b.client.AddHandler(events.MessageUpdateHandler)
	b.client.AddHandler(events.ReadyHandler)
	b.client.AddHandler(events.DisconnectHandler)
	b.client.AddHandler(commands.MessageCreateHandler)
}

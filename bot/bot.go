package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/intrntsrfr/meido/bot/commands"
	"github.com/intrntsrfr/meido/bot/events"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Bot struct {
	client    *discordgo.Session
	db        *sqlx.DB
	logger    *zap.Logger
	eh        *events.EventHandler
	config    *Config
	starttime time.Time
}

func NewBot(Config *Config) (*Bot, error) {

	// creating zap logger
	z := zap.NewDevelopmentConfig()
	z.OutputPaths = []string{"./logs.txt"}
	z.ErrorOutputPaths = []string{"./logs.txt"}
	log, err := z.Build()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer log.Sync()
	log.Info("Logger construction succeeded")

	// making client
	client, err := discordgo.New("Bot " + Config.Token)
	if err != nil {
		fmt.Println(err)
		log.Error(err.Error())
		return nil, err
	}
	log.Info("created discord client")

	// opening psql connection
	psql, err := sqlx.Connect("postgres", Config.ConnectionString)
	if err != nil {
		fmt.Println("could not connect to db " + err.Error())
		log.Error(err.Error())
		return nil, err
	}
	log.Info("Established postgres connection")

	econfig := &events.Config{
		OwoToken:      Config.OwoToken,
		DmLogChannels: Config.DmLogChannels,
		OwnerIds:      Config.OwnerIds,
	}

	chconfig := &commands.Config{
		OwoToken:      Config.OwoToken,
		DmLogChannels: Config.DmLogChannels,
		OwnerIds:      Config.OwnerIds,
	}

	eventHandler := events.NewEventHandler(client, psql, log, econfig, chconfig)
	eventHandler.Initialize()

	return &Bot{
		client:    client,
		db:        psql,
		logger:    log,
		eh:        eventHandler,
		config:    Config,
		starttime: time.Now(),
	}, nil
}

func (b *Bot) Close() {
	b.logger.Info("Shutting down bot.")
	b.db.Close()
	b.client.Close()
}

func (b *Bot) Run() error {

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	return b.client.Open()
}

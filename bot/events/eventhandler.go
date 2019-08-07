package events

import (
	"math/rand"
	"time"

	"github.com/intrntsrfr/meido/bot/commands"
	"github.com/jmoiron/sqlx"

	"github.com/bwmarrin/discordgo"

	"go.uber.org/zap"
)

type EventHandler struct {
	ch            *commands.CommandHandler
	client        *discordgo.Session
	db            *sqlx.DB
	logger        *zap.Logger
	dmLogChannels []string
	ownerIds      []string
}

func NewEventHandler(c *discordgo.Session, psql *sqlx.DB, l *zap.Logger, ec *Config, cc *commands.Config) *EventHandler {

	commhandler := commands.NewCommandHandler(c, psql, l, cc)
	commhandler.Initialize()

	return &EventHandler{
		ch:            commhandler,
		client:        c,
		db:            psql,
		logger:        l,
		dmLogChannels: ec.DmLogChannels,
		ownerIds:      ec.OwnerIds,
	}
}

func (eh *EventHandler) Initialize() {
	eh.client.AddHandler(eh.messageCreateHandler)
	eh.client.AddHandler(eh.messageUpdateHandler)
	eh.client.AddHandler(eh.guildAvailableHandler)
	eh.client.AddHandler(eh.guildRoleDeleteHandler)
	eh.client.AddHandler(eh.guildMemberAddHandler)
	eh.client.AddHandler(eh.guildMembersChunkHandler)
	eh.client.AddHandler(eh.guildMemberRemoveHandler)
	eh.client.AddHandler(eh.readyHandler)
	eh.client.AddHandler(eh.disconnectHandler)
}

const (
	dColorRed    = 13107200
	dColorOrange = 15761746
	dColorLBlue  = 6410733
	dColorGreen  = 51200
	dColorWhite  = 16777215
)

func FullHex(hex string) string {

	i := len(hex)

	for i < 6 {
		hex = "0" + hex
		i++
	}

	return hex
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

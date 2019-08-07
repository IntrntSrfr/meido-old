package commands

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"
	"github.com/intrntsrfr/owo"

	"go.uber.org/zap"
)

func NewCommandHandler(s *discordgo.Session, DB *sqlx.DB, Logger *zap.Logger, cc *Config) *CommandHandler {

	o := owo.NewOWOClient(cc.OwoToken)

	return &CommandHandler{
		client:        s,
		botStartTime:  time.Now(),
		comms:         Commandmap{},
		db:            DB,
		dmLogChannels: cc.DmLogChannels,
		ownerIds:      cc.OwnerIds,
		owo:           o,
		logger:        Logger,
	}
}

func (ch *CommandHandler) GetCommandMap() Commandmap {
	return ch.comms
}

func (cmap *Commandmap) RegisterCommand(cmd Command) {
	if cmd, ok := (*cmap)[cmd.Name]; ok {
		fmt.Println("Conflicting Commands.", cmd.Name, (*cmap)[cmd.Name])
		os.Exit(0)
	}

	(*cmap)[cmd.Name] = cmd
}

func FullHex(hex string) string {

	i := len(hex)

	for ; i < 6; i++ {
		hex = "0" + hex
	}

	return hex
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (ch *CommandHandler) SetupProfile(target *discordgo.User, ctx *service.Context, rep int) {

	dbu := &models.DiscordUser{}

	currentTime := time.Now()

	err := ch.db.Get(dbu, "SELECT * FROM discordusers WHERE userid = $1", target.ID)
	if err != nil && err != sql.ErrNoRows {
		ch.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		ch.db.Exec("INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
			target.ID,
			target.Username,
			target.Discriminator,
			0,
			currentTime,
			false,
			rep,
			currentTime,
		)
	}

	lcxp := &models.Localxp{}
	gbxp := &models.Globalxp{}

	newXp := 0

	err = ch.db.Get(lcxp, "SELECT * FROM localxp WHERE userid = $1 AND guildid = $2;", target.ID, ctx.Guild.ID)
	if err != nil && err != sql.ErrNoRows {
		ch.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		ch.db.Exec("INSERT INTO localxp(guildid, userid, xp) VALUES($1, $2, $3);", ctx.Guild.ID, target.ID, newXp)
	}

	err = ch.db.Get(gbxp, "SELECT * FROM globalxp WHERE userid = $1;", target.ID)
	if err != nil && err != sql.ErrNoRows {
		ch.logger.Error("error", zap.Error(err))
		return
	} else if err == sql.ErrNoRows {
		ch.db.Exec("INSERT INTO globalxp(userid, xp) VALUES($1, $2);", target.ID, newXp)
	}
}

func (ch *CommandHandler) HighestRole(g *discordgo.Guild, userID string) int {

	user, err := ch.client.State.Member(g.ID, userID)
	if err != nil {
		return -1
	}

	if user.User.ID == g.OwnerID {
		return math.MaxInt64
	}

	topRole := 0

	for _, val := range user.Roles {
		for _, role := range g.Roles {
			if val == role.ID {
				if role.Position > topRole {
					topRole = role.Position
				}
			}
		}
	}

	return topRole
}

func (ch *CommandHandler) UserColor(g *discordgo.Guild, userID string) int {

	member, err := ch.client.State.Member(g.ID, userID)
	if err != nil {
		return 0
	}

	roles := discordgo.Roles(g.Roles)
	sort.Sort(roles)

	for _, role := range roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID {
				if role.Color != 0 {
					return role.Color
				}
			}
		}
	}

	return 0
}

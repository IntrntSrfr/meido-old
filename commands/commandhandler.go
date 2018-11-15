package commands

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"meido-test/models"
	"meido-test/service"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	dColorRed   = 13107200
	dColorGreen = 51200
	dColorWhite = 16777215
)

type Command struct {
	Name          string
	Triggers      []string
	Description   string
	Usage         string
	RequiredPerms int
	RequiresOwner bool
	Execute       func(args []string, context *service.Context)
}

type Commandmap map[string]Command

var (
	startTime     = time.Now()
	comms         = Commandmap{}
	db            *sql.DB
	dmLogChannels []string
	ownerIds      []string
)

func Initialize(OwnerIds *[]string, DmLogChannels *[]string, DB *sql.DB) {
	comms.RegisterCommand(About)
	comms.RegisterCommand(Avatar)
	comms.RegisterCommand(Ban)
	comms.RegisterCommand(ClearAFK)
	comms.RegisterCommand(CoolNameBro)
	//comms.RegisterCommand(Filter)
	comms.RegisterCommand(Help)
	comms.RegisterCommand(Inrole)
	comms.RegisterCommand(Kick)
	comms.RegisterCommand(Lockdown)
	comms.RegisterCommand(MyRole)
	comms.RegisterCommand(Ping)
	comms.RegisterCommand(Profile)
	comms.RegisterCommand(Rep)
	comms.RegisterCommand(Repleaderboard)
	//comms.RegisterCommand(Role)
	//comms.RegisterCommand(Server)
	comms.RegisterCommand(SetUserRole)
	comms.RegisterCommand(Umr)
	comms.RegisterCommand(Unlock)
	//comms.RegisterCommand(User)
	comms.RegisterCommand(WithNick)
	comms.RegisterCommand(WithTag)
	comms.RegisterCommand(Xpleaderboard)

	comms.RegisterCommand(Test)
	comms.RegisterCommand(Dm)
	comms.RegisterCommand(Msg)

	db = DB
	dmLogChannels = *DmLogChannels
	ownerIds = *OwnerIds
}

func GetCommandMap() Commandmap {

	return comms

}

func (cmap *Commandmap) RegisterCommand(cmd Command) {

	(*cmap)[cmd.Name] = cmd
}

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	context := service.NewContext(s, m.Message)

	if m.Author.Bot {
		return
	}

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if ch.Type == discordgo.ChannelTypeDM {
		if strings.ToLower(m.Content) == "enroll me" {
			cfc, err := s.Guild("320896491596283906")
			if err != nil {
				return
			}

			var enrolledRole *discordgo.Role

			groles, err := s.GuildRoles(cfc.ID)
			if err != nil {
				return
			}

			for i := range groles {
				role := groles[i]
				if role.ID == "404333507918430212" {
					enrolledRole = role
				}
			}

			if enrolledRole == nil {
				return
			}

			for i := range cfc.Members {
				member := cfc.Members[i]

				if member.User.ID == m.Author.ID {
					err := s.GuildMemberRoleAdd(cfc.ID, m.Author.ID, enrolledRole.ID)
					if err != nil {
						return
					}
				}
			}
		} else {

			var dmembed discordgo.MessageEmbed

			if len(m.Attachments) > 0 {
				dmembed = discordgo.MessageEmbed{
					Color:       dColorWhite,
					Title:       fmt.Sprintf("Message from %v", m.Author.String()),
					Description: m.Content,
					Image:       &discordgo.MessageEmbedImage{URL: m.Attachments[0].URL},
					Footer:      &discordgo.MessageEmbedFooter{Text: m.Author.ID},
					Timestamp:   string(m.Timestamp),
				}
			} else {
				dmembed = discordgo.MessageEmbed{
					Color:       dColorWhite,
					Title:       fmt.Sprintf("Message from %v", m.Author.String()),
					Description: m.Content,
					Footer:      &discordgo.MessageEmbedFooter{Text: m.Author.ID},
					Timestamp:   string(m.Timestamp),
				}
			}

			for i := range dmLogChannels {
				dmch := dmLogChannels[i]

				_, err := s.ChannelMessageSendEmbed(dmch, &dmembed)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
		return
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return
	}

	g, err := s.Guild(ch.GuildID)
	if err != nil {
		return
	}

	perms, err := s.UserChannelPermissions(m.Author.ID, ch.ID)
	if err != nil {
		return
	}
	botPerms, err := s.UserChannelPermissions(s.State.User.ID, ch.ID)
	if err != nil {
		return
	}

	args := strings.Split(m.Content, " ")

	isIllegal := checkFilter(&context, &perms, m)
	if isIllegal {
		return
	}

	doXp(&context)

	triggerCommand := ""
	for _, val := range comms {
		for _, com := range val.Triggers {
			if strings.ToLower(args[0]) == strings.ToLower(com) {
				triggerCommand = val.Name
			}
		}
	}

	if triggerCommand != "" {

		if cmd, ok := comms[triggerCommand]; ok {

			isOwner := false

			if cmd.RequiresOwner == true {
				for _, val := range ownerIds {
					if m.Author.ID == val {
						isOwner = true
					}
				}
				if !isOwner {
					context.Send("Owner only.")
					return
				}
			}

			if !isOwner {
				if perms&cmd.RequiredPerms == 0 {
					return
				}
			}

			if botPerms&cmd.RequiredPerms == 0 {
				context.Send(fmt.Sprintf("I am missing permissions: %v", cmd.RequiredPerms))
				return
			}

			cmd.Execute(args, &context)
			db.Exec("INSERT INTO usedcommands VALUES($1, $2, $3, $4)")
			fmt.Println(fmt.Sprintf("Command executed\nCommand: %v\nUser: %v [%v]\nSource: %v [%v] - #%v [%v]\n", args, m.Author.String(), m.Author.ID, g.Name, g.ID, ch.Name, ch.ID))
		}
	}
}

func checkFilter(ctx *service.Context, perms *int, msg *discordgo.MessageCreate) bool {

	isIllegal := false

	if *perms&discordgo.PermissionManageMessages == 0 {

		rows, _ := db.Query("SELECT phrase FROM filters WHERE guildid = $1", ctx.Guild.ID)

		for rows.Next() {
			filter := models.Filter{}
			err := rows.Scan(&filter.Filter)
			if err != nil {
				continue
			}

			if strings.Contains(msg.Content, filter.Filter) {
				isIllegal = true
				break
			}
		}

		if isIllegal {
			ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
			ctx.Send(fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!", msg.Author.Mention()))
		}
	}

	return isIllegal
}

func doXp(ctx *service.Context) {

	dbu := models.DiscordUser{}

	currentTime := time.Now()
	xpTime := time.Now()
	isIgnored := false

	row := db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", ctx.User.ID)
	err := row.Scan(
		&dbu.Uid,
		&dbu.Userid,
		&dbu.Username,
		&dbu.Discriminator,
		&dbu.Xp,
		&dbu.Nextxpgaintime,
		&dbu.Xpexcluded,
		&dbu.Reputation,
		&dbu.Cangivereptime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			db.Exec("INSERT INTO discordusers(userid, username, discriminator, xp, nextxpgaintime, xpexcluded, reputation, cangivereptime) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
				ctx.User.ID,
				ctx.User.Username,
				ctx.User.Discriminator,
				0,
				currentTime,
				false,
				0,
				currentTime,
			)
		}
	} else {
		isIgnored = dbu.Xpexcluded
		xpTime = dbu.Nextxpgaintime
	}

	diff := xpTime.Sub(currentTime)

	if diff <= 0 {

		igch := models.Xpignoredchannel{}
		//igu := models.Xpignoreduser{}
		lcxp := models.Localxp{}
		gbxp := models.Globalxp{}

		newXp := Random(15, 26)

		rows, err := db.Query("SELECT * FROM xpignoredchannels WHERE guildid = $1;", ctx.Guild.ID)
		if err != nil {
			return
		}

		for rows.Next() {
			rows.Scan(
				&igch.Uid,
				&igch.Channelid,
			)
			if igch.Channelid == ctx.Channel.ID {
				newXp = 0
			}
		}
		/*
			row = db.QueryRow("SELECT * FROM xpignoreduser WHERE userid = $1;", ctx.User.ID)
			err = row.Scan(
				&igu.Uid,
				&igu.Userid,
			)

			if err != nil {
				if igu.Userid == ctx.User.ID {
					newXp = 0
				}
			} */

		if isIgnored {
			newXp = 0
		}

		row = db.QueryRow("SELECT * FROM localxp WHERE userid = $1 AND guildid = $2;", ctx.User.ID, ctx.Guild.ID)
		err = row.Scan(
			&lcxp.Uid,
			&lcxp.Guildid,
			&lcxp.Userid,
			&lcxp.Xp,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				db.Exec("INSERT INTO localxp(guildid, userid, xp) VALUES($1, $2, $3);", ctx.Guild.ID, ctx.User.ID, newXp)
			}
		} else {
			if newXp != 0 {
				db.Exec("UPDATE localxp SET xp = $1 WHERE guildid = $2 AND userid = $3;", lcxp.Xp+newXp, ctx.Guild.ID, ctx.User.ID)
			}
		}
		row = db.QueryRow("SELECT * FROM globalxp WHERE userid = $1;", ctx.User.ID)
		err = row.Scan(
			&gbxp.Uid,
			&gbxp.Userid,
			&gbxp.Xp,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				db.Exec("INSERT INTO globalxp(userid, xp) VALUES($1, $2);", ctx.User.ID, newXp)
			}
		} else {
			if newXp != 0 {
				db.Exec("UPDATE globalxp SET xp = $1 WHERE userid = $3;", lcxp.Xp+newXp, ctx.User.ID)
			}
		}

		db.Exec("UPDATE discordusers SET nextxpgaintime = $1 WHERE userid = $2;", currentTime.Add(time.Minute*time.Duration(2)), ctx.User.ID)
	}
}

func HighestRole(g *discordgo.Guild, u *discordgo.Member) int {

	if u.User.ID == g.OwnerID {
		return math.MaxInt64
	}

	topRole := 0

	for _, val := range u.Roles {
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

func HighestColor(g *discordgo.Guild, u *discordgo.Member) int {

	topRole := 0
	topColor := 0
	groles := discordgo.Roles(g.Roles)
	userroles := []*discordgo.Role{}

	for _, grole := range groles {
		for _, urole := range u.Roles {
			if urole == grole.ID {
				userroles = append(userroles, grole)
			}
		}
	}

	groles = discordgo.Roles(userroles)
	sort.Sort(groles)
	sort.Sort(sort.Reverse(groles))

	for _, grole := range groles {

		if grole.Position > topRole {

			topRole = grole.Position

			if grole.Color != 0 {
				topColor = grole.Color
			}
		}
	}

	return topColor
}

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

package commands

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/intrntsrfr/meido/bot/models"
	"github.com/intrntsrfr/meido/bot/service"
	"github.com/intrntsrfr/owo"

	"github.com/bwmarrin/discordgo"
)

const (
	dColorRed    = 13107200
	dColorOrange = 15761746
	dColorLBlue  = 6410733
	dColorGreen  = 51200
	dColorWhite  = 16777215
)

type CommandType string

const (
	Filter     CommandType = "Filter"
	Strikes    CommandType = "Strikes"
	Moderation CommandType = "Moderation"
	Fun        CommandType = "Fun"
	Utility    CommandType = "Utility"
	Profile    CommandType = "Profile"
	Owner      CommandType = "Owner"
)

type Command struct {
	Name          string
	Triggers      []string
	Description   string
	Usage         string
	Category      CommandType
	RequiredPerms int
	RequiresOwner bool
	Execute       func(args []string, context *service.Context)
}

type Commandmap map[string]Command

var permMap = map[int]string{
	1:          "create instant invite",
	2:          "kick members",
	4:          "ban members",
	8:          "administrator",
	16:         "manage channels",
	32:         "manage server",
	64:         "add reactions",
	128:        "view audit log",
	256:        "priority speaker",
	1024:       "view channel",
	2048:       "send messages",
	4096:       "send tts messages",
	8192:       "manage messages",
	16384:      "embed links",
	32768:      "attach files",
	65536:      "read message history",
	131072:     "mention everyone",
	262144:     "use external emojis",
	1048576:    "connect",
	2097152:    "speak",
	4194304:    "mute members",
	8388608:    "deafen members",
	16777216:   "move members",
	33554432:   "use VAD",
	67108864:   "change nickname",
	134217728:  "manage nicknames",
	268435456:  "manage roles",
	536870912:  "manage webhooks",
	1073741824: "manage emojis",
}

var verificationMap = map[int]string{
	0: "Unrestricted.",
	1: "Email verification.",
	2: "Email verification and account must be at least 5 minutes old.",
	3: "Email verification, account must be at least 5 minutes old and user must have been on server for 10 minutes.",
	4: "Verified phone to Discord account.",
}

var (
	client        *discordgo.Session
	botStartTime  = time.Now()
	comms         = Commandmap{}
	db            *sql.DB
	dmLogChannels []string
	ownerIds      []string
	OWOApi        *owo.OWOClient
	logger        *zap.Logger
)

func Initialize(s *discordgo.Session, OwnerIds *[]string, DmLogChannels *[]string, DB *sql.DB, owo *owo.OWOClient, Logger *zap.Logger) {

	// Filter
	comms.RegisterCommand(FilterWord)
	comms.RegisterCommand(FilterWordList)
	comms.RegisterCommand(FilterInfo)
	comms.RegisterCommand(FilterIgnoreChannel)
	comms.RegisterCommand(ClearFilter)

	// Strikes
	comms.RegisterCommand(UseStrikes)
	comms.RegisterCommand(SetMaxStrikes)
	comms.RegisterCommand(ClearStrikes)
	comms.RegisterCommand(Warn)
	comms.RegisterCommand(StrikeLog)
	//comms.RegisterCommand(StrikeLogAll)
	comms.RegisterCommand(RemoveStrike)

	// Moderation
	comms.RegisterCommand(Ban)
	comms.RegisterCommand(Hackban)
	comms.RegisterCommand(Unban)
	//comms.RegisterCommand(ClearAFK)
	comms.RegisterCommand(CoolNameBro)
	comms.RegisterCommand(NiceNameBro)
	comms.RegisterCommand(Kick)
	comms.RegisterCommand(Lockdown)
	comms.RegisterCommand(Unlock)
	comms.RegisterCommand(SetUserRole)

	// Utility
	comms.RegisterCommand(About)
	comms.RegisterCommand(Avatar)
	comms.RegisterCommand(Ping)
	comms.RegisterCommand(Help)
	comms.RegisterCommand(Inrole)
	comms.RegisterCommand(WithNick)
	comms.RegisterCommand(WithTag)
	//comms.RegisterCommand(Role)
	comms.RegisterCommand(Server)
	//comms.RegisterCommand(User)
	//comms.RegisterCommand(ListRoles)
	comms.RegisterCommand(ListUserRoles)
	comms.RegisterCommand(Invite)
	comms.RegisterCommand(Feedback)
	comms.RegisterCommand(MyRole)

	// Fun
	comms.RegisterCommand(Img)

	// Profile
	comms.RegisterCommand(ShowProfile)
	comms.RegisterCommand(Rep)
	comms.RegisterCommand(Repleaderboard)
	comms.RegisterCommand(XpLeaderboard)
	comms.RegisterCommand(GlobalXpLeaderboard)
	comms.RegisterCommand(XpIgnoreChannel)

	// Owner
	//comms.RegisterCommand(Test)
	comms.RegisterCommand(Dm)
	comms.RegisterCommand(Msg)

	client = s
	db = DB
	dmLogChannels = *DmLogChannels
	ownerIds = *OwnerIds
	OWOApi = owo
	logger = Logger
}

func GetCommandMap() Commandmap {
	return comms
}

func (cmap *Commandmap) RegisterCommand(cmd Command) {
	if cmd, ok := comms[cmd.Name]; ok {
		fmt.Println("Conflicting Commands.", cmd, comms[cmd.Name])
		os.Exit(0)
	}

	(*cmap)[cmd.Name] = cmd
}

//var id = 0

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	startTime := time.Now()

	context, err := service.NewContext(s, m.Message, startTime)
	if err != nil {
		return
	}

	//fmt.Println(fmt.Sprintf("[%v] - context in %v", id, time.Now().Sub(startTime)))

	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if ch.Type == discordgo.ChannelTypeDM {
		dmembed := discordgo.MessageEmbed{
			Color:       dColorWhite,
			Title:       fmt.Sprintf("Message from %v", m.Author.String()),
			Description: m.Content,
			Footer:      &discordgo.MessageEmbedFooter{Text: m.Author.ID},
			Timestamp:   string(m.Timestamp),
		}

		if len(m.Attachments) > 0 {
			dmembed.Image = &discordgo.MessageEmbedImage{URL: m.Attachments[0].URL}
		}

		for i := range dmLogChannels {
			dmch := dmLogChannels[i]

			_, err := s.ChannelMessageSendEmbed(dmch, &dmembed)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		return
	}

	if ch.Type != discordgo.ChannelTypeGuildText {
		return
	}

	perms, err := s.State.UserChannelPermissions(m.Author.ID, ch.ID)
	if err != nil {
		return
	}
	botPerms, err := s.State.UserChannelPermissions(s.State.User.ID, ch.ID)
	if err != nil {
		return
	}

	args := strings.Split(m.Content, " ")

	isIllegal := checkFilter(&context, &perms, m)
	if isIllegal {
		return
	}
	//fmt.Println(fmt.Sprintf("[%v] - filter in %v", id, time.Now().Sub(startTime)))

	doXp(&context)
	//fmt.Println(fmt.Sprintf("[%v] - xp in %v", id, time.Now().Sub(startTime)))

	triggerCommand := ""
	for _, val := range comms {
		for _, com := range val.Triggers {
			if strings.ToLower(args[0]) == strings.ToLower(com) {
				triggerCommand = val.Name
			}
		}
	}
	//fmt.Println(fmt.Sprintf("[%v] - checked command in %v", id, time.Now().Sub(startTime)))

	if triggerCommand != "" {

		if cmd, ok := comms[triggerCommand]; ok {

			isOwner := false

			for _, val := range ownerIds {
				if m.Author.ID == val {
					isOwner = true
				}
			}
			if cmd.RequiresOwner {
				if !isOwner {
					context.Send("Owner only.")
					return
				}
			}
			/*
				if !cmd.RequiresOwner {
				}
			*/
			if perms&cmd.RequiredPerms == 0 && perms&discordgo.PermissionAdministrator == 0 {
				//fmt.Println(perms, cmd.RequiredPerms, permMap[cmd.RequiredPerms], perms&cmd.RequiredPerms)
				return
			}

			if botPerms&cmd.RequiredPerms == 0 && perms&discordgo.PermissionAdministrator == 0 {
				context.Send(fmt.Sprintf("I am missing permissions: %v", permMap[cmd.RequiredPerms]))
				return
			}

			go cmd.Execute(args, &context)
			//fmt.Println(fmt.Sprintf("[%v] - executed command in %v\n", id, time.Now().Sub(startTime)))
			db.Exec("INSERT INTO commandlog(command, args, userid, guildid, channelid, messageid, tstamp) VALUES($1, $2, $3, $4, $5, $6, $7)", cmd.Name, strings.Join(args, " "), m.Author.ID, context.Guild.ID, ch.ID, m.ID, time.Now())
			fmt.Println(fmt.Sprintf("\nCommand executed\nCommand: %v\nUser: %v [%v]\nSource: %v [%v] - #%v [%v]", args, m.Author.String(), m.Author.ID, context.Guild.Name, context.Guild.ID, ch.Name, ch.ID))
		}
	}
	//id++
}

func checkFilter(ctx *service.Context, perms *int, msg *discordgo.MessageCreate) bool {

	isIllegal := false
	trigger := ""

	if *perms&discordgo.PermissionManageMessages == 0 {

		var count int

		row := db.QueryRow("SELECT COUNT(*) FROM filterignorechannels WHERE channelid = $1;", ctx.Channel.ID)
		err := row.Scan(&count)
		if err != nil {
			return false
		}

		if count > 0 {
			return true
		}

		rows, err := db.Query("SELECT phrase FROM filters WHERE guildid = $1", ctx.Guild.ID)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer rows.Close()

		for rows.Next() {
			filter := models.Filter{}
			err := rows.Scan(&filter.Filter)
			if err != nil {
				continue
			}

			if strings.Contains(strings.ToLower(msg.Content), strings.ToLower(filter.Filter)) {
				trigger = filter.Filter
				isIllegal = true
				break
			}
		}

		if isIllegal {
			ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)

			row := db.QueryRow("SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid = $1;", ctx.Guild.ID)

			dbg := models.DiscordGuild{}

			err := row.Scan(&dbg.UseStrikes, &dbg.MaxStrikes)
			if err != nil {
				return false
			}

			if dbg.UseStrikes {

				reason := fmt.Sprintf("Triggering filter: %v", trigger)

				strikeCount := 0

				row := db.QueryRow("SELECT COUNT(*) FROM strikes WHERE guildid = $1 AND userid = $2;", ctx.Guild.ID, ctx.User.ID)
				err := row.Scan(&strikeCount)
				if err != nil {
					return false
				}

				if strikeCount+1 >= dbg.MaxStrikes {
					//ban
					userch, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
					ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for acquiring %v strikes.\nLast warning was: %v", ctx.Guild.Name, dbg.MaxStrikes, reason))
					err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, ctx.User.ID, fmt.Sprintf("Acquired %v strikes.", dbg.MaxStrikes), 0)
					if err != nil {
						return false
					}

					ctx.Send(fmt.Sprintf("%v has been banned after acquiring too many strikes. Miss them.", ctx.User.Mention()))
					_, err := db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", ctx.User.ID, ctx.Guild.ID)
					if err != nil {
						fmt.Println(err)
					}
				} else {
					//insert warn
					userch, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
					ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been warned in %v.\nWarned for: %v", ctx.Guild.Name, reason))
					ctx.Send(fmt.Sprintf("%v has been warned\nThey are currently at strike %v/%v", ctx.User.Mention(), strikeCount+1, dbg.MaxStrikes))
					db.Exec("INSERT INTO strikes(guildid, userid, reason, executorid, tstamp) VALUES ($1, $2, $3, $4, $5);", ctx.Guild.ID, ctx.User.ID, reason, ctx.Session.State.User.ID, time.Now())
				}

			} else {
				ctx.Send(fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!", msg.Author.Mention()))
			}
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
		//igu := models.Xpignoreduser{}
		lcxp := models.Localxp{}
		gbxp := models.Globalxp{}

		newXp := Random(15, 26)

		row := db.QueryRow("SELECT COUNT(*) FROM xpignoredchannels WHERE channelid = $1;", ctx.Channel.ID)

		count := 0
		err := row.Scan(
			&count,
		)
		if count > 0 {
			newXp = 0
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
			if !isIgnored {
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
			if !isIgnored {
				db.Exec("UPDATE globalxp SET xp = $1 WHERE userid = $2;", gbxp.Xp+newXp, ctx.User.ID)
			}
		}

		db.Exec("UPDATE discordusers SET nextxpgaintime = $1 WHERE userid = $2;", currentTime.Add(time.Minute*time.Duration(2)), ctx.User.ID)
	}
}

func SetupProfile(target *discordgo.User, ctx *service.Context, rep int) {

	dbu := models.DiscordUser{}

	currentTime := time.Now()

	row := db.QueryRow("SELECT * FROM discordusers WHERE userid = $1", target.ID)
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
	}

	lcxp := models.Localxp{}
	gbxp := models.Globalxp{}

	newXp := 0

	row = db.QueryRow("SELECT * FROM localxp WHERE userid = $1 AND guildid = $2;", target.ID, ctx.Guild.ID)
	err = row.Scan(
		&lcxp.Uid,
		&lcxp.Guildid,
		&lcxp.Userid,
		&lcxp.Xp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			db.Exec("INSERT INTO localxp(guildid, userid, xp) VALUES($1, $2, $3);", ctx.Guild.ID, target.ID, newXp)
		}
	}
	row = db.QueryRow("SELECT * FROM globalxp WHERE userid = $1;", target.ID)
	err = row.Scan(
		&gbxp.Uid,
		&gbxp.Userid,
		&gbxp.Xp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			db.Exec("INSERT INTO globalxp(userid, xp) VALUES($1, $2);", target.ID, newXp)
		}
	}
}

func HighestRole(g *discordgo.Guild, userID string) int {

	user, err := client.State.Member(g.ID, userID)
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

func UserColor(g *discordgo.Guild, userID string) int {

	member, err := client.State.Member(g.ID, userID)
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

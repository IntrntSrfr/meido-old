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
)

func Initialize(s *discordgo.Session, OwnerIds *[]string, DmLogChannels *[]string, DB *sql.DB) {

	comms.RegisterCommand(FilterWord)
	comms.RegisterCommand(FilterWordList)
	comms.RegisterCommand(FilterInfo)
	comms.RegisterCommand(FilterIgnoreChannel)
	comms.RegisterCommand(ClearFilter)
	comms.RegisterCommand(UseStrikes)
	comms.RegisterCommand(SetMaxStrikes)
	comms.RegisterCommand(ClearStrikes)

	comms.RegisterCommand(Ban)
	comms.RegisterCommand(Unban)
	comms.RegisterCommand(ClearAFK)
	comms.RegisterCommand(CoolNameBro)
	comms.RegisterCommand(Kick)
	comms.RegisterCommand(Lockdown)
	comms.RegisterCommand(Unlock)
	comms.RegisterCommand(SetUserRole)

	comms.RegisterCommand(About)
	comms.RegisterCommand(Avatar)
	comms.RegisterCommand(MyRole)
	comms.RegisterCommand(Ping)
	comms.RegisterCommand(Umr)
	comms.RegisterCommand(Help)
	comms.RegisterCommand(Inrole)
	comms.RegisterCommand(WithNick)
	comms.RegisterCommand(WithTag)
	//comms.RegisterCommand(Role)
	comms.RegisterCommand(Server)
	//comms.RegisterCommand(User)

	comms.RegisterCommand(Profile)
	comms.RegisterCommand(Rep)
	comms.RegisterCommand(Repleaderboard)
	comms.RegisterCommand(XpLeaderboard)
	comms.RegisterCommand(GlobalXpLeaderboard)
	comms.RegisterCommand(XpIgnoreChannel)

	comms.RegisterCommand(Test)
	comms.RegisterCommand(Dm)
	comms.RegisterCommand(Msg)
	comms.RegisterCommand(Warn)

	client = s
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

//var id = 0

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.Bot {
		return
	}

	startTime := time.Now()

	context := service.NewContext(s, m.Message, startTime)

	//fmt.Println(fmt.Sprintf("[%v] - context in %v", id, time.Now().Sub(startTime)))

	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}

	if ch.Type == discordgo.ChannelTypeDM {
		if strings.ToLower(m.Content) == "enroll me" {
			cfc, err := s.State.Guild("320896491596283906")
			if err != nil {
				return
			}

			var enrolledRole *discordgo.Role

			groles := cfc.Roles

			for i := range groles {
				role := groles[i]
				if role.ID == "404333507918430212" {
					enrolledRole = role
				}
			}

			if enrolledRole == nil {
				return
			}

			mem, err := s.State.Member(cfc.ID, m.Author.ID)
			if err != nil {
				return
			}

			if mem.User.ID == m.Author.ID {
				err := s.GuildMemberRoleAdd(cfc.ID, m.Author.ID, enrolledRole.ID)
				if err != nil {
					return
				}
				s.ChannelMessageSend(m.ChannelID, "You're now enrolled, congratulations for passing the IQ test.")
			}
		} else {
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

			if botPerms&cmd.RequiredPerms == 0 {
				context.Send(fmt.Sprintf("I am missing permissions: %v", permMap[cmd.RequiredPerms]))
				return
			}

			go cmd.Execute(args, &context)
			//fmt.Println(fmt.Sprintf("[%v] - executed command in %v\n", id, time.Now().Sub(startTime)))
			db.Exec("INSERT INTO commandlog(command, args, userid, guildid, channelid, messageid, tstamp) VALUES($1, $2, $3, $4, $5, $6, $7)", cmd.Name, strings.Join(args, " "), m.Author.ID, g.ID, ch.ID, m.ID, time.Now())
			fmt.Println(fmt.Sprintf("\nCommand executed\nCommand: %v\nUser: %v [%v]\nSource: %v [%v] - #%v [%v]", args, m.Author.String(), m.Author.ID, g.Name, g.ID, ch.Name, ch.ID))
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

		rows, _ := db.Query("SELECT phrase FROM filters WHERE guildid = $1", ctx.Guild.ID)

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
			row := db.QueryRow("SELECT usestrikes, maxstrikes FROM discordguilds WHERE guildid = $1;", ctx.Guild.ID)

			dbg := models.DiscordGuild{}

			err := row.Scan(&dbg.UseStrikes, &dbg.MaxStrikes)
			if err != nil {
				return false
			}

			if dbg.UseStrikes {
				dbs := models.Strikes{}

				row := db.QueryRow("SELECT * FROM strikes WHERE guildid = $1 AND userid = $2;", ctx.Guild.ID, ctx.User.ID)
				err := row.Scan(&dbs.Uid, &dbs.Guildid, &dbs.Userid, &dbs.Strikes)
				if err != nil {
					if err == sql.ErrNoRows {
						if dbg.MaxStrikes < 2 {
							ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
							userch, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
							ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for triggering the filter.\n- %v", ctx.Guild.Name, trigger))
							err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, ctx.User.ID, fmt.Sprintf("Triggering filter: %v", trigger), 0)
							if err != nil {
								ctx.Send(err.Error())
								return true
							}

							embed := &discordgo.MessageEmbed{
								Title:       "User banned",
								Description: "Filter triggered",
								Fields: []*discordgo.MessageEmbedField{
									{
										Name:   "Username",
										Value:  fmt.Sprintf("%v", ctx.User.Mention()),
										Inline: true,
									},
									{
										Name:   "ID",
										Value:  fmt.Sprintf("%v", ctx.User.ID),
										Inline: true,
									},
								},
								Color: dColorRed,
							}

							ctx.SendEmbed(embed)

						} else {
							ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
							ctx.Send(fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!\nYou are currently at strike %v/%v", msg.Author.Mention(), dbs.Strikes+1, dbg.MaxStrikes))
							db.Exec("INSERT INTO strikes(guildid, userid, strikes) VALUES ($1, $2, $3);", ctx.Guild.ID, ctx.User.ID, 1)
						}
					}
				} else {
					if dbs.Strikes+1 >= dbg.MaxStrikes {
						ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
						userch, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
						ctx.Session.ChannelMessageSend(userch.ID, fmt.Sprintf("You have been banned from %v for triggering the filter.\n- %v", ctx.Guild.Name, trigger))
						err = ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, ctx.User.ID, fmt.Sprintf("Triggering filter: %v", trigger), 0)
						if err != nil {
							ctx.Send(err.Error())
							return true
						}

						embed := &discordgo.MessageEmbed{
							Title:       "User banned",
							Description: "Filter triggered",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Username",
									Value:  fmt.Sprintf("%v", ctx.User.Mention()),
									Inline: true,
								},
								{
									Name:   "ID",
									Value:  fmt.Sprintf("%v", ctx.User.ID),
									Inline: true,
								},
							},
							Color: dColorRed,
						}

						ctx.SendEmbed(embed)

						_, err := db.Exec("DELETE FROM strikes WHERE userid = $1 AND guildid = $2;", ctx.User.ID, ctx.Guild.ID)
						if err != nil {
							fmt.Println(err)
						}

					} else {
						ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
						ctx.Send(fmt.Sprintf("%v, you are not allowed to use a banned word/phrase!\nYou are currently at strike %v/%v", msg.Author.Mention(), dbs.Strikes+1, dbg.MaxStrikes))
						db.Exec("UPDATE strikes SET strikes = $1 WHERE userid = $2 AND guildid = $3;", dbs.Strikes+1, ctx.User.ID, ctx.Guild.ID)
					}
				}

			} else {
				ctx.Session.ChannelMessageDelete(ctx.Channel.ID, msg.ID)
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

	diff := xpTime.Sub(currentTime.Add(time.Hour * 1))

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

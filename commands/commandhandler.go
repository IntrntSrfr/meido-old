package commands

import (
	"database/sql"
	"fmt"
	"meido-test/models"
	"meido-test/service"
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
	comms.RegisterCommand(ClearAFK)
	comms.RegisterCommand(CoolNameBro)
	comms.RegisterCommand(Help)
	comms.RegisterCommand(Inrole)
	comms.RegisterCommand(MyRole)
	comms.RegisterCommand(Ping)
	comms.RegisterCommand(Server)
	comms.RegisterCommand(SetUserRole)
	comms.RegisterCommand(Test)
	comms.RegisterCommand(WithNick)
	comms.RegisterCommand(WithTag)

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

	doLocalXp()
	doGlobalXp()

	triggerCommand := ""
	for _, val := range comms {
		for _, com := range val.Triggers {
			if args[0] == com {
				triggerCommand = val.Name
			}
		}
	}

	if triggerCommand != "" {

		if cmd, ok := comms[triggerCommand]; ok {

			if cmd.RequiresOwner == true {
				isOwner := false
				for _, val := range ownerIds {
					if m.Author.ID == val {
						isOwner = true
					}
				}
				if !isOwner {
					return
				}
			}

			if perms&cmd.RequiredPerms == 0 {
				return
			}
			if botPerms&cmd.RequiredPerms == 0 {
				context.Send(fmt.Sprintf("I am missing permissions: %v", cmd.RequiredPerms))
				return
			}

			cmd.Execute(args, &context)
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

func doLocalXp() {

}

func doGlobalXp() {

}

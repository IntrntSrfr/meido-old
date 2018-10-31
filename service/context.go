package service

import (
	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Session *discordgo.Session
	Message *discordgo.Message
	//Db      *sql.DB
}

var ctx = Context{}

func NewContext(s *discordgo.Session, m *discordgo.Message) Context {
	return Context{Session: s, Message: m}
}

func (c *Context) Send(input string) {
	c.Session.ChannelMessageSend(ctx.Message.ChannelID, input)
}

func SendEmbed(embed discordgo.MessageEmbed) {
	ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, &embed)
}

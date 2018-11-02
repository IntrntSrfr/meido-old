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
	_, err := c.Session.ChannelMessageSend(c.Message.ChannelID, input)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
	}
}

func (c *Context) SendEmbed(embed discordgo.MessageEmbed) {
	_, err := c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, &embed)
	if err != nil {
		c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
	}
}

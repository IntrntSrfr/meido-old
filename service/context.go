package service

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Session *discordgo.Session
	Message *discordgo.Message
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Db      *sql.DB
}

var ctx = Context{}

func NewContext(s *discordgo.Session, m *discordgo.Message) Context {
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		ch = nil
	}

	g, err := s.Guild(ch.GuildID)
	if err != nil {
		g = nil
	}

	u, err := s.User(m.Author.ID)
	if err != nil {
		u = nil
	}

	return Context{
		Session: s,
		Message: m,
		User:    u,
		Channel: ch,
		Guild:   g,
	}
}

func (c *Context) Send(input string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Message.ChannelID, input)
}

func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)
}

package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Context struct {
	Session   *discordgo.Session
	Message   *discordgo.Message
	Guild     *discordgo.Guild
	Channel   *discordgo.Channel
	User      *discordgo.User
	Db        *sql.DB
	StartTime time.Time
}

var ctx = Context{}

func NewContext(s *discordgo.Session, m *discordgo.Message, t time.Time) Context {
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
		Session:   s,
		Message:   m,
		User:      u,
		Channel:   ch,
		Guild:     g,
		StartTime: t,
	}
}

func (c *Context) Send(a ...interface{}) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("%v", a...))
}

func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)
}

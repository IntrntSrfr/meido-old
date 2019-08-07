package service

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	Session   *discordgo.Session
	Message   *discordgo.Message
	Guild     *discordgo.Guild
	Channel   *discordgo.Channel
	User      *discordgo.User
	Db        *sqlx.DB
	StartTime time.Time
}

var ctx = Context{}

func NewContext(s *discordgo.Session, m *discordgo.Message, t time.Time) (Context, error) {
	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return Context{}, err
	}

	g := &discordgo.Guild{}

	if ch.Type == discordgo.ChannelTypeGuildText {
		g, err = s.State.Guild(ch.GuildID)
		if err != nil {
			return Context{}, err
		}
	}
	/*
		u, err := s.User(m.Author.ID)
		if err != nil {
			u = nil
		} */

	return Context{
		Session:   s,
		Message:   m,
		User:      m.Author,
		Channel:   ch,
		Guild:     g,
		StartTime: t,
	}, nil
}

func (c *Context) Send(a ...interface{}) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprint(a...))
}

func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)
}

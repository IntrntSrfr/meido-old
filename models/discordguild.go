package models

type DiscordGuild struct {
	Uid                  int
	Guildid              string
	UseStrikes           bool
	MaxStrikes           int
	IgnoredChannels      string
	WelcomeChannel       string
	MsgDeleteLogChannel  string
	MsgEditLogChannel    string
	UserJoinedLogChannel string
	UserLeftLogChannel   string
	BanLogChannel        string
	VoiceLogChannel      string
}

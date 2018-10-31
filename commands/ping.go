package commands

import (
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Ping = Command{
	Name:          "jeff",
	Description:   "jeffe",
	Aliases:       []string{"jeffer", "jeffette"},
	Usage:         "m?jeff",
	RequiredPerms: discordgo.PermissionManageMessages,
	Function: func(context *service.Context) {

		comms := GetCommandMap()

		list := "```\n"
		for _, val := range *comms {
			list += val.Name + "\n"
		}
		list += "```"
		context.Session.ChannelMessageSend(context.Message.ChannelID, list)
	},
}

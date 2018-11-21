package commands

import (
	"fmt"
	"meido-test/service"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Help = Command{
	Name:          "help",
	Description:   "Shows info about commands.",
	Triggers:      []string{"m?help", "m?h"},
	Usage:         "m?help <optional command name>",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		comms := GetCommandMap()

		if len(args) < 2 {

			list := "```\n"
			for _, val := range comms {
				list += val.Name + "\n"
			}
			list += "```"

			ctx.Send(list)
		} else {

			comm := args[1]

			if cmd, ok := comms[comm]; ok {
				embed := discordgo.MessageEmbed{
					Title:       cmd.Name,
					Description: cmd.Description,
					Color:       dColorWhite,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Usage",
							Value: cmd.Usage,
						},
						{
							Name:  "Triggers",
							Value: strings.Join(cmd.Triggers, ", "),
						},
						{
							Name:  "Required permissions",
							Value: fmt.Sprintf("%v", permMap[cmd.RequiredPerms]),
						},
					},
				}
				ctx.SendEmbed(&embed)
			} else {
				ctx.Send("Command not found")
			}
		}
	},
}

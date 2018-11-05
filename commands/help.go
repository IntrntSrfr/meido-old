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
	Aliases:       []string{},
	Usage:         "m?help <optional command name>",
	RequiredPerms: discordgo.PermissionManageMessages,
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
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Usage",
							Value: cmd.Usage,
						},
						{
							Name:  "Aliases",
							Value: aliases(cmd.Aliases),
						},
						{
							Name:  "Required permissions",
							Value: fmt.Sprintf("%v", cmd.RequiredPerms),
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

func aliases(list []string) string {
	if len(list) < 1 {
		return "no aliases"
	} else {
		return strings.Join(list, ", ")
	}
}

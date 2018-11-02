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
	Function: func(context *service.Context) {

		args := strings.Split(context.Message.Content, " ")

		comms := GetCommandMap()

		if len(args) < 2 {

			list := "```\n"
			for _, val := range comms {
				list += val.Name + "\n"
			}
			list += "```"

			context.Send(list)
		} else {

			comm := args[1]

			if cmd, ok := comms[comm]; ok {
				embed := discordgo.MessageEmbed{
					Title:       cmd.Name,
					Description: cmd.Description,
					Fields: []*discordgo.MessageEmbedField{
						&discordgo.MessageEmbedField{
							Name:  "Usage",
							Value: cmd.Usage,
						},
						&discordgo.MessageEmbedField{
							Name:  "Aliases",
							Value: aliases(cmd.Aliases),
						},
						&discordgo.MessageEmbedField{
							Name:  "Required permissions",
							Value: fmt.Sprintf("%v", cmd.RequiredPerms),
						},
					},
				}
				context.SendEmbed(embed)
			} else {
				context.Send("Command not found")
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

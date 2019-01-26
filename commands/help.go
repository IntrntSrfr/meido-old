package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/service"

	"github.com/bwmarrin/discordgo"
)

var Help = Command{
	Name:          "help",
	Description:   "Shows info about commands.",
	Triggers:      []string{"m?help", "m?h"},
	Usage:         "m?help <optional command name>",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) < 2 {

			list := "```css\nList of commands:\n"
			for _, val := range comms {
				t := strings.Join(val.Triggers, ", ")
				if val.RequiresOwner {
					list += fmt.Sprintf("%v - [%v] (OWNER ONLY)\n", val.Name, t)
				} else if val.RequiredPerms == discordgo.PermissionSendMessages {
					list += fmt.Sprintf("%v - [%v]\n", val.Name, t)
				} else {
					list += fmt.Sprintf("%v - [%v] (%v)\n", val.Name, t, permMap[val.RequiredPerms])
				}
			}
			list += "```"

			_, err := ctx.Send(list)
			if err != nil {
				ctx.Send(err)
			}
		} else {

			comm := args[1]

			triggerCommand := ""
			for _, val := range comms {

				if strings.ToLower(comm) == strings.ToLower(val.Name) {
					triggerCommand = val.Name
					break
				}

				if triggerCommand == "" {
					for _, com := range val.Triggers {
						if strings.ToLower(comm) == strings.ToLower(com) {
							triggerCommand = val.Name
						}
					}
				}
			}

			if cmd, ok := comms[triggerCommand]; ok {
				perm := ""
				if cmd.RequiredPerms == discordgo.PermissionSendMessages {
					perm = "None"
				} else {
					perm = permMap[cmd.RequiredPerms]
				}
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
							Value: fmt.Sprintf("%v", perm),
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

package commands

import (
	"fmt"
	"strings"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var Help = Command{
	Name:          "Help",
	Description:   "Shows info about commands.",
	Triggers:      []string{"m?help", "m?h"},
	Usage:         "m?help <optional command name>\nm?h ban\nm?h .b",
	Category:      Utility,
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		if len(args) < 2 {

			listFilter := strings.Builder{}
			listStrikes := strings.Builder{}
			listMod := strings.Builder{}
			listFun := strings.Builder{}
			listUtil := strings.Builder{}
			listProfile := strings.Builder{}
			listOwner := strings.Builder{}

			for _, val := range comms {
				switch val.Category {
				case Filter:
					listFilter.WriteString(fmt.Sprintf("%v\t", val.Triggers[0]))
					if len(val.Triggers) > 1 {
						for _, trig := range val.Triggers[1:] {
							listFilter.WriteString(fmt.Sprintf("[%v] ", trig))
						}
					}
					listFilter.WriteString("\n")
				case Strikes:
					listStrikes.WriteString(fmt.Sprintf("%v\t", val.Triggers[0]))
					if len(val.Triggers) > 1 {
						for _, trig := range val.Triggers[1:] {
							listStrikes.WriteString(fmt.Sprintf("[%v] ", trig))
						}
					}
					listStrikes.WriteString("\n")
				case Moderation:
					listMod.WriteString(fmt.Sprintf("%v\t", val.Triggers[0]))
					if len(val.Triggers) > 1 {
						for _, trig := range val.Triggers[1:] {
							listMod.WriteString(fmt.Sprintf("[%v] ", trig))
						}
					}
					listMod.WriteString("\n")
				case Fun:
					listFun.WriteString(fmt.Sprintf("%v\t", val.Triggers[0]))
					if len(val.Triggers) > 1 {
						for _, trig := range val.Triggers[1:] {
							listFun.WriteString(fmt.Sprintf("[%v] ", trig))
						}
					}
					listFun.WriteString("\n")
				case Utility:
					listUtil.WriteString(fmt.Sprintf("%v\t", val.Triggers[0]))
					if len(val.Triggers) > 1 {
						for _, trig := range val.Triggers[1:] {
							listUtil.WriteString(fmt.Sprintf("[%v] ", trig))
						}
					}
					listUtil.WriteString("\n")
				case Profile:
					listProfile.WriteString(fmt.Sprintf("%v\t", val.Triggers[0]))
					if len(val.Triggers) > 1 {
						for _, trig := range val.Triggers[1:] {
							listProfile.WriteString(fmt.Sprintf("[%v] ", trig))
						}
					}
					listProfile.WriteString("\n")
				case Owner:
					listOwner.WriteString(fmt.Sprintf("%v\t", val.Triggers[0]))
					if len(val.Triggers) > 1 {
						for _, trig := range val.Triggers[1:] {
							listOwner.WriteString(fmt.Sprintf("[%v] ", trig))
						}
					}
					listOwner.WriteString("\n")
				default:
				}
			}

			embed := &discordgo.MessageEmbed{
				Title: "Commands and aliases",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Filter",
						Value:  fmt.Sprintf("```ini\n%v\n```", listFilter.String()),
						Inline: false,
					},
					{
						Name:   "Strike",
						Value:  fmt.Sprintf("```ini\n%v\n```", listStrikes.String()),
						Inline: false,
					},
					{
						Name:   "Moderation",
						Value:  fmt.Sprintf("```ini\n%v\n```", listMod.String()),
						Inline: false,
					},
					{
						Name:   "Fun",
						Value:  fmt.Sprintf("```ini\n%v\n```", listFun.String()),
						Inline: false,
					},
					{
						Name:   "Utility",
						Value:  fmt.Sprintf("```ini\n%v\n```", listUtil.String()),
						Inline: false,
					},
					{
						Name:   "Profile",
						Value:  fmt.Sprintf("```ini\n%v\n```", listProfile.String()),
						Inline: false,
					},
					{
						Name:   "Owner",
						Value:  fmt.Sprintf("```ini\n%v\n```", listOwner.String()),
						Inline: false,
					},
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Do `m?help <command name/alias>` to get further information about a command.",
				},
			}
			/*
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
			*/
			_, err := ctx.SendEmbed(embed)
			if err != nil {
				ctx.Send(err)
			}
		} else {

			comm := args[1:]

			scomm := strings.Join(comm, " ")

			triggerCommand := ""
			for _, val := range comms {

				if strings.ToLower(scomm) == strings.ToLower(val.Name) {
					triggerCommand = val.Name
					break
				}

				if triggerCommand == "" {
					for _, com := range val.Triggers {
						if strings.ToLower(scomm) == strings.ToLower(com) {
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

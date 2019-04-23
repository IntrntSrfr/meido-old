package commands

import (
	"os"
	"path/filepath"

	"github.com/intrntsrfr/meido/bot/service"

	"github.com/bwmarrin/discordgo"
)

var Img = Command{
	Name:          "Img",
	Description:   "Easter eggs",
	Triggers:      []string{"m?img"},
	Usage:         "m?img umr",
	Category:      Fun,
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {
		if len(args) >= 2 {
			var path string
			switch args[1] {
			case "umr":
				path, _ = filepath.Abs("../meido/bot/misc/umr.jpg")
				file, err := os.Open(path)
				if err != nil {
					ctx.Send(err.Error())
					return
				}
				defer file.Close()

				ctx.Session.ChannelFileSend(ctx.Channel.ID, "umr.jpg", file)
			case "hamster":
				path, _ = filepath.Abs("../meido/bot/misc/hamster.png")
				file, err := os.Open(path)
				if err != nil {
					ctx.Send(err.Error())
					return
				}
				defer file.Close()

				ctx.Session.ChannelFileSend(ctx.Channel.ID, "hamster.png", file)
			default:
				return
			}
		}
	},
}

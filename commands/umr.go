package commands

import (
	"meido/service"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

var Img = Command{
	Name:          "img",
	Description:   "easter eggs",
	Triggers:      []string{"m?img"},
	Usage:         "m?img umr",
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {
		if len(args) >= 2 {
			var path string
			switch args[1] {
			case "umr":
				path, _ = filepath.Abs("../meido/stuff/umr.jpg")
				file, err := os.Open(path)
				if err != nil {
					ctx.Send(err.Error())
					return
				}
				defer file.Close()

				ctx.Session.ChannelFileSend(ctx.Channel.ID, "umr.jpg", file)
			case "hamster":
				path, _ = filepath.Abs("../meido/stuff/hamster.png")
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

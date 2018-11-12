package commands

import (
	"meido-test/service"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

var Umr = Command{
	Name:          "UMR",
	Description:   "UMR UMR UMR UMR UMR UMR UMR UMR UMR",
	Triggers:      []string{"m?umr"},
	Usage:         "m?umr",
	RequiredPerms: discordgo.PermissionManageMessages,
	Execute: func(args []string, ctx *service.Context) {
		path, _ := filepath.Abs("../meido-test/stuff/umr.jpg")
		file, err := os.Open(path)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		defer file.Close()

		ctx.Session.ChannelFileSend(ctx.Channel.ID, "umr.jpg", file)
	},
}

package commands

/*
import (
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Chainban = Command{
	Name:          "ban",
	Description:   "Bans a user. Reason and prune days is optional.",
	Triggers:      []string{"m?ban", "m?b", ".b"},
	Usage:         ".b @internet surfer#0001\n.b 163454407999094786\n.b 163454407999094786 being very mean\n.b 163454407999094786 1 being very mean\n.b 163454407999094786 1",
	RequiredPerms: discordgo.PermissionSendMessages,
	Execute: func(args []string, ctx *service.Context) {

		var cnt int
		row := db.QueryRow("SELECT COUNT(*) FROM cbwhitelist WHERE userid = $1;", ctx.User.ID)
		err := row.Scan(&cnt)
		if err != nil {
			ctx.Send("Some error occured")
			return
		}

		if cnt < 1 {
			ctx.Send("You dont have permission to do that.")
			return
		}

		goodbans := 0
		badbans := 0

		for _, g := range ctx.Session.State.Guilds {

		}
	},
}
*/

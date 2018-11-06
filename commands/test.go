package commands

import (
	"fmt"
	"meido-test/models"
	"meido-test/service"

	"github.com/bwmarrin/discordgo"
)

var Test = Command{
	Name:          "test",
	Description:   "does epic testing.",
	Aliases:       []string{},
	Usage:         "m?test",
	RequiredPerms: discordgo.PermissionVoiceMoveMembers,
	Execute: func(args []string, ctx *service.Context) {
		rows, err := ctx.Db.Query("SELECT * FROM discordusers ORDER BY xp LIMIT 10;")
		if err != nil {
			return
		}

		list := []models.Discorduser{}

		for rows.Next() {
			dbu := models.Discorduser{}

			err = rows.Scan(
				&dbu.Uid,
				&dbu.Userid,
				&dbu.Username,
				&dbu.Discriminator,
				&dbu.Xp,
				&dbu.Nextxpgaintime,
				&dbu.Xpexcluded,
				&dbu.Reputation,
				&dbu.Cangivereptime)

			if err != nil {
				continue
			}
			list = append(list, dbu)
		}

		block := "```\n"

		for _, val := range list {
			block += fmt.Sprintf("%v#%v, %v reputation, %v experience\n", val.Username, val.Discriminator, val.Reputation, val.Xp)
		}

		block += "```"
		ctx.Send(block)
	},
}

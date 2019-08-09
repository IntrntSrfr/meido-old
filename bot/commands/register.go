package commands

import "github.com/bwmarrin/discordgo"

func (ch *CommandHandler) Initialize() {

	// Filter
	var FilterWord = Command{
		Name:          "Filter word",
		Description:   "Adds a word or phrase to the bots filter, making the bot automatically delete said words or phrases if posted, ignored if user has manage messages permission.",
		Triggers:      []string{"m?filterword", "m?fw"},
		Usage:         "m?fw jeff\nm?filterword jeff",
		Category:      Filter,
		RequiredPerms: discordgo.PermissionManageMessages,
		//RequiresOwner: true,
		Execute: ch.filterWord,
	}

	var FilterWordList = Command{
		Name:          "Filter word list",
		Description:   "Shows filtered words or phrases.",
		Triggers:      []string{"m?filterwordlist", "m?fwl"},
		Usage:         "m?filterwordlist\nm?fwl",
		Category:      Filter,
		RequiredPerms: discordgo.PermissionSendMessages,
		//RequiresOwner: true,
		Execute: ch.filterWordList,
	}

	var FilterInfo = Command{
		Name:          "Filter info",
		Description:   "Shows filter info.",
		Triggers:      []string{"m?filterinfo", "m?fi"},
		Usage:         "m?filterinfo\nm?fi",
		Category:      Filter,
		RequiredPerms: discordgo.PermissionManageMessages,
		//RequiresOwner: true,
		Execute: ch.filterInfo,
	}
	var FilterIgnoreChannel = Command{
		Name:          "Filter ignore channel",
		Description:   "Sets a channel to be ignored by filter.",
		Triggers:      []string{"m?filterignorechannel", "m?figch"},
		Usage:         "m?figch\nm?figch 393558442977263619\nm?filterignorechannel #gamers",
		Category:      Filter,
		RequiredPerms: discordgo.PermissionManageMessages,
		//RequiresOwner: true,
		Execute: ch.filterIgnoreChannel,
	}
	var ClearFilter = Command{
		Name:          "Clear filter",
		Description:   "Clears the list of filtered words.",
		Triggers:      []string{"m?clearfilter"},
		Usage:         "m?clearfilter",
		Category:      Filter,
		RequiredPerms: discordgo.PermissionManageMessages,
		//RequiresOwner: true,
		Execute: ch.clearFilter,
	}

	// Strikes
	var UseStrikes = Command{
		Name:          "Use strikes",
		Description:   "Toggles strike system.",
		Triggers:      []string{"m?usestrikes"},
		Usage:         "m?usestrikes",
		Category:      Strikes,
		RequiredPerms: discordgo.PermissionManageMessages,
		//RequiresOwner: true,
		Execute: ch.useStrikes,
	}
	var SetMaxStrikes = Command{
		Name:          "Set max strikes",
		Description:   "Sets max strikes. Max 10.",
		Triggers:      []string{"m?maxstrikes"},
		Usage:         "m?maxstrikes 5",
		Category:      Strikes,
		RequiredPerms: discordgo.PermissionManageMessages,
		//RequiresOwner: true,
		Execute: ch.setMaxStrikes,
	}
	var ClearStrikes = Command{
		Name:          "Clear strikes",
		Description:   "Clears the strikes on a user.",
		Triggers:      []string{"m?clearstrikes", "m?cs"},
		Usage:         "m?clearstrikes @internet surfer#0001\nm?cs 163454407999094786",
		Category:      Strikes,
		RequiredPerms: discordgo.PermissionManageMessages,
		//RequiresOwner: true,
		Execute: ch.clearStrikes,
	}
	var Warn = Command{
		Name:          "Warn",
		Description:   "Warns a user, adding a strike. Does not work if strike system is disabled.",
		Triggers:      []string{"m?warn", ".warn"},
		Usage:         "m?warn 163454407999094786\n.warn @internet surfer#0001",
		Category:      Strikes,
		RequiredPerms: discordgo.PermissionBanMembers,
		Execute:       ch.warn,
	}
	var StrikeLog = Command{
		Name:          "Strike log",
		Description:   "Shows a users strikes.",
		Triggers:      []string{"m?strikelog"},
		Usage:         "m?strikelog 163454407999094786\nm?strikelog @internet surfer#0001",
		Category:      Strikes,
		RequiredPerms: discordgo.PermissionBanMembers,
		//RequiresOwner: true,
		Execute: ch.strikeLog,
	}
	var RemoveStrike = Command{
		Name:          "Remove strike",
		Description:   "Removes a strike from a user. Use strikelog to check a users strike ids. Use the strike ids provided to find the strike you want to remove.",
		Triggers:      []string{"m?removestrike", "m?rmstrike"},
		Usage:         "m?removestrike [strike id]\nm?rmstrike 123",
		Category:      Strikes,
		RequiredPerms: discordgo.PermissionBanMembers,
		Execute:       ch.removeStrike,
	}

	// Moderation
	var Ban = Command{
		Name:          "Ban",
		Description:   "Bans a user. Reason and prune days is optional.",
		Triggers:      []string{"m?ban", "m?b", ".b", ".ban"},
		Usage:         ".b @internet surfer#0001\n.b 163454407999094786\n.b 163454407999094786 being very mean\n.b 163454407999094786 1 being very mean\n.b 163454407999094786 1",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionBanMembers,
		Execute:       ch.ban,
	}
	var Hackban = Command{
		Name:          "Hackban",
		Description:   "Hackbans one or several users. Prunes 7 days.",
		Triggers:      []string{"m?hackban", "m?hb"},
		Usage:         "m?hb 123 123 12 31 23 123 ",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionBanMembers,
		Execute:       ch.hackban,
	}
	var Unban = Command{
		Name:          "Unban",
		Description:   "Unbans a user.",
		Triggers:      []string{"m?unban", "m?ub", ".ub", ".unban"},
		Usage:         ".unban 163454407999094786",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionBanMembers,
		Execute:       ch.unban,
	}
	/*
		var ClearAFK = Command{
			Name:          "Clearafk",
			Description:   "Moves AFK users to AFK channel, if there is one.",
			Triggers:      []string{"m?clearafk"},
			Usage:         "m?clearafk",
			Category:      Utility,
			RequiredPerms: discordgo.PermissionVoiceMoveMembers,
			Execute:       ch.clearAFK,
		}
	*/
	var CoolNameBro = Command{
		Name:          "Cool name bro",
		Description:   "Renames attentionseeking nick- or usernames.",
		Triggers:      []string{"m?coolnamebro", "m?cnb"},
		Usage:         "m?coolnamebro my name is shit",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionManageNicknames,
		Execute:       ch.coolNamebro,
	}
	var Kick = Command{
		Name:          "Kick",
		Description:   "Kick a user. Reason and prune days is optional.",
		Triggers:      []string{"m?kick", "m?k", ".kick", ".k"},
		Usage:         "m?k @internet surfer#0001\n.k 163454407999094786",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionKickMembers,
		Execute:       ch.kick,
	}
	var Lockdown = Command{
		Name:          "Lockdown",
		Description:   "Locks down the current channel, denying the everyonerole send message perms.",
		Triggers:      []string{"m?lockdown"},
		Usage:         "m?lockdown",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionManageRoles,
		Execute:       ch.lockdown,
	}
	var Unlock = Command{
		Name:          "Unlock",
		Description:   "Unlocks a previously locked channel, setting the everyone roles send message permissions to default.",
		Triggers:      []string{"m?unlock"},
		Usage:         "m?unlock",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionManageRoles,
		Execute:       ch.unlock,
	}
	var SetUserRole = Command{
		Name:          "Set userrole",
		Description:   "Sets a users custom role. First provide the user, followed by the role.",
		Triggers:      []string{"m?setuserrole"},
		Usage:         "m?setuserrole 163454407999094786 kumiko",
		Category:      Moderation,
		RequiredPerms: discordgo.PermissionManageRoles,
		Execute:       ch.setUserRole,
	}

	// Utility
	var About = Command{
		Name:          "About",
		Description:   "Shows info about Meido.",
		Triggers:      []string{"m?about"},
		Usage:         "m?about",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.about,
	}
	var Avatar = Command{
		Name:          "Avatar",
		Description:   "Displays a users profile picture.",
		Triggers:      []string{"m?avatar", ">av", "m?av"},
		Usage:         ">av\n>av @internet surfer#0001\n>av 163454407999094786",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.avatar,
	}
	var Ping = Command{
		Name:          "Ping",
		Description:   "Displays bot latency.",
		Triggers:      []string{"m?ping"},
		Usage:         "m?ping",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.ping,
	}
	var Help = Command{
		Name:          "Help",
		Description:   "Shows info about commands.",
		Triggers:      []string{"m?help", "m?h"},
		Usage:         "m?help <optional command name>\nm?h ban\nm?h .b",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.help,
	}
	var Inrole = Command{
		Name:          "Inrole",
		Description:   "Shows a list of who and how many users who are in a specified role.",
		Triggers:      []string{"m?inrole"},
		Usage:         "m?inrole gamers",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.inrole,
	}
	var WithNick = Command{
		Name:          "With nick",
		Description:   "Shows how many has an input user- or nickname.",
		Triggers:      []string{"m?withnick"},
		Usage:         "m?withnick meido",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.withNick,
	}
	var WithTag = Command{
		Name:          "With tag",
		Description:   "Shows how many has an input discriminator.",
		Triggers:      []string{"m?withtag"},
		Usage:         "m?withtag <0001/#0001>",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.withTag,
	}
	var Server = Command{
		Name:          "Server",
		Description:   "Shows information about the current server.",
		Triggers:      []string{"m?server", "m?s"},
		Usage:         "m?server",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.server,
	}
	var ListUserRoles = Command{
		Name:          "List userroles",
		Description:   "Sets a users custom role. First provide the user, followed by the role.",
		Triggers:      []string{"m?listuserroles"},
		Usage:         "m?listuserroles",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionManageRoles,
		Execute:       ch.listUserRoles,
	}
	var Invite = Command{
		Name:          "Invite",
		Description:   "Sends bot invite link and support server invite.",
		Triggers:      []string{"m?invite", "m?inv"},
		Usage:         "m?invite",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.invite,
	}
	var Feedback = Command{
		Name:          "Feedback",
		Description:   "Sends your very nice and helpful feedback to the Meido Café.",
		Triggers:      []string{"m?feedback", "m?fb"},
		Usage:         "m?feedback wow what a really COOL and NICE bot that works flawlessly",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.feedback,
	}
	var MyRole = Command{
		Name:          "Myrole",
		Description:   "Gets information about a custom role, or lets the owner of the role edit its name or color.",
		Triggers:      []string{"m?myrole"},
		Usage:         "m?myrole\nm?myrole 163454407999094786\nm?myrole color c0ff33\nm?myrole name kumiko",
		Category:      Utility,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.myRole,
	}

	// Fun
	var Img = Command{
		Name:          "Img",
		Description:   "Easter eggs",
		Triggers:      []string{"m?img"},
		Usage:         "m?img umr",
		Category:      Fun,
		RequiredPerms: discordgo.PermissionManageMessages,
		Execute:       ch.img,
	}

	// Profile
	var ShowProfile = Command{
		Name:          "Profile",
		Description:   "Shows a user profile.",
		Triggers:      []string{"m?profile", "m?p"},
		Usage:         "m?profile\nm?profile @internet surfer#0001\nm?profile 163454407999094786",
		Category:      Profile,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.showProfile,
	}
	var Rep = Command{
		Name:          "Rep",
		Description:   "Gives a user a reputation point or checks whether you can give it or not.",
		Triggers:      []string{"m?rep"},
		Usage:         "m?rep\nm?rep @internet surfer#0001\nm?rep 163454407999094786",
		Category:      Profile,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.rep,
	}
	var Repleaderboard = Command{
		Name:          "Rep leaderboard",
		Description:   "Checks the reputation leaderboard.",
		Triggers:      []string{"m?rplb"},
		Usage:         "m?rplb",
		Category:      Profile,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.repleaderboard,
	}
	var XpLeaderboard = Command{
		Name:          "XP leaderboard",
		Description:   "Checks local leaderboard.",
		Triggers:      []string{"m?xplb"},
		Usage:         "m?xplb",
		Category:      Profile,
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.xpLeaderboard,
	}
	var GlobalXpLeaderboard = Command{
		Name:          "Global XP Leaderboard",
		Description:   "Checks the global xp leaderboard.",
		Triggers:      []string{"m?gxplb"},
		Usage:         "m?gxplb",
		RequiredPerms: discordgo.PermissionSendMessages,
		Execute:       ch.globalXpLeaderboard,
	}
	var XpIgnoreChannel = Command{
		Name:          "XP ignore channel",
		Description:   "Adds or removes a channel to or from the xp ignored list.",
		Triggers:      []string{"m?xpignorechannel", "m?xpigch"},
		Usage:         "m?xpigch\nm?xpigch 123123123123",
		Category:      Profile,
		RequiredPerms: discordgo.PermissionManageChannels,
		Execute:       ch.xpIgnoreChannel,
	}

	// Owner
	var Dm = Command{
		Name:          "DM",
		Description:   "Sends a direct message. Owner only.",
		Triggers:      []string{"m?dm"},
		Usage:         "m?dm 163454407999094786 jeff",
		Category:      Owner,
		RequiredPerms: discordgo.PermissionSendMessages,
		RequiresOwner: true,
		Execute:       ch.dm,
	}
	var Msg = Command{
		Name:          "Msg",
		Description:   "Sends a message to a channel. Owner only.",
		Triggers:      []string{"m?msg"},
		Usage:         "m?msg 497106582144942101 jeff",
		Category:      Owner,
		RequiredPerms: discordgo.PermissionSendMessages,
		RequiresOwner: true,
		Execute:       ch.msg,
	}
	var Refresh = Command{
		Name:          "Refresh",
		Description:   "Force refreshes db stuff in case its stuck. Owner only",
		Triggers:      []string{"m?refresh"},
		Usage:         "m?refresh",
		Category:      Owner,
		RequiredPerms: discordgo.PermissionSendMessages,
		RequiresOwner: true,
		Execute:       ch.refresh,
	}

	// Filter
	ch.comms.RegisterCommand(FilterWord)
	ch.comms.RegisterCommand(FilterWordList)
	ch.comms.RegisterCommand(FilterInfo)
	ch.comms.RegisterCommand(FilterIgnoreChannel)
	ch.comms.RegisterCommand(ClearFilter)

	// Strikes
	ch.comms.RegisterCommand(UseStrikes)
	ch.comms.RegisterCommand(SetMaxStrikes)
	ch.comms.RegisterCommand(ClearStrikes)
	ch.comms.RegisterCommand(Warn)
	ch.comms.RegisterCommand(StrikeLog)
	//ch.RegisterCommand(StrikeLogAll)
	ch.comms.RegisterCommand(RemoveStrike)

	// Moderation
	ch.comms.RegisterCommand(Ban)
	ch.comms.RegisterCommand(Hackban)
	ch.comms.RegisterCommand(Unban)
	//comms.RegisterCommand(ClearAFK)
	ch.comms.RegisterCommand(CoolNameBro)
	//ch.comms.RegisterCommand(NiceNameBro)
	ch.comms.RegisterCommand(Kick)
	ch.comms.RegisterCommand(Lockdown)
	ch.comms.RegisterCommand(Unlock)
	ch.comms.RegisterCommand(SetUserRole)

	// Utility
	ch.comms.RegisterCommand(About)
	ch.comms.RegisterCommand(Avatar)
	ch.comms.RegisterCommand(Ping)
	ch.comms.RegisterCommand(Help)
	ch.comms.RegisterCommand(Inrole)
	ch.comms.RegisterCommand(WithNick)
	ch.comms.RegisterCommand(WithTag)
	//comms.RegisterCommand(Role)
	ch.comms.RegisterCommand(Server)
	//comms.RegisterCommand(User)
	//comms.RegisterCommand(ListRoles)
	ch.comms.RegisterCommand(ListUserRoles)
	ch.comms.RegisterCommand(Invite)
	ch.comms.RegisterCommand(Feedback)
	ch.comms.RegisterCommand(MyRole)

	// Fun
	ch.comms.RegisterCommand(Img)

	// Profile
	ch.comms.RegisterCommand(ShowProfile)
	ch.comms.RegisterCommand(Rep)
	ch.comms.RegisterCommand(Repleaderboard)
	ch.comms.RegisterCommand(XpLeaderboard)
	ch.comms.RegisterCommand(GlobalXpLeaderboard)
	ch.comms.RegisterCommand(XpIgnoreChannel)

	// Owner
	//comms.RegisterCommand(Test)
	ch.comms.RegisterCommand(Dm)
	ch.comms.RegisterCommand(Msg)
	ch.comms.RegisterCommand(Refresh)
}

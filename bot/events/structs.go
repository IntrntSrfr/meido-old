package events

type Config struct {
	OwoToken      string
	DmLogChannels []string
	OwnerIds      []string
}

var permMap = map[int]string{
	1:          "create instant invite",
	2:          "kick members",
	4:          "ban members",
	8:          "administrator",
	16:         "manage channels",
	32:         "manage server",
	64:         "add reactions",
	128:        "view audit log",
	256:        "priority speaker",
	1024:       "view channel",
	2048:       "send messages",
	4096:       "send tts messages",
	8192:       "manage messages",
	16384:      "embed links",
	32768:      "attach files",
	65536:      "read message history",
	131072:     "mention everyone",
	262144:     "use external emojis",
	1048576:    "connect",
	2097152:    "speak",
	4194304:    "mute members",
	8388608:    "deafen members",
	16777216:   "move members",
	33554432:   "use VAD",
	67108864:   "change nickname",
	134217728:  "manage nicknames",
	268435456:  "manage roles",
	536870912:  "manage webhooks",
	1073741824: "manage emojis",
}

var verificationMap = map[int]string{
	0: "Unrestricted.",
	1: "Email verification.",
	2: "Email verification and account must be at least 5 minutes old.",
	3: "Email verification, account must be at least 5 minutes old and user must have been on server for 10 minutes.",
	4: "Verified phone to Discord account.",
}

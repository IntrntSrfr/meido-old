package bot

type Config struct {
	Token            string   `json:"Token"`
	OwoAPIKey        string   `json:"OWOToken"`
	ConnectionString string   `json:"Connectionstring"`
	DmLogChannels    []string `json:"DmLogChannels"`
	OwnerIds         []string `json:"OwnerIds"`
}

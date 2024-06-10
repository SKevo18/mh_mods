package src

const (
	AppName = "Moorhuhn Mods"
	AppID = "ac.svit.mh-mods"
	AppVersion = "0.1.0"
)

var SupportedGames = map[string]string{
	"mhk_1": "Moorhuhn Kart 1",
	"mhk_2": "Moorhuhn Kart 2",
	"mhk_3": "Moorhuhn Kart 3",
	"mhk_4": "Moorhuhn Kart 4",
}

func GetGames() []string {
	var games []string
	for _, game := range SupportedGames {
		games = append(games, game)
	}
	return games
}

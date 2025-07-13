package core

type TeamID string

const (
	TeamPlayers TeamID = "Players"
	TeamEnemies TeamID = "Enemies"
)

func TeamFromString(s string) TeamID {
	switch s {
	case "players":
		return TeamPlayers
	case "enemies":
		return TeamEnemies
	default:
		return TeamEnemies
	}
}

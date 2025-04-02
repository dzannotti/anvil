package core

type TeamID string

const (
	TeamPlayers TeamID = "Players"
	TeamEnemies TeamID = "Enemies"
)

type Team struct {
	Name    string
	Members []*Actor
}

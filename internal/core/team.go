package core

type TeamID string

const (
	TeamPlayers TeamID = "Players"
	TeamEnemies TeamID = "Enemies"
	TeamNeutral TeamID = "Neutral"
	TeamGaea    TeamID = "Gaea"
)

type Team struct {
	Name    string
	Members []*Actor
}

package team

type Team int

const (
	None Team = iota
	Player
	Enemy
)

func (t Team) String() string {
	switch t {
	case None:
		return "None"
	case Player:
		return "Player"
	case Enemy:
		return "Enemy"
	default:
		return "Unknown"
	}
}

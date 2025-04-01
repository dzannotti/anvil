package definition

type Encounter interface {
	IsOver() bool
	Teams() []Team
}

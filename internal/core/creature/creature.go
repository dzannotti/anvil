package creature

import (
	"anvil/internal/core/definition"
	"anvil/internal/log"
)

type Creature struct {
	log          *log.EventLog
	name         string
	hitPoints    int
	maxHitPoints int
	actions      []definition.Action
}

func New(log *log.EventLog, name string, hitPoints int) *Creature {
	return &Creature{
		log:          log,
		name:         name,
		hitPoints:    hitPoints,
		maxHitPoints: hitPoints,
	}
}

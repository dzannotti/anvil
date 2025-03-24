package creature

import "anvil/internal/log"

type Creature struct {
	log          *log.EventLog
	name         string
	hitPoints    int
	maxHitPoints int
}

func New(log *log.EventLog, name string, hitPoints int) *Creature {
	return &Creature{
		log:          log,
		name:         name,
		hitPoints:    hitPoints,
		maxHitPoints: hitPoints,
	}
}

package ai

import "anvil/internal/ai/simple"

type AI interface {
	Play()
}

var NewSimple = simple.New

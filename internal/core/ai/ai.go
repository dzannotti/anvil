package ai

import "anvil/internal/core/ai/simple"

type AI interface {
	Play()
}

var NewSimple = simple.New

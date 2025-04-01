package ai

import "anvil/internal/ai/simple"

type AI interface {
	Play()
}

type Simple = simple.Simple

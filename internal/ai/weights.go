package ai

func NewBerserkerWeights() *Weights {
	return &Weights{
		Weights: map[string]float32{
			"damage_enemy":        2.0,
			"friendly_fire":       0.5,
			"survival_threat":     0.3,
			"kill_potential":      1.5,
			"enemy_proximity":     0.2,
			"threat_priority":     0.8,
			"low_health_bonus":    1.8,
			"tactical_value":      0.6,
			"movement_efficiency": 0.4,
		},
	}
}

func NewDefensiveWeights() *Weights {
	return &Weights{
		Weights: map[string]float32{
			"damage_enemy":        1.0,
			"friendly_fire":       2.0,
			"survival_threat":     2.0,
			"kill_potential":      0.8,
			"enemy_proximity":     1.8,
			"threat_priority":     1.5,
			"low_health_bonus":    1.0,
			"tactical_value":      1.3,
			"movement_efficiency": 1.6,
		},
	}
}

func NewDefaultWeights() *Weights {
	return &Weights{
		Weights: map[string]float32{
			"damage_enemy":        1.0,
			"friendly_fire":       1.5,
			"survival_threat":     1.0,
			"kill_potential":      1.2,
			"enemy_proximity":     1.0,
			"threat_priority":     1.2,
			"low_health_bonus":    1.4,
			"tactical_value":      1.0,
			"movement_efficiency": 1.0,
		},
	}
}

package ai

// NewBerserkerWeights creates weights for aggressive, damage-focused AI behavior.
func NewBerserkerWeights() *Weights {
	return &Weights{
		DamageEnemy:        2.0,
		FriendlyFire:       0.5,
		SurvivalThreat:     0.3,
		KillPotential:      1.5,
		EnemyProximity:     0.2,
		ThreatPriority:     0.8,
		LowHealthBonus:     1.8,
		TacticalValue:      0.6,
		MovementEfficiency: 0.4,
	}
}

// NewDefensiveWeights creates weights for cautious, survival-focused AI behavior.
func NewDefensiveWeights() *Weights {
	return &Weights{
		DamageEnemy:        1.0,
		FriendlyFire:       2.0,
		SurvivalThreat:     2.0,
		KillPotential:      0.8,
		EnemyProximity:     1.8,
		ThreatPriority:     1.5,
		LowHealthBonus:     1.0,
		TacticalValue:      1.3,
		MovementEfficiency: 1.6,
	}
}

// NewDefaultWeights creates balanced weights for standard AI behavior.
func NewDefaultWeights() *Weights {
	return &Weights{
		DamageEnemy:        1.0,
		FriendlyFire:       1.5,
		SurvivalThreat:     1.0,
		KillPotential:      1.2,
		EnemyProximity:     1.0,
		ThreatPriority:     1.2,
		LowHealthBonus:     1.4,
		TacticalValue:      1.0,
		MovementEfficiency: 1.0,
	}
}

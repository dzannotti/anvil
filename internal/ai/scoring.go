package ai

// Weights represents the complete set of weights for AI decision making.
// Each field corresponds to a specific aspect of combat evaluation.
type Weights struct {
	// Combat weights
	DamageEnemy     float32 // Priority for dealing damage to enemies
	FriendlyFire    float32 // Penalty for risking damage to allies
	KillPotential   float32 // Bonus for actions that could kill an enemy
	
	// Positioning weights  
	SurvivalThreat     float32 // Penalty for dangerous positions
	EnemyProximity     float32 // Consideration of distance to enemies
	MovementEfficiency float32 // Preference for efficient movement
	
	// Target selection weights
	ThreatPriority   float32 // Priority for high-threat targets
	LowHealthBonus   float32 // Bonus for targeting low-health enemies
	TacticalValue    float32 // General tactical value assessment
}

// Scores contains raw metric values before weight application.
type Scores struct {
	// Combat scores
	DamageEnemy   int // Expected damage to enemies
	FriendlyFire  int // Risk of damaging allies
	KillPotential int // Likelihood of killing an enemy
	
	// Positioning scores
	SurvivalThreat     int // Danger level of position
	EnemyProximity     int // Distance consideration score
	MovementEfficiency int // Movement cost evaluation
	
	// Target selection scores
	ThreatPriority   int // Target threat level
	LowHealthBonus   int // Low health target bonus
	TacticalValue    int // Overall tactical benefit
}

// WeightedScores contains the final weighted scores used for decision making.
type WeightedScores struct {
	// Combat scores
	DamageEnemy   int // Weighted damage potential
	FriendlyFire  int // Weighted friendly fire risk
	KillPotential int // Weighted kill opportunity
	
	// Positioning scores
	SurvivalThreat     int // Weighted survival consideration
	EnemyProximity     int // Weighted proximity factor
	MovementEfficiency int // Weighted movement preference
	
	// Target selection scores
	ThreatPriority   int // Weighted threat prioritization
	LowHealthBonus   int // Weighted low health targeting
	TacticalValue    int // Weighted tactical assessment
}

// Total calculates the sum of all weighted scores for final decision making.
func (w *WeightedScores) Total() int {
	return w.DamageEnemy + w.FriendlyFire + w.KillPotential +
		w.SurvivalThreat + w.EnemyProximity + w.MovementEfficiency +
		w.ThreatPriority + w.LowHealthBonus + w.TacticalValue
}

// ApplyWeights transforms raw scores into weighted scores using the provided weights.
func (s *Scores) ApplyWeights(weights *Weights) *WeightedScores {
	return &WeightedScores{
		DamageEnemy:        int(float32(s.DamageEnemy) * weights.DamageEnemy),
		FriendlyFire:       int(float32(s.FriendlyFire) * weights.FriendlyFire),
		KillPotential:      int(float32(s.KillPotential) * weights.KillPotential),
		SurvivalThreat:     int(float32(s.SurvivalThreat) * weights.SurvivalThreat),
		EnemyProximity:     int(float32(s.EnemyProximity) * weights.EnemyProximity),
		MovementEfficiency: int(float32(s.MovementEfficiency) * weights.MovementEfficiency),
		ThreatPriority:     int(float32(s.ThreatPriority) * weights.ThreatPriority),
		LowHealthBonus:     int(float32(s.LowHealthBonus) * weights.LowHealthBonus),
		TacticalValue:      int(float32(s.TacticalValue) * weights.TacticalValue),
	}
}

// mapToScores converts legacy map[string]int metrics to structured Scores.
// This function bridges the gap between the old metric system and new struct-based scoring.
func mapToScores(metricsMap map[string]int) *Scores {
	return &Scores{
		DamageEnemy:        metricsMap["damage_enemy"],
		FriendlyFire:       metricsMap["friendly_fire"],
		KillPotential:      metricsMap["kill_potential"],
		SurvivalThreat:     metricsMap["survival_threat"],
		EnemyProximity:     metricsMap["enemy_proximity"],
		MovementEfficiency: metricsMap["movement_efficiency"],
		ThreatPriority:     metricsMap["threat_priority"],
		LowHealthBonus:     metricsMap["low_health_bonus"],
		TacticalValue:      metricsMap["tactical_value"],
	}
}
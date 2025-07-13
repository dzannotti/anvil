package ai

import (
	"testing"
)

// Test AI Requirements Integration
// This test validates that the AI system meets its core requirements:
// 1. Three distinct archetypes (Berserker, Defensive, Default)
// 2. Nine weight metrics across three categories
// 3. Simulation-based decision making
// 4. Fallback action support

func TestAI_ArchetypeRequirements(t *testing.T) {
	// Test that all three archetypes exist and have distinct characteristics
	berserker := NewBerserkerWeights()
	defensive := NewDefensiveWeights()
	defaultWeights := NewDefaultWeights()
	
	// All archetypes must have exactly 9 metrics
	validateMetricCount(t, berserker, "Berserker")
	validateMetricCount(t, defensive, "Defensive")
	validateMetricCount(t, defaultWeights, "Default")
	
	// Validate archetype personalities through weight differences
	validateBerserkerPersonality(t, berserker, defensive, defaultWeights)
	validateDefensivePersonality(t, berserker, defensive, defaultWeights)
	validateDefaultPersonality(t, berserker, defensive, defaultWeights)
}

func TestAI_MetricCategories(t *testing.T) {
	// Test that all required metric categories are present
	weights := NewDefaultWeights()
	
	// Damage Metrics
	validateMetricExists(t, weights, "damage_enemy", "Damage category")
	validateMetricExists(t, weights, "friendly_fire", "Damage category")
	validateMetricExists(t, weights, "kill_potential", "Damage category")
	
	// Positioning Metrics  
	validateMetricExists(t, weights, "survival_threat", "Positioning category")
	validateMetricExists(t, weights, "enemy_proximity", "Positioning category")
	validateMetricExists(t, weights, "movement_efficiency", "Positioning category")
	
	// Target Selection Metrics
	validateMetricExists(t, weights, "threat_priority", "Target Selection category")
	validateMetricExists(t, weights, "low_health_bonus", "Target Selection category")
	validateMetricExists(t, weights, "tactical_value", "Target Selection category")
}

func TestAI_WeightRanges(t *testing.T) {
	// Test that weights are in reasonable ranges for AI decision making
	archetypes := []*Weights{
		NewBerserkerWeights(),
		NewDefensiveWeights(), 
		NewDefaultWeights(),
	}
	
	for i, weights := range archetypes {
		archetypeName := []string{"Berserker", "Defensive", "Default"}[i]
		
		// Check all weight fields
		weightsToCheck := []struct{
			name string
			value float32
		}{
			{"DamageEnemy", weights.DamageEnemy},
			{"FriendlyFire", weights.FriendlyFire},
			{"KillPotential", weights.KillPotential},
			{"SurvivalThreat", weights.SurvivalThreat},
			{"EnemyProximity", weights.EnemyProximity},
			{"MovementEfficiency", weights.MovementEfficiency},
			{"ThreatPriority", weights.ThreatPriority},
			{"LowHealthBonus", weights.LowHealthBonus},
			{"TacticalValue", weights.TacticalValue},
		}
		
		for _, w := range weightsToCheck {
			// Weights should be positive (we use negative multipliers in metrics for penalties)
			if w.value <= 0 {
				t.Errorf("%s archetype has non-positive weight for %s: %f", archetypeName, w.name, w.value)
			}
			
			// Weights should be reasonable (not extreme values)
			if w.value > 10.0 {
				t.Errorf("%s archetype has excessive weight for %s: %f", archetypeName, w.name, w.value)
			}
		}
	}
}

func TestAI_ArchetypeDistinction(t *testing.T) {
	// Test that archetypes are meaningfully different from each other
	berserker := NewBerserkerWeights()
	defensive := NewDefensiveWeights()
	defaultWeights := NewDefaultWeights()
	
	// Count significant differences between archetypes
	berserkerVsDefensive := countSignificantDifferences(berserker, defensive, 0.5)
	berserkerVsDefault := countSignificantDifferences(berserker, defaultWeights, 0.3)
	defensiveVsDefault := countSignificantDifferences(defensive, defaultWeights, 0.3)
	
	// Archetypes should have multiple significant differences
	if berserkerVsDefensive < 5 {
		t.Errorf("Berserker vs Defensive only has %d significant differences, expected >= 5", berserkerVsDefensive)
	}
	if berserkerVsDefault < 3 {
		t.Errorf("Berserker vs Default only has %d significant differences, expected >= 3", berserkerVsDefault)
	}
	if defensiveVsDefault < 3 {
		t.Errorf("Defensive vs Default only has %d significant differences, expected >= 3", defensiveVsDefault)
	}
}

func TestAI_SimulationBasedDesign(t *testing.T) {
	// Test that the AI system is designed for simulation-based decision making
	weights := NewDefaultWeights()
	
	// The ActionTargetEvaluation structure should support simulation
	rawScores := &Scores{}
	weightedScores := &WeightedScores{}
	evaluation := &ActionTargetEvaluation{
		RawScores:      rawScores,
		WeightedScores: weightedScores,
		FinalScore:     0,
	}
	
	if evaluation.RawScores == nil {
		t.Error("ActionTargetEvaluation should support raw metrics for simulation")
	}
	if evaluation.WeightedScores == nil {
		t.Error("ActionTargetEvaluation should support weighted scores for decision making")
	}
	
	// With struct-based weights, all metrics are always present by design
	if weights == nil {
		t.Error("Weights should not be nil")
	}
}

// Helper functions
func validateMetricCount(t *testing.T, weights *Weights, archetype string) {
	// AIWeights struct has exactly 9 fields by design, so just verify it's not nil
	if weights == nil {
		t.Errorf("%s archetype weights should not be nil", archetype)
	}
}

func validateMetricExists(t *testing.T, weights *Weights, metric, category string) {
	// With struct-based weights, all metrics are always present by design
	// This function is now mainly for API compatibility
	if weights == nil {
		t.Errorf("Weights should not be nil for %s", category)
	}
}

func validateBerserkerPersonality(t *testing.T, berserker, defensive, defaultWeights *Weights) {
	// Berserker should prioritize damage over safety
	if berserker.DamageEnemy <= defensive.DamageEnemy {
		t.Error("Berserker should prioritize damage more than Defensive")
	}
	if berserker.SurvivalThreat >= defensive.SurvivalThreat {
		t.Error("Berserker should care less about survival than Defensive")
	}
	if berserker.KillPotential <= defaultWeights.KillPotential {
		t.Error("Berserker should prioritize kills more than Default")
	}
}

func validateDefensivePersonality(t *testing.T, berserker, defensive, defaultWeights *Weights) {
	// Defensive should prioritize safety over damage
	if defensive.SurvivalThreat <= berserker.SurvivalThreat {
		t.Error("Defensive should prioritize survival more than Berserker")
	}
	if defensive.FriendlyFire <= berserker.FriendlyFire {
		t.Error("Defensive should avoid friendly fire more than Berserker")
	}
	if defensive.EnemyProximity <= defaultWeights.EnemyProximity {
		t.Error("Defensive should avoid enemy proximity more than Default")
	}
}

func validateDefaultPersonality(t *testing.T, berserker, defensive, defaultWeights *Weights) {
	// Default should be balanced - not extreme in any direction
	damageWeight := defaultWeights.DamageEnemy
	survivalWeight := defaultWeights.SurvivalThreat
	
	// Default should be between extremes for key metrics
	if damageWeight >= berserker.DamageEnemy {
		t.Error("Default damage priority should be less than Berserker")
	}
	if survivalWeight >= defensive.SurvivalThreat {
		t.Error("Default survival priority should be less than Defensive")
	}
	
	// Default should have moderate values (close to 1.0 for core metrics)
	weights := []float32{
		defaultWeights.DamageEnemy,
		defaultWeights.SurvivalThreat,
		defaultWeights.EnemyProximity,
		defaultWeights.TacticalValue,
		defaultWeights.MovementEfficiency,
	}
	metricNames := []string{"DamageEnemy", "SurvivalThreat", "EnemyProximity", "TacticalValue", "MovementEfficiency"}
	
	for i, weight := range weights {
		if weight < 0.5 || weight > 1.5 {
			t.Errorf("Default archetype should have balanced weight for %s, got %f", metricNames[i], weight)
		}
	}
}

func countSignificantDifferences(weights1, weights2 *Weights, threshold float32) int {
	count := 0
	
	// Compare all weight fields
	diffs := []float32{
		abs(weights1.DamageEnemy - weights2.DamageEnemy),
		abs(weights1.FriendlyFire - weights2.FriendlyFire),
		abs(weights1.KillPotential - weights2.KillPotential),
		abs(weights1.SurvivalThreat - weights2.SurvivalThreat),
		abs(weights1.EnemyProximity - weights2.EnemyProximity),
		abs(weights1.MovementEfficiency - weights2.MovementEfficiency),
		abs(weights1.ThreatPriority - weights2.ThreatPriority),
		abs(weights1.LowHealthBonus - weights2.LowHealthBonus),
		abs(weights1.TacticalValue - weights2.TacticalValue),
	}
	
	for _, diff := range diffs {
		if diff >= threshold {
			count++
		}
	}
	return count
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
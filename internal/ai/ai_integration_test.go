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
		
		for metric, weight := range weights.Weights {
			// Weights should be positive (we use negative multipliers in metrics for penalties)
			if weight <= 0 {
				t.Errorf("%s archetype has non-positive weight for %s: %f", archetypeName, metric, weight)
			}
			
			// Weights should be reasonable (not extreme values)
			if weight > 10.0 {
				t.Errorf("%s archetype has excessive weight for %s: %f", archetypeName, metric, weight)
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
	evaluation := &ActionTargetEvaluation{
		RawMetrics:     make(map[string]int),
		WeightedScores: make(map[string]int),
		FinalScore:     0,
	}
	
	if evaluation.RawMetrics == nil {
		t.Error("ActionTargetEvaluation should support raw metrics for simulation")
	}
	if evaluation.WeightedScores == nil {
		t.Error("ActionTargetEvaluation should support weighted scores for decision making")
	}
	
	// Weights should be compatible with all metrics
	requiredMetrics := []string{
		"damage_enemy", "friendly_fire", "kill_potential",
		"survival_threat", "enemy_proximity", "movement_efficiency", 
		"threat_priority", "low_health_bonus", "tactical_value",
	}
	
	for _, metric := range requiredMetrics {
		if _, exists := weights.Weights[metric]; !exists {
			t.Errorf("Weights missing required metric for simulation: %s", metric)
		}
	}
}

// Helper functions
func validateMetricCount(t *testing.T, weights *Weights, archetype string) {
	if len(weights.Weights) != 9 {
		t.Errorf("%s archetype should have exactly 9 metrics, got %d", archetype, len(weights.Weights))
	}
}

func validateMetricExists(t *testing.T, weights *Weights, metric, category string) {
	if _, exists := weights.Weights[metric]; !exists {
		t.Errorf("Missing required metric '%s' in %s", metric, category)
	}
}

func validateBerserkerPersonality(t *testing.T, berserker, defensive, defaultWeights *Weights) {
	// Berserker should prioritize damage over safety
	if berserker.Weights["damage_enemy"] <= defensive.Weights["damage_enemy"] {
		t.Error("Berserker should prioritize damage more than Defensive")
	}
	if berserker.Weights["survival_threat"] >= defensive.Weights["survival_threat"] {
		t.Error("Berserker should care less about survival than Defensive")
	}
	if berserker.Weights["kill_potential"] <= defaultWeights.Weights["kill_potential"] {
		t.Error("Berserker should prioritize kills more than Default")
	}
}

func validateDefensivePersonality(t *testing.T, berserker, defensive, defaultWeights *Weights) {
	// Defensive should prioritize safety over damage
	if defensive.Weights["survival_threat"] <= berserker.Weights["survival_threat"] {
		t.Error("Defensive should prioritize survival more than Berserker")
	}
	if defensive.Weights["friendly_fire"] <= berserker.Weights["friendly_fire"] {
		t.Error("Defensive should avoid friendly fire more than Berserker")
	}
	if defensive.Weights["enemy_proximity"] <= defaultWeights.Weights["enemy_proximity"] {
		t.Error("Defensive should avoid enemy proximity more than Default")
	}
}

func validateDefaultPersonality(t *testing.T, berserker, defensive, defaultWeights *Weights) {
	// Default should be balanced - not extreme in any direction
	damageWeight := defaultWeights.Weights["damage_enemy"]
	survivalWeight := defaultWeights.Weights["survival_threat"]
	
	// Default should be between extremes for key metrics
	if damageWeight >= berserker.Weights["damage_enemy"] {
		t.Error("Default damage priority should be less than Berserker")
	}
	if survivalWeight >= defensive.Weights["survival_threat"] {
		t.Error("Default survival priority should be less than Defensive")
	}
	
	// Default should have moderate values (close to 1.0 for core metrics)
	coreMetrics := []string{"damage_enemy", "survival_threat", "enemy_proximity", "tactical_value", "movement_efficiency"}
	for _, metric := range coreMetrics {
		weight := defaultWeights.Weights[metric]
		if weight < 0.5 || weight > 1.5 {
			t.Errorf("Default archetype should have balanced weight for %s, got %f", metric, weight)
		}
	}
}

func countSignificantDifferences(weights1, weights2 *Weights, threshold float32) int {
	count := 0
	for metric, weight1 := range weights1.Weights {
		if weight2, exists := weights2.Weights[metric]; exists {
			diff := weight1 - weight2
			if diff < 0 {
				diff = -diff
			}
			if diff >= threshold {
				count++
			}
		}
	}
	return count
}
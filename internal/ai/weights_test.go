package ai

import (
	"testing"
)

func TestNewBerserkerWeights(t *testing.T) {
	weights := NewBerserkerWeights()
	
	// Validate berserker characteristics
	assertWeight(t, weights, "damage_enemy", 2.0, "Berserker should prioritize damage")
	assertWeight(t, weights, "kill_potential", 1.5, "Berserker should prioritize kills")
	assertWeight(t, weights, "low_health_bonus", 1.8, "Berserker should target low health enemies")
	
	// Should have lower defensive priorities
	assertWeight(t, weights, "survival_threat", 0.3, "Berserker should care less about survival")
	assertWeight(t, weights, "friendly_fire", 0.5, "Berserker should care less about friendly fire")
	assertWeight(t, weights, "enemy_proximity", 0.2, "Berserker should care less about enemy proximity")
	
	// Should have all 9 required metrics
	assertAllMetricsPresent(t, weights)
}

func TestNewDefensiveWeights(t *testing.T) {
	weights := NewDefensiveWeights()
	
	// Validate defensive characteristics
	assertWeight(t, weights, "survival_threat", 2.0, "Defensive should prioritize survival")
	assertWeight(t, weights, "friendly_fire", 2.0, "Defensive should avoid friendly fire")
	assertWeight(t, weights, "enemy_proximity", 1.8, "Defensive should avoid enemy proximity")
	assertWeight(t, weights, "movement_efficiency", 1.6, "Defensive should value positioning")
	
	// Should have lower aggressive priorities
	assertWeight(t, weights, "damage_enemy", 1.0, "Defensive should be moderate on damage")
	assertWeight(t, weights, "kill_potential", 0.8, "Defensive should care less about kills")
	
	// Should have all 9 required metrics
	assertAllMetricsPresent(t, weights)
}

func TestNewDefaultWeights(t *testing.T) {
	weights := NewDefaultWeights()
	
	// Validate balanced characteristics - most weights around 1.0
	assertWeight(t, weights, "damage_enemy", 1.0, "Default should be balanced on damage")
	assertWeight(t, weights, "survival_threat", 1.0, "Default should be balanced on survival")
	assertWeight(t, weights, "enemy_proximity", 1.0, "Default should be balanced on proximity")
	assertWeight(t, weights, "tactical_value", 1.0, "Default should be balanced on tactics")
	assertWeight(t, weights, "movement_efficiency", 1.0, "Default should be balanced on movement")
	
	// Some slight preferences for safety and targeting
	assertWeight(t, weights, "friendly_fire", 1.5, "Default should prefer avoiding friendly fire")
	assertWeight(t, weights, "low_health_bonus", 1.4, "Default should prefer low health targets")
	assertWeight(t, weights, "kill_potential", 1.2, "Default should slightly prefer kills")
	assertWeight(t, weights, "threat_priority", 1.2, "Default should slightly prefer threats")
	
	// Should have all 9 required metrics
	assertAllMetricsPresent(t, weights)
}

func TestWeightsArchetypeDistinction(t *testing.T) {
	berserker := NewBerserkerWeights()
	defensive := NewDefensiveWeights()
	defaultW := NewDefaultWeights()
	
	// Berserker should be most aggressive
	if berserker.Weights["damage_enemy"] <= defensive.Weights["damage_enemy"] {
		t.Error("Berserker should have higher damage priority than defensive")
	}
	if berserker.Weights["damage_enemy"] <= defaultW.Weights["damage_enemy"] {
		t.Error("Berserker should have higher damage priority than default")
	}
	
	// Defensive should be most cautious
	if defensive.Weights["survival_threat"] <= berserker.Weights["survival_threat"] {
		t.Error("Defensive should have higher survival priority than berserker")
	}
	if defensive.Weights["friendly_fire"] <= berserker.Weights["friendly_fire"] {
		t.Error("Defensive should have higher friendly fire avoidance than berserker")
	}
	
	// Default should be between the extremes for damage
	damageRank := []float32{
		berserker.Weights["damage_enemy"],
		defaultW.Weights["damage_enemy"],
		defensive.Weights["damage_enemy"],
	}
	if !isDescendingOrder(damageRank) {
		t.Error("Damage priority should be: Berserker > Default > Defensive")
	}
	
	// Default should be between the extremes for survival
	survivalRank := []float32{
		defensive.Weights["survival_threat"],
		defaultW.Weights["survival_threat"],
		berserker.Weights["survival_threat"],
	}
	if !isDescendingOrder(survivalRank) {
		t.Error("Survival priority should be: Defensive > Default > Berserker")
	}
}

func TestWeightsStructure(t *testing.T) {
	weights := NewDefaultWeights()
	
	if weights == nil {
		t.Fatal("Weights should not be nil")
	}
	
	if weights.Weights == nil {
		t.Fatal("Weights.Weights map should not be nil")
	}
	
	// Should contain exactly 9 metrics
	if len(weights.Weights) != 9 {
		t.Errorf("Expected 9 weight metrics, got %d", len(weights.Weights))
	}
}

// Helper functions
func assertWeight(t *testing.T, weights *Weights, metric string, expected float32, description string) {
	actual, exists := weights.Weights[metric]
	if !exists {
		t.Errorf("Missing required metric: %s", metric)
		return
	}
	if actual != expected {
		t.Errorf("%s: expected %s = %f, got %f", description, metric, expected, actual)
	}
}

func assertAllMetricsPresent(t *testing.T, weights *Weights) {
	requiredMetrics := []string{
		"damage_enemy",
		"friendly_fire", 
		"survival_threat",
		"kill_potential",
		"enemy_proximity",
		"threat_priority",
		"low_health_bonus",
		"tactical_value",
		"movement_efficiency",
	}
	
	for _, metric := range requiredMetrics {
		if _, exists := weights.Weights[metric]; !exists {
			t.Errorf("Missing required metric: %s", metric)
		}
	}
	
	if len(weights.Weights) != len(requiredMetrics) {
		t.Errorf("Expected exactly %d metrics, got %d", len(requiredMetrics), len(weights.Weights))
	}
}

func isDescendingOrder(values []float32) bool {
	for i := 1; i < len(values); i++ {
		if values[i-1] < values[i] {
			return false
		}
	}
	return true
}
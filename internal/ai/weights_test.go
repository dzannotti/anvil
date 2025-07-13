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
	if berserker.DamageEnemy <= defensive.DamageEnemy {
		t.Error("Berserker should have higher damage priority than defensive")
	}
	if berserker.DamageEnemy <= defaultW.DamageEnemy {
		t.Error("Berserker should have higher damage priority than default")
	}
	
	// Defensive should be most cautious
	if defensive.SurvivalThreat <= berserker.SurvivalThreat {
		t.Error("Defensive should have higher survival priority than berserker")
	}
	if defensive.FriendlyFire <= berserker.FriendlyFire {
		t.Error("Defensive should have higher friendly fire avoidance than berserker")
	}
	
	// Default should be between the extremes for damage
	damageRank := []float32{
		berserker.DamageEnemy,
		defaultW.DamageEnemy,
		defensive.DamageEnemy,
	}
	if !isDescendingOrder(damageRank) {
		t.Error("Damage priority should be: Berserker > Default > Defensive")
	}
	
	// Default should be between the extremes for survival
	survivalRank := []float32{
		defensive.SurvivalThreat,
		defaultW.SurvivalThreat,
		berserker.SurvivalThreat,
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
	
	// With struct-based weights, all 9 fields are always present by design
	// Test that the struct has reasonable values
	if weights.DamageEnemy <= 0 || weights.DamageEnemy > 10 {
		t.Errorf("DamageEnemy should be reasonable, got %f", weights.DamageEnemy)
	}
	if weights.SurvivalThreat <= 0 || weights.SurvivalThreat > 10 {
		t.Errorf("SurvivalThreat should be reasonable, got %f", weights.SurvivalThreat)
	}
}

// Helper functions
func assertWeight(t *testing.T, weights *Weights, metric string, expected float32, description string) {
	var actual float32
	
	switch metric {
	case "damage_enemy":
		actual = weights.DamageEnemy
	case "friendly_fire":
		actual = weights.FriendlyFire
	case "survival_threat":
		actual = weights.SurvivalThreat
	case "kill_potential":
		actual = weights.KillPotential
	case "enemy_proximity":
		actual = weights.EnemyProximity
	case "threat_priority":
		actual = weights.ThreatPriority
	case "low_health_bonus":
		actual = weights.LowHealthBonus
	case "tactical_value":
		actual = weights.TacticalValue
	case "movement_efficiency":
		actual = weights.MovementEfficiency
	default:
		t.Errorf("Unknown metric: %s", metric)
		return
	}
	
	if actual != expected {
		t.Errorf("%s: expected %s = %f, got %f", description, metric, expected, actual)
	}
}

func assertAllMetricsPresent(t *testing.T, weights *Weights) {
	// With struct-based weights, all metrics are always present by design
	// This function is kept for API compatibility but just checks for nil
	if weights == nil {
		t.Error("Weights should not be nil")
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
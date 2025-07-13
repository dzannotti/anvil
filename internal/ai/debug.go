package ai

import "fmt"

var aiDebugEnabled = true

func logAIDecision(message string) {
	if aiDebugEnabled {
		fmt.Println(message)
	}
}

func logActionBreakdown(evaluation *ActionTargetEvaluation) {
	if !aiDebugEnabled || evaluation == nil {
		return
	}

	fmt.Printf("  📊 Raw Metrics: ")
	for metric, value := range evaluation.RawMetrics {
		fmt.Printf("%s:%d ", metric, value)
	}
	fmt.Println()

	fmt.Printf("  ⚖️  Weighted Scores: ")
	for metric, value := range evaluation.WeightedScores {
		fmt.Printf("%s:%d ", metric, value)
	}
	fmt.Println()

	if len(evaluation.Movement) > 0 {
		fmt.Printf("  🚶 Movement: %v -> %v\n", evaluation.Position, evaluation.Target)
	}
}

func getArchetypeName(weights *Weights) string {
	damageWeight := weights.Weights["damage_enemy"]
	survivalWeight := weights.Weights["survival_threat"]

	switch {
	case damageWeight >= 2.0 && survivalWeight <= 0.5:
		return "BERSERKER"
	case survivalWeight >= 1.8 && damageWeight <= 1.2:
		return "DEFENSIVE"
	default:
		return "DEFAULT"
	}
}

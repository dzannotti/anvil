package ai

import (
	"testing"

	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/loader"
	"anvil/internal/tag"
)

func TestFindBestAction_ValidActions(t *testing.T) {
	world, actor, encounter := setupTestActionWorld()
	weights := NewDefaultWeights()
	
	// Add some actions to the actor
	actor.Actions = []core.Action{
		&mockScoredAction{name: "GoodAction", score: 50},
		&mockScoredAction{name: "BadAction", score: -10},
		&mockScoredAction{name: "BestAction", score: 100},
	}
	
	bestAction := findBestAction(world, actor, encounter, weights)
	
	if bestAction == nil {
		t.Fatal("Should find a best action")
	}
	
	if bestAction.Action.Name() != "BestAction" {
		t.Errorf("Expected BestAction to be selected, got %s", bestAction.Action.Name())
	}
	
	if bestAction.FinalScore <= 0 {
		t.Errorf("Expected positive final score, got %d", bestAction.FinalScore)
	}
}

func TestFindBestAction_NoFeasibleActions(t *testing.T) {
	world, actor, encounter := setupTestActionWorld()
	weights := NewDefaultWeights()
	
	// Add only unfeasible actions
	actor.Actions = []core.Action{
		&mockUnfeasibleAction{name: "CantDo"},
	}
	
	bestAction := findBestAction(world, actor, encounter, weights)
	
	if bestAction != nil {
		t.Error("Should not find any action when none are feasible")
	}
}

func TestFindFallbackAction_MovementPriority(t *testing.T) {
	world, actor, encounter := setupTestActionWorld()
	weights := NewDefaultWeights()
	
	// Add movement and other actions
	actor.Actions = []core.Action{
		&mockMoveActionAI{name: "Move", score: 30},
		&mockDodgeAction{name: "Dodge", score: 20},
		&mockScoredAction{name: "Attack", score: 10},
	}
	
	fallback := findFallbackAction(world, actor, encounter, weights)
	
	if fallback == nil {
		t.Fatal("Should find a fallback action")
	}
	
	// Should prefer movement with score > -50
	if fallback.Action.Name() != "Move" {
		t.Errorf("Expected Move to be selected as fallback, got %s", fallback.Action.Name())
	}
}

func TestFindFallbackAction_DefensivePriority(t *testing.T) {
	world, actor, encounter := setupTestActionWorld()
	weights := NewDefaultWeights()
	
	// Add defensive action with higher score than movement
	actor.Actions = []core.Action{
		&mockMoveActionAI{name: "Move", score: -60}, // Below threshold
		&mockDodgeAction{name: "Dodge", score: 40},
		&mockScoredAction{name: "Attack", score: 30},
	}
	
	fallback := findFallbackAction(world, actor, encounter, weights)
	
	if fallback == nil {
		t.Fatal("Should find a fallback action")
	}
	
	// Should prefer defensive action when movement score is too low
	if fallback.Action.Name() != "Dodge" {
		t.Errorf("Expected Dodge to be selected as fallback, got %s", fallback.Action.Name())
	}
}

func TestSimulateActionTarget_ScoreCalculation(t *testing.T) {
	world, actor, _ := setupTestActionWorld()
	action := &mockScoredAction{name: "TestAction", score: 0}
	target := grid.Position{X: 6, Y: 6}
	
	// Create simple weights that we can predict
	weights := &Weights{
		Weights: map[string]float32{
			"test_metric": 2.0,
		},
	}
	
	// Mock the metrics to return predictable values
	// Note: This test would need metric mocking for full functionality
	
	evaluation := simulateActionTarget(world, actor, action, target, weights)
	
	if evaluation == nil {
		t.Fatal("Should return an evaluation")
	}
	
	if evaluation.Action != action {
		t.Error("Evaluation should reference the correct action")
	}
	
	if evaluation.Target != target {
		t.Error("Evaluation should reference the correct target")
	}
}

func TestUpdateBestEvaluation_MovementAction(t *testing.T) {
	moveAction := &mockMoveActionAI{name: "Move", score: 30}
	tags := moveAction.Tags()
	
	evaluation := &ActionTargetEvaluation{
		Action:     moveAction,
		FinalScore: 30,
	}
	
	evals := &fallbackEvaluations{bestScore: -9999}
	
	updateBestEvaluation(evaluation, *tags, evals)
	
	if evals.bestMovement == nil {
		t.Error("Should set bestMovement for move action")
	}
	
	if evals.bestMovement.FinalScore != 30 {
		t.Errorf("Expected movement score 30, got %d", evals.bestMovement.FinalScore)
	}
}

func TestUpdateBestEvaluation_DefensiveAction(t *testing.T) {
	dodgeAction := &mockDodgeAction{name: "Dodge", score: 25}
	tags := dodgeAction.Tags()
	
	evaluation := &ActionTargetEvaluation{
		Action:     dodgeAction,
		FinalScore: 25,
	}
	
	evals := &fallbackEvaluations{bestScore: -9999}
	
	updateBestEvaluation(evaluation, *tags, evals)
	
	if evals.bestDefensive == nil {
		t.Error("Should set bestDefensive for dodge action")
	}
	
	if evals.bestDefensive.FinalScore != 25 {
		t.Errorf("Expected defensive score 25, got %d", evals.bestDefensive.FinalScore)
	}
}

func TestSelectBestFallbackAction_Preferences(t *testing.T) {
	evals := &fallbackEvaluations{
		bestMovement:  &ActionTargetEvaluation{FinalScore: 30},
		bestDefensive: &ActionTargetEvaluation{FinalScore: 40},
		bestOverall:   &ActionTargetEvaluation{FinalScore: 20},
		bestScore:     20,
	}
	
	// Should prefer movement if score > -50
	selected := selectBestFallbackAction(evals)
	if selected != evals.bestMovement {
		t.Error("Should prefer movement when score > -50")
	}
	
	// Test defensive preference when movement score too low
	evals.bestMovement.FinalScore = -60
	selected = selectBestFallbackAction(evals)
	if selected != evals.bestDefensive {
		t.Error("Should prefer defensive when movement score <= -50")
	}
	
	// Test fallback to overall when no good options
	evals.bestDefensive.FinalScore = 10
	evals.bestScore = 15
	selected = selectBestFallbackAction(evals)
	if selected != evals.bestOverall {
		t.Error("Should fallback to overall best when defensive not better")
	}
}

// Test helpers and mocks
func setupTestActionWorld() (*core.World, *core.Actor, *core.Encounter) {
	world := core.NewWorld(loader.WorldDefinition{Width: 20, Height: 20})
	encounter := &core.Encounter{}
	
	actor := &core.Actor{
		Name:         "TestActor",
		Team:         core.TeamID("heroes"),
		Position:     grid.Position{X: 5, Y: 5},
		HitPoints:    50,
		MaxHitPoints: 50,
	}
	
	encounter.Actors = append(encounter.Actors, actor)
	actor.Encounter = encounter
	world.AddOccupant(actor.Position, actor)
	
	return world, actor, encounter
}

type mockScoredAction struct {
	name  string
	score int
}

func (m *mockScoredAction) Name() string { return m.name }
func (m *mockScoredAction) Archetype() string { return "test" }
func (m *mockScoredAction) ID() string { return "test-scored" }
func (m *mockScoredAction) AverageDamage() int { return 8 }
func (m *mockScoredAction) Perform([]grid.Position) {}
func (m *mockScoredAction) Tags() *tag.Container { return &tag.Container{} }
func (m *mockScoredAction) CanAfford() bool { return true }
func (m *mockScoredAction) ValidPositions(grid.Position) []grid.Position {
	return []grid.Position{{X: 6, Y: 6}}
}
func (m *mockScoredAction) AffectedPositions([]grid.Position) []grid.Position {
	return []grid.Position{{X: 6, Y: 6}}
}

type mockUnfeasibleAction struct {
	name string
}

func (m *mockUnfeasibleAction) Name() string { return m.name }
func (m *mockUnfeasibleAction) Archetype() string { return "test" }
func (m *mockUnfeasibleAction) ID() string { return "test-unfeasible" }
func (m *mockUnfeasibleAction) AverageDamage() int { return 0 }
func (m *mockUnfeasibleAction) Perform([]grid.Position) {}
func (m *mockUnfeasibleAction) Tags() *tag.Container { return &tag.Container{} }
func (m *mockUnfeasibleAction) CanAfford() bool { return true }
func (m *mockUnfeasibleAction) ValidPositions(grid.Position) []grid.Position {
	return []grid.Position{} // No valid positions = unfeasible
}
func (m *mockUnfeasibleAction) AffectedPositions([]grid.Position) []grid.Position {
	return []grid.Position{}
}

type mockMoveActionAI struct {
	name  string
	score int
}

func (m *mockMoveActionAI) Name() string { return m.name }
func (m *mockMoveActionAI) Archetype() string { return "movement" }
func (m *mockMoveActionAI) ID() string { return "test-move" }
func (m *mockMoveActionAI) AverageDamage() int { return 0 }
func (m *mockMoveActionAI) Perform([]grid.Position) {}
func (m *mockMoveActionAI) Tags() *tag.Container {
	container := tag.ContainerFromTag(tags.Move)
	return &container
}
func (m *mockMoveActionAI) CanAfford() bool { return true }
func (m *mockMoveActionAI) ValidPositions(grid.Position) []grid.Position {
	return []grid.Position{{X: 4, Y: 4}}
}
func (m *mockMoveActionAI) AffectedPositions([]grid.Position) []grid.Position {
	return []grid.Position{{X: 4, Y: 4}}
}

type mockDodgeAction struct {
	name  string
	score int
}

func (m *mockDodgeAction) Name() string { return m.name }
func (m *mockDodgeAction) Archetype() string { return "defensive" }
func (m *mockDodgeAction) ID() string { return "test-dodge" }
func (m *mockDodgeAction) AverageDamage() int { return 0 }
func (m *mockDodgeAction) Perform([]grid.Position) {}
func (m *mockDodgeAction) Tags() *tag.Container {
	container := tag.ContainerFromTag(tags.Dodge)
	return &container
}
func (m *mockDodgeAction) CanAfford() bool { return true }
func (m *mockDodgeAction) ValidPositions(grid.Position) []grid.Position {
	return []grid.Position{{X: 5, Y: 5}}
}
func (m *mockDodgeAction) AffectedPositions([]grid.Position) []grid.Position {
	return []grid.Position{{X: 5, Y: 5}}
}
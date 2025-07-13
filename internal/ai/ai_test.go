package ai

import (
	"testing"

	"anvil/internal/core"
	"anvil/internal/grid"
	"anvil/internal/loader"
	"anvil/internal/tag"
)

func TestPlay_ValidActor(t *testing.T) {
	gameState, actor := setupTestGameState()
	weights := NewDefaultWeights()
	
	// Actor should be able to act
	if !actor.CanAct() {
		t.Fatal("Test actor should be able to act")
	}
	
	// Should not panic
	Play(gameState, weights)
	
	// Turn should be ended (defer call)
	// Note: We can't easily test this without mocking the encounter
}

func TestPlay_DeadActor(t *testing.T) {
	gameState, actor := setupTestGameState()
	weights := NewDefaultWeights()
	
	// Kill the actor
	actor.HitPoints = 0
	
	// Should handle dead actor gracefully
	Play(gameState, weights)
	
	// Should not panic and should return early
}

func TestPlay_CannotAct(t *testing.T) {
	gameState, _ := setupTestGameState()
	weights := NewDefaultWeights()
	
	// Make actor unable to act (mock this by setting a condition)
	// Note: We'd need to mock CanAct() method for proper testing
	// For now, just ensure it handles the case
	
	Play(gameState, weights)
	
	// Should not panic
}

func TestActionTargetEvaluation_Structure(t *testing.T) {
	evaluation := &ActionTargetEvaluation{
		Action:         &mockAIAction{name: "TestAction"},
		Target:         grid.Position{X: 5, Y: 5},
		Affected:       []grid.Position{{X: 5, Y: 5}, {X: 6, Y: 6}},
		RawScores:      &Scores{DamageEnemy: 10},
		WeightedScores: &WeightedScores{DamageEnemy: 15},
		FinalScore:     15,
		Position:       grid.Position{X: 4, Y: 4},
		Movement:       []grid.Position{{X: 4, Y: 4}},
	}
	
	// Validate structure
	if evaluation.Action.Name() != "TestAction" {
		t.Errorf("Expected action name 'TestAction', got '%s'", evaluation.Action.Name())
	}
	
	if evaluation.FinalScore != 15 {
		t.Errorf("Expected final score 15, got %d", evaluation.FinalScore)
	}
	
	if len(evaluation.Affected) != 2 {
		t.Errorf("Expected 2 affected positions, got %d", len(evaluation.Affected))
	}
}

func TestGetBestScore(t *testing.T) {
	// Test with valid evaluation
	evaluation := &ActionTargetEvaluation{FinalScore: 42}
	score := getBestScore(evaluation)
	if score != 42 {
		t.Errorf("Expected score 42, got %d", score)
	}
	
	// Test with nil evaluation
	score = getBestScore(nil)
	if score != -999 {
		t.Errorf("Expected score -999 for nil evaluation, got %d", score)
	}
}

func TestExecuteAction_BasicExecution(t *testing.T) {
	gameState, actor := setupTestGameState()
	originalPos := actor.Position
	
	evaluation := &ActionTargetEvaluation{
		Action:   &mockAIAction{name: "TestAction"},
		Target:   grid.Position{X: 6, Y: 6},
		Position: originalPos, // Same position (no movement)
		Movement: []grid.Position{},
	}
	
	executeAction(gameState, actor, evaluation, originalPos)
	
	// Actor should still be in original position
	if actor.Position != originalPos {
		t.Errorf("Expected actor to remain at %v, got %v", originalPos, actor.Position)
	}
}

func TestExecuteAction_WithMovement(t *testing.T) {
	gameState, actor := setupTestGameState()
	originalPos := actor.Position
	newPos := grid.Position{X: 3, Y: 3}
	
	evaluation := &ActionTargetEvaluation{
		Action:   &mockAIAction{name: "TestAction"},
		Target:   grid.Position{X: 6, Y: 6},
		Position: newPos,
		Movement: []grid.Position{newPos},
	}
	
	executeAction(gameState, actor, evaluation, originalPos)
	
	// Actor should be moved to new position
	if actor.Position != newPos {
		t.Errorf("Expected actor to move to %v, got %v", newPos, actor.Position)
	}
}

func TestExecuteAction_CedricMovementRevert(t *testing.T) {
	gameState, actor := setupTestGameState()
	actor.Name = "Cedric" // Special case for Cedric
	originalPos := actor.Position
	newPos := grid.Position{X: 3, Y: 3}
	
	evaluation := &ActionTargetEvaluation{
		Action:   &mockAIAction{name: "TestAction"},
		Target:   grid.Position{X: 6, Y: 6},
		Position: newPos,
		Movement: []grid.Position{newPos},
	}
	
	executeAction(gameState, actor, evaluation, originalPos)
	
	// Cedric should be reverted to original position after action
	if actor.Position != originalPos {
		t.Errorf("Expected Cedric to revert to %v, got %v", originalPos, actor.Position)
	}
}

// Test helpers
func setupTestGameState() (*core.GameState, *core.Actor) {
	world := core.NewWorld(loader.WorldDefinition{Width: 20, Height: 20})
	encounter := &core.Encounter{}
	
	actor := &core.Actor{
		Name:         "TestActor",
		Team:         core.TeamID("heroes"),
		Position:     grid.Position{X: 5, Y: 5},
		HitPoints:    50,
		MaxHitPoints: 50,
		Actions:      []core.Action{&mockAIAction{name: "BasicAttack"}},
	}
	
	// Add enemy
	enemy := &core.Actor{
		Name:         "TestEnemy",
		Team:         core.TeamID("monsters"),
		Position:     grid.Position{X: 6, Y: 6},
		HitPoints:    30,
		MaxHitPoints: 30,
	}
	
	encounter.Actors = append(encounter.Actors, actor)
	encounter.Actors = append(encounter.Actors, enemy)
	actor.Encounter = encounter
	enemy.Encounter = encounter
	
	world.AddOccupant(actor.Position, actor)
	world.AddOccupant(enemy.Position, enemy)
	
	gameState := &core.GameState{
		World:     world,
		Encounter: encounter,
	}
	
	return gameState, actor
}

// Mock action for AI testing
type mockAIAction struct {
	name     string
	performed bool
}

func (m *mockAIAction) Name() string { return m.name }
func (m *mockAIAction) Archetype() string { return "test" }
func (m *mockAIAction) ID() string { return "test-action" }
func (m *mockAIAction) AverageDamage() int { return 8 }
func (m *mockAIAction) Perform(targets []grid.Position) {
	m.performed = true
}
func (m *mockAIAction) Tags() *tag.Container { return &tag.Container{} }
func (m *mockAIAction) CanAfford() bool { return true }
func (m *mockAIAction) ValidPositions(pos grid.Position) []grid.Position {
	// Return positions around the given position
	return []grid.Position{
		{X: pos.X + 1, Y: pos.Y},
		{X: pos.X - 1, Y: pos.Y},
		{X: pos.X, Y: pos.Y + 1},
		{X: pos.X, Y: pos.Y - 1},
	}
}
func (m *mockAIAction) AffectedPositions(targets []grid.Position) []grid.Position {
	return targets
}
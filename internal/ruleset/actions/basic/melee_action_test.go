package basic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/loader"
	"anvil/internal/tag"
)

func createTestSetup() (*core.World, *core.Encounter, *eventbus.Dispatcher) {
	dispatcher := &eventbus.Dispatcher{}
	world := core.NewWorld(loader.WorldDefinition{Width: 10, Height: 10})
	encounter := &core.Encounter{
		Dispatcher: dispatcher,
		World:      world,
		Actors:     []*core.Actor{},
	}
	return world, encounter, dispatcher
}

func createTestActor(world *core.World, encounter *core.Encounter, dispatcher *eventbus.Dispatcher, pos grid.Position, name string, team core.TeamID) *core.Actor {
	attributes := stats.Attributes{
		Strength:     16,
		Dexterity:    14,
		Constitution: 15,
		Intelligence: 8,
		Wisdom:       12,
		Charisma:     10,
	}
	proficiencies := stats.Proficiencies{Bonus: 2}
	resources := core.Resources{Max: map[tag.Tag]int{
		tags.ResourceAction:    1,
		tags.ResourceWalkSpeed: 5,
	}}

	actor := &core.Actor{
		Dispatcher:    dispatcher,
		Encounter:     encounter,
		Position:      pos,
		World:         world,
		Name:          name,
		Team:          team,
		HitPoints:     20,
		MaxHitPoints:  20,
		Attributes:    attributes,
		Proficiencies: proficiencies,
		Resources:     resources,
	}
	actor.Resources.LongRest()
	world.AddOccupant(pos, actor)
	encounter.Actors = append(encounter.Actors, actor)
	return actor
}

// MockWeapon implements the weapon interface for testing
type MockWeapon struct {
	name       string
	damage     *expression.Expression
	damageType tag.Tag
}

func (m MockWeapon) Damage() *expression.Expression {
	return m.damage
}

func (m MockWeapon) Name() string {
	return m.name
}

func (m MockWeapon) Tags() *tag.Container {
	container := tag.NewContainer(m.damageType)
	return &container
}

func createTestWeapon(name string, _ string, damageType tag.Tag, _ int) *MockWeapon {
	expr := expression.FromDamageDice(1, 8, name, tag.NewContainer(damageType))
	return &MockWeapon{
		name:       name,
		damage:     &expr,
		damageType: damageType,
	}
}

func TestMeleeAction_Creation(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	actor := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(actor, "Attack with Sword", weapon, 1, tag.NewContainer(tags.Melee), cost)

	assert.NotNil(t, action)
	assert.Equal(t, "Attack with Sword", action.Name())
	assert.Equal(t, "attack", action.Archetype())
	assert.Equal(t, 1, action.Reach())
	assert.True(t, action.Tags().MatchTag(tags.Attack))
	assert.True(t, action.Tags().MatchTag(tags.Melee))
	assert.Equal(t, cost, action.Cost())
	assert.Equal(t, actor, action.Owner())
}

func TestMeleeAction_ValidPositions_ReachConstraint(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	target := createTestActor(world, encounter, dispatcher, grid.Position{X: 6, Y: 5}, "Orc", core.TeamEnemies)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)
	valid := action.ValidPositions(attacker.Position)

	// Should be able to attack adjacent enemy
	assert.Contains(t, valid, target.Position)

	// Move target out of reach
	world.RemoveOccupant(target.Position, target)
	target.Position = grid.Position{X: 7, Y: 5}
	world.AddOccupant(target.Position, target)

	valid = action.ValidPositions(attacker.Position)
	// Should NOT be able to attack enemy 2 squares away with reach 1
	assert.NotContains(t, valid, target.Position)
}

func TestMeleeAction_ValidPositions_ExtendedReach(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	target := createTestActor(world, encounter, dispatcher, grid.Position{X: 7, Y: 5}, "Orc", core.TeamEnemies)
	weapon := createTestWeapon("Spear", "1d8", tags.Piercing, 2)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 2, tag.NewContainer(tags.Melee), cost)
	valid := action.ValidPositions(attacker.Position)

	// Should be able to attack enemy 2 squares away with reach 2
	assert.Contains(t, valid, target.Position)
}

func TestMeleeAction_ValidPositions_OnlyTargetsEnemies(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	ally := createTestActor(world, encounter, dispatcher, grid.Position{X: 6, Y: 5}, "Wizard", core.TeamPlayers)
	enemy := createTestActor(world, encounter, dispatcher, grid.Position{X: 4, Y: 5}, "Orc", core.TeamEnemies)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)
	valid := action.ValidPositions(attacker.Position)

	// Should be able to attack enemy
	assert.Contains(t, valid, enemy.Position)
	// Should NOT be able to attack ally
	assert.NotContains(t, valid, ally.Position)
	// Should NOT be able to attack empty squares
	assert.NotContains(t, valid, grid.Position{X: 5, Y: 6})
}

func TestMeleeAction_ValidPositions_SkipsDeadEnemies(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	enemy := createTestActor(world, encounter, dispatcher, grid.Position{X: 6, Y: 5}, "Orc", core.TeamEnemies)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)

	// Enemy is alive - should be valid target
	valid := action.ValidPositions(attacker.Position)
	assert.Contains(t, valid, enemy.Position)

	// Kill the enemy
	enemy.AddCondition(tags.Dead, &core.Effect{Name: "Test"})

	// Dead enemy should not be valid target
	valid = action.ValidPositions(attacker.Position)
	assert.NotContains(t, valid, enemy.Position)
}

func TestMeleeAction_ValidPositions_RequiresAffordableCost(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	enemy := createTestActor(world, encounter, dispatcher, grid.Position{X: 6, Y: 5}, "Orc", core.TeamEnemies)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)

	// With action available - should be valid
	valid := action.ValidPositions(attacker.Position)
	assert.Contains(t, valid, enemy.Position)

	// Consume the action
	attacker.ConsumeResource(tags.ResourceAction, 1)

	// Without action available - should be empty
	valid = action.ValidPositions(attacker.Position)
	assert.Empty(t, valid)
}

func TestMeleeAction_AffectedPositions(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)
	target := grid.Position{X: 6, Y: 5}
	affected := action.AffectedPositions([]grid.Position{target})

	// Melee attacks only affect the target square
	assert.Len(t, affected, 1)
	assert.Equal(t, target, affected[0])
}

func TestMeleeAction_TagCombination(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	actionTags := tag.NewContainer(tags.Melee, tags.MartialWeapon)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, actionTags, cost)
	combinedTags := action.Tags()

	// Should have action tags
	assert.True(t, combinedTags.MatchTag(tags.Melee))
	assert.True(t, combinedTags.MatchTag(tags.MartialWeapon))
	// Should have automatically added Attack tag
	assert.True(t, combinedTags.MatchTag(tags.Attack))
	// Should have weapon damage tags
	assert.True(t, combinedTags.MatchTag(tags.Slashing))
}

func TestMeleeAction_AverageDamage(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)
	avgDamage := action.AverageDamage()

	// 1d8 has average of 4.5, should be rounded/calculated appropriately
	assert.Greater(t, avgDamage, 0)
	assert.LessOrEqual(t, avgDamage, 8) // Max damage of 1d8
}

func TestMeleeAction_ConfigurableCost(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)

	// Test standard action cost
	standardCost := map[tag.Tag]int{tags.ResourceAction: 1}
	standardAction := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), standardCost)
	assert.Equal(t, standardCost, standardAction.Cost())

	// Test AOO cost (reaction instead of action)
	aooCost := map[tag.Tag]int{tags.ResourceReaction: 1}
	aooAction := NewMeleeAction(attacker, "Attack of Opportunity", weapon, 1, tag.NewContainer(tags.Melee), aooCost)
	assert.Equal(t, aooCost, aooAction.Cost())

	// Test bonus action cost
	bonusCost := map[tag.Tag]int{tags.ResourceBonusAction: 1}
	bonusAction := NewMeleeAction(attacker, "Offhand Attack", weapon, 1, tag.NewContainer(tags.Melee), bonusCost)
	assert.Equal(t, bonusCost, bonusAction.Cost())
}

func TestMeleeAction_Perform_ExecutesAttackSequence(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	target := createTestActor(world, encounter, dispatcher, grid.Position{X: 6, Y: 5}, "Orc", core.TeamEnemies)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)

	// Should have action available before attack
	assert.True(t, action.CanAfford())
	assert.Equal(t, 1, attacker.Resources.Remaining(tags.ResourceAction))

	// Perform the attack
	action.Perform([]grid.Position{target.Position})

	// Should consume the action
	assert.False(t, action.CanAfford())
	assert.Equal(t, 0, attacker.Resources.Remaining(tags.ResourceAction))
}

func TestMeleeAction_Perform_ConsumesCorrectResource(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	target := createTestActor(world, encounter, dispatcher, grid.Position{X: 6, Y: 5}, "Orc", core.TeamEnemies)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)

	// Add reaction resource for AOO test
	attacker.Resources.Max[tags.ResourceReaction] = 1
	attacker.Resources.LongRest()

	// Test action cost
	actionCost := map[tag.Tag]int{tags.ResourceAction: 1}
	actionAttack := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), actionCost)
	actionAttack.Perform([]grid.Position{target.Position})
	assert.Equal(t, 0, attacker.Resources.Remaining(tags.ResourceAction))
	assert.Equal(t, 1, attacker.Resources.Remaining(tags.ResourceReaction)) // Should not consume reaction

	// Reset and test reaction cost (AOO)
	attacker.Resources.LongRest()
	reactionCost := map[tag.Tag]int{tags.ResourceReaction: 1}
	aooAttack := NewMeleeAction(attacker, "AOO", weapon, 1, tag.NewContainer(tags.Melee), reactionCost)
	aooAttack.Perform([]grid.Position{target.Position})
	assert.Equal(t, 1, attacker.Resources.Remaining(tags.ResourceAction)) // Should not consume action
	assert.Equal(t, 0, attacker.Resources.Remaining(tags.ResourceReaction))
}

func TestMeleeAction_ValidPositions_CircularReach(t *testing.T) {
	world, encounter, dispatcher := createTestSetup()
	attacker := createTestActor(world, encounter, dispatcher, grid.Position{X: 5, Y: 5}, "Fighter", core.TeamPlayers)
	weapon := createTestWeapon("Sword", "1d8", tags.Slashing, 1)
	cost := map[tag.Tag]int{tags.ResourceAction: 1}

	action := NewMeleeAction(attacker, "Attack", weapon, 1, tag.NewContainer(tags.Melee), cost)

	// Place enemies at orthogonally adjacent positions (reach 1)
	orthogonalPositions := []grid.Position{
		{X: 4, Y: 5}, // West
		{X: 6, Y: 5}, // East
		{X: 5, Y: 4}, // North
		{X: 5, Y: 6}, // South
	}

	// Place enemies at diagonally adjacent positions (outside reach 1)
	diagonalPositions := []grid.Position{
		{X: 4, Y: 4}, // Northwest
		{X: 6, Y: 4}, // Northeast
		{X: 4, Y: 6}, // Southwest
		{X: 6, Y: 6}, // Southeast
	}

	for i, pos := range orthogonalPositions {
		createTestActor(world, encounter, dispatcher, pos, fmt.Sprintf("OrthEnemy%d", i), core.TeamEnemies)
	}
	for i, pos := range diagonalPositions {
		createTestActor(world, encounter, dispatcher, pos, fmt.Sprintf("DiagEnemy%d", i), core.TeamEnemies)
	}

	valid := action.ValidPositions(attacker.Position)

	// Should be able to attack orthogonally adjacent enemies (reach 1)
	for _, pos := range orthogonalPositions {
		assert.Contains(t, valid, pos, "Should be able to attack orthogonally adjacent enemy at %v", pos)
	}

	// Should NOT be able to attack diagonally adjacent enemies with reach 1 (D&D 5e rules)
	for _, pos := range diagonalPositions {
		assert.NotContains(t, valid, pos, "Should NOT be able to attack diagonally adjacent enemy at %v with reach 1", pos)
	}

	// Should NOT be able to attack enemies 2 squares away
	farPositions := []grid.Position{
		{X: 3, Y: 5}, // 2 West
		{X: 7, Y: 5}, // 2 East
		{X: 5, Y: 3}, // 2 North
		{X: 5, Y: 7}, // 2 South
	}
	for _, pos := range farPositions {
		assert.NotContains(t, valid, pos, "Should NOT be able to attack enemy at %v with reach 1", pos)
	}
}

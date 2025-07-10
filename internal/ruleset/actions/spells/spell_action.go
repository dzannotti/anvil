package spells

import (
	"anvil/internal/core"
	"anvil/internal/core/shapes"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

// SavingThrow holds the details for a spell's saving throw requirement.
type SavingThrow struct {
	Attribute    tag.Tag // e.g., tags.Dexterity
	DC           int     // The calculated Difficulty Class
	EffectOnSave tag.Tag // e.g., tags.HalfDamage, tags.NegateEffect
}

// SpellAction is a generic, configurable action that can represent various types of spells.
type SpellAction struct {
	owner *core.Actor
	id    string
	name  string
	tags  tag.Container
	cost  map[tag.Tag]int

	// --- Targeting & Delivery ---
	spellRange   int
	aoeShape     func(grid.Position) []grid.Position
	isAttackRoll bool

	// --- Payload (What the spell does) ---
	damageExpression *expression.Expression // The dice expression for damage. Nil if no damage.
	damageTags       tag.Container          // Tags for the damage (e.g., "damage.fire").
	effectsToApply   []*core.Effect         // Effects to apply on hit/fail.

	// --- Defense (How targets resist) ---
	savingThrow *SavingThrow
}

// NewSpellAction creates and configures a new SpellAction.
func NewSpellAction(
	owner *core.Actor,
	name string,
	cost map[tag.Tag]int,
	spellRange int,
	aoeShape func(grid.Position) []grid.Position,
	isAttackRoll bool,
	damageExpression *expression.Expression,
	damageTags tag.Container,
	effectsToApply []*core.Effect,
	savingThrow *SavingThrow,
) *SpellAction {
	return &SpellAction{
		owner:            owner,
		id:               uuid.New().String(),
		name:             name,
		cost:             cost,
		spellRange:       spellRange,
		aoeShape:         aoeShape,
		isAttackRoll:     isAttackRoll,
		damageExpression: damageExpression,
		damageTags:       damageTags,
		effectsToApply:   effectsToApply,
		savingThrow:      savingThrow,
	}
}

//
// core.DamageSource interface implementation
//

// Name returns the spell's name.
func (s *SpellAction) Name() string {
	return s.name
}

// Damage returns a clone of the spell's damage expression.
func (s *SpellAction) Damage() *expression.Expression {
	if s.damageExpression == nil {
		// Return an empty expression to avoid nil pointer issues downstream.
		return &expression.Expression{}
	}
	cloned := s.damageExpression.Clone()
	return &cloned
}

// Tags returns the damage tags for the spell.
func (s *SpellAction) Tags() *tag.Container {
	return &s.damageTags
}

//
// core.Action interface implementation
//

func (s *SpellAction) Owner() *core.Actor {
	return s.owner
}

func (s *SpellAction) Archetype() string {
	return "spell"
}

func (s *SpellAction) ID() string {
	return s.id
}

func (s *SpellAction) Cost() map[tag.Tag]int {
	return s.cost
}

func (s *SpellAction) CanAfford() bool {
	return s.owner.Resources.CanAfford(s.cost)
}

func (s *SpellAction) Commit() {
	if !s.CanAfford() {
		panic("Attempt to commit action without affording cost")
	}
	for t, amount := range s.cost {
		s.owner.ConsumeResource(t, amount)
	}
}

func (s *SpellAction) ValidPositions(from grid.Position) []grid.Position {
	if !s.CanAfford() {
		return []grid.Position{}
	}

	// Generate a circle representing the maximum range.
	positionsInRange := shapes.Circle(from, s.spellRange)
	validPositions := make([]grid.Position, 0)

	if s.isAttackRoll {
		// For attack rolls, valid positions are occupied by enemies.
		for _, pos := range positionsInRange {
			if !s.owner.World.IsValidPosition(pos) {
				continue
			}
			actor := s.owner.World.ActorAt(pos)
			if actor != nil && s.owner.Team != actor.Team && !actor.IsDead() {
				validPositions = append(validPositions, pos)
			}
		}
	} else {
		// For AoE or targeted effects, any valid grid cell in range is a valid target point.
		for _, pos := range positionsInRange {
			if s.owner.World.IsValidPosition(pos) {
				validPositions = append(validPositions, pos)
			}
		}
	}
	return validPositions
}

func (s *SpellAction) AffectedPositions(targetPos []grid.Position) []grid.Position {
	if len(targetPos) == 0 {
		return []grid.Position{}
	}
	pointOfOrigin := targetPos[0]

	if s.aoeShape != nil {
		// If there's an AoE shape, calculate all positions affected by it.
		return s.aoeShape(pointOfOrigin)
	}

	// If no AoE, the only affected position is the target itself.
	return []grid.Position{pointOfOrigin}
}

func (s *SpellAction) Perform(pos []grid.Position) {
	s.owner.Dispatcher.Begin(core.UseActionEvent{Action: s, Source: s.owner, Target: pos})
	defer s.owner.Dispatcher.End()

	s.Commit()

	affectedActors := []*core.Actor{}
	for _, p := range s.AffectedPositions(pos) {
		actor := s.owner.World.ActorAt(p)
		if actor != nil {
			affectedActors = append(affectedActors, actor)
		}
	}
	s.owner.Dispatcher.Emit(core.TargetEvent{Target: affectedActors})

	if s.isAttackRoll {
		// Path A: Spell Attack
		target := s.owner.World.ActorAt(pos[0])
		if target == nil {
			return // Target might have moved or died
		}
		attackResult := s.owner.AttackRoll(target, *s.Tags())
		if attackResult.Success {
			s.applyPayload(target, attackResult.Critical, true, false)
		}
	} else {
		// Path B: Automatic Hit / Saving Throw
		for _, target := range affectedActors {
			if target.IsDead() {
				continue
			}

			saveSuccessful := false
			if s.savingThrow != nil {
				saveExpr := expression.FromD20("Saving Throw")
				target.Evaluate(&core.PreSavingThrow{
					Source:          target,
					Expression:      &saveExpr,
					DifficultyClass: s.savingThrow.DC,
					Attribute:       s.savingThrow.Attribute,
				})
				saveResult := saveExpr.Evaluate()
				target.Evaluate(&core.PostSavingThrow{
					Source:          target,
					Result:          saveResult,
					DifficultyClass: s.savingThrow.DC,
					Attribute:       s.savingThrow.Attribute,
				})
				saveSuccessful = saveResult.Value >= s.savingThrow.DC
			}

			s.applyPayload(target, false, true, saveSuccessful)
		}
	}
}

// applyPayload is a helper to apply damage and effects based on the outcome of an attack or save.
func (s *SpellAction) applyPayload(target *core.Actor, isCritical bool, isHit bool, saveSuccessful bool) {
	if !isHit {
		return
	}

	// Apply Damage
	if s.damageExpression != nil {
		damageToDeal := s.owner.DamageRoll(s, isCritical)

		// Handle half damage on save
		if saveSuccessful && s.savingThrow != nil && s.savingThrow.EffectOnSave.Match(tags.EffectSaveHalfDamage) {
			// Halve the total damage value
			halvedValue := damageToDeal.Value / 2
			damageToDeal.ReplaceWith(halvedValue, "Successful Save")
		}

		target.TakeDamage(*damageToDeal)
	}

	// Apply Effects
	if len(s.effectsToApply) > 0 {
		if s.savingThrow != nil && s.savingThrow.EffectOnSave.Match(tags.EffectSaveNegate) && saveSuccessful {
			// Do not apply effects if the save negates them
		} else {
			for _, effect := range s.effectsToApply {
				target.AddEffect(effect.Clone()) // Clone effect to give a unique instance
			}
		}
	}
}

func (s *SpellAction) AverageDamage() int {
	if s.damageExpression == nil {
		return 0
	}
	return s.damageExpression.ExpectedValue()
}

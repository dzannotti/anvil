package spells

import (
	"testing"

	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/expression"
	"anvil/internal/tag"

	"github.com/stretchr/testify/assert"
)

// TestSpellAction_Creation ensures that a SpellAction can be created with all its components.
func TestSpellAction_Creation(t *testing.T) {
	dispatcher := &eventbus.Dispatcher{}
	world := core.NewWorld(10, 10)
	owner := &core.Actor{
		Dispatcher: dispatcher,
		World:      world,
		Attributes: stats.Attributes{Intelligence: 16},
	}

	cost := map[tag.Tag]int{tags.Action: 1}
	damageExpr := expression.FromDamageDice(8, 6, "Fireball", tag.NewContainer(tags.Fire))
	savingThrow := &SavingThrow{
		Attribute:    tags.Dexterity,
		DC:           15,
		EffectOnSave: tags.EffectSaveHalfDamage,
	}

	action := NewSpellAction(
		owner,
		"Fireball",
		cost,
		30,       // range
		nil,      // aoeShape
		false,    // isAttackRoll
		&damageExpr,
		tag.NewContainer(tags.Fire),
		nil,      // effectsToApply
		savingThrow,
	)

	assert.NotNil(t, action)
	assert.Equal(t, "Fireball", action.Name())
	assert.Equal(t, "spell", action.Archetype())
	assert.Equal(t, owner, action.Owner())
	assert.Equal(t, cost, action.Cost())
	assert.Equal(t, 15, action.savingThrow.DC)
	assert.True(t, action.Tags().MatchTag(tags.Fire))
	assert.NotNil(t, action.Damage())
}

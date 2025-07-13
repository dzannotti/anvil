package expression

import "anvil/internal/tag"

var Primary = tag.FromString("primary")

var (
	DamageAcid        = tag.FromString("damage.acid")
	DamageBludgeoning = tag.FromString("damage.bludgeoning")
	DamageCold        = tag.FromString("damage.cold")
	DamageFire        = tag.FromString("damage.fire")
	DamageForce       = tag.FromString("damage.force")
	DamageLightning   = tag.FromString("damage.lightning")
	DamageNecrotic    = tag.FromString("damage.necrotic")
	DamagePiercing    = tag.FromString("damage.piercing")
	DamagePoison      = tag.FromString("damage.poison")
	DamagePsychic     = tag.FromString("damage.psychic")
	DamageRadiant     = tag.FromString("damage.radiant")
	DamageSlashing    = tag.FromString("damage.slashing")
	DamageThunder     = tag.FromString("damage.thunder")
)

type Expression struct {
	Value      int
	Components []Component
	Rng        Roller
}

func (e *Expression) Evaluate() *Expression {
	e.Value = 0
	ctx := &Context{Rng: e.Rng}
	for _, component := range e.Components {
		e.Value += component.Evaluate(ctx)
	}
	return e
}

func (e *Expression) Clone() *Expression {
	clone := *e
	clone.Components = make([]Component, len(e.Components))
	for i, component := range e.Components {
		clone.Components[i] = component.Clone()
	}
	return &clone
}

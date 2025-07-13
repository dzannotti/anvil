package basic

import (
	"fmt"
	"regexp"
	"strconv"

	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/loader"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

type Weapon struct {
	archetype string
	id        string
	name      string
	damage    expression.Expression
	tags      tag.Container
	reach     int
}

func NewWeapon(archetype, id, name string, damage expression.Expression, weaponTags tag.Container, reach int) *Weapon {
	return &Weapon{
		archetype: archetype,
		id:        id,
		name:      name,
		damage:    damage,
		tags:      weaponTags,
		reach:     reach,
	}
}

func NewWeaponFromDefinition(def loader.WeaponDefinition) *Weapon {
	damageExpr := expression.Expression{Rng: expression.NewRngRoller()}
	for _, dmg := range def.Damage {
		if err := parseDamageFormula(dmg.Formula, def.Name, dmg.Kind, &damageExpr); err != nil {
			panic(fmt.Sprintf("invalid damage formula '%s' for weapon '%s': %v", dmg.Formula, def.Archetype, err))
		}
	}

	weaponTags := make([]tag.Tag, len(def.Tags))
	for i, tagStr := range def.Tags {
		weaponTags[i] = tag.FromString(tagStr)
	}

	return &Weapon{
		archetype: def.Archetype,
		id:        uuid.New().String(),
		name:      def.Name,
		damage:    damageExpr,
		tags:      tag.NewContainer(weaponTags...),
		reach:     def.Reach,
	}
}

func parseDamageFormula(formula, weaponName, kind string, expr *expression.Expression) error {
	diceRe := regexp.MustCompile(`^(\d+)d(\d+)$`)
	if matches := diceRe.FindStringSubmatch(formula); len(matches) == 3 {
		times, err := strconv.Atoi(matches[1])
		if err != nil {
			return fmt.Errorf("invalid number of dice: %s", matches[1])
		}

		sides, err := strconv.Atoi(matches[2])
		if err != nil {
			return fmt.Errorf("invalid number of sides: %s", matches[2])
		}

		expr.AddDamageDice(times, sides, tag.NewContainer(tag.FromString(kind)), weaponName)
		return nil
	}

	if constant, err := strconv.Atoi(formula); err == nil {
		expr.AddDamageConstant(constant, tag.NewContainer(tag.FromString(kind)), weaponName)
		return nil
	}

	return fmt.Errorf("invalid damage formula format: %s (expected format like '1d4', '2d6', or '5')", formula)
}

func (w Weapon) Archetype() string {
	return w.archetype
}

func (w Weapon) ID() string {
	return w.id
}

func (w Weapon) Name() string {
	return w.name
}

func (w Weapon) Tags() *tag.Container {
	tags := w.tags.Clone()
	return &tags
}

func (w Weapon) OnEquip(a *core.Actor) {
	cost := map[tag.Tag]int{tags.ResourceAction: 1}
	a.AddAction(NewMeleeAction(a, fmt.Sprintf("Attack with %s", w.name), &w, w.reach, w.tags, cost))
}

func (w Weapon) Damage() *expression.Expression {
	dmg := w.damage.Clone()
	return dmg
}

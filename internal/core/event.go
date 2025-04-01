package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/serialize"
	"anvil/internal/expression"
	"anvil/internal/tag"
	"encoding/json"
)

func NewEncounterEvent(encounter definition.Encounter) (string, []byte) {
	teams := make([]map[string]any, 0, len(encounter.Teams()))
	for i := range encounter.Teams() {
		teams = append(teams, snapshotTeam(encounter.Teams()[i]))
	}
	data := serialize.ToJSON(map[string]any{
		"Teams": teams,
		"World": snapshotWorld(encounter.World()),
	})
	return "encounter", data
}

func NewRoundEvent(round int, c []definition.Creature) (string, []byte) {
	creatures := make([]map[string]any, 0, len(c))
	for i := range c {
		creatures = append(creatures, snapshotCreature(c[i]))
	}
	data := serialize.ToJSON(map[string]any{
		"Round":     round,
		"Creatures": creatures,
	})

	return "round", data
}

func NewTurnEvent(turn int, src definition.Creature) (string, []byte) {
	data := serialize.ToJSON(map[string]any{
		"Turn":     turn,
		"Creature": snapshotCreature(src),
	})

	return "turn", data
}

func NewTakeDamageEvent(src definition.Creature, damage int) (string, []byte) {
	data := serialize.ToJSON(map[string]any{
		"Creature": snapshotCreature(src),
		"Damage":   damage,
	})

	return "take_damage", data
}

func NewAttackRollEvent(src definition.Creature, dst definition.Creature) (string, []byte) {
	data := serialize.ToJSON(map[string]any{
		"creature": snapshotCreature(src),
		"target":   snapshotCreature(dst),
	})

	return "attack_roll", data
}

func NewExpressionResultEvent(expr expression.Expression) (string, []byte) {
	data := serialize.ToJSON(map[string]any{
		"Expression": expr,
	})

	return "expression_result", data
}

func NewAttributeCalculationEvent(tag tag.Tag, expr expression.Expression) (string, []byte) {
	data, _ := json.Marshal(map[string]any{
		"Expression": expr,
	})
	return "attribute_calculation", data
}

func NewCheckResultEvent(value int, against int, critical bool, success bool) (string, []byte) {
	data, _ := json.Marshal(map[string]any{
		"Value":    value,
		"Against":  against,
		"Critical": critical,
		"Success":  success,
	})
	return "check_result", data
}

func NewUseActionEvent(a definition.Action, src definition.Creature, target definition.Creature) (string, []byte) {
	data, _ := json.Marshal(map[string]any{
		"Action":  snapshotAction(a),
		"Source":  snapshotCreature(src),
		"Target":  snapshotCreature(target),
		"Success": true,
	})
	return "use_action", data
}

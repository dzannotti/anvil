package ai

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

func findPotentialTargets(_ *core.World, actor *core.Actor, action core.Action, encounter *core.Encounter) []grid.Position {
	actionTags := action.Tags()

	if actionTags.HasTag(tags.Move) || actionTags.HasTag(tags.Dash) {
		return findMovementTargets(actor, action)
	}

	hostileActors := encounter.HostileActors(actor)
	var potentialTargets []grid.Position
	for _, hostileActor := range hostileActors {
		target := hostileActor.Position
		if canReachTargetFromAnyPosition(actor, action, target) {
			potentialTargets = append(potentialTargets, target)
		}
	}

	return potentialTargets
}

func findMovementTargets(actor *core.Actor, action core.Action) []grid.Position {
	validPositions := action.ValidPositions(actor.Position)
	if len(validPositions) == 0 {
		return []grid.Position{}
	}

	return selectStrategicMovementPositions(actor, validPositions)
}

func selectStrategicMovementPositions(actor *core.Actor, validPositions []grid.Position) []grid.Position {
	if len(validPositions) <= 8 {
		return validPositions
	}

	var strategic []grid.Position
	currentPos := actor.Position
	strategic = append(strategic, currentPos)

	nearestEnemy := findNearestEnemy(actor)
	if nearestEnemy == nil {
		return limitPositions(strategic, validPositions, 8)
	}

	enemyPos := nearestEnemy.Position
	currentDistance := calculateDistance(currentPos, enemyPos)

	closerPositions := make([]grid.Position, 0, 3)
	fartherPositions := make([]grid.Position, 0, 3)

	for _, pos := range validPositions {
		if pos == currentPos {
			continue
		}

		newDistance := calculateDistance(pos, enemyPos)
		switch {
		case newDistance < currentDistance && len(closerPositions) < 3:
			closerPositions = append(closerPositions, pos)
		case newDistance > currentDistance && len(fartherPositions) < 3:
			fartherPositions = append(fartherPositions, pos)
		}
	}

	strategic = append(strategic, closerPositions...)
	strategic = append(strategic, fartherPositions...)

	return limitPositions(strategic, validPositions, 8)
}

func findNearestEnemy(actor *core.Actor) *core.Actor {
	if actor.Encounter == nil {
		return nil
	}

	var nearestEnemy *core.Actor
	minDistance := 999
	currentPos := actor.Position

	for _, enemy := range actor.Encounter.HostileActors(actor) {
		if enemy.IsDead() {
			continue
		}
		distance := calculateDistance(currentPos, enemy.Position)
		if distance < minDistance {
			minDistance = distance
			nearestEnemy = enemy
		}
	}

	return nearestEnemy
}

func limitPositions(strategic []grid.Position, validPositions []grid.Position, maxCount int) []grid.Position {
	if len(strategic) >= maxCount {
		return strategic[:maxCount]
	}

	addedPositions := make(map[grid.Position]bool)
	for _, pos := range strategic {
		addedPositions[pos] = true
	}

	for _, pos := range validPositions {
		if len(strategic) >= maxCount {
			break
		}
		if !addedPositions[pos] {
			strategic = append(strategic, pos)
			addedPositions[pos] = true
		}
	}

	return strategic
}

func checkActionFeasibility(action core.Action, actor *core.Actor) bool {
	if !actor.CanAct() || actor.IsDead() {
		return false
	}

	if !action.CanAfford() {
		return false
	}

	actionTags := action.Tags()
	return actionTags.HasTag(tags.Attack) ||
		actionTags.HasTag(tags.Spell) ||
		actionTags.HasTag(tags.Move) ||
		actionTags.HasTag(tags.Dash) ||
		actionTags.HasTag(tags.Dodge) ||
		actionTags.HasTag(tags.Help) ||
		actionTags.HasTag(tags.Teleport)
}

func canReachTargetFromAnyPosition(actor *core.Actor, action core.Action, target grid.Position) bool {
	validPositions := action.ValidPositions(actor.Position)
	for _, validPos := range validPositions {
		if validPos == target {
			return true
		}
	}
	return false
}

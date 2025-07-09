# AI System Implementation: Simulation-Based Decision Making

*Status: In Progress*

## Overview

This document outlines the implementation of a simulation-based AI system that follows a structured decision-making flow to evaluate actions, determine optimal positioning, and select the best course of action each turn.

## Target AI Decision Flow

Our AI system will follow a 7-step decision process:

### 1. Initial Action Feasibility Check
- **Purpose**: Verify each action can be performed
- **Checks**: Resource costs, cooldowns, conditions, status effects
- **Outcome**: Filter out impossible actions to avoid wasted computation

### 2. Target Evaluation & Simulation
- **Purpose**: Evaluate potential targets for each viable action
- **Process**: For each action, identify all valid targets (characters, terrain, positions)
- **Simulation**: Run "what if" scenarios for each target
- **Scoring**: Calculate damage scores, health impacts, buff/debuff effects
- **Output**: Per-target simulation results

### 3. Action Scoring & Aggregation
- **Purpose**: Transform raw metric values into weighted behavioral scores, then aggregate to final action score
- **Process**: Apply archetype weights to raw metrics, then sum/subtract/multiply to get final score
- **Input**: Raw metric values (e.g., `damage_dealt: 25`, `kill_potential: 15`)
- **Weighting**: Apply creature configuration weights to metrics
- **Output**: Single "final score" per action-target combination

### 4. Position Optimization
- **Purpose**: Determine optimal casting/execution position for each action
- **Evaluation**: Consider surface conditions, nearby allies/enemies, line of sight, range
- **Scoring**: Calculate "PositionScores" for different execution positions
- **Output**: Best execution position per action

### 5. Movement Assessment
- **Purpose**: Calculate movement costs and path efficiency
- **Evaluation**: Path difficulty, movement costs, opportunity attacks
- **Scoring**: Calculate "MovementScore" for reaching optimal positions
- **Output**: Movement efficiency rating per action

### 6. Final Action Selection
- **Purpose**: Select highest-scoring action with complete context
- **Combines**: Action score + Position score + Movement score
- **Determines**: Which skill/attack, which target, which position, which movement path
- **Output**: Complete action plan

### 7. Fallback Mechanism
- **Purpose**: Handle cases where no positive actions exist
- **Fallback**: Find optimal positioning without specific action
- **Movement**: Move to best available tactical position
- **Output**: Movement-only turn if no good offensive actions

## Current System Issues

### Structural Problems
- **Overly simplistic flow**: Calculate all scores → pick highest → execute
- **No target simulation**: Actions evaluated in isolation
- **No position optimization**: Actions evaluated from current position only
- **No fallback mechanism**: AI gets stuck when no good actions exist
- **No feasibility checking**: Wastes computation on impossible actions

### Implementation Problems
- **Hardcoded scoring**: BaseDamageScore = 20, planWeight = 0.8
- **Simple additive scoring**: All metrics just add to a total score
- **No behavioral differentiation**: All creatures behave identically
- **No situational awareness**: Same logic regardless of context

## New Architecture Design

### Core Principle: Metrics as Simulation Engine
The AI will simulate potential outcomes for each action-target combination by running metrics:

```go
// Raw metrics are produced by running our metric system
// Example: damageMetric.Evaluate() returns:
// {
//     "damage_dealt": 25,
//     "kill_potential": 15,
//     "threat_elimination": 8,
//     "overkill_waste": -3,
// }
```

### Target-Centric Evaluation
Instead of evaluating actions in isolation, evaluate action-target combinations:

```go
type ActionTargetEvaluation struct {
    Action       core.Action
    Target       grid.Position
    RawMetrics   map[string]int  // Raw metric values from simulation
    Position     grid.Position   // Optimal casting position
    Movement     []grid.Position // Path to position
    FinalScore   int             // After applying archetype weights
}
```

## Simulation-Based Decision System

### Action Feasibility Evaluation
The AI starts by filtering actions based on immediate feasibility:

```go
func (ai *AI) checkActionFeasibility(action core.Action, actor *core.Actor) bool {
    // Check resource costs (AP, spell slots, etc.)
    if !actor.Resources.CanAfford(action.Cost()) {
        return false
    }
    
    // Check cooldowns and conditions
    if action.OnCooldown() || actor.HasBlockingCondition(action) {
        return false
    }
    
    // Check status effects that prevent this action type
    if actor.Conditions.BlocksAction(action) {
        return false
    }
    
    return true
}
```

### Target Simulation Process
For each viable action, the AI identifies all potential targets and runs metrics to simulate outcomes:

```go
func (ai *AI) simulateActionTarget(action core.Action, target grid.Position, fromPosition grid.Position) map[string]int {
    affected := action.AffectedPositions([]grid.Position{target})
    
    // "Simulation" is achieved by running our metrics - these produce raw metric values
    damageResults := ai.damageMetric.Evaluate(ai.world, ai.actor, action, target, affected)
    controlResults := ai.controlMetric.Evaluate(ai.world, ai.actor, action, target, affected)
    positionResults := ai.positionMetric.Evaluate(ai.world, ai.actor, action, target, affected)
    
    // Combine all raw metric values into a single map
    allMetrics := make(map[string]int)
    for k, v := range damageResults {
        allMetrics[k] = v
    }
    for k, v := range controlResults {
        allMetrics[k] = v
    }
    for k, v := range positionResults {
        allMetrics[k] = v
    }
    
    return allMetrics
}
```

### Position Optimization
The AI evaluates different positions from which to execute each action by running metrics from each position:

```go
func (ai *AI) optimizePosition(action core.Action, target grid.Position) (grid.Position, int) {
    validPositions := action.ValidPositions(ai.actor.Position)
    bestPosition := ai.actor.Position
    bestScore := 0
    
    for _, pos := range validPositions {
        // Temporarily move actor to evaluate position
        originalPos := ai.actor.Position
        ai.actor.Position = pos
        
        // Run metrics from this position to "simulate" what would happen
        affected := action.AffectedPositions([]grid.Position{target})
        positionResults := ai.positionMetric.Evaluate(ai.world, ai.actor, action, target, affected)
        
        // Calculate position score from metric results
        score := positionResults["tactical_advantage"] + positionResults["threat_distance"]
        
        if score > bestScore {
            bestScore = score
            bestPosition = pos
        }
        
        // Restore original position
        ai.actor.Position = originalPos
    }
    
    return bestPosition, bestScore
}
```

### Movement Cost Assessment
The AI calculates the cost and efficiency of reaching optimal positions using metrics:

```go
func (ai *AI) assessMovement(from, to grid.Position) int {
    path, found := ai.world.FindPath(from, to)
    if !found {
        return math.MinInt // Impossible movement
    }
    
    // Use movement metric to evaluate the path
    // Create a mock movement action for this path
    moveAction := createMoveAction(from, to)
    
    // Run metrics to simulate movement consequences
    movementResults := ai.movementMetric.Evaluate(ai.world, ai.actor, moveAction, to, []grid.Position{to})
    
    // Movement score includes path cost, opportunity attacks, surface effects
    return movementResults["movement_efficiency"] + movementResults["opportunity_cost"]
}
```

### Archetype-Based Behavior Modification
Different creature types will have different weights for each metric:

```go
type AIArchetype struct {
    Name string
    Weights map[string]float32  // Weight for each metric (e.g., "damage_dealt": 2.0)
}

var BerserkerArchetype = AIArchetype{
    Name: "berserker",
    Weights: map[string]float32{
        "damage_dealt":       2.0,  // 2x damage preference
        "kill_potential":     1.5,  // Loves finishing enemies
        "threat_elimination": 1.0,  // Normal threat assessment
        "survival_chance":    0.3,  // Ignores danger
        "tactical_advantage": 0.5,  // Doesn't care about positioning
    },
}

func (ai *AI) applyArchetypeWeights(rawMetrics map[string]int, archetype AIArchetype) int {
    finalScore := 0
    
    for metricName, rawValue := range rawMetrics {
        weight := archetype.Weights[metricName]
        if weight == 0 {
            weight = 1.0  // Default weight if not specified
        }
        
        weightedScore := int(float32(rawValue) * weight)
        finalScore += weightedScore
    }
    
    return finalScore
}
```

## Implementation Approach

### Metrics ARE the Simulation Engine
Our metrics system IS how we simulate "what if" scenarios. When the AI needs to know "What will happen if I cast fireball at position X?", it runs the damage metric, control metric, and position metric to find out.

### Target-First Evaluation
Actions are meaningless without targets. The AI evaluates action-target combinations by running metrics for each combination.

### Position-Aware Planning
The AI considers where to execute actions by temporarily moving the actor to different positions and running metrics from those positions.

### Fallback Behavior
When no actions score positively, the AI falls back to pure positioning - using the movement metric to find the best available tactical position.

### Behavioral Differentiation Through Weights
Different creature types (berserker, defensive, mage, etc.) are differentiated through weight multipliers applied to metric results.

## Key Advantages

1. **Contextual Decision Making**: Actions are evaluated with full context of target, position, and movement
2. **Realistic Simulation**: AI considers actual outcomes rather than abstract scores
3. **Flexible Targeting**: Can handle character targets, terrain targets, area targets
4. **Intelligent Positioning**: AI considers optimal casting positions, not just current position
5. **Graceful Fallback**: AI has sensible behavior when no good actions exist
6. **Extensible Architecture**: Easy to add new action types and scoring considerations

## Weight and Metrics System

For detailed specifications of weights, metrics, and archetype configurations, see [AI Weights and Metrics](12_AI_WEIGHTS_AND_METRICS.md).

**Summary**:
- **Phase 1**: 3 core weights (`damage_enemy`, `friendly_fire`, `survival_threat`) from 2 metrics
- **Metrics simulate outcomes**: Raw values like `damage_enemy: 25`, `survival_threat: -12`
- **Weights create behavior**: Berserker 2.0x damage focus, Defensive 2.0x survival focus
- **Future expansion**: 12+ weights across 5+ metrics for sophisticated behavior

## Next Steps

Start with Phase 1 implementation:
1. Implement damage and positioning simulation metrics
2. Create basic archetype weight system  
3. Test behavioral differences between archetypes
4. Build foundation for future metric expansion
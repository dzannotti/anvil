# AI Weights and Metrics System

*Status: Design Phase*

## Overview

This document defines the weight system and metrics that will drive AI decision-making. Metrics simulate "what if" scenarios by producing raw values, which are then transformed into behavioral preferences through archetype-specific weights.

## Weight System Analysis

### Larian's Weight Categories

Based on Larian's system, we can identify these key weight categories:

#### 1. Damage Weights
- `MULTIPLIER_DAMAGE_ENEMY_POS`: 1.0 (damage dealt to enemies)
- `MULTIPLIER_DAMAGE_ALLY_NEG`: 1.5 (friendly fire penalty)
- `MULTIPLIER_KILL_ENEMY`: 2.50 (bonus for killing enemies)

#### 2. Healing Weights
- `MULTIPLIER_HEAL_SELF_POS`: 1.0 (self-healing value)
- `MULTIPLIER_HEAL_ALLY_POS`: 1.0 (ally healing value)

#### 3. Control Weights
- `MULTIPLIER_CONTROL_ENEMY_POS`: 1.0 (applying control effects to enemies)
- `MULTIPLIER_CONTROL_SELF_NEG`: 2.0 (avoiding control effects on self)

#### 4. Target Priority Weights
- `MULTIPLIER_TARGET_MY_ENEMY`: 1.50 (targeting enemies)
- `MULTIPLIER_TARGET_MY_HOSTILE`: 3.00 (targeting hostile creatures)
- `MULTIPLIER_TARGET_AGGRO_MARKED`: 5.00 (targeting marked/prioritized enemies)

#### 5. Positioning Weights
- `MULTIPLIER_ENDPOS_FLANKED`: 0.05 (avoid ending in flanked positions)
- `MULTIPLIER_ENDPOS_HEIGHT_DIFFERENCE`: 0.002 (height advantage)

#### 6. Special Situation Weights
- `MULTIPLIER_RESURRECT`: 4.00 (resurrection actions)
- `MULTIPLIER_STATUS_REMOVE`: 1.00 (removing status effects)
- `MULTIPLIER_CHARMED`: 2.50 (charm effects)

## Our Weight System Design

### Phase 1: Core Combat Weights (Priority Implementation)
We'll start with the 3 most important and simple weights:

1. **damage_enemy** (from damage metric) - How much damage we deal to enemies
2. **friendly_fire** (from damage metric) - Penalty for damaging allies  
3. **survival_threat** (from positioning metric) - How much danger we're in

### Phase 2: Advanced Combat Weights
4. **kill_potential** (from damage metric) - Bonus for finishing off enemies
5. **tactical_position** (from positioning metric) - Positional advantage
6. **movement_efficiency** (from movement metric) - Cost/benefit of movement

### Phase 3: Complex Behavioral Weights
7. **crowd_control** (from control metric) - Applying/avoiding status effects
8. **target_priority** (from target metric) - Intelligent target selection
9. **resource_conservation** (from resource metric) - Efficient use of abilities

### Phase 4: Advanced Tactical Weights
10. **area_denial** (from control metric) - Battlefield control
11. **team_coordination** (from team metric) - Supporting allies
12. **threat_elimination** (from combined metrics) - Removing biggest threats

## Detailed Metric Specifications

### Damage Simulation Metric

**Purpose**: Simulate combat damage outcomes for action-target combinations

**Inputs**: 
- `action`: The action being evaluated (attack, spell, etc.)
- `target`: Target position
- `affected`: All positions affected by the action

**Simulation Process**:
1. **Identify affected actors**: Find all actors at affected positions
2. **Line of sight check**: Only consider targets with clear line of sight
3. **Calculate raw damage**: Use action's average damage for each target
4. **Assess kill potential**: Check if damage >= target's current HP
5. **Evaluate overkill**: Calculate wasted damage beyond target HP
6. **Friendly fire assessment**: Identify ally damage vs enemy damage

**Weight Outputs**:
```go
{
    "damage_enemy":    25,  // Raw damage to enemy actors
    "friendly_fire":   -8,  // Negative value for ally damage
    "kill_potential":  15,  // Bonus if we can finish off enemies
    "overkill_waste":  -3,  // Penalty for excessive damage
}
```

**Detailed Logic**:
```go
func (d DamageMetric) Evaluate(world, actor, action, target, affected) map[string]int {
    results := map[string]int{}
    
    for _, pos := range affected {
        targetActor := world.ActorAt(pos)
        if targetActor == nil { continue }
        
        if !world.HasLineOfSight(actor.Position, pos) { continue }
        
        damage := action.AverageDamage()
        
        if actor.IsHostileTo(targetActor) {
            // Enemy damage
            actualDamage := min(damage, targetActor.HitPoints)
            results["damage_enemy"] += actualDamage
            
            // Kill potential bonus
            if damage >= targetActor.HitPoints {
                results["kill_potential"] += 15
            }
            
            // Overkill penalty
            if damage > targetActor.HitPoints {
                results["overkill_waste"] -= (damage - targetActor.HitPoints) / 2
            }
        } else {
            // Friendly fire penalty
            actualDamage := min(damage, targetActor.HitPoints)
            results["friendly_fire"] -= actualDamage * 2  // 2x penalty
        }
    }
    
    return results
}
```

### Positioning Simulation Metric

**Purpose**: Simulate tactical positioning and threat assessment

**Inputs**:
- `action`: Action being evaluated (including movement actions)
- `target`: Target position (where we're considering moving or casting from)
- `affected`: Positions affected by the action

**Simulation Process**:
1. **Threat assessment**: Evaluate immediate danger at target position
2. **Tactical advantage**: Check for height, flanking, cover opportunities
3. **Escape routes**: Count available retreat paths
4. **Movement cost**: Calculate efficiency of reaching position

**Weight Outputs**:
```go
{
    "survival_threat":     -12,  // Negative for dangerous positions
    "tactical_position":    8,   // Positive for advantageous positions
    "escape_routes":        3,   // Number of available retreat paths
    "movement_efficiency": -5,   // Cost of reaching this position
}
```

**Detailed Logic**:
```go
func (p PositioningMetric) Evaluate(world, actor, action, target, affected) map[string]int {
    results := map[string]int{}
    
    // Threat assessment - how dangerous is this position?
    nearbyEnemies := world.ActorsInRange(target, 2, func(a *Actor) bool {
        return a.IsHostileTo(actor)
    })
    
    threatLevel := 0
    for _, enemy := range nearbyEnemies {
        // Estimate threat based on enemy damage potential
        threatLevel += enemy.EstimatedDamage() 
    }
    results["survival_threat"] = -threatLevel  // Negative because threat is bad
    
    // Tactical positioning
    tacticalScore := 0
    
    // Height advantage
    if p.hasHeightAdvantage(target, nearbyEnemies) {
        tacticalScore += 10
    }
    
    // Flanking opportunities
    flankingTargets := p.countFlankingOpportunities(target, nearbyEnemies)
    tacticalScore += flankingTargets * 5
    
    // Cover availability
    if p.hasGoodCover(world, target) {
        tacticalScore += 8
    }
    
    results["tactical_position"] = tacticalScore
    
    // Escape routes
    escapeCount := p.countEscapeRoutes(world, target)
    results["escape_routes"] = escapeCount
    
    // Movement efficiency
    if action.Tags().HasTag(tags.Move) {
        path, found := world.FindPath(actor.Position, target)
        if found {
            results["movement_efficiency"] = -int(path.TotalCost)
        } else {
            results["movement_efficiency"] = -999  // Impossible movement
        }
    }
    
    return results
}
```

### Control Simulation Metric (Future)

**Purpose**: Simulate crowd control and status effect outcomes

**Simulation Process**:
1. **Status effect application**: Evaluate control effects on targets
2. **Status removal**: Check for removing debuffs from allies
3. **Area denial**: Assess battlefield control potential
4. **Duration and effectiveness**: Consider how long effects last

**Weight Outputs**:
```go
{
    "crowd_control":   20,  // Value of control effects applied
    "status_removal":  15,  // Value of removing ally debuffs
    "area_denial":     10,  // Battlefield control value
}
```

### Target Priority Metric (Future)

**Purpose**: Simulate intelligent target selection

**Simulation Process**:
1. **Threat assessment**: Identify most dangerous enemies
2. **Vulnerability detection**: Find weakest/wounded targets
3. **Role identification**: Prioritize healers, casters, etc.
4. **Focus fire coordination**: Encourage targeting same enemies

**Weight Outputs**:
```go
{
    "threat_elimination": 25,  // Value of targeting dangerous enemies
    "vulnerability_exploit": 18, // Value of targeting weak enemies
    "role_priority": 12,       // Bonus for targeting priority roles
    "focus_fire": 8,           // Coordination bonus
}
```

### Movement Simulation Metric (Future)

**Purpose**: Simulate movement efficiency and safety

**Simulation Process**:
1. **Path cost calculation**: Movement points required
2. **Opportunity attack assessment**: Risk of AoOs along path
3. **Surface effect evaluation**: Hazardous terrain penalties
4. **Positioning improvement**: Benefit of new position vs current

**Weight Outputs**:
```go
{
    "movement_efficiency": -5,  // Cost of movement
    "opportunity_risk": -8,     // Risk of opportunity attacks
    "surface_penalty": -3,      // Hazardous terrain costs
    "position_improvement": 12, // Benefit of new position
}
```

## Archetype Weight Configurations

### Phase 1: Core Archetypes (3 weights)
```go
var BerserkerArchetype = AIArchetype{
    Name: "berserker",
    Weights: map[string]float32{
        "damage_enemy":     2.0,  // 2x damage focus
        "friendly_fire":    0.5,  // Doesn't care about allies
        "survival_threat":  0.3,  // Ignores danger
    },
}

var DefensiveArchetype = AIArchetype{
    Name: "defensive", 
    Weights: map[string]float32{
        "damage_enemy":     1.0,  // Normal damage
        "friendly_fire":    2.0,  // Avoids friendly fire
        "survival_threat":  2.0,  // Prioritizes safety
    },
}

var DefaultArchetype = AIArchetype{
    Name: "default",
    Weights: map[string]float32{
        "damage_enemy":     1.0,  // Baseline
        "friendly_fire":    1.5,  // Slight friendly fire avoidance
        "survival_threat":  1.0,  // Normal survival instinct
    },
}
```

### Future: Advanced Archetypes
```go
var MageArchetype = AIArchetype{
    Name: "mage",
    Weights: map[string]float32{
        "damage_enemy":           1.2,  // Moderate damage focus
        "friendly_fire":          1.8,  // Careful with AoE
        "survival_threat":        1.8,  // Stays safe
        "kill_potential":         1.3,  // Likes finishing enemies
        "tactical_position":      1.4,  // Values positioning
        "crowd_control":          1.6,  // Uses control spells
        "area_denial":            1.5,  // Battlefield control
        "resource_conservation":  1.7,  // Efficient spell use
    },
}

var RogueArchetype = AIArchetype{
    Name: "rogue",
    Weights: map[string]float32{
        "damage_enemy":        1.5,  // Good damage focus
        "friendly_fire":       1.2,  // Somewhat careful
        "survival_threat":     1.4,  // Moderately cautious
        "kill_potential":      2.0,  // Loves finishing enemies
        "tactical_position":   1.8,  // Positioning expert
        "vulnerability_exploit": 2.2, // Targets weak enemies
        "movement_efficiency": 1.3,  // Mobile fighter
    },
}
```

## Implementation Priority

### Phase 1: Foundation (3 weights, 2 metrics)
1. **Damage Simulation Metric**: damage_enemy, friendly_fire, kill_potential  
2. **Positioning Simulation Metric**: survival_threat, tactical_position
3. **Basic Archetype System**: Berserker, Defensive, Default

**Success Criteria**: Berserker acts more aggressively, Defensive avoids danger

### Phase 2: Enhanced Tactics (6 weights, 3 metrics)
4. **Movement Simulation Metric**: movement_efficiency, escape_routes
5. **Enhanced Positioning**: Height advantage, flanking calculations
6. **Kill Focus**: Better finish-off behavior

**Success Criteria**: Smarter positioning, finishing wounded enemies

### Phase 3: Advanced Behavior (9 weights, 4 metrics)
7. **Control Simulation Metric**: crowd_control, status_removal, area_denial
8. **Target Priority Metric**: threat_elimination, vulnerability_exploit
9. **Resource Management**: Efficient ability usage

**Success Criteria**: Intelligent target selection, crowd control usage

### Phase 4: Sophisticated AI (12+ weights, 5+ metrics)
10. **Advanced Archetypes**: Mage, Rogue, Healer variants
11. **Situational Adaptation**: Context-aware weight adjustments
12. **Team Coordination**: Multi-actor tactical planning

**Success Criteria**: Rich behavioral variety, team-based tactics

## Testing and Validation

### Behavioral Validation Tests
1. **Aggression Test**: Berserker should engage more readily than Defensive
2. **Safety Test**: Defensive should retreat when overwhelmed
3. **Efficiency Test**: All archetypes should make reasonable decisions
4. **Differentiation Test**: Different archetypes should behave noticeably different

### Performance Benchmarks
- Metric evaluation should complete in <1ms per action-target combination
- Full decision process should complete in <10ms per turn
- Memory usage should remain reasonable with multiple metrics

### Balance Validation
- No archetype should be obviously superior in all situations
- All weights should contribute meaningfully to decisions
- Edge cases should be handled gracefully (no crashes, infinite loops)

## Future Extensions

### Dynamic Weight Adjustment
- Adjust weights based on health percentage
- Modify behavior when resources are low
- Adapt to enemy composition

### Situational Modifiers
- Different weights in different encounter types
- Environmental factors affecting decisions
- Time pressure modifications

### Learning and Adaptation
- Track successful vs unsuccessful decisions
- Adjust weights based on outcomes
- Player behavior adaptation

## Complete Larian Weight Mapping

This section maps every weight mentioned in Larian's system to the metrics we'll need to support them.

### Damage-Related Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_DAMAGE_SELF_POS` | 1.0 | Healing Metric | Phase 3 | Self-healing benefits |
| `MULTIPLIER_DAMAGE_SELF_NEG` | 1.0 | Damage Metric | Phase 1 | Self-damage penalties |
| `MULTIPLIER_DAMAGE_ENEMY_POS` | 1.0 | Damage Metric | Phase 1 | **Core weight: damage_enemy** |
| `MULTIPLIER_DAMAGE_ENEMY_NEG` | 1.0 | Damage Metric | Phase 2 | Avoiding damage to enemies (rare) |
| `MULTIPLIER_DAMAGE_ALLY_POS` | 1.0 | Healing Metric | Phase 3 | Healing allies |
| `MULTIPLIER_DAMAGE_ALLY_NEG` | 1.5 | Damage Metric | Phase 1 | **Core weight: friendly_fire** |

### Healing-Related Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_HEAL_SELF_POS` | 1.0 | Healing Metric | Phase 3 | Self-healing value |
| `MULTIPLIER_HEAL_SELF_NEG` | 1.0 | Healing Metric | Phase 3 | Reverse healing (damage) |
| `MULTIPLIER_HEAL_ENEMY_POS` | 1.0 | Healing Metric | Phase 4 | Healing enemies (rare) |
| `MULTIPLIER_HEAL_ENEMY_NEG` | 1.0 | Healing Metric | Phase 4 | Damaging enemies via "healing" |
| `MULTIPLIER_HEAL_ALLY_POS` | 1.0 | Healing Metric | Phase 3 | Ally healing value |
| `MULTIPLIER_HEAL_ALLY_NEG` | 1.0 | Healing Metric | Phase 3 | Damaging allies via "healing" |

### Control Status Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_CONTROL_SELF_POS` | 1.0 | Control Metric | Phase 3 | Beneficial self-effects |
| `MULTIPLIER_CONTROL_SELF_NEG` | 2.0 | Control Metric | Phase 3 | Avoiding self-debuffs |
| `MULTIPLIER_CONTROL_ENEMY_POS` | 1.0 | Control Metric | Phase 3 | **Core weight: crowd_control** |
| `MULTIPLIER_CONTROL_ENEMY_NEG` | 1.0 | Control Metric | Phase 4 | Removing enemy buffs |
| `MULTIPLIER_CONTROL_ALLY_POS` | 1.0 | Control Metric | Phase 3 | Buffing allies |
| `MULTIPLIER_CONTROL_ALLY_NEG` | 1.0 | Control Metric | Phase 3 | Removing ally debuffs |

### Target Priority Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_TARGET_MY_ENEMY` | 1.50 | Target Priority Metric | Phase 3 | **Core weight: threat_elimination** |
| `MULTIPLIER_TARGET_MY_HOSTILE` | 3.00 | Target Priority Metric | Phase 3 | High-priority hostile targets |
| `MULTIPLIER_TARGET_SUMMON` | 0.35 | Target Priority Metric | Phase 4 | De-prioritize summons |
| `MULTIPLIER_TARGET_AGGRO_MARKED` | 5.00 | Target Priority Metric | Phase 4 | Marked/taunted targets |
| `MULTIPLIER_TARGET_PREFERRED` | 20.00 | Target Priority Metric | Phase 4 | Special priority targets |

### Positioning Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_ENDPOS_FLANKED` | 0.05 | Positioning Metric | Phase 2 | **Core weight: survival_threat** |
| `MULTIPLIER_ENDPOS_HEIGHT_DIFFERENCE` | 0.002 | Positioning Metric | Phase 2 | **Core weight: tactical_position** |
| `MULTIPLIER_ENDPOS_TURNED_INVISIBLE` | 0.01 | Positioning Metric | Phase 4 | Invisibility positioning |

### Special Action Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_KILL_ENEMY` | 2.50 | Damage Metric | Phase 2 | **Core weight: kill_potential** |
| `MULTIPLIER_RESURRECT` | 4.00 | Special Actions Metric | Phase 4 | Resurrection priority |
| `MULTIPLIER_STATUS_REMOVE` | 1.00 | Control Metric | Phase 3 | **Core weight: status_removal** |
| `MULTIPLIER_CHARMED` | 2.50 | Control Metric | Phase 4 | Charm effect value |

### Action Cost Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_MOVEMENTCOST` | 1.0 | Movement Metric | Phase 2 | **Core weight: movement_efficiency** |
| `MULTIPLIER_COOLDOWN` | 1.0 | Resource Metric | Phase 3 | Cooldown impact on desirability |
| `MULTIPLIER_ITEM_USE_LOWHEALTH` | 1.0 | Resource Metric | Phase 3 | Item usage when low health |

### Surface/Environment Weights

| Larian Weight | Default Value | Our Metric | Implementation Phase | Notes |
|---------------|---------------|------------|---------------------|--------|
| `MULTIPLIER_SURFACE_CLOUD` | 1.0 | Environmental Metric | Phase 4 | Cloud surface effects |
| `MULTIPLIER_SURFACE_FIRE` | 1.0 | Environmental Metric | Phase 4 | Fire surface effects |
| `MULTIPLIER_SURFACE_WATER` | 1.0 | Environmental Metric | Phase 4 | Water surface effects |

## Metric Requirements Summary

To fully support Larian's weight system, we need these metrics:

### Phase 1 (3 weights, 2 metrics)
- **Damage Metric**: `damage_enemy`, `friendly_fire`, `kill_potential`
- **Positioning Metric**: `survival_threat`, `tactical_position`

### Phase 2 (6 weights, 3 metrics)
- **Movement Metric**: `movement_efficiency`, `opportunity_cost`
- Enhanced positioning for flanking and height

### Phase 3 (9 weights, 4 metrics)
- **Control Metric**: `crowd_control`, `status_removal`, `buff_allies`
- **Healing Metric**: `self_healing`, `ally_healing`

### Phase 4 (15+ weights, 7+ metrics)
- **Target Priority Metric**: `threat_elimination`, `vulnerability_exploit`, `role_priority`
- **Resource Metric**: `cooldown_efficiency`, `resource_conservation`
- **Environmental Metric**: `surface_effects`, `terrain_advantage`
- **Special Actions Metric**: `resurrection`, `item_usage`, `special_abilities`

### Total System Scope
- **42 Larian weights** mapped to our system
- **7 core metrics** to implement all weights
- **4 implementation phases** for gradual rollout
- **Full Larian compatibility** when complete

This mapping ensures our system can eventually support the full sophistication of Larian's proven AI architecture while starting with a manageable 3-weight foundation.
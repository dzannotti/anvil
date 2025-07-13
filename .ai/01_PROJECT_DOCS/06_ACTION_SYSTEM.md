# Action System

*Status: Implemented*

## Overview

The Action System is the **execution engine** for all player and NPC actions in D&D 2024 combat. It provides a unified interface for movement, attacks, spells, and special abilities while integrating seamlessly with the resource management, positioning, effect system, and AI decision-making components.

**Core Philosophy: "Actions are the Interface Between Intent and Effect"**

Every player decision and AI choice flows through the Action System. Actions encapsulate targeting logic, resource costs, validation rules, and execution mechanics while maintaining clean separation between what can be done and how it's accomplished.

## Core Components

### Action Interface

All actions implement a unified interface that supports targeting, validation, and execution:

```go
type Action interface {
    Name() string                                    // Human-readable action name
    Archetype() string                              // Action type identifier
    ID() string                                     // Unique instance identifier
    Tags() *tag.Container                           // Action properties and categories
    
    Perform(pos []grid.Position)                    // Execute the action
    ValidPositions(from grid.Position) []grid.Position // Get valid target positions
    AffectedPositions(target []grid.Position) []grid.Position // Get affected positions
    AverageDamage() int                             // AI scoring helper
}
```

### Action Categories

The system supports multiple action types with shared patterns:

```go
// Movement actions
type MoveAction struct {
    owner     *Actor              // Action owner
    name      string              // Display name
    tags      tag.Container       // Action properties
    cost      map[tag.Tag]int     // Resource costs
    reach     int                 // Maximum range
}

// Combat actions  
type MeleeAction struct {
    owner        *Actor           // Action owner
    name         string           // Display name
    tags         tag.Container    // Action properties
    cost         map[tag.Tag]int  // Resource costs
    reach        int              // Attack reach
    damageSource DamageSource     // Damage calculation
}
```

### Resource Management

Actions consume various resource types with validation and enforcement:

```go
// Common resource types
const (
    ResourceAction      = "resource.action"       // Main action per turn
    ResourceBonusAction = "resource.bonusaction"  // Bonus action per turn
    ResourceReaction    = "resource.reaction"     // Reaction per round
    ResourceSpeed       = "resource.speed"        // Movement speed
    ResourceWalkSpeed   = "resource.walkspeed"    // Walk speed specifically
    ResourceSpellSlot1  = "resource.spellslot.1"  // 1st level spell slots
)

// Resource cost validation
func (a *Action) CanAfford() bool {
    return a.owner.Resources.CanAfford(a.cost)
}

// Resource consumption
func (a *Action) Commit() {
    for resource, amount := range a.cost {
        a.owner.ConsumeResource(resource, amount)
    }
}
```

## API Reference

### Action Creation

Actions are typically created through factories and added to actors:

```go
// Movement action (automatically added to all actors)
moveAction := basic.NewMoveAction(actor)

// Melee attack action (added by weapons)
weaponAction := basic.NewMeleeAction(
    owner,
    "Longsword Attack",
    weapon,                                    // Damage source
    1,                                        // Reach in squares
    tag.NewContainer(tags.Melee, tags.Slashing), // Action tags
    map[tag.Tag]int{tags.ResourceAction: 1},  // Cost
)

// Add actions to actor
actor.AddAction(moveAction, weaponAction)
```

### Position-Based Targeting

Actions work with grid positions for targeting and area effects:

```go
// Get valid target positions
validPositions := action.ValidPositions(actor.Position)

// Get positions affected by action
affectedPositions := action.AffectedPositions([]grid.Position{target})

// Execute action at specific position
action.Perform([]grid.Position{targetPosition})
```

### Action Execution Flow

```go
func (a *MeleeAction) Perform(pos []grid.Position) {
    target := a.owner.World.ActorAt(pos[0])
    
    // Begin action event (triggers effects)
    a.owner.Dispatcher.Begin(core.UseActionEvent{
        Action: a,
        Source: a.owner,
        Target: pos,
    })
    defer a.owner.Dispatcher.End()
    
    // Consume resources
    a.Commit()
    
    // Execute attack sequence
    attackResult := a.owner.AttackRoll(target, *a.Tags())
    if attackResult.Success {
        damage := a.owner.DamageRoll(a, attackResult.Critical)
        target.TakeDamage(*damage)
    }
}
```

## Action Implementations

### Movement Actions

Movement actions handle pathfinding, terrain, and speed limitations:

```go
type MoveAction struct {
    cost map[tag.Tag]int // Speed cost per square
}

func (a *MoveAction) ValidPositions(from grid.Position) []grid.Position {
    speed := a.owner.Resources.Remaining(tags.ResourceWalkSpeed)
    shape := shapes.Circle(from, speed)
    
    valid := []grid.Position{}
    for _, pos := range shape {
        // Check terrain validity
        if !a.owner.World.IsValidPosition(pos) {
            continue
        }
        
        // Check for obstacles
        cell := a.owner.World.Grid.At(pos)
        if cell.IsOccupied() {
            continue
        }
        
        // Check pathfinding
        path, ok := a.owner.World.FindPath(from, pos)
        if !ok || path.Speed() > speed {
            continue
        }
        
        valid = append(valid, pos)
    }
    return valid
}

func (a *MoveAction) Perform(pos []grid.Position) {
    // Find optimal path
    path, _ := a.owner.World.FindPath(a.owner.Position, pos[0])
    
    // Emit movement event
    a.owner.Dispatcher.Begin(core.MoveEvent{
        Source: a.owner,
        From:   a.owner.Position,
        To:     pos[0],
        Path:   path,
    })
    defer a.owner.Dispatcher.End()
    
    // Execute movement step by step
    for _, position := range path.Positions()[1:] {
        a.owner.ConsumeResource(tags.ResourceWalkSpeed, 1)
        a.owner.Move(position, a)
    }
}
```

### Melee Combat Actions

Melee actions handle targeting, reach, and damage resolution:

```go
type MeleeAction struct {
    reach        int            // Attack reach
    damageSource DamageSource   // Weapon or natural attack
}

func (a *MeleeAction) ValidPositions(from grid.Position) []grid.Position {
    // Check resource availability
    if !a.CanAfford() {
        return []grid.Position{}
    }
    
    // Get positions within reach
    shape := shapes.Circle(from, a.reach)
    enemies := a.owner.Enemies()
    
    valid := []grid.Position{}
    for _, pos := range shape {
        // Must contain a living enemy
        actor := a.owner.World.ActorAt(pos)
        if actor == nil || actor.IsDead() {
            continue
        }
        
        if slices.Contains(enemies, actor) {
            valid = append(valid, pos)
        }
    }
    return valid
}

func (a *MeleeAction) Tags() *tag.Container {
    // Combine action tags with weapon tags
    combined := tag.NewContainerFromContainer(a.tags)
    combined.Add(*a.damageSource.Tags())
    return &combined
}
```

### Spell Actions (Pattern)

```go
type SpellAction struct {
    spellLevel    int
    school        tag.Tag
    castingTime   tag.Tag        // Action, bonus action, etc.
    range         int
    areaOfEffect  Shape
    damageSource  DamageSource
    saveAttribute tag.Tag        // For saving throw spells
    saveDC        int
}

func (a *SpellAction) ValidPositions(from grid.Position) []grid.Position {
    // Check spell slot availability
    slotType := tag.FromString(fmt.Sprintf("resource.spellslot.%d", a.spellLevel))
    if !a.owner.Resources.CanAfford(map[tag.Tag]int{slotType: 1}) {
        return []grid.Position{}
    }
    
    // Check range and line of sight
    shape := shapes.Circle(from, a.range)
    valid := []grid.Position{}
    for _, pos := range shape {
        if a.owner.World.HasLineOfSight(from, pos) {
            valid = append(valid, pos)
        }
    }
    return valid
}

func (a *SpellAction) Perform(pos []grid.Position) {
    // Consume spell slot
    a.Commit()
    
    // Get affected positions
    affected := a.areaOfEffect.Positions(pos[0])
    
    // Apply spell effects
    for _, affectedPos := range affected {
        target := a.owner.World.ActorAt(affectedPos)
        if target == nil {
            continue
        }
        
        if a.saveAttribute != nil {
            // Saving throw spell
            save := target.SaveThrow(a.saveAttribute, a.saveDC)
            if save.Success {
                // Half damage or no effect
            } else {
                // Full effect
            }
        } else {
            // Attack roll spell
            attack := a.owner.AttackRoll(target, *a.Tags())
            if attack.Success {
                damage := a.owner.DamageRoll(a, attack.Critical)
                target.TakeDamage(*damage)
            }
        }
    }
}
```

## Actor Integration

### Action Management

Actors maintain collections of available actions with duplicate prevention:

```go
type Actor struct {
    Actions []Action    // Available actions
}

func (a *Actor) AddAction(actions ...Action) {
    for _, action := range actions {
        // Prevent duplicates
        if a.HasAction(action) {
            continue
        }
        a.Actions = append(a.Actions, action)
    }
}

func (a *Actor) HasAction(action Action) bool {
    for _, existing := range a.Actions {
        if existing.ID() == action.ID() {
            return true
        }
    }
    return false
}
```

### Default Actions

All actors receive basic actions automatically:

```go
// Movement is universal
actor.AddAction(basic.NewMoveAction(actor))

// Weapons add their own actions
for _, weapon := range actor.Equipped {
    if weaponActions := weapon.GetActions(); len(weaponActions) > 0 {
        actor.AddAction(weaponActions...)
    }
}
```

### Best Action Queries

Actors provide query methods for AI decision-making:

```go
func (a *Actor) BestWeaponAttack() Action {
    var best Action
    bestDamage := 0
    
    for _, action := range a.Actions {
        // Only consider attack actions
        if !action.Tags().HasTag(tags.Attack) {
            continue
        }
        
        if action.AverageDamage() > bestDamage {
            best = action
            bestDamage = action.AverageDamage()
        }
    }
    return best
}
```

## Registry Integration

Actions are registered as factories for data-driven creation:

```go
// Register basic actions
registry.RegisterAction("move", func(owner *core.Actor, options map[string]interface{}) core.Action {
    return basic.NewMoveAction(owner)
})

registry.RegisterAction("attack", func(owner *core.Actor, options map[string]interface{}) core.Action {
    // Extract weapon from options
    weapon := options["weapon"].(core.DamageSource)
    reach := options["reach"].(int)
    
    return basic.NewMeleeAction(owner, "Attack", weapon, reach, 
        tag.NewContainer(tags.Melee), 
        map[tag.Tag]int{tags.ResourceAction: 1})
})

// Spells can be registered similarly
registry.RegisterAction("fireball", func(owner *core.Actor, options map[string]interface{}) core.Action {
    return spells.NewFireball(owner, options)
})
```

## AI Integration

### Action Scoring

The AI system evaluates actions using multiple metrics:

```go
type Score struct {
    Action   core.Action
    Position grid.Position
    Metrics  map[string]int
    Total    int
}

func ScoreAction(world *core.World, actor *core.Actor, action core.Action) []Score {
    validPositions := action.ValidPositions(actor.Position)
    scores := make([]Score, len(validPositions))
    
    for i, pos := range validPositions {
        score := ScorePosition(world, actor, action, pos)
        scores[i] = score
    }
    
    return scores
}

func ScorePosition(world *core.World, actor *core.Actor, action core.Action, pos grid.Position) Score {
    affected := action.AffectedPositions([]grid.Position{pos})
    
    score := Score{
        Action:   action,
        Position: pos,
        Metrics:  make(map[string]int),
    }
    
    // Apply AI metrics
    for _, metric := range metrics.Default {
        metricName := reflect.TypeOf(metric).String()
        score.Metrics[metricName] = metric.Evaluate(world, actor, action, pos, affected)
    }
    
    // Calculate total score
    for _, value := range score.Metrics {
        score.Total += value
    }
    
    return score
}
```

### AI Decision Making

```go
func Play(gameState *core.GameState) {
    actor := gameState.Encounter.ActiveActor()
    
    // Score all available actions
    allScores := []Score{}
    for _, action := range actor.Actions {
        actionScores := ScoreAction(gameState.World, actor, action)
        allScores = append(allScores, actionScores...)
    }
    
    // Pick best action
    best := PickBestAction(allScores)
    if best != nil {
        best.Action.Perform([]grid.Position{best.Position})
    }
    
    // End turn
    gameState.Encounter.EndTurn()
}
```

## Event System Integration

### Action Events

Actions emit comprehensive events for external systems:

```go
// Action usage events
type UseActionEvent struct {
    Action core.Action
    Source *Actor
    Target []grid.Position
}

// Movement events
type MoveEvent struct {
    World  *World
    Source *Actor
    From   grid.Position
    To     grid.Position
    Path   *pathfinding.Result
}

type MoveStepEvent struct {
    Source *Actor
    From   grid.Position
    To     grid.Position
}
```

### Effect System Integration

Actions trigger effect evaluation through actor methods:

```go
func (a *Actor) AttackRoll(target *Actor, tags tag.Container) CheckResult {
    // Create attack expression
    expression := expression.FromD20("Attack Roll")
    
    // Trigger pre-attack effects
    a.Evaluate(&PreAttackRoll{
        Source:     a,
        Target:     target,
        Expression: &expression,
        Tags:       tags,
    })
    
    // Evaluate and check success
    result := expression.Evaluate()
    targetAC := target.Attribute(tags.ArmorClass).Value
    
    return CheckResult{
        Expression: result,
        Success:    result.Value >= targetAC,
        Critical:   result.IsCriticalSuccess(),
    }
}
```

## Usage Patterns

### When to Create Actions

**✅ Create Actions for:**
- Player-triggered abilities (attacks, spells, movement)
- Class features that consume resources
- Spell effects with targeting
- Monster special abilities
- Equipment-granted abilities

**❌ Don't Create Actions for:**
- Passive effects (use Effect System)
- Automatic triggers (use Effect System)
- UI state changes
- Pure calculations

### Action Design Guidelines

**Resource Management:**
```go
// Clear resource costs
cost := map[tag.Tag]int{
    tags.ResourceAction:    1,    // Main action
    tags.ResourceSpellSlot3: 1,   // 3rd level spell slot
}

// Validate before execution
if !action.CanAfford() {
    return // Cannot perform action
}
```

**Position Validation:**
```go
// Always validate targeting
func (a *Action) ValidPositions(from grid.Position) []grid.Position {
    // Check resource availability
    if !a.CanAfford() {
        return []grid.Position{}
    }
    
    // Check range and line of sight
    // Check for valid targets
    // Apply action-specific rules
}
```

**Event Integration:**
```go
// Always emit action events
func (a *Action) Perform(pos []grid.Position) {
    a.owner.Dispatcher.Begin(core.UseActionEvent{...})
    defer a.owner.Dispatcher.End()
    
    // Action implementation
}
```

## Performance Characteristics

- **Position Validation**: O(r²) where r is action range (circular area)
- **Action Execution**: O(1) for single-target, O(n) for area effects
- **AI Scoring**: O(a × p × m) where a=actions, p=positions, m=metrics
- **Memory**: Minimal overhead, actions share damage sources
- **Event Processing**: O(e) where e is number of active effects

## Design Rationale

### Interface-Based Design
- **Polymorphism**: All actions work through same interface
- **Extensibility**: New action types integrate seamlessly  
- **Testing**: Easy to mock and test actions in isolation
- **AI Integration**: Consistent scoring and evaluation

### Position-Based Targeting
- **Grid Combat**: Natural fit for tactical combat
- **Area Effects**: Easy to represent spell areas and weapon reach
- **Pathfinding**: Movement actions use sophisticated pathfinding
- **Line of Sight**: Spells and ranged attacks check visibility

### Resource Management
- **Authenticity**: Matches D&D action economy
- **Balance**: Prevents action spam and maintains game balance
- **Flexibility**: Supports different resource types and costs
- **Validation**: Prevents invalid actions before execution

### Event-Driven Execution
- **Effect Integration**: Actions trigger effects naturally
- **External Communication**: UI and logging systems receive events
- **Debugging**: Complete audit trail of action execution
- **Replay**: Events enable combat replay and analysis

## Integration with Other Systems

### Expression System Integration
Actions use expressions for attack rolls and damage calculations:
```go
// Attack roll with bonuses from effects
attackRoll := expression.FromD20("Attack Roll")
a.owner.Evaluate(&PreAttackRoll{Expression: &attackRoll, ...})

// Damage roll with critical hit support
damageRoll := a.damageSource.Damage()
if critical {
    damageRoll.DoubleDice("Critical Hit")
}
```

### Tag System Integration
Actions use tags for categorization and effect targeting:
```go
// Action properties
actionTags := tag.NewContainer(tags.Attack, tags.Melee, tags.Weapon)

// Effect targeting
if actionTags.HasTag(tags.Spell) {
    // Apply spell-specific effects
}
```

### Pathfinding Integration
Movement actions use sophisticated pathfinding:
```go
// Check reachability
path, canReach := world.FindPath(from, to)
if !canReach || path.Speed() > availableSpeed {
    // Cannot move to position
}
```

## File Locations

- **Core Interface**: `/internal/core/action.go`
- **Basic Actions**: `/internal/ruleset/actions/basic/`
- **Class Actions**: `/internal/ruleset/actions/classes/`
- **Registry**: `/internal/ruleset/registry.go`, `seed.go`
- **AI Integration**: `/internal/ai/ai.go`
- **Tests**: `/internal/ruleset/actions/basic/*_test.go`

The Action System provides the foundation for all interactive gameplay while maintaining clean separation between targeting logic, resource management, and execution mechanics. This architecture enables complex D&D combat while keeping the codebase maintainable and extensible.


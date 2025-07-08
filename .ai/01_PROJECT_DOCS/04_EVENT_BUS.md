# Event Bus

## Overview

The EventBus is the **central nervous system** of the Anvil game engine. It serves as the primary communication mechanism between the game core and all external interfaces (CLI, GUI, logging, AI systems). Every significant game state change, action, and interaction flows through the event system, making it the single source of truth for what's happening in the game.

**Core Philosophy: "All interactions with the outside world flow through events"**

The EventBus implements a **hierarchical observer pattern** where:
- Game logic produces events describing what happened
- External systems (UI, logging, AI) consume events reactively  
- Complex game actions create nested event trees
- No direct coupling between game logic and presentation layers

## Core Components

### Event Structure

```go
type Event struct {
    Data  any     // Typed payload (e.g., AttackRollEvent, TakeDamageEvent)
    Kind  string  // Event type identifier (e.g., "attackRoll", "takeDamage") 
    Depth int     // Nesting level for hierarchical display
    End   bool    // True for closing events in hierarchical pairs
}
```

### Event Dispatcher Interface

```go
type EventDispatcher interface {
    Begin(kind string, data any)  // Start hierarchical event
    End()                         // Close hierarchical event  
    Emit(kind string, data any)   // Simple event (Begin + End)
}
```

### Event Handler

```go
type EventHandler func(event Event)
```

### Dispatcher Implementation

```go
type Dispatcher struct {
    events      []Event         // Event history
    stack       Stack[Event]    // Event stack for depth tracking
    subscribers []EventHandler  // Registered event handlers
}
```

## API Reference

### Basic Usage

```go
// Create dispatcher
dispatcher := eventbus.Dispatcher{}

// Subscribe to ALL events
dispatcher.SubscribeAll(func(event eventbus.Event) {
    prettyprint.Print(os.Stdout, event)
})

// Subscribe to SPECIFIC event types
dispatcher.Subscribe(core.TakeDamageType, func(event eventbus.Event) {
    damageEvent := event.Data.(core.TakeDamageEvent)
    ui.ShowDamageNumber(damageEvent.Target, damageEvent.Damage.Value)
})

dispatcher.Subscribe(core.ConditionChangedType, func(event eventbus.Event) {
    conditionEvent := event.Data.(core.ConditionChangedEvent)
    ui.UpdateStatusEffects(conditionEvent.Source, conditionEvent.Condition)
})

// Simple events (atomic operations)
actor.Dispatcher.Emit(core.ConditionChangedType, core.ConditionChangedEvent{
    Source:    actor,
    Condition: tags.Poisoned,
    Added:     true,
})

// Complex events (hierarchical operations)
actor.Dispatcher.Begin(core.UseActionType, core.UseActionEvent{
    Action: fireball,
    Source: caster,
    Target: targets,
})
defer actor.Dispatcher.End() // ALWAYS use defer

// Nested events happen inside...
actor.Dispatcher.Emit(core.TargetType, core.TargetEvent{Target: enemies})
```

### Subscription Methods

The EventBus provides two subscription approaches:

#### SubscribeAll() - Universal Event Listening
```go
// Receives ALL events - useful for logging, debugging, metrics
dispatcher.SubscribeAll(func(event eventbus.Event) {
    // This handler gets every single event
    prettyprint.Print(os.Stdout, event)  // Log everything
    metrics.RecordEvent(event.Kind)      // Count all event types
})
```

**Use SubscribeAll() for:**
- ✅ Comprehensive logging systems
- ✅ Debug event tracing  
- ✅ Metrics collection
- ✅ Event replay/recording systems

#### Subscribe() - Event-Specific Listening  
```go
// Only receives specific event types - cleaner and more performant
dispatcher.Subscribe(core.TakeDamageType, func(event eventbus.Event) {
    damageEvent := event.Data.(core.TakeDamageEvent)
    ui.ShowDamageNumber(damageEvent.Target, damageEvent.Damage.Value)
})

dispatcher.Subscribe(core.ConditionChangedType, func(event eventbus.Event) {
    conditionEvent := event.Data.(core.ConditionChangedEvent)
    ui.UpdateStatusIcons(conditionEvent.Source, conditionEvent.Condition)
})
```

**Use Subscribe() for:**
- ✅ UI component updates (only care about specific events)
- ✅ Game system reactions (AI, sound effects, etc.)
- ✅ Feature-specific handlers
- ✅ Better performance (no unnecessary event processing)

#### Performance Comparison
```go
// ❌ Old way: Process all events, filter manually
dispatcher.SubscribeAll(func(event eventbus.Event) {
    switch event.Kind {
    case core.TakeDamageType:
        // Handle damage
    case core.ConditionChangedType:
        // Handle condition  
    default:
        return // Ignore most events
    }
})

// ✅ New way: Only receive relevant events
dispatcher.Subscribe(core.TakeDamageType, handleDamage)
dispatcher.Subscribe(core.ConditionChangedType, handleCondition)
// No switch statements, no wasted processing
```

## Event Types and Data

All events are strongly typed with dedicated structures in `/internal/core/event.go`:

### Game Flow Events
- `EncounterType` - Combat encounter starts (`EncounterEvent`)
- `RoundType` - New combat round (`RoundEvent`)
- `TurnType` - Actor's turn begins (`TurnEvent`)

### Combat Events
- `UseActionType` - Action being performed (`UseActionEvent`)
- `AttackRollType` - Attack roll made (`AttackRollEvent`)
- `DamageRollType` - Damage being calculated (`DamageRollEvent`)
- `TakeDamageType` - Actor taking damage (`TakeDamageEvent`)

### Character Events
- `AttributeCalculationType` - Attribute being calculated (`AttributeCalculationEvent`)
- `ConditionChangedType` - Status effect added/removed (`ConditionChangedEvent`)
- `SpendResourceType` - Resource consumed (`SpendResourceEvent`)

### Movement Events
- `MoveType` - Actor movement with full path (`MoveEvent`)
- `MoveStepType` - Individual movement step (`MoveStepEvent`)

## Usage Patterns

### When to Use Emit() vs Begin()/End()

**Use `Emit()` for simple, atomic operations:**
```go
// ✅ Instant state changes
actor.Dispatcher.Emit(core.ConditionChangedType, core.ConditionChangedEvent{...})

// ✅ Resource consumption  
actor.Dispatcher.Emit(core.SpendResourceType, core.SpendResourceEvent{...})

// ✅ Simple notifications
actor.Dispatcher.Emit(core.DeathType, core.DeathEvent{Actor: actor})
```

**Use `Begin()/End()` for complex, hierarchical operations:**
```go
// ✅ Complex actions with multiple steps
actor.Dispatcher.Begin(core.UseActionType, core.UseActionEvent{...})
defer actor.Dispatcher.End() // CRITICAL: Always use defer

// ✅ Operations that can fail partway through
actor.Dispatcher.Begin(core.DamageRollType, core.DamageRollEvent{...})
// ... nested events for calculation steps
actor.Dispatcher.End()

// ✅ Multi-step processes with intermediate results
actor.Dispatcher.Begin(core.AttributeCalculationType, event)
// ... emit intermediate calculation events
actor.Dispatcher.End()
```

### The Action Pattern

Most game actions follow this structure:

```go
func (a *FireballAction) Perform(targets []grid.Position) {
    // 1. Begin action event
    a.owner.Dispatcher.Begin(core.UseActionType, core.UseActionEvent{
        Action: a,
        Source: a.owner, 
        Target: targets,
    })
    defer a.owner.Dispatcher.End()
    
    // 2. Target selection
    affectedActors := a.findTargetsInArea(targets[0])
    a.owner.Dispatcher.Emit(core.TargetType, core.TargetEvent{
        Target: affectedActors,
    })
    
    // 3. Damage calculation (nested hierarchical event)
    a.owner.Dispatcher.Begin(core.DamageRollType, core.DamageRollEvent{
        Source: a.owner,
        DamageSource: a.damageSource,
    })
    
    damageExpression := a.calculateDamage()
    damageExpression.Evaluate()
    
    a.owner.Dispatcher.Emit(core.ExpressionResultType, core.ExpressionResultEvent{
        Expression: &damageExpression,
    })
    a.owner.Dispatcher.End()
    
    // 4. Apply effects (each creates more events)
    for _, target := range affectedActors {
        target.TakeDamage(damageExpression.Clone())
    }
}
```

## External Interface Integration

### CLI Integration (Hierarchical Logging)

```go
func main() {
    dispatcher := eventbus.Dispatcher{}
    
    // Subscribe for pretty-printed hierarchical logging
    dispatcher.Subscribe(func(event eventbus.Event) {
        prettyprint.Print(os.Stdout, event)
    })
    
    gameState := demo.New(&dispatcher)
    encounter := gameState.Encounter
    encounter.Start() // Events automatically flow to console
}

// Output example:
// 🎯 encounter: Combat begins (3 actors)
//   ├─ 🎯 round: Round 1 begins
//   │   ├─ 🎯 turn: Cedric's turn
//   │   │   ├─ 🎯 useAction: Fireball
//   │   │   │   ├─ 🎯 target: 2 zombies targeted
//   │   │   │   ├─ 🎲 damageRoll: 8d6 fire damage
//   │   │   │   │   ├─ 📊 expressionResult: 35 fire damage
//   │   │   │   │   └─ ⚡ (end)
//   │   │   │   └─ ⚡ (end)
//   │   │   └─ ⚡ (end)
//   │   └─ ⚡ (end)
//   └─ ⚡ (end)
```

### GUI Integration (Reactive UI Updates)

```go
func main() {
    dispatcher := eventbus.Dispatcher{}
    ui := InitializeUI()
    
    // Event-specific subscriptions (much cleaner than switch statements!)
    dispatcher.Subscribe(core.TakeDamageType, func(event eventbus.Event) {
        // Show floating damage numbers
        dmgEvent := event.Data.(core.TakeDamageEvent)
        ui.ShowDamageNumber(dmgEvent.Target.Position, dmgEvent.Damage.Value)
    })
    
    dispatcher.Subscribe(core.ConditionChangedType, func(event eventbus.Event) {
        // Update status effect icons
        condEvent := event.Data.(core.ConditionChangedEvent)
        ui.UpdateStatusEffects(condEvent.Source, condEvent.Condition, condEvent.Added)
    })
    
    dispatcher.Subscribe(core.MoveStepType, func(event eventbus.Event) {
        // Animate movement
        moveEvent := event.Data.(core.MoveStepEvent)
        ui.AnimateMovement(moveEvent.Source, moveEvent.From, moveEvent.To)
    })
    
    dispatcher.Subscribe(core.DeathType, func(event eventbus.Event) {
        // Play death animation
        deathEvent := event.Data.(core.DeathEvent)
        ui.PlayDeathAnimation(deathEvent.Actor)
    })
    
    // Optional: Subscribe to all events for debugging/logging
    dispatcher.SubscribeAll(func(event eventbus.Event) {
        log.Debug("Event: %s", event.Kind)
    })
    
    // Same game logic, different presentation
    gameState := demo.New(&dispatcher)
    
    for ui.IsRunning() {
        ui.ProcessInput()
        gameState.Update()
        ui.Render()
    }
}
```

### Testing Integration

```go
func TestCombatSequence(t *testing.T) {
    // Record all events for verification
    var events []eventbus.Event
    dispatcher := eventbus.Dispatcher{}
    
    dispatcher.Subscribe(func(event eventbus.Event) {
        events = append(events, event)
    })
    
    // Set up test scenario
    attacker := createTestActor(&dispatcher, "Attacker", 20)
    target := createTestActor(&dispatcher, "Target", 10)
    
    // Perform attack
    attackAction := base.NewAttackAction(attacker)
    attackAction.Perform([]grid.Position{target.Position})
    
    // Verify complete event flow
    assert.True(t, hasEvent(events, core.UseActionType))
    assert.True(t, hasEvent(events, core.AttackRollType))
    
    // Verify hierarchical structure
    useActionStart := findEvent(events, core.UseActionType, false)
    useActionEnd := findEvent(events, core.UseActionType, true)
    attackRoll := findEvent(events, core.AttackRollType, false)
    
    assert.True(t, useActionStart.Index < attackRoll.Index)
    assert.True(t, attackRoll.Index < useActionEnd.Index)
}
```

## When to Use EventBus

### ✅ **Good Uses**

**Game State Changes:**
- Actor attributes changing (HP, conditions, resources)
- Turn-based game flow (encounters, rounds, turns)
- Combat actions and their results
- Movement and positioning changes

**External Communication:**
- UI updates in response to game events
- Logging game actions for debugging/replay
- AI systems observing game state
- Analytics and metrics collection

**Complex Operations:**
- Multi-step actions (spells with multiple effects)
- Operations that can fail partway through
- Calculations with intermediate results
- Any operation benefiting from audit trails

### ❌ **Bad Uses**

**Internal State Management:**
```go
// DON'T: Use events for internal object state
actor.Dispatcher.Emit("internalFieldChanged", field) // ❌ Too granular

// DO: Use events for meaningful state changes
actor.Dispatcher.Emit(core.ConditionChangedType, condition) // ✅ External significance
```

**Performance-Critical Paths:**
```go
// DON'T: Events in tight loops
for _, cell := range world.AllCells() {
    cell.Dispatcher.Emit("cellProcessed", cell) // ❌ Performance killer
}

// DO: Batch operations, emit summary
world.Dispatcher.Emit(core.WorldProcessedType, summary) // ✅ Efficient
```

**Implementation Details:**
```go
// DON'T: Expose internal algorithms
calculator.Dispatcher.Emit("intermediateCalculation", step) // ❌ Internal detail

// DO: Emit meaningful results
calculator.Dispatcher.Emit(core.AttributeCalculationType, result) // ✅ External meaning
```

## Event Design Guidelines

### 1. Event Naming
Use **noun-based event types** describing what happened:

```go
const (
    AttackRollType    = "attackRoll"     // ✅ What happened
    TakeDamageType    = "takeDamage"     // ✅ What happened  
    ConditionChanged  = "conditionChanged" // ✅ What happened
)
```

### 2. Event Data Design
Create **strongly-typed event structures**:

```go
// ✅ Well-designed event
type AttackRollEvent struct {
    Source *Actor   // Who attacked
    Target *Actor   // Who was attacked
    Roll   int      // Dice result
    Bonus  int      // Attack bonus
    Total  int      // Final result
}
```

### 3. Event Granularity
Balance **meaningful information** vs **event noise**:

```go
// ✅ Right granularity - meaningful to external systems
actor.Dispatcher.Emit(core.TakeDamageType, core.TakeDamageEvent{
    Target: actor,
    Damage: damageExpression,
})

// ❌ Too granular - internal implementation detail
actor.Dispatcher.Emit("hitPointsDecremented", 1) 
```

## Common Patterns and Best Practices

### 1. Always Use Defer for Begin/End

```go
// ✅ Always use defer for Begin/End pairs
actor.Dispatcher.Begin(eventType, eventData)
defer actor.Dispatcher.End()

// ❌ Manual End() calls can be missed
actor.Dispatcher.Begin(eventType, eventData)
if condition { return } // Oops! End() not called
actor.Dispatcher.End()
```

### 2. Keep Events Immutable

```go
// ✅ Event data should be immutable snapshots
event := core.TakeDamageEvent{
    Target: actor,
    Damage: damageExpression.Clone(), // Clone mutable data
}

// ❌ Don't pass mutable references
event := core.TakeDamageEvent{
    Target: actor,
    Damage: &mutableExpression, // Could change after event
}
```

### 3. Handle Subscriber Errors

```go
// Events should not fail due to subscriber errors
dispatcher.Subscribe(func(event eventbus.Event) {
    defer func() {
        if r := recover(); r != nil {
            log.Error("Event subscriber panic: %v", r)
        }
    }()
    // Subscriber logic
})
```

## Debugging and Observability

### Event Tracing
The hierarchical nature enables powerful debugging:

```bash
# Run CLI to see hierarchical event output
make run-cli

# Events display as tree structure showing the complete flow
```

### Performance Monitoring
Track event volume and types for performance tuning:

```go
eventCounts := make(map[string]int)
dispatcher.Subscribe(func(event eventbus.Event) {
    eventCounts[event.Kind]++
    if eventCounts[event.Kind] > 1000 {
        log.Warn("High event volume for %s", event.Kind)
    }
})
```

## Architecture Benefits

The EventBus enables **safe system evolution**:

1. **New Features:** Add new event types without changing existing code
2. **UI Changes:** Modify subscribers without touching game logic  
3. **Logging Updates:** Change logging format by updating subscribers
4. **Testing:** Add test subscribers to verify complex flows
5. **Analytics:** Add metrics subscribers for telemetry

The event-driven architecture ensures the game core remains decoupled from its presentation and external interfaces, making the system highly maintainable and extensible.

## File Locations

- **Core Types**: `/internal/eventbus/` (Dispatcher, Event, EventHandler)
- **Interface**: `/internal/core/logger.go` (EventDispatcher interface)
- **Event Types**: `/internal/core/event.go` (all event structures and constants)
- **Pretty Printing**: `/internal/prettyprint/` (CLI event formatting)
- **Usage Examples**: Domain objects in `/internal/core/` and `/internal/ruleset/`

## Integration Checklist

When adding EventBus to new code:

- [ ] Inject `EventDispatcher` into domain objects
- [ ] Use `Emit()` for simple state changes
- [ ] Use `Begin()/defer End()` for complex operations
- [ ] Create typed event structures in `core/event.go`
- [ ] Add event type constants
- [ ] Update pretty print formatters if needed
- [ ] Test event flows in unit tests
- [ ] Verify events appear correctly in CLI output
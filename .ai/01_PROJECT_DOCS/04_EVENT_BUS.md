# Event Bus

## Overview

The Event Bus is a hierarchical event communication system that manages event distribution with depth tracking and type-safe data interfaces. It provides external communication capabilities for narrative events, logging, UI updates, and system coordination.

## Core Components

### Event Structure

```go
type Event struct {
    Data  any       // Event payload data
    Kind  tag.Tag   // Event type identifier (tag-based)
    Depth int       // Hierarchical depth for nested events
    End   bool      // Indicates event completion
}
```

### Event Handler

```go
type EventHandler func(data Event)
```

Event handlers receive all events dispatched through the system and can process them based on event kind and data.

### Dispatcher

```go
type Dispatcher struct {
    events      []Event                    // Event history
    stack       *stack.Stack[Event]        // Event stack for depth tracking
    subscribers []EventHandler             // Registered event handlers
}
```

## API Reference

### Dispatcher Creation

```go
// Create new dispatcher
dispatcher := eventbus.NewDispatcher()
```

### Event Subscription

```go
// Subscribe to all events
dispatcher.Subscribe(func(event Event) {
    fmt.Printf("Event: %s (depth: %d, end: %t)\n", 
        event.Kind, event.Depth, event.End)
})
```

### Event Dispatching

#### Simple Events

```go
// Dispatch a simple event (start + immediate end)
dispatcher.Add(core.AttackType, attackData)
```

#### Hierarchical Events

```go
// Start a hierarchical event
dispatcher.Start(core.TurnType, turnData)

// Nested events (automatically inherit depth)
dispatcher.Add(core.ActionType, actionData)
dispatcher.Add(core.MovementType, movementData)

// End the hierarchical event
dispatcher.End()
```

### Event Stack Management

The dispatcher automatically manages event depth using an internal stack:

```go
dispatcher.Start(core.EncounterType, encounterData) // Depth: 0
    dispatcher.Start(core.TurnType, turnData)       // Depth: 1
        dispatcher.Add(core.AttackType, attackData) // Depth: 2
    dispatcher.End()                                // End turn
dispatcher.End()                                    // End encounter
```

## Type-Safe Data Interfaces

The event bus provides interface definitions that other packages can implement to maintain type safety without circular dependencies.

### Expression Data Interface

```go
// ExpressionData represents data from expression evaluation
type ExpressionData interface {
    GetValue() int
    GetComponents() []ComponentData
    GetIsCritical() int
}

// ComponentData represents data from individual components
type ComponentData interface {
    GetValue() int
    GetSource() string
    GetType() string
}
```

### Tag Data Interface

```go
// TagData represents tag information
type TagData interface {
    AsString() string
    AsStrings() []string
}
```

## Usage Patterns

### Combat Event Broadcasting

```go
dispatcher := eventbus.NewDispatcher()

// Subscribe to combat events
dispatcher.Subscribe(func(event Event) {
    switch event.Kind.AsString() {
    case core.AttackType.AsString():
        if attackData, ok := event.Data.(AttackData); ok {
            logCombatAttack(attackData)
        }
    case core.DamageType.AsString():
        if dmgData, ok := event.Data.(DamageData); ok {
            updateHealthBars(dmgData)
        }
    }
})

// Broadcast combat encounter
dispatcher.Start(core.EncounterType, EncounterData{
    Participants: participants,
    Location: "Goblin Cave",
})

// Individual combat actions
dispatcher.Add(core.InitiativeType, InitiativeData{
    Order: initiativeOrder,
})

dispatcher.Start(core.TurnType, TurnData{
    Actor: player,
    Actions: availableActions,
})

dispatcher.Add(core.AttackType, AttackData{
    Attacker: player,
    Target: goblin,
    Result: attackExpression,
})

dispatcher.End() // End turn
dispatcher.End() // End encounter
```

### Expression Evaluation Events

```go
// Integration with expression system
expr := expression.New()
expr.AddDice(1, 20, "Attack roll")
expr.AddConstant(5, "Proficiency bonus")

result := expr.Evaluate()

// Broadcast evaluation result
dispatcher.Add(core.ExpressionType, struct{
    Expression eventbus.ExpressionData
    Context    string
}{
    Expression: result, // Implements ExpressionData interface
    Context:    "attack roll",
})
```

### UI Update Events

```go
// UI synchronization
dispatcher.Subscribe(func(event Event) {
    if strings.HasPrefix(event.Kind.AsString(), "ui.") {
        updateInterface(event)
    }
})

// Broadcast UI updates
dispatcher.Add(core.UIHealthType, HealthUpdateData{
    Entity: player,
    OldHP: 25,
    NewHP: 18,
})

dispatcher.Add(core.UIInventoryType, InventoryData{
    Entity: player,
    Items: updatedInventory,
})
```

### Logging and Audit Trail

```go
// Complete event logging
dispatcher.Subscribe(func(event Event) {
    logger.Info("Event", 
        "kind", event.Kind.AsString(),
        "depth", event.Depth,
        "end", event.End,
        "data", event.Data,
    )
})

// Events automatically create audit trail through depth tracking
```

## Event Categories

### Combat Events

- `combat.encounter` - Combat encounter start/end
- `combat.turn` - Individual turn management
- `combat.attack` - Attack attempts and results
- `combat.damage` - Damage dealt and taken
- `combat.healing` - Healing applied

### Expression Events

- `expression.evaluated` - Expression evaluation results
- `expression.modified` - Expression modifications (advantage, etc.)
- `expression.critical` - Critical hit/failure detection

### System Events

- `system.error` - Error conditions
- `system.warning` - Warning conditions
- `system.debug` - Debug information

### UI Events

- `ui.health.update` - Health changes
- `ui.inventory.update` - Inventory modifications
- `ui.status.update` - Status effect changes

## Performance Characteristics

- **Event Dispatch**: O(n) where n is number of subscribers
- **Stack Operations**: O(1) for push/pop operations
- **Memory**: Linear growth with event history retention
- **Depth Tracking**: O(1) depth calculation using stack size

## Design Rationale

The Event Bus implements several key design decisions:

1. **Hierarchical Events**: Depth tracking enables nested event structures
2. **Type Safety**: Interface definitions prevent circular dependencies
3. **Complete History**: Event retention enables replay and audit capabilities
4. **Flexible Subscription**: Single handler type receives all events
5. **Stack-Based Depth**: Automatic depth management for nested events
6. **Immutable Events**: Events are created and dispatched without modification

## Integration Points

### Expression System Integration

```go
// Expression implements EventBus interfaces
type Expression struct {
    // ... fields
}

func (e *Expression) GetValue() int { return e.Value }
func (e *Expression) GetComponents() []eventbus.ComponentData {
    // Convert internal components to interface
}
func (e *Expression) GetIsCritical() int { return e.IsCritical }
```

### Tag System Integration

```go
// Tag implements EventBus interfaces
func (t Tag) AsString() string { return t.value }
func (t Tag) AsStrings() []string { return t.AsStrings() }
```

This design enables loose coupling between systems while maintaining type safety and complete event traceability throughout the game engine.
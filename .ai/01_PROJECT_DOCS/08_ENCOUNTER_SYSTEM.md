# Encounter System

*Status: Implemented*

## Overview

The Encounter System is the core orchestrator for combat and scene management. It coordinates turn-based combat, manages actor lifecycles, and provides event-driven coordination between all game systems.

## Core Components

### Encounter Structure

```go
type Encounter struct {
    EventBus        *eventbus.Dispatcher  // Event coordination
    World           *World                // Game world reference
    Actors          []*Actor              // All participants
    InitiativeOrder []*Actor              // Turn order
    Turn            int                   // Current turn index
    Round           int                   // Current round number
}
```

### Actor Structure

```go
type Actor struct {
    Name        string               // Actor identifier
    Team        TeamID               // Team affiliation
    HitPoints   int                  // Current health
    MaxHitPoints int                 // Maximum health
    Damage      int                  // Attack damage
    EventBus    *eventbus.Dispatcher // Event dispatcher
    Encounter   *Encounter           // Parent encounter
    World       *World               // World reference
}
```

## API Reference

### Encounter Creation

```go
// Create new encounter
world := &core.World{Width: 10, Height: 10}
dispatcher := eventbus.NewDispatcher()
actors := []*core.Actor{player, enemy}
encounter := core.NewEncounter(dispatcher, world, actors)
```

### Combat Orchestration

```go
// Start combat encounter
encounter.Start()

// Combat loop
for !encounter.IsOver() {
    // Get current actor's AI
    actorAI := aiMap[encounter.ActiveActor()]
    
    // AI makes decision and executes action
    actorAI.Play(encounter)
    
    // System automatically advances turn
}

// End encounter
encounter.End()
```

### Turn Management

```go
// Manual turn advancement (normally handled automatically)
encounter.EndTurn()

// Check current actor
currentActor := encounter.ActiveActor()

// Check turn/round state
turn := encounter.Turn
round := encounter.Round
```

## Event Integration

The encounter system is fully event-driven and emits comprehensive events:

### Encounter Events

```go
// Encounter start
core.EncounterEvent{
    Actors: []*Actor,
    World:  *World,
}
```

### Round Events

```go
// Round progression
core.RoundEvent{
    Round:  int,
    Actors: []*Actor,
}
```

### Turn Events

```go
// Turn progression
core.TurnEvent{
    Turn:  int,
    Actor: *Actor,
}
```

### Combat Events

```go
// Attack events
core.AttackEvent{
    Attacker: *Actor,
    Target:   *Actor,
    Damage:   int,
}

// Damage events
core.DamageEvent{
    Target:    *Actor,
    Amount:    int,
    NewHealth: int,
}

// Death events
core.DeathEvent{
    Target: *Actor,
}
```

## Actor System Integration

### Combat Actions

```go
// Actor performs attack
attacker.Attack(target)
// Emits: AttackEvent (start), DamageEvent, AttackEvent (end)

// Actor takes damage
actor.TakeDamage(amount)
// Emits: DamageEvent, potentially DeathEvent

// Check actor state
isDead := actor.IsDead()
isAlive := !actor.IsDead()
```

### Team Management

```go
// Team-based victory conditions
type TeamID string

const (
    PlayersTeam TeamID = "Players"
    EnemiesTeam TeamID = "Enemies"
)

// Check winner
winner, hasWinner := encounter.Winner()
if hasWinner {
    fmt.Printf("Team %s wins!\n", winner)
}
```

## AI System Integration

```go
// AI brain interface
type Brain interface {
    Play(encounter *core.Encounter)
}

// AI decision making
type SimpleBrain struct {
    actor *core.Actor
}

func (b *SimpleBrain) Play(encounter *core.Encounter) {
    // Select target
    target := b.selectTarget(encounter)
    
    // Execute action
    b.actor.Attack(target)
    
    // End turn
    encounter.EndTurn()
}
```

## Usage Patterns

### Complete Combat Example

```go
// Setup
dispatcher := eventbus.NewDispatcher()
world := &core.World{Width: 10, Height: 10}

// Create actors
player := core.NewActor(dispatcher, world, "Hero", core.PlayersTeam, 15, 15, 6)
enemy := core.NewActor(dispatcher, world, "Goblin", core.EnemiesTeam, 8, 8, 4)

// Create AI
playerAI := ai.NewSimpleBrain(player)
enemyAI := ai.NewSimpleBrain(enemy)
aiMap := map[*core.Actor]*ai.Brain{
    player: playerAI,
    enemy:  enemyAI,
}

// Create encounter
encounter := core.NewEncounter(dispatcher, world, []*core.Actor{player, enemy})

// Event logging
dispatcher.Subscribe(func(event eventbus.Event) {
    prettyprint.Print(os.Stdout, event)
})

// Execute combat
encounter.Start()
for !encounter.IsOver() {
    actorAI := aiMap[encounter.ActiveActor()]
    actorAI.Play(encounter)
}
encounter.End()

// Check results
winner, _ := encounter.Winner()
fmt.Printf("Winner: %s\n", winner)
```

### Event Subscription for UI Updates

```go
// Subscribe to specific events for UI updates
dispatcher.Subscribe(func(event eventbus.Event) {
    switch event.Kind.AsString() {
    case core.DamageType.AsString():
        if !event.End {
            data := event.Data.(core.DamageEvent)
            updateHealthBar(data.Target, data.NewHealth)
        }
    case core.DeathType.AsString():
        if !event.End {
            data := event.Data.(core.DeathEvent)
            showDeathAnimation(data.Target)
        }
    }
})
```

## Design Rationale

The Encounter System implements several key design decisions:

1. **Event-Driven Architecture**: All state changes emit events for loose coupling
2. **Turn-Based Coordination**: Clear turn/round progression with automatic management
3. **Team-Based Victory**: Flexible team system supporting multiplayer scenarios
4. **AI Integration**: Clean interface for pluggable AI decision making
5. **Hierarchical Events**: Nested event structure (encounter → round → turn → action)
6. **Immutable Event Data**: Events carry complete state snapshots

This design enables complete combat orchestration while maintaining clean separation between game logic, AI decision making, and presentation layers.
# Effect System

## Overview

The Effect System is the primary mechanism for implementing game rules, status effects, and conditional logic within the Anvil engine. It provides a robust, event-driven framework where effects can listen for specific game states and modify them in a predictable, prioritized order. This system is the foundation for everything from simple attribute bonuses to complex, multi-stage abilities like Undead Fortitude or Attacks of Opportunity.

The core philosophy is to decouple rule logic from the core game systems. Instead of hardcoding checks like `if (target.hasCondition("poisoned"))`, the system evaluates a state object (e.g., `PreAttackRoll`) against an actor's active effects, allowing each effect to apply its logic independently.

## Core Components

### Effect

The fundamental building block of the system. An `Effect` is a container for handlers that react to specific game states.

```go
type Effect struct {
	Archetype string
	ID        string
	Name      string
	Handlers  Handlers // Map of event names to handler functions
	Priority  Priority // Determines execution order
}
```

### EffectContainer

A collection of effects attached to a game entity (like an Actor). It manages adding, removing, and evaluating effects in the correct priority order.

```go
type EffectContainer struct {
	effects []*Effect
}
```

### State Objects

Plain Go structs that represent a specific point-in-time event or calculation in the game. These are the "events" that effects listen for. They are defined in `internal/core/effect_state.go`.

```go
// Example: State object for calculations before an attack roll
type PreAttackRoll struct {
	Source     *Actor
	Target     *Actor
	Expression *expression.Expression
	Tags       tag.Container
}

// Example: State object for when a turn starts
type TurnStarted struct {
	Source *Actor
}
```

### Priority

An integer value that determines the execution order of effects. Lower numbers run earlier. This ensures that foundational effects (like calculating base attributes) run before multiplicative or conditional effects.

```go
const (
	PriorityNormal       Priority = iota // 0
	PriorityEarly        Priority = -20
	PriorityBase         Priority = -60
	PriorityBaseOverride Priority = -40
	PriorityLate         Priority = 20
	PriorityLast         Priority = 40
)
```

## Evaluation Flow

1.  **State Creation**: A core system (e.g., combat logic) creates a state object, like `PreAttackRoll{...}`.
2.  **Evaluation**: The system passes this state object to an entity's `EffectContainer.Evaluate()`.
3.  **Prioritized Iteration**: The container iterates through its sorted list of effects.
4.  **Handler Matching**: Each effect checks if it has a handler registered for the type of the state object (e.g., a handler for `*PreAttackRoll`).
5.  **Execution**: If a handler exists, it's executed, potentially modifying the state object.
6.  **Completion**: After all effects have been evaluated, the core system uses the (potentially modified) state object to continue its logic.

## API Reference

### EffectContainer

```go
// Add effects to the container. They will be sorted automatically by priority.
func (c *EffectContainer) Add(effect ...*Effect)

// Remove an effect from the container.
func (c *EffectContainer) Remove(effect *Effect)

// Evaluate a state object against all contained effects.
func (c *EffectContainer) Evaluate(state any)
```

### Effect Creation and Handler Registration

The primary way to create an effect is to define a factory function that returns a new `*Effect` and registers its handlers.

The `On()` method provides a type-safe, reflection-based way to register a handler. It infers the state object type from the handler function's parameter.

```go
// Registers a handler. The event name is derived from the state type.
func (e *Effect) On(handler any)

// Example of a handler function signature
func onPreAttack(state *core.PreAttackRoll) {
    // logic to modify the attack roll expression
    state.Expression.AddConstant(2, "Bless")
}

// Registering the handler
myEffect.On(onPreAttack)
```

## Usage Patterns

### Example 1: Simple Attribute Modifier (Bless Spell)

This effect grants a +2 bonus on all attack rolls.

**Effect Definition (`/internal/ruleset/effects/basic/effect_bless.go`):**
```go
func NewBless() *core.Effect {
	effect := &core.Effect{
		Name:     "Bless",
		Priority: core.PriorityNormal,
	}

	effect.On(func(state *core.PreAttackRoll) {
		state.Expression.AddConstant(2, "Bless Bonus")
	})

	effect.On(func(state *core.PreSavingThrow) {
		state.Expression.AddConstant(2, "Bless Bonus")
	})

	return effect
}
```

**Applying and Using the Effect:**
```go
// In the game logic...
blessedActor := world.GetActor("Cedric")
blessEffect := effects.NewBless()
blessedActor.Effects.Add(blessEffect)

// Later, during an attack...
attackRollState := &core.PreAttackRoll{
    Source:     blessedActor,
    Target:     world.GetActor("Goblin"),
    Expression: expression.FromD20("Attack"),
}

// The container evaluates the effect, and the handler adds the bonus.
blessedActor.Effects.Evaluate(attackRollState)

// attackRollState.Expression now includes the +2 bonus.
result := attackRollState.Expression.Evaluate()
```

### Example 2: Undead Fortitude

This effect allows a zombie to make a saving throw to avoid being knocked to 0 HP.

**Effect Definition (`/internal/ruleset/effects/shared/effect_undead_fortitude.go`):**
```go
func NewUndeadFortitude() *core.Effect {
    effect := &core.Effect{
        Name:     "Undead Fortitude",
        Priority: core.PriorityLate, // Runs after damage is calculated but before it's applied.
    }

    effect.On(func(state *core.PreTakeDamage) {
        // If damage would kill the actor...
        if state.Source.HP-state.Expression.Value <= 0 {
            // ...make a Constitution saving throw.
            dc := 5 + state.Expression.Value
            savingThrow := expression.FromD20("Undead Fortitude Save")
            savingThrow.AddConstant(state.Source.Attributes.Constitution.Modifier(), "CON Mod")
            
            if savingThrow.Evaluate().Value >= dc {
                // If successful, set the damage to leave the actor at 1 HP.
                finalDamage := state.Source.HP - 1
                state.Expression.ReplaceWith(finalDamage, "Undead Fortitude")
            }
        }
    })

    return effect
}
```

## Design Rationale

1.  **Decoupling**: The primary goal is to decouple rule implementation from core game logic. The combat system doesn't need to know about "Bless" or "Undead Fortitude"; it only needs to `Evaluate` the correct state objects at the correct times.
2.  **Prioritization**: Rules in a TTRPG have a strict order of operations (e.g., base values are determined before multipliers). The `Priority` system enforces this, preventing bugs and ensuring consistent calculations.
3.  **Type Safety**: Using `Effect.On()` and specific state structs allows for compile-time checks (via `go vet`) and reduces runtime errors. Handlers are strongly typed, so there's no need for type assertions on the event payload.
4.  **Discoverability**: All possible hooks for effects are explicitly defined as structs in `effect_state.go`. This serves as a clear, centralized registry of all extension points in the game logic.
5.  **Performance**: While reflection has a reputation for being slow, it is only used once during handler registration. The critical path—evaluation—involves a simple map lookup by string, which is highly performant.

This design creates a system that is scalable, maintainable, and easy to test, as each effect can be tested in isolation by simply feeding it the state objects it's designed to handle.
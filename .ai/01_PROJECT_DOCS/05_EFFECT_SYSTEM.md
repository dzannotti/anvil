# Effect System

*Status: Implemented*

## Overview

The Effect System is the **internal rule resolution engine** for D&D 2024 mechanics. It provides event-driven callbacks that modify calculations, apply conditions, and trigger rule-based behaviors. Effects are the mechanism by which spells, feats, class features, magical items, and monster abilities interact with the core game systems.

**Core Philosophy: "Rules are Effects, Effects are Rules"**

Every D&D rule that modifies gameplay - from basic attribute modifiers to complex spell interactions - is implemented as an Effect. This creates a unified, extensible system for rule implementation that maintains D&D's layered complexity while providing clean code organization.

## Core Components

### Effect Structure

```go
type Effect struct {
    Archetype string              // Effect type identifier (e.g., "fighting-style-defense")
    ID        string              // Unique instance identifier 
    Name      string              // Human-readable name
    Handlers  map[string]func()   // Event type -> handler function mapping
    Priority  Priority            // Execution order relative to other effects
}
```

### Priority System

Effects execute in strict priority order to ensure correct rule layering:

```go
const (
    PriorityBase         Priority = -60  // Core rules (attribute modifiers)
    PriorityBaseOverride Priority = -40  // Base rule overrides
    PriorityEarly        Priority = -20  // Early modifications
    PriorityNormal       Priority = 0    // Standard effects
    PriorityLate         Priority = 20   // Late modifications (critical hits)
    PriorityLast         Priority = 40   // Final adjustments
)
```

### Event-Driven Architecture

Effects respond to specific game events using reflection-based handler registration:

```go
// Effects register handlers for specific event types
effect.On(func(s *core.PreAttackRoll) {
    // Modify attack roll calculation
    s.Expression.AddConstant(bonus, "Effect Name")
})
```

## API Reference

### Effect Creation

```go
// Create effect with priority
effect := &core.Effect{
    Name:     "Fighting Style: Defense",
    Priority: core.PriorityNormal,
}

// Register event handlers using reflection
effect.On(func(s *core.AttributeCalculation) {
    if s.Attribute.MatchExact(tags.ArmorClass) {
        s.Expression.AddConstant(1, "Defense Fighting Style")
    }
})
```

### Handler Registration

The `On()` method uses reflection to automatically map handler functions to event types:

```go
// Handler parameter type determines which events trigger it
effect.On(func(s *core.PreDamageRoll) {
    // Automatically called for PreDamageRoll events
})

effect.On(func(s *core.TurnStarted) {
    // Automatically called for TurnStarted events  
})
```

### Effect Container Management

```go
// Add effects to actors (automatically sorted by priority)
actor.AddEffect(attributeModifierEffect, proficiencyEffect)

// Remove specific effects
actor.RemoveEffect(temporaryEffect)

// Effects evaluate in priority order
actor.Evaluate(&core.PreAttackRoll{...})
```

## Event Types and States

The system defines comprehensive event types for all game mechanics:

### Combat Events

```go
// Attack roll modification
type PreAttackRoll struct {
    Source     *Actor              // Attacker
    Target     *Actor              // Target being attacked
    Expression *expression.Expression // Attack roll calculation
    Tags       tag.Container       // Attack properties (melee, ranged, etc.)
}

type PostAttackRoll struct {
    Source *Actor              // Attacker
    Target *Actor              // Target
    Result *expression.Expression // Final attack roll result
    Tags   tag.Container       // Attack properties
}
```

### Damage Events

```go
// Damage calculation modification
type PreDamageRoll struct {
    Expression *expression.Expression // Damage calculation
    Source     *Actor              // Damage source
    Tags       tag.Container       // Damage types and properties
}

type PreTakeDamage struct {
    Expression *expression.Expression // Incoming damage
    Source     *Actor              // Target taking damage
}
```

### Character Events

```go
// Attribute calculation (AC, saves, etc.)
type AttributeCalculation struct {
    Source     *Actor              // Actor whose attribute is calculated
    Expression *expression.Expression // Attribute calculation
    Attribute  tag.Tag             // Which attribute (AC, Strength, etc.)
}

// Saving throw modification
type PreSavingThrow struct {
    Expression      *expression.Expression // Save calculation
    Source          *Actor              // Actor making save
    Attribute       tag.Tag             // Save type (Dexterity, Wisdom, etc.)
    DifficultyClass int                 // Target DC
}
```

### Condition and Turn Events

```go
// Condition changes (poisoned, charmed, etc.)
type ConditionChanged struct {
    Source    *Actor    // Actor affected
    Condition tag.Tag   // Condition type
    From      *Effect   // Effect that caused the change
}

// Turn lifecycle
type TurnStarted struct {
    Source *Actor // Actor whose turn started
}

type TurnEnded struct {
    Source *Actor // Actor whose turn ended
}
```

### Movement Events

```go
// Movement validation and modification
type PreMoveStep struct {
    Source  *Actor        // Actor moving
    Action  Action        // Movement action
    From    grid.Position // Starting position
    To      grid.Position // Destination position
    CanMove bool          // Whether movement is allowed
}
```

## Effect Implementation Patterns

### Basic Attribute Modifiers

```go
func NewAttributeModifierEffect() *core.Effect {
    fx := &core.Effect{Name: "Attribute Modifier", Priority: core.PriorityBase}
    
    // Apply STR/DEX to attack rolls
    fx.On(func(s *core.PreAttackRoll) {
        if s.Tags.HasTag(tags.Finesse) || s.Tags.HasTag(tags.Ranged) {
            dexMod := stats.AttributeModifier(s.Source.Attribute(tags.Dexterity).Value)
            s.Expression.AddConstant(dexMod, "Dexterity Modifier")
        } else {
            strMod := stats.AttributeModifier(s.Source.Attribute(tags.Strength).Value)
            s.Expression.AddConstant(strMod, "Strength Modifier")
        }
    })
    
    return fx
}
```

### Conditional Effects

```go
func NewFightingStyleDefense() *core.Effect {
    fx := &core.Effect{Name: "Fighting Style: Defense"}
    
    fx.On(func(s *core.AttributeCalculation) {
        // Only apply to AC calculations
        if !s.Attribute.MatchExact(tags.ArmorClass) {
            return
        }
        
        // Only if wearing armor
        armorTypes := tag.NewContainer(tags.LightArmor, tags.MediumArmor, tags.HeavyArmor)
        for _, item := range s.Source.Equipped {
            if item.Tags().HasAny(armorTypes) {
                s.Expression.AddConstant(1, fx.Name)
                return
            }
        }
    })
    
    return fx
}
```

### Critical Hit Effects

```go
func NewCritEffect() *core.Effect {
    fx := &core.Effect{Name: "Critical Hit", Priority: core.PriorityLate}
    
    fx.On(func(s *core.PreDamageRoll) {
        if s.Expression.IsCriticalSuccess() {
            s.Expression.DoubleDice("Critical Hit")
        }
    })
    
    return fx
}
```

### Turn-Based Effects

```go
func NewConcentrationEffect() *core.Effect {
    fx := &core.Effect{Name: "Concentration"}
    
    fx.On(func(s *core.TurnStarted) {
        // Check concentration at start of turn
        if s.Source.HasCondition(tags.Concentrating) {
            // Trigger concentration check if needed
        }
    })
    
    fx.On(func(s *core.PostTakeDamage) {
        // Constitution save when taking damage
        if s.Source.HasCondition(tags.Concentrating) {
            dc := max(10, s.ActualDamage/2)
            // Trigger concentration save
        }
    })
    
    return fx
}
```

## Actor Integration

### Default Effects

All actors receive core effects automatically through the registry:

```go
// Universal effects (applied to all actors)
actor.AddEffect(registry.NewEffect("attribute-modifier", nil))
actor.AddEffect(registry.NewEffect("proficiency-modifier", nil))
actor.AddEffect(registry.NewEffect("critical", nil))
actor.AddEffect(registry.NewEffect("attack-of-opportunity", nil))

// Conditional effects based on creature type
if !actor.HasCondition(tags.Undead) {
    actor.AddEffect(registry.NewEffect("death-saving-throw", nil))
}
actor.AddEffect(registry.NewEffect("death", nil))
```

### Manual Effect Addition

```go
// Class features
fighter.AddEffect(registry.NewEffect("fighting-style-defense", nil))

// Spell effects
target.AddEffect(registry.NewEffect("bless", nil))

// Magic item effects
actor.Equip(enchantedWeapon) // Items can add effects via OnEquip()
```

### Effect Evaluation

Effects are automatically triggered through the actor's evaluation system:

```go
// Actor methods trigger effect evaluation
func (a *Actor) StartTurn() {
    a.Resources.Reset()
    a.Evaluate(&TurnStarted{Source: a})  // Triggers turn start effects
}

func (a *Actor) AddCondition(condition tag.Tag, source *Effect) {
    a.Conditions.Add(condition, source)
    a.Evaluate(&ConditionChanged{...})   // Triggers condition effects
}
```

## Registry Integration

Effects are registered as factories in the central registry:

```go
// Basic effects available to all actors
registry.RegisterEffect("attribute-modifier", func(_ map[string]interface{}) *core.Effect {
    return effectsBasic.NewAttributeModifierEffect()
})

registry.RegisterEffect("critical", func(_ map[string]interface{}) *core.Effect {
    return effectsBasic.NewCritEffect()
})

// Class-specific effects
registry.RegisterEffect("fighting-style-defense", func(_ map[string]interface{}) *core.Effect {
    return effectsFighter.NewFightingStyleDefense()
})

// Race-specific effects  
registry.RegisterEffect("undead-fortitude", func(_ map[string]interface{}) *core.Effect {
    return effectsShared.NewUndeadFortitudeEffect()
})
```

## Usage Patterns

### When to Create Effects

**✅ Create Effects for:**
- Class features and feats
- Spell effects and magical abilities
- Racial traits and monster abilities
- Magic item enchantments
- Environmental conditions
- Rule modifications and exceptions

**❌ Don't Create Effects for:**
- One-time action resolutions
- Simple data lookups
- UI state changes
- Logging or metrics

### Effect Granularity

**Good Granularity:**
```go
// One effect per distinct game feature
NewBlessEffect()           // +1d4 to attack rolls and saves
NewShieldOfFaithEffect()   // +2 AC
NewRageEffect()            // Damage bonus, advantage, resistance
```

**Poor Granularity:**
```go
// Too specific
NewBlessAttackBonusEffect()    // ❌ Part of bless
NewBlessSaveBonusEffect()      // ❌ Part of bless

// Too broad  
NewSpellEffect()               // ❌ Not specific enough
```

### Priority Guidelines

- **PriorityBase**: Core D&D rules (attribute modifiers, proficiency)
- **PriorityNormal**: Most class features, spells, items
- **PriorityLate**: Critical hits, damage doubling
- **PriorityLast**: Final overrides, special exceptions

## Performance Characteristics

- **Handler Registration**: O(1) using reflection-based type mapping
- **Effect Evaluation**: O(n) where n is number of effects on actor
- **Priority Sorting**: O(n log n) when effects are added
- **Memory**: Minimal overhead, handlers stored as function pointers
- **Concurrency**: Thread-safe evaluation using goroutines and wait groups

## Design Rationale

### Event-Driven Architecture
- **Decoupling**: Effects don't know about each other, only events
- **Extensibility**: New effects just register for relevant events
- **Debugging**: Clear audit trail of which effects triggered
- **Testing**: Easy to test effects in isolation

### Priority-Based Execution
- **Rule Layering**: Matches D&D's rule priority system
- **Predictability**: Consistent execution order
- **Conflict Resolution**: Clear precedence for conflicting effects
- **Extensibility**: New priorities can be added without breaking existing effects

### Reflection-Based Handlers
- **Type Safety**: Compile-time checking of handler signatures
- **Clean API**: No manual event type registration
- **Performance**: Reflection cost paid once at registration
- **Maintainability**: Handler changes don't require registration updates

### Immutable Event Data
- **Safety**: Effects can't interfere with each other's data
- **Debugging**: Event state is consistent throughout evaluation
- **Testing**: Predictable behavior for test scenarios
- **Concurrency**: Safe for parallel effect evaluation

## Integration with Other Systems

### Expression System Integration
Effects modify expressions by adding components:
```go
// Effects add components to expressions
s.Expression.AddConstant(bonus, "Effect Name")
s.Expression.AddDice(1, 4, "Bless")
```

### Tag System Integration
Effects use tags for conditional logic:
```go
// Check attack properties
if s.Tags.HasTag(tags.Ranged) {
    // Apply ranged attack bonuses
}

// Check damage types
if s.Tags.HasTag(tags.Fire) {
    // Apply fire resistance/vulnerability
}
```

### Event Bus Integration
Effects operate on internal events, while the Event Bus handles external communication:
```go
// Internal effect evaluation (not on event bus)
actor.Evaluate(&PreAttackRoll{...})

// External event notification (on event bus)  
actor.Dispatcher.Emit(AttackRollEvent{...})
```

This creates clean separation between internal rule resolution and external system communication.

## File Locations

- **Core System**: `/internal/core/effect.go`, `effect_container.go`, `effect_state.go`
- **Basic Effects**: `/internal/ruleset/effects/basic/`
- **Class Effects**: `/internal/ruleset/effects/classes/`
- **Shared Effects**: `/internal/ruleset/effects/shared/`
- **Registry**: `/internal/ruleset/seed.go` (effect registration)
- **Tests**: `/internal/core/effect_test.go`, `effect_container_test.go`

The Effect System provides the foundation for implementing D&D 2024's complex rule interactions while maintaining clean, testable, and extensible code architecture.

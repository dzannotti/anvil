# Expression System

## Overview

The Expression System is a sophisticated calculation engine that maintains complete audit trails through recursive component structures. It provides formula-driven calculations with deferred evaluation, supporting complex D&D mechanics like dice rolling, damage calculation, advantage/disadvantage, and critical hits.

## Core Architecture

### Expression Structure

```go
type Expression struct {
    Components []Component  // List of calculation components
    Value      int         // Evaluated result
    IsCritical int         // Critical hit state (1=success, -1=failure, 0=normal)
    rng        Roller      // Random number generator
}
```

### Component Structure

```go
type Component struct {
    Kind          ComponentKind    // Evaluation behavior (strategy pattern)
    Source        string          // Human-readable source description
    Type          Kind            // Component type identifier
    Tags          tag.Container   // Associated tags for rule matching
    Value         int             // Evaluated result
    SubComponents []Component     // Nested components for complex calculations
    Values        []int           // Individual roll values (for dice)
}
```

### ComponentKind Interface

The system uses the Strategy pattern for different calculation behaviors:

```go
type ComponentKind interface {
    Evaluate(ctx EvaluationContext) int
    Expected() int
    Clone() ComponentKind
    Values() []int
}
```

## Component Types

### Basic Components

#### Constant Components
```go
// Fixed value components
expr := expression.New()
expr.AddConstant(5, "Proficiency bonus")
// Result: 5 (always)
```

#### Dice Components
```go
// Standard dice rolling
expr := expression.New()
expr.AddDice(2, 6, "Short sword damage") // 2d6
// Result: 2-12 (random)
```

#### D20 Components
```go
// Special D20 handling with advantage/disadvantage
expr := expression.New()
expr.AddD20("Attack roll")
// Result: 1-20 (random)
```

### Damage Components

#### Damage Constants
```go
// Damage with type tags
tags := tag.ContainerFromString("damage.fire")
expr := expression.New()
expr.AddDamageConstant(3, "Fire damage bonus", tags)
```

#### Damage Dice
```go
// Dice damage with type tags
tags := tag.ContainerFromString("damage.fire,damage.elemental")
expr := expression.New()
expr.AddDamageDice(2, 6, "Fireball damage", tags)
```

### Modifier Components

#### Halved Components
```go
// Halve specific damage types
expr := expression.New()
expr.AddDamageDice(4, 6, "Fire damage", tag.ContainerFromString("damage.fire"))
expr.HalveDamage(tag.New("damage.fire"))
// Fire damage is halved
```

#### Max Dice Components
```go
// Maximize dice (typically for critical hits)
expr := expression.New()
expr.AddDice(2, 6, "Sword damage")
expr.MaxDice(tag.New("damage.weapon"))
// Adds maximum possible dice value
```

#### Replaced Components
```go
// Replace entire expression with constant
expr := expression.New()
expr.AddDice(1, 20, "Attack roll")
expr.ReplaceWith(15, "Reliable talent")
// Result: 15 (always)
```

## Advanced Features

### Advantage and Disadvantage

```go
// D20 advantage/disadvantage system
expr := expression.New()
expr.AddD20("Attack roll")

// Apply advantage (roll twice, take higher)
exprAdv := expr.WithAdvantage()

// Apply disadvantage (roll twice, take lower)
exprDisadv := expr.WithDisadvantage()

// Advantage and disadvantage cancel out
exprNormal := expr.WithAdvantage().WithDisadvantage()
```

### Critical Hit Detection

```go
// Automatic critical detection
expr := expression.New()
expr.AddD20("Attack roll")
result := expr.Evaluate()

// Check critical state
if expr.IsCritical == 1 {
    // Critical success (natural 20)
} else if expr.IsCritical == -1 {
    // Critical failure (natural 1)
}
```

### Expression Grouping

```go
// Group components by damage type
expr := expression.New()
expr.AddDamageDice(2, 6, "Weapon damage", tag.ContainerFromString("damage.slashing"))
expr.AddDamageConstant(3, "Strength bonus", tag.ContainerFromString("damage.slashing"))
expr.AddDamageDice(1, 4, "Frost damage", tag.ContainerFromString("damage.cold"))

grouped := expr.EvaluateGrouped()
// Returns map[string]Expression with damage types as keys
```

### Double Dice (Critical Hits)

```go
// Double dice components for critical hits
expr := expression.New()
expr.AddDice(1, 8, "Longsword damage")
expr.DoubleDice(tag.New("damage.weapon"))
// Adds duplicate dice components
```

## Evaluation Context

### Custom Rollers

```go
// Custom random number generation
type TestRoller struct {
    values []int
    index  int
}

func (r *TestRoller) Roll(sides int) int {
    value := r.values[r.index]
    r.index++
    return value
}

expr := expression.FromD20("Attack roll")
expr.SetRoller(&TestRoller{values: []int{20}}) // Always roll 20
```

### Evaluation Process

```go
// Evaluation creates complete audit trail
expr := expression.New()
expr.AddDice(2, 6, "Damage roll")
expr.AddConstant(3, "Strength bonus")

result := expr.Evaluate()
// result.Value contains total
// result.Components contains detailed breakdown
```

## Factory Functions

### Expression Creation

```go
// Factory functions for common patterns
attackRoll := expression.FromD20("Attack roll")
damage := expression.FromDice(1, 8, "Sword damage")
constant := expression.FromConstant(5, "Proficiency bonus")

// Damage with tags
fireDamage := expression.FromDamageDice(3, 6, "Fireball", tag.ContainerFromString("damage.fire"))
```

### Component Factories

```go
// Direct component creation
constantComp := expression.NewConstantComponent(5, "Bonus")
diceComp := expression.NewDiceComponent(2, 6, "Damage")
d20Comp := expression.NewD20Component("Attack")
```

## Usage Patterns

### Attack Roll Calculation

```go
// Complete attack roll with proficiency and modifiers
attack := expression.New()
attack.AddD20("Attack roll")
attack.AddConstant(5, "Proficiency bonus")
attack.AddConstant(3, "Strength modifier")

// With advantage
attackAdv := attack.WithAdvantage()

result := attackAdv.Evaluate()
if result.IsCritical == 1 {
    // Critical hit!
}
```

### Damage Calculation

```go
// Weapon damage with multiple types
damage := expression.New()
damage.AddDamageDice(1, 8, "Longsword", tag.ContainerFromString("damage.slashing"))
damage.AddConstant(3, "Strength modifier", tag.ContainerFromString("damage.slashing"))
damage.AddDamageDice(1, 6, "Frost enchantment", tag.ContainerFromString("damage.cold"))

// For critical hits
critDamage := damage.Clone()
critDamage.DoubleDice(tag.New("damage.weapon"))
```

### Saving Throw with Modifiers

```go
// Saving throw with various modifiers
save := expression.New()
save.AddD20("Constitution save")
save.AddConstant(2, "Constitution modifier")
save.AddConstant(5, "Proficiency bonus")

// Apply disadvantage if poisoned
if isPoisoned {
    save = save.WithDisadvantage()
}
```

### Spell Damage with Resistance

```go
// Spell damage with potential resistance
spell := expression.New()
spell.AddDamageDice(8, 6, "Fireball", tag.ContainerFromString("damage.fire"))

// Apply resistance (half damage)
if hasFireResistance {
    spell.HalveDamage(tag.New("damage.fire"))
}
```

## Audit Trail System

### Component Hierarchy

```go
// Nested component structure maintains full audit trail
expr := expression.New()
expr.AddDice(2, 6, "Base damage")
expr.AddConstant(3, "Strength bonus")
expr.MaxDice(tag.New("damage.weapon")) // Critical hit

result := expr.Evaluate()
// result.Components contains:
// - Original 2d6 component with individual roll values
// - Strength bonus component
// - Max dice component showing maximum added
```

### Evaluation Result

```go
type EvaluationResult struct {
    Value      int         // Total result
    Components []Component // Detailed breakdown
}

// Each component tracks:
// - Source description
// - Individual values (for dice)
// - Sub-components for complex calculations
// - Tags for rule matching
```

## Performance Characteristics

- **Evaluation**: O(n) where n is number of components
- **Cloning**: O(n) deep copy of all components
- **Grouping**: O(n) with tag-based sorting
- **Memory**: Minimal overhead with component reuse

## Design Rationale

The Expression System implements several key design decisions:

1. **Deferred Evaluation**: Expressions are formulas until explicitly evaluated
2. **Complete Audit Trail**: Every calculation step is preserved
3. **Strategy Pattern**: ComponentKind interface enables extensible behaviors
4. **Immutable Results**: Evaluation creates new data, doesn't modify expressions
5. **Tag Integration**: Deep integration with tag system for rule matching
6. **Composable Operations**: Methods return new expressions for chaining

This design enables complex D&D calculations while maintaining full traceability and extensibility.
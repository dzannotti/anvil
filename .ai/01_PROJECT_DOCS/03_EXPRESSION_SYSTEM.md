# Expression System

## Overview

The Expression System is a specialized calculation engine for D&D 2024 mechanics, handling dice rolls, damage calculations, advantage/disadvantage, and critical hits. It maintains complete audit trails through hierarchical component structures while providing a clean, type-safe API for game mechanics.

## Core Architecture

### Expression Structure

```go
type Expression struct {
    Components []Component  // List of calculation components
    Value      int         // Total evaluated result
    Rng        DiceRoller  // Random number generator interface
}
```

### Component Structure

```go
type Component struct {
    Type            tag.Tag       // Component type identifier
    Value           int          // Final calculated value
    Source          string       // Human-readable audit
    Values          []int        // Individual roll results (for dice)
    Times           int          // Number of dice/multiplier
    Sides           int          // Die sides
    HasAdvantage    []string     // Sources granting advantage
    HasDisadvantage []string      // Sources imposing disadvantage
    Tags            tag.Container // Damage types and metadata
    Components      []Component   // Nested components for complex calculations
    IsCritical      int           // Critical state (CriticalSuccess=1, CriticalFailure=-1, CriticalNone=0)
}
```

### Component Types

The system uses hierarchical tag-based component types:

```go
// Basic types
Constant       = tag.FromString("Component.Type.Constant")
Dice           = tag.FromString("Component.Type.Dice")
D20            = tag.FromString("Component.Type.Dice.D20")

// Damage-specific types (extend base types)
DamageConstant = tag.FromString("Component.Type.Constant.Damage") 
DamageDice     = tag.FromString("Component.Type.Dice.Damage")
```

### Critical Hit Constants

```go
const (
    CriticalSuccess = 1   // Natural 20 or forced critical
    CriticalFailure = -1  // Natural 1 or forced failure
    CriticalNone    = 0   // Normal result
)
```

## Evaluation Architecture

### Root-Level Evaluation Only

**Critical Design Principle**: The expression system evaluates **only** the root-level components in `Expression.Components[]`. Nested components in `Component.Components[]` are **never evaluated or re-evaluated** - they exist purely for audit trail purposes.

```go
expression := FromDice(2, 6, "Sword Damage")
expression.HalveDamage(tags.Slashing, "Resistance")

// After HalveDamage:
// Root level: [Component{Type: Constant, Value: 6, Source: "Halved (Resistance) Sword Damage"}]
// Nested:     Component.Components = [Component{Type: Dice, Values: [4,8], Source: "Sword Damage"}]
//
// Only the root Constant component contributes to Expression.Value
// The nested Dice component preserves the original roll history
```

### Audit Trail Preservation

Methods that transform expressions preserve calculation history by moving original components into nested structures:

```go
// ReplaceWith example
spell := FromDice(1, 20, "Spell Attack")
spell.ReplaceWith(15, "Reliable Talent")

// Result:
// Root: [Component{Type: Constant, Value: 15, Source: "Reliable Talent"}]
// Nested: Component.Components = [Component{Type: D20, Value: 12, Source: "Spell Attack"}]
```

**Methods Leveraging This Architecture:**
- **`HalveDamage()`**: Creates constant with halved value, preserves original in nested
- **`ReplaceWith()`**: Replaces all components with constant, preserves originals 
- **`DoubleDice()`**: Duplicates dice components, maintains separate audit trails
- **`MaxDice()`**: Adds maximized constants, preserves original dice

### Primary Damage Collation

In damage expressions, the first component establishes the "primary" damage type. Components added with empty tags automatically inherit this primary type, enabling clean ability modifier categorization:

```go
// Flaming longsword: 1d8 slashing + 2d4 fire + STR modifier
damage := FromDamageDice(1, 8, "Longsword", tag.NewContainer(tags.Slashing))     // Primary damage
damage.AddDamageDice(2, 4, "Flame Tongue", tag.NewContainer(tags.Fire))         // Explicit fire
damage.AddDamageConstant(3, "STR Modifier", tag.NewContainer())                 // Gets slashing (primary)

// Before grouping:
// Components[0]: 1d8 slashing (primary)
// Components[1]: 2d4 fire (explicit)  
// Components[2]: +3 slashing (inherited from primary)

grouped := damage.EvaluateGroup()
// Results in two separate root components:
// grouped.Components[0]: DamageConstant{Value: 11, Tags: slashing} // 1d8(8) + 3 STR = 11 slashing
// grouped.Components[1]: DamageConstant{Value: 6, Tags: fire}      // 2d4(6) = 6 fire
//
// Each root component contains the total for that damage type
// Original components are preserved in Component.Components[] for audit trail
```

**Key Behavior**: Ability modifiers, enhancement bonuses, and other non-typed damage automatically attach to the weapon's primary damage type, which is exactly what D&D rules expect.

The `primaryTags()` method handles this inheritance:
- If input tags are empty or contain "primary", returns `Components[0].Tags`
- Otherwise returns the provided tags unchanged
- This ensures modifiers group with the main weapon damage while preserving explicit damage types

## Public API

### Factory Functions

Create expressions with specific component types:

```go
// Basic expressions
FromConstant(value int, source string, components ...Component) Expression
FromDice(times int, sides int, source string, components ...Component) Expression
FromD20(source string, components ...Component) Expression

// Damage expressions with tag support
FromDamageConstant(value int, source string, tags tag.Container, components ...Component) Expression
FromDamageDice(times int, sides int, source string, tags tag.Container, components ...Component) Expression
FromDamageResult(res Expression) Expression // Create copy for damage
```

### Component Addition Methods

Add components to existing expressions:

```go
// Basic components
AddConstant(value int, source string, components ...Component)
AddDice(times int, sides int, source string, components ...Component)
AddD20(source string, components ...Component)

// Damage components with tags
AddDamageConstant(value int, source string, tags tag.Container, components ...Component)
AddDamageDice(times int, sides int, source string, tags tag.Container, components ...Component)

// D20 modifiers (only works with D20 components)
WithAdvantage(source string)
WithDisadvantage(source string)
```

### Evaluation Methods

```go
Evaluate() *Expression           // Calculate final result
EvaluateGroup() *Expression      // Group by tags, then evaluate
Clone() Expression               // Create deep copy
```

### Critical Hit Detection

```go
IsCriticalSuccess() bool                    // Check for critical success
IsCriticalFailure() bool                    // Check for critical failure
SetCriticalSuccess(source string)          // Force critical success
SetCriticalFailure(source string)          // Force critical failure
```

### Damage Manipulation

```go
HalveDamage(tag tag.Tag, source string)    // Apply resistance/half damage
DoubleDice(source string)                  // Double dice (critical hits)
MaxDice(source string)                     // Maximize dice (critical hits)
ReplaceWith(value int, source string)      // Replace with constant value (items or effects that reset calculation)
HasDamageType(t tag.Tag) bool              // Check for damage type
```

## Usage Examples

### Basic Attack Roll

```go
// Standard attack roll with ability modifier
attack := FromD20("Longsword Attack")
attack.AddConstant(5, "Strength Modifier")
attack.AddConstant(3, "Proficiency Bonus")

result := attack.Evaluate()
fmt.Printf("Attack roll: %d\n", result.Value)

// Check for critical
if result.IsCriticalSuccess() {
    fmt.Println("Critical hit!")
}
```

### Advantage and Disadvantage

```go
// Advantage from spell
attack := FromD20("Attack with Bless")
attack.WithAdvantage("Bless Spell")

// Disadvantage from condition
attack.WithDisadvantage("Poisoned")
// Advantage and disadvantage cancel out

result := attack.Evaluate()
```

### Damage Calculation

```go
// Multi-type damage expression
damage := FromDamageDice(1, 8, "Longsword", tag.NewContainer(tags.Slashing))
damage.AddDamageConstant(3, "Strength Modifier", tag.NewContainer(tags.Slashing))
damage.AddDamageDice(1, 6, "Flame Tongue", tag.NewContainer(tags.Fire))

result := damage.Evaluate()
fmt.Printf("Total damage: %d\n", result.Value)
```

### Critical Hit Damage

```go
// Standard damage
damage := FromDamageDice(1, 8, "Longsword", tag.NewContainer(tags.Slashing))
damage.AddDamageConstant(3, "Strength Modifier", tag.NewContainer(tags.Slashing))

// Apply critical hit (double dice only)
damage.DoubleDice("Critical Hit")

result := damage.Evaluate()
```

### Damage Resistance

```go
// Spell damage
spell := FromDamageDice(8, 6, "Fireball", tag.NewContainer(tags.Fire))

// Apply fire resistance (half damage)
spell.HalveDamage(tags.Fire, "Fire Resistance")

result := spell.Evaluate()
```

### Damage Grouping

```go
// Complex multi-type damage
damage := FromDamageDice(1, 8, "Weapon", tag.NewContainer(tags.Slashing))
damage.AddDamageConstant(3, "STR Mod", tag.NewContainer(tags.Slashing))
damage.AddDamageDice(1, 6, "Poison", tag.NewContainer(tags.Poison))
damage.AddDamageDice(1, 4, "Cold", tag.NewContainer(tags.Cold))

// Group by damage type
grouped := damage.EvaluateGroup()
// Results in separate damage constants for each type
```

## Advanced Features

### Custom Dice Rollers

```go
type TestRoller struct {
    values []int
    index  int
}

func (r *TestRoller) Roll(sides int) int {
    value := r.values[r.index%len(r.values)]
    r.index++
    return value
}

// Use custom roller for testing
expr := FromDice(2, 6, "Test Roll")
expr.Rng = &TestRoller{values: []int{6, 6}} // Always max
result := expr.Evaluate()
```

### Audit Trail Access

```go
damage := FromDamageDice(2, 6, "Fireball", tag.NewContainer(tags.Fire))
result := damage.Evaluate()

// Access individual components
for i, comp := range result.Components {
    fmt.Printf("Component %d: %s = %d\n", i, comp.Source, comp.Value)
    if len(comp.Values) > 0 {
        fmt.Printf("  Individual rolls: %v\n", comp.Values)
    }
}
```

### Complex Damage Scenarios

```go
// Sneak attack with multiple damage types
attack := FromDamageDice(1, 6, "Shortbow", tag.NewContainer(tags.Piercing))
attack.AddDamageConstant(4, "DEX Modifier", tag.NewContainer(tags.Piercing))
attack.AddDamageDice(3, 6, "Sneak Attack", tag.NewContainer(tags.Piercing))
attack.AddDamageDice(1, 6, "Poison Arrow", tag.NewContainer(tags.Poison))

// Apply resistances selectively
if hasResistance {
    attack.HalveDamage(tags.Poison, "Poison Resistance")
}

result := attack.Evaluate()
```

## Component Type Matching

The system uses hierarchical type matching:

```go
// Type hierarchy
component.Type.Match(Dice)           // Matches: Dice, D20, DamageDice
component.Type.Match(D20)            // Matches: D20 only
component.Type.Match(Constant)       // Matches: Constant, DamageConstant
```

This enables targeted operations:

```go
// DoubleDice only affects dice components
damage.DoubleDice("Critical")  // Doubles Dice, D20, DamageDice

// HalveDamage works on any component with matching tags
damage.HalveDamage(tags.Fire, "Resistance")  // Affects any fire damage
```

## Design Rationale

### Type Safety with Flexibility
- **Tag-based typing**: Compile-time safety with runtime flexibility
- **Hierarchical types**: Specialized behavior while maintaining compatibility
- **No interface overhead**: Direct type checking for performance

### Mutable State for Game Mechanics
- **In-place modification**: Natural for building complex expressions
- **Builder pattern**: Chainable operations for readability
- **Deep cloning**: Safe copying when needed

### Complete Audit Trails
- **Nested components**: Tracks all calculation steps
- **Source tracking**: Human-readable descriptions
- **Value preservation**: Individual dice rolls maintained

### D&D-Specific Optimizations
- **Advantage/disadvantage**: Built-in dual-roll mechanics
- **Critical hit support**: Automatic detection and manual setting
- **Damage type grouping**: Handles resistance/vulnerability correctly
- **Flexible tagging**: Supports complex rule interactions

## Performance Characteristics

- **Evaluation**: O(n) where n is total components (including nested)
- **Cloning**: O(n) deep copy with proper slice handling
- **Grouping**: O(n²) worst case for tag comparison
- **Memory**: Minimal overhead, components reused efficiently

## When to Use What

### Expression Types
- **`FromD20()`**: Attack rolls, saves, ability checks
- **`FromDice()`**: Simple damage, random tables
- **`FromConstant()`**: Fixed bonuses, modifiers
- **`FromDamageDice/Constant()`**: Any damage that needs type tracking

### Modification Methods
- **`WithAdvantage/Disadvantage()`**: Only for D20 components
- **`DoubleDice()`**: Critical hits (affects dice only)
- **`MaxDice()`**: Alternative critical hit mechanic
- **`HalveDamage()`**: Resistance, vulnerabilities
- **`ReplaceWith()`**: Special abilities that override rolls

### Evaluation Methods
- **`Evaluate()`**: Standard calculation
- **`EvaluateGroup()`**: When damage types need separation

This system provides the perfect balance of type safety, performance, and flexibility needed for D&D 2024 mechanics while maintaining complete calculation transparency.
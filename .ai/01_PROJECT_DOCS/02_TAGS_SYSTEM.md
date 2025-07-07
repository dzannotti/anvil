# Tags System

## Overview

The Tags System is a hierarchical string-based identification system inspired by UE5 GameplayTags. It provides robust categorization and querying capabilities for game entities, damage types, spell schools, and rule conditions.

## Core Components

### Tag Structure

Tags use dot-separated hierarchical strings:
- `item.weapon.martial.melee.axe`
- `damage.fire`
- `spell.school.evocation`
- `creature.type.undead`

### Input Normalization

The system provides robust input sanitization:

```go
// All inputs are normalized automatically
tag := tag.New("ITEM.Weapon@#$%Martial  MELEE") 
// Results in: "item.weapon.martial.melee"
```

**Normalization Rules:**
- **Case**: All lowercase conversion
- **Whitespace**: Complete removal (`\s\v` patterns)
- **Special Characters**: Removal of `[@#$%\-^&*]`
- **Non-ASCII**: Removal of unicode characters
- **Multiple Dots**: Consolidation (`...` becomes `.`)
- **Boundary Dots**: Removal of leading/trailing dots

## API Reference

### Tag Creation

```go
// Multiple creation methods
tag1 := tag.New("item.weapon.martial")
tag2 := tag.FromString("item.weapon.martial")

// Both are equivalent
```

### Tag Querying

```go
weapon := tag.New("item.weapon.martial.melee.axe")
item := tag.New("item.weapon.martial")

// Exact matching
weapon.MatchExact(item) // false

// Hierarchical matching (item matches weapon)
weapon.Match(item) // true
item.Match(weapon) // false
```

### Tag Utilities

```go
tag := tag.New("item.weapon.martial")

// String representations
tag.AsString()  // "item.weapon.martial"
tag.AsStrings() // ["item", "weapon", "martial"]

// Validation
tag.IsEmpty()   // false
tag.IsValid()   // true
```

## Container System

### Overview

The Container system manages collections of tags with advanced querying capabilities.

### Container Creation

```go
// Multiple creation patterns
container1 := tag.NewContainer()
container2 := tag.NewContainerFromString("damage.fire,damage.cold")
container3 := tag.ContainerFromString("spell.school.evocation")

// From individual tags
weaponTag := tag.New("item.weapon.martial")
container4 := tag.NewContainerFromTag(weaponTag)
```

### Container Management

```go
container := tag.NewContainer()

// Add tags
container.AddTag(tag.New("damage.fire"))
container.AddTag(tag.New("damage.cold"))

// Remove tags
container.RemoveTag(tag.New("damage.fire"))

// Duplicate prevention is automatic
container.AddTag(tag.New("damage.cold")) // No duplicate added
```

### Container Querying

```go
damageTypes := tag.ContainerFromString("damage.fire,damage.cold,damage.lightning")
queryTag := tag.New("damage.fire")

// Basic presence check
damageTypes.HasTag(queryTag) // true

// Multiple tag queries
fireAndCold := tag.ContainerFromString("damage.fire,damage.cold")
damageTypes.HasAny(fireAndCold) // true
damageTypes.HasAll(fireAndCold) // true
```

### Hierarchical Matching

```go
weaponTags := tag.ContainerFromString("item.weapon.martial.melee.axe,item.armor.heavy")
weaponQuery := tag.New("item.weapon")

// Hierarchical matching
weaponTags.MatchTag(weaponQuery) // true (axe matches weapon)

// Multiple hierarchical queries
queries := tag.ContainerFromString("item.weapon,item.armor")
weaponTags.MatchAny(queries) // true
weaponTags.MatchAll(queries) // true
```

### Container Utilities

```go
container := tag.ContainerFromString("damage.fire,damage.cold")

// Metadata
container.Len()      // 2
container.IsEmpty()  // false
container.AsStrings() // ["damage.fire", "damage.cold"]

// Unique identifier (deterministic, sorted)
container.ID() // "damage.cold.damage.fire"

// Safe copying
clone := container.Clone()
container.Add(otherContainer) // Merge containers
```

## Matching Logic

### Exact vs Hierarchical

```go
weapon := tag.New("item.weapon.martial.melee.axe")
martial := tag.New("item.weapon.martial")
gameplay := tag.New("gameplay.weapon") // Different hierarchy

// Exact matching
weapon.MatchExact(martial) // false
weapon.MatchExact(weapon)  // true

// Hierarchical matching
weapon.Match(martial) // true (martial is parent of weapon)
martial.Match(weapon) // false (weapon is not parent of martial)
weapon.Match(gameplay) // false (different hierarchy)
```

### Safety Features

The system prevents partial string matches:

```go
weaponTag := tag.New("gameplay.weapon")
weaponryTag := tag.New("gameplay.weaponry")

// This correctly returns false (not a partial string match)
weaponTag.Match(weaponryTag) // false
```

## Usage Patterns

### Damage Type System

```go
// Damage type categorization
fireDamage := tag.ContainerFromString("damage.fire,damage.elemental")
coldDamage := tag.ContainerFromString("damage.cold,damage.elemental")
physicalDamage := tag.ContainerFromString("damage.physical,damage.bludgeoning")

// Resistance checking
elementalResistance := tag.New("damage.elemental")
fireDamage.MatchTag(elementalResistance) // true
physicalDamage.MatchTag(elementalResistance) // false
```

### Item Categorization

```go
// Weapon categorization
longbow := tag.ContainerFromString("item.weapon.martial.ranged.bow")
greatsword := tag.ContainerFromString("item.weapon.martial.melee.sword")
dagger := tag.ContainerFromString("item.weapon.simple.melee.dagger")

// Proficiency checking
martialWeapons := tag.New("item.weapon.martial")
longbow.MatchTag(martialWeapons) // true
dagger.MatchTag(martialWeapons) // false
```

### Spell School System

```go
// Spell categorization
fireball := tag.ContainerFromString("spell.school.evocation,spell.level.3")
magicMissile := tag.ContainerFromString("spell.school.evocation,spell.level.1")
charm := tag.ContainerFromString("spell.school.enchantment,spell.level.1")

// School-based effects
evocationSpells := tag.New("spell.school.evocation")
fireball.MatchTag(evocationSpells) // true
charm.MatchTag(evocationSpells) // false
```

## Performance Characteristics

- **Normalization**: O(n) where n is input string length
- **Exact Matching**: O(1) string comparison
- **Hierarchical Matching**: O(m) where m is prefix length
- **Container Operations**: O(n) where n is number of tags in container
- **Memory**: Minimal overhead with string interning through normalization

## Design Rationale

The Tags System implements several key design decisions:

1. **Immutable Tags**: Tags are immutable after creation, preventing accidental modification
2. **Automatic Normalization**: Robust input handling prevents user errors
3. **Hierarchical Matching**: Enables flexible rule queries without complex logic
4. **Container Abstraction**: Simplifies multi-tag operations and queries
5. **Safety First**: Prevents partial matches and validates input integrity

This design enables flexible, safe, and performant categorization throughout the game engine.
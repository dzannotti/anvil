# DamageSource Interface Refactor Plan

## Overview

Convert `DamageSource` from a struct containing dice data to an interface returning expressions. This enables weapons, actions, and spells to be damage sources directly, eliminating redundant data and clarifying responsibilities.

## Current State Analysis

### Current DamageSource Struct

```go
type DamageSource struct {
    Times  int
    Sides  int
    Source string
    Tags   tag.Container
}
```

### Current Usage Patterns

- **Actions**: Contain `[]core.DamageSource` arrays
- **Weapons**: Contain `[]core.DamageSource` arrays
- **Actor.DamageRoll()**: Takes `[]DamageSource` parameter
- **Events**: `DamageRollEvent` contains `[]DamageSource`
- **Constructors**: Create struct literals everywhere

### Files Requiring Changes

- `internal/core/damage_source.go` - Interface definition
- `internal/core/actor_calculation.go` - DamageRoll method
- `internal/core/event.go` - Event structures
- `internal/expression/expression.go` - Add ExpectedValue() method
- `internal/ruleset/base/base_action.go` - Action base class
- `internal/ruleset/base/action_attack.go` - Attack actions
- `internal/ruleset/shared/action_fireball.go` - Spell actions
- `internal/ruleset/item/weapon/weapon.go` - Weapon implementation
- `internal/ruleset/item/weapon/weapons.go` - Weapon constructors
- `internal/ruleset/monster/undead/zombie/zombie.go` - Monster attacks

## Refactor Steps

### Phase 1: Interface Foundation

#### Step 1.1: Create New DamageSource Interface

- [x] Create new interface in `damage_source.go`
- [x] Keep old struct temporarily as `LegacyDamageSource`
- [x] Add interface methods:

  ```go
  type DamageSource interface {
      Name() string
      Damage() *expression.Expression
      Tags() *tag.Container
  }
  ```

#### Step 1.2: Add ExpectedValue() to Expression

- [x] Add `ExpectedValue() float64` method to `expression.Expression`
- [x] Calculate expected value from constants and dice averages
- [x] Implementation: constants as-is, dice as `times * (sides + 1) / 2.0`
- [x] Handle nested components recursively

#### Step 1.3: Create Legacy Adapter

- [x] Make `LegacyDamageSource` implement `DamageSource` interface
- [x] Convert Times/Sides to expression in `Damage()` method
- [x] Ensure backward compatibility during transition

### Phase 2: Core System Updates

#### Step 2.1: Update Actor.DamageRoll Method

- [ ] Change signature from `DamageRoll([]DamageSource, bool)` to `DamageRoll(DamageSource, bool)`
- [ ] Update implementation to call `source.Damage()` instead of building expression from struct
- [ ] Update event emission to use single source

#### Step 2.2: Update Events

- [ ] Change `DamageRollEvent.DamageSource` from `[]DamageSource` to `DamageSource`
- [ ] Update any other events that reference damage sources
- [ ] Ensure event handlers still work

#### Step 2.3: Update Base Action Class

- [ ] Remove `damage []core.DamageSource` field from `base.Action`
- [ ] Remove `Damage()` accessor method
- [ ] Actions will implement `DamageSource` interface directly

### Phase 3: Weapon Implementation

#### Step 3.1: Make Weapons Implement DamageSource

- [ ] Add interface methods to `Weapon` struct:
  - `Damage() *expression.Expression` - convert current damage data
  - `Tags() *tag.Container` - return weapon tags
  - `Name() string` - return weapon name
- [ ] Remove `damage []core.DamageSource` field from weapon

#### Step 3.2: Update Weapon Constructors

- [ ] Update `NewDagger()`, `NewGreatAxe()` constructors
- [ ] Remove DamageSource struct creation
- [ ] Store dice data directly in weapon fields or compute in `Damage()` method

#### Step 3.3: Update Weapon OnEquip

- [ ] Update attack action creation to pass weapon as damage source
- [ ] Remove intermediate DamageSource array creation

### Phase 4: Action Implementation

#### Step 4.1: Make Actions Implement DamageSource

- [ ] Add `DamageSource` interface methods to action types:
  - `AttackAction` - delegate to weapon, add modifiers
  - `FireballAction` - return spell damage expression
  - Natural weapon attacks (zombie slam) - return built-in damage

#### Step 4.2: Update Attack Actions

- [ ] `AttackAction.Damage()` - get weapon damage, add strength modifier
- [ ] `AttackAction.Tags()` - combine weapon tags + action tags
- [ ] `AttackAction.Name()` - return action name (not weapon name)

#### Step 4.3: Update Spell Actions  

- [ ] `FireballAction.Damage()` - return spell damage expression
- [ ] Remove embedded DamageSource array
- [ ] Update spell constructors

#### Step 4.4: Update Monster Actions

- [ ] Natural weapons (zombie slam) implement interface directly
- [ ] Remove DamageSource struct creation in constructors

### Phase 5: Melee Attack Resolution

#### Step 5.1: Design Composite Pattern

- [ ] Decide how `AttackAction` combines weapon + action data:

  ```go
  func (a *AttackAction) Damage() *expression.Expression {
      expr := a.weapon.Damage().Clone()
      expr.AddConstant(a.owner.StrengthModifier(), "Strength")
      return expr
  }
  
  func (a *AttackAction) Tags() *tag.Container {
      combined := tag.NewContainer(tags.Melee, tags.Attack)
      combined.Add(*a.weapon.Tags())
      return &combined
  }
  ```

#### Step 5.2: Update AttackAction Constructor

- [ ] Pass weapon reference instead of copying damage data
- [ ] Store weapon for delegation

### Phase 6: Integration & Cleanup

#### Step 6.1: Update All Call Sites

- [ ] Update places that call `Actor.DamageRoll()` to pass single source
- [ ] Update action `Perform()` methods
- [ ] Update AI damage calculation logic to use `ExpectedValue()`

#### Step 6.2: Remove Legacy Code

- [ ] Remove `LegacyDamageSource` struct
- [ ] Remove old constructors and helper methods
- [ ] Clean up any remaining `[]DamageSource` references

#### Step 6.3: Update Expression Integration

- [ ] Ensure all damage sources return properly configured expressions
- [ ] Test critical hits, resistance, advantage/disadvantage
- [ ] Verify damage type grouping still works

### Phase 7: Testing & Validation

#### Step 7.1: Build Validation

- [ ] Ensure code compiles at each step
- [ ] Fix any type mismatches

#### Step 7.2: Functional Testing  

- [ ] Run CLI demo to verify combat works
- [ ] Test weapon attacks, spell attacks, natural weapons
- [ ] Verify damage calculations are correct

#### Step 7.3: Edge Case Testing

- [ ] Test multiple damage types (flaming sword)
- [ ] Test resistance/vulnerability interactions  
- [ ] Test critical hit damage doubling
- [ ] Test complex spell damage (fireball with saves)

## Potential Challenges

### Challenge 1: AI Integration

**Problem**: AI needs to calculate average damage for decision making
**Solution**: Use `expression.ExpectedValue()` method to calculate averages from dice and constants


## Success Criteria

- [ ] All weapons implement DamageSource interface
- [ ] All actions implement DamageSource interface  
- [ ] No more embedded `[]DamageSource` arrays
- [ ] `Actor.DamageRoll()` takes single DamageSource parameter
- [ ] Melee attacks properly combine weapon + action data
- [ ] `expression.ExpectedValue()` method works for AI calculations
- [ ] All tests pass
- [ ] CLI demo runs successfully
- [ ] Damage calculations remain accurate
- [ ] Expression system integration works correctly

## Rollback Plan

If refactor fails, abandon this git branch and do not merge. The plan is designed to be all-or-nothing - either we complete the full refactor successfully or we don't merge the changes at all.

This refactor will significantly clean up the damage system and eliminate redundant data while providing much more flexibility for complex damage scenarios.

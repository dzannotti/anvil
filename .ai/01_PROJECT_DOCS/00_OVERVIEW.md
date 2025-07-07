# Spellbinder Architecture Overview

## Overview

Spellbinder is a modular D&D 2024 rules engine designed for clean architecture with no spaghetti dependencies. Each system is self-contained and communicates through well-defined interfaces.

## Core Design Principles

- **Clean Architecture**: No circular dependencies between systems
- **Modularity**: Each action/spell/feat is standalone, only referencing basic framework
- **No Silent Mutations**: All state changes are explicit and traceable
- **Minimal Code**: Less code = simpler system
- **Audit Trail**: Every calculation and action is fully traceable

## System Architecture

The core systems are documented in individual files:

- [Tags System](02_TAGS_SYSTEM.md) - Hierarchical identification and categorization
- [Expression System](03_EXPRESSION_SYSTEM.md) - Calculation engine with audit trails
- [Event Bus](04_EVENT_BUS.md) - Hierarchical event communication system
- [Effect System](05_EFFECT_SYSTEM.md) - Rule resolution callbacks _(future)_
- [Action System](06_ACTION_SYSTEM.md) - Game actions with metadata _(future)_
- [Entity System](07_ENTITY_SYSTEM.md) - Game objects with state _(future)_
- [Encounter System](08_ENCOUNTER_SYSTEM.md) - Combat/scene management _(future)_
- [AI System](09_AI_SYSTEM.md) - Decision making for NPCs _(future)_

## Data Flow

### Combat Data Flow

```
User/AI Input → Combat System → Expression + Effect Systems → Event Bus → UI/Logging
      ↓              ↓                    ↓                      ↓
   Decision      Action Execution    Calculation/Rules      Notification
```

### AI Decision Flow

```
AI System → Combat System → Expression + Effect Systems → Utility Score → Action Decision
    ↓             ↓                 ↓                         ↓              ↓
  Query      Execute Query    Calculate Results         Score Result    Select Action
```

## Technical Notes

- **No Concurrency**: Systems are designed for single-threaded operation
- **No Performance Testing**: Focus on correctness, optimize later if needed
- **Clean Interfaces**: Each system has clear, testable boundaries
- **Minimal Dependencies**: Each milestone depends only on previous ones

## CLI vs GUI vs GAME Strategy

**CLI**: Primary development interface

- Log output for all events
- Text-based state visualization
- Debugging and testing focus

**GUI**: User experience layer

- 2D rendering for debug representation
- Interactive dialogs

**GAME**: User experience layer

- 3D rendering for visual representation

The CLI provides complete functionality for development and testing, while GUI adds visual polish for end users.

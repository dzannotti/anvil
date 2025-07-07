# AI System

*Status: Future Implementation*

## Overview

Utility AI-based decision making system for non-player entities. **Does NOT perform combat calculations** - relies on Combat System for all mathematical operations.

## Architecture

**Pattern**: Utility AI pattern with weighted metrics  
**Information Access**: Same data available to human players (no cheating)  
**Decision Process**: Calculate utility scores for available actions using Combat System results

## Components

- **Utility Functions**: Calculate numeric scores for potential actions
- **Metric Weights**: Per-agent behavior configuration
- **Action Evaluator**: Queries Combat System to score available actions
- **Decision Engine**: Selects highest utility action

## Behavior Customization

- Range preference (melee vs ranged combat)
- Role preference (damage vs healing vs support)
- Risk tolerance (aggressive vs defensive)
- Target priority (focus fire vs spread damage)

## Example Configuration

```go
// Example AI agent configuration
aiAgent := AIAgent{
    Weights: map[string]float64{
        "damage_potential": 0.8,
        "self_preservation": 0.6,
        "healing_allies": 0.4,
        "range_preference": 0.7,  // Prefers ranged combat
    },
    Metrics: []UtilityMetric{
        DamageMetric{},
        SurvivalMetric{},
        HealingMetric{},
        PositionMetric{},
    },
}
```

## Integration Points

- Receives same event bus information as players
- Uses Combat System for action evaluation (not Expression System directly)
- Uses standard action system for execution
- Respects all game rules and limitations
- Provides reasoning for decisions (via utility scores, not expressions)

## Implementation Notes

This system will be implemented in a future milestone to provide intelligent NPC behavior that integrates with all other systems while maintaining game balance.
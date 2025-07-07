# Action System

*Status: Future Implementation*

## Overview

Function + metadata containers for all game actions.

## Components

- **Metadata**: cost, range, requirements, etc.
- **Function**: trigger effects and emit events
- **Validation**: precondition checks

## Flow

1. Validate action requirements
2. Trigger `before*` effect events
3. Execute core action logic
4. Trigger `after*` effect events
5. Emit event bus notifications

## Implementation Notes

This system will be implemented in a future milestone to provide structured game action execution that integrates with the Expression, Tag, and Effect systems.
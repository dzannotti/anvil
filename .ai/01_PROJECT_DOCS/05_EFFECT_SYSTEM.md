# Effect System

*Status: Future Implementation*

## Overview

Internal rule resolution system using event-driven callbacks attached to entities.

## Purpose

Modify calculations and game state based on entity conditions.

## Scope

Internal engine events only (not external communication).

## Event Types

- `AttributeCalculation`: Modify attribute expressions
- `BeforeAttackRoll`: Adjust attack calculations
- `DamageDealt`: Modify damage expressions
- `TurnStart`: Trigger beginning-of-turn effects

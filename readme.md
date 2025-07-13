# Anvil

A modular D&D 2024 rules engine with clean architecture and comprehensive system integration.

## Overview

Anvil (formerly Spellbinder) is a sophisticated rules engine implementing D&D 2024 mechanics through a series of interconnected systems. The architecture emphasizes clean interfaces, complete audit trails, and extensible design patterns.

### Core Systems

- **Tags System** - Hierarchical identification and categorization
- **Expression System** - Calculation engine with complete audit trails  
- **Event Bus** - Hierarchical event communication system
- **Effect System** - Rule resolution callbacks for D&D mechanics
- **Action System** - Game actions with resource management and targeting
- **AI System** - Decision making for NPCs _(in development)_

### Architecture Principles

- **Clean Architecture** - No circular dependencies between systems
- **Event-Driven** - All interactions flow through well-defined events  
- **Modular** - Each component is self-contained and testable
- **Audit Trail** - Every calculation and action is fully traceable

## Documentation

Comprehensive system documentation is available in `.ai/01_PROJECT_DOCS/`:

- [00_OVERVIEW.md](.ai/01_PROJECT_DOCS/00_OVERVIEW.md) - Architecture overview and design principles
- [01_CODING_STANDARDS.md](.ai/01_PROJECT_DOCS/01_CODING_STANDARDS.md) - Code style and conventions  
- [02_TAGS_SYSTEM.md](.ai/01_PROJECT_DOCS/02_TAGS_SYSTEM.md) - Hierarchical tagging system
- [03_EXPRESSION_SYSTEM.md](.ai/01_PROJECT_DOCS/03_EXPRESSION_SYSTEM.md) - Calculation and dice rolling engine
- [04_EVENT_BUS.md](.ai/01_PROJECT_DOCS/04_EVENT_BUS.md) - Event communication system
- [05_EFFECT_SYSTEM.md](.ai/01_PROJECT_DOCS/05_EFFECT_SYSTEM.md) - Rule resolution and effect callbacks
- [06_ACTION_SYSTEM.md](.ai/01_PROJECT_DOCS/06_ACTION_SYSTEM.md) - Game actions and resource management
- [07_AI_SYSTEM.md](.ai/01_PROJECT_DOCS/07_AI_SYSTEM.md) - AI decision making _(planned)_

Start with the [Overview](.ai/01_PROJECT_DOCS/00_OVERVIEW.md) for architecture principles and system interactions.

## Development Commands

| Command              | Description                                 |
| -------------------- | ------------------------------------------- |
| `make build`         | Build both CLI and GUI binaries             |
| `make cli`           | Build CLI binary only                       |
| `make gui`           | Build GUI binary only                       |
| `make test`          | Run all tests                               |
| `make tdd`           | Run tests in watch mode with concise output |
| `make test-coverage` | Run tests with coverage report              |
| `make fmt`           | Format all Go code                          |
| `make fmt-check`     | Check if code is formatted (CI)             |
| `make lint`          | Run linter on all code                      |
| `make lint-fix`      | Run linter and auto-fix issues              |
| `make clean`         | Remove build artifacts                      |
| `make clean-cache`   | Clean Go cache and module cache             |
| `make deps`          | Install development tools                   |
| `make hooks`         | Install pre-commit hooks                    |
| `make ci`            | Run full CI pipeline locally                |
| `make help`          | Show this help message                      |

## Project Structure

```text
├── cmd/
│   ├── cli/          # CLI application
│   └── gui/          # GUI application
├── internal/         # Internal packages
├── bin/              # Built binaries
└── Makefile          # Build automation
```

## Getting Started

**One-time setup after cloning:**

```bash
make setup
```

This installs all tools and sets up git hooks automatically.

**Daily development:**

```bash
make tdd    # Test-driven development mode
make ci     # Run full CI pipeline locally
```

## TODO Tracker 

- [ ] make actions definition based like items
- [ ] make creatures definition based like items and actions
- [ ] leverage registry to replace demo package
- [ ] rewrite ai
- [ ] finesse
- [ ] unarmed strike
- [ ] fire bolt
- [ ] prone
- [ ] instant death (overkill)
- [ ] resistance/vulnerability
- [ ] consider/poc using ids instead of references
- [ ] something with temp hit points
- [ ] lucky?
- [ ] action that gives poison
- [ ] action that uses start/end turn
- [ ] web
- [ ] dash
- [ ] dodge
- [ ] help
- [ ] vulnerability
- [ ] resistances
- [ ] up casting
- [ ] bark skin
- [ ] bane
- [ ] hex
- [ ] aid
- [ ] armor of agathys
- [ ] blade ward
- [ ] bless
- [ ] chill touch
- [ ] power word stun
- [ ] hold person
- [ ] hold monster
- [ ] hunter's mark
- [ ] prismatic spray
- [ ] tasha's laughter
- [ ] thunder wave
- [ ] phantasmal killer
- [ ] invisibility
- [ ] Elementals (can move thru other people)
- [ ] ability check

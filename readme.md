# Anvil

A game project focused on core game logic implementation.

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

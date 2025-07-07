# Coding Standards

## Comments

- **NO unnecessary comments** - Code should be self-documenting
- **NO comments that repeat what the code obviously does**
- Complex code should be broken down into smaller functions, not commented
- Only add comments for complex business logic or "why" explanations

## Control Flow

- **Use early returns and continues** instead of nested happy paths
- **Always add blank line after closing brackets of if statements**
- Invert conditions to fail fast: `if !condition { return/continue }`

## Examples

### ❌ Bad - Nested happy path with unnecessary comment

```go
for _, component := range components {
    if component.IsValid() { // Check if component is valid
        doSomething(component)
    }
}
```

### ✅ Good - Early continue with blank line after if

```go
for _, component := range components {
    if !component.IsValid() {
        continue
    }

    doSomething(component)
}
```

### ❌ Bad - Nested conditions

```go
if len(components) > 0 {
    if component.IsDamage() {
        processComponent(component)
    }
}
```

### ✅ Good - Early returns

```go
if len(components) == 0 {
    return
}

if !component.IsDamage() {
    return
}

processComponent(component)
```

## File Organization

- One struct/interface per file when possible
- Keep files under 200 lines (prefer ~100 lines)
- Split by feature/functionality, not just by struct
- File names should match their primary content

## Naming

- Use concise, clear names: `user` not `currentUserObject`
- Interface names should describe capability: `ComponentKind` not `ComponentBehavior`
- Field names should be consistent across similar structs

## Function/Method Design

- Functions should do one thing well
- Use early returns to reduce nesting
- Happy path should be unindented
- Prefer explicit over implicit
- No side effects unless clearly indicated by name

## Go-Specific Standards

### Package Management

- **Always use npm-style dependency management**: Use `go get` and `go mod` commands
- **Never manually edit go.mod** - use CLI tools

### Error Handling

- Handle errors explicitly
- Use early returns for error conditions
- Clear error messages that indicate what went wrong

### Testing

- Use table-driven tests when appropriate
- Test public APIs, not implementation details, assuming we'll throw away the whole implementation and want unit test to tell us what we missed
- Tests should be deterministic and fast

## Git Etiquette

- We commit often, after each checkpoint
- Before pushing make sure make ci is clean with no errors and no warnings
- We do not co-author commits

## Enforcement

These standards should be applied consistently across the codebase. When refactoring, always apply these rules to touched code.

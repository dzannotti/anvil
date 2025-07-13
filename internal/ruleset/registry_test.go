package ruleset

import (
	"testing"

	"anvil/internal/core"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/loader"
	"anvil/internal/tag"

	"github.com/stretchr/testify/assert"
)

func newEmptyRegistry() *Registry {
	return &Registry{
		actions: make(map[string]ActionFactory),
		effects: make(map[string]EffectFactory),
		items:   make(map[string]ItemFactory),
		actors:  make(map[string]ActorFactory),
	}
}

func TestRegistry_EmptyRegistry(t *testing.T) {
	registry := newEmptyRegistry()

	assert.NotNil(t, registry)
	assert.False(t, registry.HasAction("test"))
	assert.False(t, registry.HasEffect("test"))
	assert.False(t, registry.HasItem("test"))
	assert.False(t, registry.HasActor("test"))
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()

	assert.NotNil(t, registry)

	// Check that basic actions are registered
	assert.True(t, registry.HasAction("move"))

	// Check that basic effects are registered
	assert.True(t, registry.HasEffect("critical"))
	assert.True(t, registry.HasEffect("death"))
	assert.True(t, registry.HasEffect("death-saving-throw"))
	assert.True(t, registry.HasEffect("attack-of-opportunity"))
	assert.True(t, registry.HasEffect("proficiency-modifier"))
	assert.True(t, registry.HasEffect("attribute-modifier"))
	assert.True(t, registry.HasEffect("undead-fortitude"))
	assert.True(t, registry.HasEffect("fighting-style-defense"))

	// Check that basic items are registered
	assert.True(t, registry.HasItem("chainmail"))

	// Check that actors are registered
	assert.True(t, registry.HasActor("zombie"))

	// Check that some weapons are loaded from YAML
	assert.True(t, registry.HasItem("dagger"))
	assert.True(t, registry.HasItem("greataxe"))
}

func TestRegistry_ActionRegistration(t *testing.T) {
	registry := newEmptyRegistry()

	// Register a test action
	registry.RegisterAction("test-action", func(_ *core.Actor, _ map[string]interface{}) core.Action {
		return &MockAction{name: "test-action"}
	})

	assert.True(t, registry.HasAction("test-action"))

	// Create an action
	actor := &core.Actor{}
	action := registry.NewAction("test-action", actor, nil)

	assert.NotNil(t, action)
	mockAction := action.(*MockAction)
	assert.Equal(t, "test-action", mockAction.name)
}

func TestRegistry_EffectRegistration(t *testing.T) {
	registry := newEmptyRegistry()

	// Register a test effect
	registry.RegisterEffect("test-effect", func(_ map[string]interface{}) *core.Effect {
		return &core.Effect{Name: "test-effect"}
	})

	assert.True(t, registry.HasEffect("test-effect"))

	// Create an effect
	effect := registry.NewEffect("test-effect", nil)

	assert.NotNil(t, effect)
	assert.Equal(t, "test-effect", effect.Name)
}

func TestRegistry_ItemRegistration(t *testing.T) {
	registry := newEmptyRegistry()

	// Register a test item
	registry.RegisterItem("test-item", func(_ map[string]interface{}) core.Item {
		return &MockItem{name: "test-item"}
	})

	assert.True(t, registry.HasItem("test-item"))

	// Create an item
	item := registry.NewItem("test-item", nil)

	assert.NotNil(t, item)
	mockItem := item.(*MockItem)
	assert.Equal(t, "test-item", mockItem.name)
}

func TestRegistry_ActorRegistration(t *testing.T) {
	registry := newEmptyRegistry()

	// Register a test actor
	registry.RegisterActor("test-actor", func(_ map[string]interface{}) *core.Actor {
		return &core.Actor{Name: "test-actor"}
	})

	assert.True(t, registry.HasActor("test-actor"))

	// Create an actor
	actor := registry.NewActor("test-actor", nil)

	assert.NotNil(t, actor)
	assert.Equal(t, "test-actor", actor.Name)
}

func TestRegistry_PanicOnMissingArchetype(t *testing.T) {
	registry := newEmptyRegistry()
	actor := &core.Actor{}

	// Test panics for missing archetypes
	assert.Panics(t, func() {
		registry.NewAction("missing-action", actor, nil)
	})

	assert.Panics(t, func() {
		registry.NewEffect("missing-effect", nil)
	})

	assert.Panics(t, func() {
		registry.NewItem("missing-item", nil)
	})

	assert.Panics(t, func() {
		registry.NewActor("missing-actor", nil)
	})
}

func TestRegistry_ZombieCreation(t *testing.T) {
	registry := NewRegistry()

	dispatcher := &eventbus.Dispatcher{}
	world := core.NewWorld(loader.WorldDefinition{Width: 10, Height: 10})
	pos := grid.Position{X: 5, Y: 5}
	name := "Test Zombie"

	options := map[string]interface{}{
		"dispatcher": dispatcher,
		"world":      world,
		"position":   pos,
		"name":       name,
	}

	zombie := registry.NewActor("zombie", options)

	assert.NotNil(t, zombie)
	assert.Equal(t, name, zombie.Name)
}

func TestRegistry_ZombieCreationWithoutName(t *testing.T) {
	registry := NewRegistry()

	dispatcher := &eventbus.Dispatcher{}
	world := core.NewWorld(loader.WorldDefinition{Width: 10, Height: 10})
	pos := grid.Position{X: 5, Y: 5}

	options := map[string]interface{}{
		"dispatcher": dispatcher,
		"world":      world,
		"position":   pos,
	}

	zombie := registry.NewActor("zombie", options)

	assert.NotNil(t, zombie)
	assert.Equal(t, "Zombie", zombie.Name)
}

func TestRegistry_ZombieCreationPanicsOnMissingOptions(t *testing.T) {
	registry := NewRegistry()

	dispatcher := &eventbus.Dispatcher{}
	world := core.NewWorld(loader.WorldDefinition{Width: 10, Height: 10})
	pos := grid.Position{X: 5, Y: 5}

	// Test panic when missing dispatcher
	assert.Panics(t, func() {
		options := map[string]interface{}{
			"world":    world,
			"position": pos,
		}
		registry.NewActor("zombie", options)
	})

	// Test panic when missing world
	assert.Panics(t, func() {
		options := map[string]interface{}{
			"dispatcher": dispatcher,
			"position":   pos,
		}
		registry.NewActor("zombie", options)
	})

	// Test panic when missing position
	assert.Panics(t, func() {
		options := map[string]interface{}{
			"dispatcher": dispatcher,
			"world":      world,
		}
		registry.NewActor("zombie", options)
	})
}

func TestRegistry_WeaponCreation(t *testing.T) {
	registry := NewRegistry()

	// Test dagger creation
	dagger := registry.NewItem("dagger", nil)
	assert.NotNil(t, dagger)

	// Test greataxe creation
	greataxe := registry.NewItem("greataxe", nil)
	assert.NotNil(t, greataxe)
}

// Mock types for testing
type MockAction struct {
	name string
}

func (m *MockAction) Name() string                                        { return m.name }
func (m *MockAction) Archetype() string                                   { return "mock-action" }
func (m *MockAction) ID() string                                          { return "mock-id" }
func (m *MockAction) Tags() *tag.Container                                { tags := tag.ContainerFromTag(); return &tags }
func (m *MockAction) Perform(_ []grid.Position)                           {}
func (m *MockAction) ValidPositions(_ grid.Position) []grid.Position      { return nil }
func (m *MockAction) AffectedPositions(_ []grid.Position) []grid.Position { return nil }
func (m *MockAction) AverageDamage() int                                  { return 0 }
func (m *MockAction) CanAfford() bool                                     { return true }

type MockItem struct {
	name string
}

func (m *MockItem) Name() string         { return m.name }
func (m *MockItem) Archetype() string    { return "mock-item" }
func (m *MockItem) ID() string           { return "mock-id" }
func (m *MockItem) Tags() *tag.Container { tags := tag.ContainerFromTag(); return &tags }
func (m *MockItem) OnEquip(*core.Actor)  {}

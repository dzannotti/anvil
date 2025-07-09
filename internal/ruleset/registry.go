package ruleset

import (
	"fmt"

	"anvil/internal/core"
)

// ActionFactory creates actions for actors with options
type ActionFactory func(owner *core.Actor, options map[string]interface{}) core.Action

// EffectFactory creates effects with options
type EffectFactory func(options map[string]interface{}) *core.Effect

// ItemFactory creates items with options
type ItemFactory func(options map[string]interface{}) core.Item

// CreatureFactory creates creature actors with options
type CreatureFactory func(options map[string]interface{}) *core.Actor

// Registry serves as a unified database for all D&D ruleset content
type Registry struct {
	actions   map[string]ActionFactory
	effects   map[string]EffectFactory
	items     map[string]ItemFactory
	creatures map[string]CreatureFactory
}

// NewRegistry creates a new empty registry
func NewRegistry() *Registry {
	return &Registry{
		actions:   make(map[string]ActionFactory),
		effects:   make(map[string]EffectFactory),
		items:     make(map[string]ItemFactory),
		creatures: make(map[string]CreatureFactory),
	}
}

// RegisterAction registers an action factory for the given archetype
func (r *Registry) RegisterAction(archetype string, factory ActionFactory) {
	r.actions[archetype] = factory
}

func (r *Registry) NewAction(archetype string, owner *core.Actor, options map[string]interface{}) (core.Action, error) {
	factory, exists := r.actions[archetype]
	if !exists {
		return nil, fmt.Errorf("action archetype '%s' not found", archetype)
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(owner, options), nil
}

func (r *Registry) ListActions() []string {
	archetypes := make([]string, 0, len(r.actions))
	for archetype := range r.actions {
		archetypes = append(archetypes, archetype)
	}
	return archetypes
}

func (r *Registry) HasAction(archetype string) bool {
	_, exists := r.actions[archetype]
	return exists
}

// RegisterEffect registers an effect factory for the given archetype
func (r *Registry) RegisterEffect(archetype string, factory EffectFactory) {
	r.effects[archetype] = factory
}

func (r *Registry) NewEffect(archetype string, options map[string]interface{}) (*core.Effect, error) {
	factory, exists := r.effects[archetype]
	if !exists {
		return nil, fmt.Errorf("effect archetype '%s' not found", archetype)
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(options), nil
}

func (r *Registry) ListEffects() []string {
	archetypes := make([]string, 0, len(r.effects))
	for archetype := range r.effects {
		archetypes = append(archetypes, archetype)
	}
	return archetypes
}

func (r *Registry) HasEffect(archetype string) bool {
	_, exists := r.effects[archetype]
	return exists
}

// RegisterItem registers an item factory for the given archetype
func (r *Registry) RegisterItem(archetype string, factory ItemFactory) {
	r.items[archetype] = factory
}

func (r *Registry) NewItem(archetype string, options map[string]interface{}) (core.Item, error) {
	factory, exists := r.items[archetype]
	if !exists {
		return nil, fmt.Errorf("item archetype '%s' not found", archetype)
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(options), nil
}

func (r *Registry) ListItems() []string {
	archetypes := make([]string, 0, len(r.items))
	for archetype := range r.items {
		archetypes = append(archetypes, archetype)
	}
	return archetypes
}

func (r *Registry) HasItem(archetype string) bool {
	_, exists := r.items[archetype]
	return exists
}

// RegisterCreature registers a creature factory for the given archetype
func (r *Registry) RegisterCreature(archetype string, factory CreatureFactory) {
	r.creatures[archetype] = factory
}

func (r *Registry) NewCreature(archetype string, options map[string]interface{}) (*core.Actor, error) {
	factory, exists := r.creatures[archetype]
	if !exists {
		return nil, fmt.Errorf("creature archetype '%s' not found", archetype)
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(options), nil
}

func (r *Registry) ListCreatures() []string {
	archetypes := make([]string, 0, len(r.creatures))
	for archetype := range r.creatures {
		archetypes = append(archetypes, archetype)
	}
	return archetypes
}

func (r *Registry) HasCreature(archetype string) bool {
	_, exists := r.creatures[archetype]
	return exists
}

// DefaultRegistry is the global registry instance
var DefaultRegistry = NewRegistry()

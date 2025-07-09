package ruleset

import (
	"fmt"

	"anvil/internal/core"
)

type RegistryReader interface {
	NewAction(archetype string, owner *core.Actor, options map[string]interface{}) core.Action
	NewEffect(archetype string, options map[string]interface{}) *core.Effect
	NewItem(archetype string, options map[string]interface{}) core.Item
	NewCreature(archetype string, options map[string]interface{}) *core.Actor
	HasAction(archetype string) bool
	HasEffect(archetype string) bool
	HasItem(archetype string) bool
	HasCreature(archetype string) bool
}

type ActionFactory func(owner *core.Actor, options map[string]interface{}) core.Action
type EffectFactory func(options map[string]interface{}) *core.Effect
type ItemFactory func(options map[string]interface{}) core.Item
type CreatureFactory func(options map[string]interface{}) *core.Actor

type Registry struct {
	actions   map[string]ActionFactory
	effects   map[string]EffectFactory
	items     map[string]ItemFactory
	creatures map[string]CreatureFactory
}

func (r *Registry) RegisterAction(archetype string, factory ActionFactory) {
	r.actions[archetype] = factory
}

func (r *Registry) NewAction(archetype string, owner *core.Actor, options map[string]interface{}) core.Action {
	factory, exists := r.actions[archetype]
	if !exists {
		panic(fmt.Sprintf("action archetype '%s' not found", archetype))
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(owner, options)
}

func (r *Registry) RegisterEffect(archetype string, factory EffectFactory) {
	r.effects[archetype] = factory
}

func (r *Registry) NewEffect(archetype string, options map[string]interface{}) *core.Effect {
	factory, exists := r.effects[archetype]
	if !exists {
		panic(fmt.Sprintf("effect archetype '%s' not found", archetype))
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(options)
}

func (r *Registry) RegisterItem(archetype string, factory ItemFactory) {
	r.items[archetype] = factory
}

func (r *Registry) NewItem(archetype string, options map[string]interface{}) core.Item {
	factory, exists := r.items[archetype]
	if !exists {
		panic(fmt.Sprintf("item archetype '%s' not found", archetype))
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(options)
}

func (r *Registry) RegisterCreature(archetype string, factory CreatureFactory) {
	r.creatures[archetype] = factory
}

func (r *Registry) NewCreature(archetype string, options map[string]interface{}) *core.Actor {
	factory, exists := r.creatures[archetype]
	if !exists {
		panic(fmt.Sprintf("creature archetype '%s' not found", archetype))
	}

	if options == nil {
		options = make(map[string]interface{})
	}

	return factory(options)
}

func (r *Registry) HasAction(archetype string) bool {
	_, exists := r.actions[archetype]
	return exists
}

func (r *Registry) HasEffect(archetype string) bool {
	_, exists := r.effects[archetype]
	return exists
}

func (r *Registry) HasItem(archetype string) bool {
	_, exists := r.items[archetype]
	return exists
}

func (r *Registry) HasCreature(archetype string) bool {
	_, exists := r.creatures[archetype]
	return exists
}

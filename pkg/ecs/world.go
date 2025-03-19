package ecs

import "log"

// World is the main ECS container that holds all entities, components, and systems
type World struct {
	EntityManager    *EntityManager
	ComponentManager *ComponentManager
	systems          []System
	eventQueue       []EventInterface // Simple event queue for communication
	eventHandlers    map[EventType][]func(EventInterface)
	Logger           *log.Logger
}

func NewWorld(logger *log.Logger) *World {
	return &World{
		EntityManager:    NewEntityManager(),
		ComponentManager: NewComponentManager(),
		systems:          []System{},
		eventQueue:       []EventInterface{},
		eventHandlers:    make(map[EventType][]func(EventInterface)),
		Logger:           logger,
	}
}

func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) RemoveEntity(entity Entity) {
	w.EntityManager.RemoveEntity(entity)
	w.ComponentManager.RemoveAllComponents(entity)
}

func (w *World) Update() {
	for _, system := range w.systems {
		system.Update(w)
	}

	// Process events after all systems have updated
	w.processEvents()
}

// Simple event system for communication between ECS and external systems
type EventType string

type EventInterface interface {
	Type() EventType
	Entity() Entity
	Data() any
}

func (w *World) RegisterEventHandler(eventType EventType, handler func(EventInterface)) {
	w.eventHandlers[eventType] = append(w.eventHandlers[eventType], handler)
}

func (w *World) QueueEvent(event EventInterface) {
	w.eventQueue = append(w.eventQueue, event)
}

func (w *World) processEvents() {
	// Process all events in the queue
	for _, event := range w.eventQueue {
		if handlers, exists := w.eventHandlers[event.Type()]; exists {
			for _, handler := range handlers {
				handler(event)
			}
		}
	}

	// Clear queue
	w.eventQueue = w.eventQueue[:0]
}

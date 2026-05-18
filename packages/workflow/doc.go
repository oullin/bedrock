// Package workflow provides a Petri-Net based workflow and state-machine engine.
//
// The root package exposes the engine (Workflow[T], NewStateMachine), definitions
// (Definition, DefinitionBuilder, Transition, Place), and lifecycle events.
//
// Specialized helpers live in subpackages: audit (Trail recording), store
// (SingleState/MultiState marking stores), registry (multi-workflow registry),
// validator (extra definition checks), config (Definition loader on top of
// packages/config.Repository), events (typed Petri-Net event dispatcher), and
// multisteps (chevere-style DAG job orchestrator).
package workflow

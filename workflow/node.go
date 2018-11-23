package workflow

import "github.com/satori/go.uuid"

// MaxNumDep is the maximum number of dependents/children a node may have.
const MaxNumDep = 1000

// Signal is used for cross-node dependency communication.
type Signal struct {
	ID   uuid.UUID
	Pass bool
}

// Node composes an execution graph.
type Node interface {
	// Getters
	ID() uuid.UUID
	Name() string
	Ready() <-chan Signal
	Done() <-chan Signal
	Parents() []Node
	Children() []Node
	IsConditional() bool

	AddChild(Node) error
	AddParent(Node) error
	Activate()
	Execute() error
}

// Condition represents a conditional statement.
type Condition func() bool

// TruthyCondition always returns true.
var TruthyCondition = func() bool {
	return true
}

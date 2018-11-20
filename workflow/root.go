package workflow

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// NewRoot returns a Root that satisfies the Node interface.
func NewRoot(name string) Node {
	r := &Root{
		id:       uuid.NewV1(),
		name:     name,
		ready:    make(chan Signal, 1),
		done:     make(chan Signal, MaxNumDep),
		parents:  make(map[uuid.UUID]Node),
		children: make(map[uuid.UUID]Node),
	}

	return r
}

// Root implements Node. It serves as the starting point of an execution graph.
type Root struct {
	id       uuid.UUID
	name     string
	parents  map[uuid.UUID]Node
	children map[uuid.UUID]Node
	ready    chan Signal
	done     chan Signal
}

// ID returns Root's unique identifier.
func (r *Root) ID() uuid.UUID {
	return r.id
}

// Name returns Node's name.
func (r *Root) Name() string {
	return r.name
}

// Ready returns a channel that emits ready signal.
func (r *Root) Ready() <-chan Signal {
	return r.ready
}

// Done returns a channel that emits done signal.
func (r *Root) Done() <-chan Signal {
	return r.done
}

// Parents is a getter for a Node's dependency.
func (r *Root) Parents() []Node {
	return make([]Node, 0)
}

// Children is a getter for a Node's dependents.
func (r *Root) Children() []Node {
	nodes := make([]Node, 0, len(r.children))
	for _, n := range r.children {
		nodes = append(nodes, n)
	}

	return nodes
}

// AddChild inserts a child/dependent node to current node.
func (r *Root) AddChild(child Node) error {
	if len(r.children) == MaxNumDep {
		return errors.New("maximum number of children reached")
	}

	r.children[child.ID()] = child
	return nil
}

// Activate turns a node on and actively checks whether dependencies are met.
func (r *Root) Activate() {
	r.ready <- Signal{ID: r.id, Pass: true}
}

// Execute performs an action.
func (r *Root) Execute() error {
	log.Debugf("started %s", r.name)
	for i := 0; i < len(r.children); i++ {
		r.done <- Signal{ID: r.id, Pass: true}
	}
	log.Debugf("completed %s", r.name)
	return nil
}

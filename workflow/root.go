package workflow

import (
	"errors"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// NewRoot returns a Root that satisfies the Node interface.
func NewRoot(name string) Node {
	r := &Root{
		id:        uuid.NewV1(),
		name:      name,
		activated: false,
		mutex:     &sync.Mutex{},
		ready:     make(chan Signal, 1),
		done:      make(chan Signal, MaxNumDep),
		parents:   make(map[uuid.UUID]Node),
		children:  make(map[uuid.UUID]Node),
	}

	return r
}

// Root implements Node. It serves as the starting point of an execution graph.
type Root struct {
	id        uuid.UUID
	name      string
	activated bool
	mutex     *sync.Mutex

	// Means to communicate with other nodes
	ready chan Signal
	done  chan Signal

	parents  map[uuid.UUID]Node
	children map[uuid.UUID]Node
}

// ID returns Node's unique identifier.
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

// AddParent adds a dependency to current node.
func (r *Root) AddParent(n Node) error {
	return errors.New("root cannot have any parent")
}

// Children is a getter for a Node's dependents.
func (r *Root) Children() []Node {
	nodes := make([]Node, 0, len(r.children))
	for _, n := range r.children {
		nodes = append(nodes, n)
	}

	return nodes
}

// AddChild adds a child/dependent node to current node.
func (r *Root) AddChild(n Node) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.activated {
		return errors.New("node has been locked down, cannot modify its parent/child")
	}

	if len(r.children) == MaxNumDep {
		return errors.New("maximum number of children reached")
	}

	r.children[n.ID()] = n
	return nil
}

// Activate turns a node on and actively checks whether dependencies are met.
func (r *Root) Activate() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.activated = true
	r.ready <- Signal{ID: r.id, Pass: true}
}

// Execute performs an action.
func (r *Root) Execute() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.activated {
		return errors.New("must activate a node before execution")
	}

	log.Debugf("started %s", r.name)
	time.Sleep(100 * time.Millisecond) // Do something
	log.Debugf("completed %s", r.name)

	for i := 0; i < len(r.children); i++ {
		r.done <- Signal{ID: r.id, Pass: true}
	}

	return nil
}

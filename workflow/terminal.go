package workflow

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// NewTerminal returns a Terminal that satisfies the Node interface.
func NewTerminal(dependencies []Node, name string) Node {
	t := &Terminal{
		id:        uuid.NewV1(),
		name:      name,
		activated: false,
		mutex:     &sync.Mutex{},
		ready:     make(chan Signal, 1),
		done:      make(chan Signal, 1),
		parents:   make(map[uuid.UUID]Node),
		children:  make(map[uuid.UUID]Node),
	}

	for _, dep := range dependencies {
		t.AddParent(dep)
		dep.AddChild(t)
	}

	return t
}

// Terminal implements Node. It is inserted as leaf node of an execution graph.
type Terminal struct {
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
func (t *Terminal) ID() uuid.UUID {
	return t.id
}

// Name returns Node's name.
func (t *Terminal) Name() string {
	return t.name
}

// Ready returns a channel that emits ready signal.
func (t *Terminal) Ready() <-chan Signal {
	return t.ready
}

// Done returns a channel that emits done signal.
func (t *Terminal) Done() <-chan Signal {
	return t.done
}

// Parents is a getter for a Node's dependency.
func (t *Terminal) Parents() []Node {
	nodes := make([]Node, 0, len(t.parents))
	for _, n := range t.parents {
		nodes = append(nodes, n)
	}

	return nodes
}

// AddParent adds a dependency to current node.
func (t *Terminal) AddParent(n Node) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.activated {
		return errors.New("node has been locked down, cannot modify its parent/child")
	}

	t.parents[n.ID()] = n
	return nil
}

// Children is a getter for a Node's dependents.
func (t *Terminal) Children() []Node {
	return make([]Node, 0)
}

// AddChild adds a child/dependent node to current node.
func (t *Terminal) AddChild(child Node) error {
	return errors.New("terminal cannot have any child")
}

// IsConditional indicates whether a Node is conditional.
func (t *Terminal) IsConditional() bool {
	return false
}

// Activate turns a node on and actively checks whether dependencies are met.
func (t *Terminal) Activate() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	mux := make(chan Signal, len(t.parents))
	met := make(map[uuid.UUID]struct{})
	for _, dep := range t.parents {
		go func(id uuid.UUID, mux chan<- Signal, done <-chan Signal) {
			mux <- <-done
		}(dep.ID(), mux, dep.Done())
	}

	for sig := range mux {
		met[sig.ID] = struct{}{}

		if len(met) == len(t.parents) {
			t.activated = true
			t.ready <- Signal{ID: t.id, Pass: true}
			return
		}
	}
}

// Execute performs an action.
func (t *Terminal) Execute() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.activated {
		return errors.New("must activate a node before execution")
	}

	log.Infof("terminal node %s is reached", t.name)

	t.done <- Signal{ID: t.id, Pass: true}

	return nil
}

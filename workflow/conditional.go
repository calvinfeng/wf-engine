package workflow

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// NewConditional returns a Conditional that satisfies the Node interface.
func NewConditional(dependencies []Node, name string) Node {
	c := &Conditional{
		id:         uuid.NewV1(),
		name:       name,
		activated:  false,
		mutex:      &sync.Mutex{},
		ready:      make(chan Signal, 1),
		done:       make(chan Signal, MaxNumDep),
		parents:    make(map[uuid.UUID]Node),
		children:   make(map[uuid.UUID]Node),
		conditions: make(map[uuid.UUID]Condition),
	}

	for _, dep := range dependencies {
		c.AddParent(dep)
		dep.AddChild(c)
	}

	return c
}

// Condition represents a conditional statement.
type Condition func() bool

// Conditional implements Node. It will only activate children that satisfy
// the given condition.
type Conditional struct {
	id        uuid.UUID
	name      string
	activated bool
	mutex     *sync.Mutex

	// Means to communicate with other nodes
	ready chan Signal
	done  chan Signal

	parents    map[uuid.UUID]Node
	children   map[uuid.UUID]Node
	conditions map[uuid.UUID]Condition
}

// ID returns Node's unique identifier.
func (c *Conditional) ID() uuid.UUID {
	return c.id
}

// Name returns Node's name.
func (c *Conditional) Name() string {
	return c.name
}

// Parents is a getter for a Node's dependency.
func (c *Conditional) Parents() []Node {
	nodes := make([]Node, 0, len(c.parents))
	for _, n := range c.parents {
		nodes = append(nodes, n)
	}

	return nodes
}

// AddParent adds a dependency to current node.
func (c *Conditional) AddParent(n Node) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.activated {
		return errors.New("node has been locked down, cannot modify its parent/child")
	}

	c.parents[n.ID()] = n
	return nil
}

// Children is a getter for a Node's dependents.
func (c *Conditional) Children() []Node {
	nodes := make([]Node, 0, len(c.children))
	for _, n := range c.children {
		if c.conditions[n.ID()]() {
			nodes = append(nodes, n)
		}
	}

	return nodes
}

// AddChild adds a child/dependent node to current node.
func (c *Conditional) AddChild(n Node) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.activated {
		return errors.New("node has been locked down, cannot modify its parent/child")
	}

	if len(c.children) == MaxNumDep {
		return errors.New("maximum number of children reached")
	}

	// TODO: Caller should be able to pass in conditions.
	c.children[n.ID()] = n
	c.conditions[n.ID()] = func() bool {
		return true
	}

	return nil
}

// Ready returns a channel that emits ready signal.
func (c *Conditional) Ready() <-chan Signal {
	return c.ready
}

// Done returns a channel that emits done signal.
func (c *Conditional) Done() <-chan Signal {
	return c.done
}

// Activate turns a node on and actively checks whether dependencies are met.
func (c *Conditional) Activate() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	mux := make(chan Signal, len(c.parents))
	met := make(map[uuid.UUID]struct{})
	for _, dep := range c.parents {
		go func(id uuid.UUID, mux chan<- Signal, done <-chan Signal) {
			mux <- <-done
		}(dep.ID(), mux, dep.Done())
	}

	for sig := range mux {
		met[sig.ID] = struct{}{}

		if len(met) == len(c.parents) {
			c.activated = true
			c.ready <- Signal{ID: c.id, Pass: true}
			return
		}
	}
}

// Execute performs an action.
func (c *Conditional) Execute() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.activated {
		return errors.New("must activate a node before execution")
	}

	logrus.Infof("conditional node %s has started", c.name)
	logrus.Infof("conditional node %s is done", c.name)

	// TODO: Don't send done to children that does not satisfy condition.
	for i := 0; i < len(c.children); i++ {
		c.done <- Signal{ID: c.id, Pass: true}
	}

	return nil
}

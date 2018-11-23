package workflow

import (
	"errors"
	"fmt"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewConditional returns a Conditional that satisfies the Node interface.
func NewConditional(dependencies []Node, name string) Node {
	c := &Conditional{
		id:        uuid.NewV1(),
		name:      name,
		activated: false,
		mutex:     &sync.Mutex{},
		ready:     make(chan Signal, 1),
		done:      make(chan Signal, MaxNumDep),
		parents:   make(map[uuid.UUID]Node),
		children:  make(map[uuid.UUID]Node),
		cond:      false,
	}

	for _, dep := range dependencies {
		c.AddParent(dep)
		dep.AddChild(c)
	}

	return c
}

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

	parents  map[uuid.UUID]Node
	children map[uuid.UUID]Node

	cond bool
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

	// If condition is not satisfied, return no children.
	if !c.cond {
		return nodes
	}

	for _, n := range c.children {
		nodes = append(nodes, n)
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

	c.children[n.ID()] = n

	return nil
}

// IsConditional indicates whether a Node is conditional.
func (c *Conditional) IsConditional() bool {
	return true
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

	// Wait a little bit for the global state to poll server, because I didn't use websocket.
	time.Sleep(viper.GetDuration("conditional.wait_duration"))

	c.cond = true
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("freight%d", i)
		robot := requestIDLERobot(name)
		if robot == nil {
			c.cond = false
			break
		}

		if robot.CurrentPose.X != 10 || robot.CurrentPose.Y != 10 {
			c.cond = false
			break
		}
	}

	if c.cond {
		logrus.Infof("conditional node %s has been satisfied ", c.name)
	} else {
		logrus.Infof("conditional node %s has NOT been satisfied ", c.name)
	}

	// TODO: Don't send done to children that does not satisfy condition.
	for i := 0; i < len(c.children); i++ {
		c.done <- Signal{ID: c.id, Pass: true}
	}

	return nil
}

package workflow

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// NewJob returns an Job that satisfies the Node interface.
func NewJob(dependencies []Node, name string) Node {
	j := &Job{
		id:        uuid.NewV1(),
		name:      name,
		activated: false,
		mutex:     &sync.Mutex{},
		ready:     make(chan Signal, 1),
		done:      make(chan Signal, MaxNumDep),
		parents:   make(map[uuid.UUID]Node),
		children:  make(map[uuid.UUID]Node),
	}

	for _, dep := range dependencies {
		j.AddParent(dep)
		dep.AddChild(j)
	}

	return j
}

// Job implements Node. It performs an action through a task.
type Job struct {
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
func (j *Job) ID() uuid.UUID {
	return j.id
}

// Name returns Node's name.
func (j *Job) Name() string {
	return j.name
}

// Parents is a getter for a Node's dependency.
func (j *Job) Parents() []Node {
	nodes := make([]Node, 0, len(j.parents))
	for _, n := range j.parents {
		nodes = append(nodes, n)
	}

	return nodes
}

// AddParent adds a dependency to current node.
func (j *Job) AddParent(n Node) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.activated {
		return errors.New("node has been locked down, cannot modify its parent/child")
	}

	j.parents[n.ID()] = n
	return nil
}

// Children is a getter for a Node's dependents.
func (j *Job) Children() []Node {
	nodes := make([]Node, 0, len(j.children))
	for _, n := range j.children {
		nodes = append(nodes, n)
	}

	return nodes
}

// AddChild inserts a child/dependent node to current node.
func (j *Job) AddChild(n Node) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.activated {
		return errors.New("node has been locked down, cannot modify its parent/child")
	}

	if len(j.children) == MaxNumDep {
		return errors.New("maximum number of children reached")
	}

	j.children[n.ID()] = n
	return nil
}

// Ready returns a channel that emits ready signal.
func (j *Job) Ready() <-chan Signal {
	return j.ready
}

// Done returns a channel that emits done signal.
func (j *Job) Done() <-chan Signal {
	return j.done
}

// Activate turns a node on and actively checks whether dependencies are met.
func (j *Job) Activate() {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	mux := make(chan Signal, len(j.parents))
	met := make(map[uuid.UUID]struct{})
	for _, dep := range j.parents {
		go func(id uuid.UUID, mux chan<- Signal, done <-chan Signal) {
			mux <- <-done
			logrus.Debugf("job %s's dependency %s is met\n", j.name, id)
		}(dep.ID(), mux, dep.Done())
	}

	for sig := range mux {
		met[sig.ID] = struct{}{}

		if len(met) == len(j.parents) {
			logrus.Debugf("job %s is emitting ready signal", j.name)
			j.activated = true
			j.ready <- Signal{ID: j.id, Pass: true}
			return
		}
	}
}

// Execute performs an action.
func (j *Job) Execute() error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if !j.activated {
		return errors.New("must activate a node before execution")
	}

	logrus.Debugf("job %s has started", j.name)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Do something
	logrus.Debugf("job %s has completed", j.name)

	for i := 0; i < len(j.children); i++ {
		j.done <- Signal{ID: j.id, Pass: true}
	}

	return nil
}

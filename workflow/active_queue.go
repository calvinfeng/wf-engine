package workflow

import uuid "github.com/satori/go.uuid"

// NewActiveQueue returns an active queue.
func NewActiveQueue() *ActiveQueue {
	return &ActiveQueue{
		mux: make(chan Signal),
		set: make(map[uuid.UUID]Node),
	}
}

// ActiveQueue implements first-ready-first-out policy. A node is ready when it has all its
// dependencies met.
type ActiveQueue struct {
	mux chan Signal
	set map[uuid.UUID]Node
}

func (q *ActiveQueue) next() Node {
	sig := <-q.mux

	n, ok := q.set[sig.ID]
	if !ok {
		panic("node is not found in queue")
	}

	delete(q.set, sig.ID)
	return n
}

func (q *ActiveQueue) has(n Node) bool {
	_, ok := q.set[n.ID()]
	return ok
}

func (q *ActiveQueue) add(n Node) {
	q.set[n.ID()] = n
	go n.Activate()
	go func(id uuid.UUID, mux chan<- Signal, ready <-chan Signal) {
		mux <- <-ready
	}(n.ID(), q.mux, n.Ready())
}

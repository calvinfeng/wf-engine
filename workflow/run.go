package workflow

import (
	"errors"
)

// Run starts an executation graph.
func Run(root Node) error {
	if len(root.Parents()) != 0 {
		return errors.New("root node cannot have any dependency")
	}

	queue := NewActiveQueue()
	queue.add(root)
	for len(queue.set) > 0 {
		node := queue.next()
		go node.Execute()
		for _, child := range node.Children() {
			if queue.has(child) {
				continue
			}

			queue.add(child)
		}
	}

	return nil
}

package workflow

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	log.SetLevel(log.ErrorLevel)
}

func TestSingleDependency(t *testing.T) {
	root := NewRoot("Start")
	A := NewJob([]Node{root}, "Job A")
	B := NewJob([]Node{A}, "Job B")
	NewTerminal([]Node{B}, "End")

	err := Run(root)
	if err != nil {
		t.Error(err)
	}
}

func TestMultipleDependency(t *testing.T) {
	root := NewRoot("Start")
	A := NewJob([]Node{root}, "Job A")
	B := NewJob([]Node{root}, "Job B")
	C := NewJob([]Node{A, B}, "Job C")
	NewTerminal([]Node{C}, "End")

	err := Run(root)
	if err != nil {
		t.Error(err)
	}
}

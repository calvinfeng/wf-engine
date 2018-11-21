package main

import (
	"testing"
	"wf-engine/workflow"

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
	root := workflow.NewRoot("Start")
	A := workflow.NewJob([]workflow.Node{root}, "Job A")
	B := workflow.NewJob([]workflow.Node{A}, "Job B")
	workflow.NewTerminal([]workflow.Node{B}, "End")

	err := workflow.Run(root)
	if err != nil {
		t.Error(err)
	}
}

func TestMultipleDependency(t *testing.T) {
	root := workflow.NewRoot("Start")
	A := workflow.NewJob([]workflow.Node{root}, "Job A")
	B := workflow.NewJob([]workflow.Node{root}, "Job B")
	C := workflow.NewJob([]workflow.Node{A, B}, "Job C")
	workflow.NewTerminal([]workflow.Node{C}, "End")

	err := workflow.Run(root)
	if err != nil {
		t.Error(err)
	}
}

package main

import (
	"fmt"
	"wf-engine/workflow"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	log.SetLevel(log.DebugLevel)
}

func main() {
	R := workflow.NewRoot("Start")
	A := workflow.NewJob([]workflow.Node{R}, "A")
	B := workflow.NewJob([]workflow.Node{R}, "B")
	C := workflow.NewJob([]workflow.Node{A, B}, "C")
	T := workflow.NewTerminal([]workflow.Node{C}, "End")

	for _, n := range []workflow.Node{R, A, B, C, T} {
		fmt.Printf("%s -> %s\n", n.Name(), n.ID())
	}

	fmt.Printf("\n=======================\n\n")

	err := workflow.Run(R)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Graph is completed")
}

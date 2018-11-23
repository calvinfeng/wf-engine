package main

import (
	"context"
	"wf-engine/fleet"
	"wf-engine/global"

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
	// R := workflow.NewRoot("Start")
	// A := workflow.NewJob([]workflow.Node{R}, "A")
	// B := workflow.NewJob([]workflow.Node{R}, "B")
	// C := workflow.NewJob([]workflow.Node{A, B}, "C")
	// D := workflow.NewConditional([]workflow.Node{C}, "Conditional")
	// T1 := workflow.NewTerminal([]workflow.Node{D}, "Ending One")
	// T2 := workflow.NewTerminal([]workflow.Node{D}, "Ending Two")

	// for _, n := range []workflow.Node{R, A, B, C, T1, T2} {
	// 	fmt.Printf("%s -> %s\n", n.Name(), n.ID())
	// }

	// fmt.Printf("\n=======================\n\n")

	// err := workflow.Run(R)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Info("Graph is completed")

	ctx := context.Background()
	go global.State.Activate(ctx)

	log.Info("server is listening on 8000")
	if err := fleet.RunServer(); err != nil {
		log.Fatal(err)
	}
}

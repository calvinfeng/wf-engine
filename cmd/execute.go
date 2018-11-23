package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"
	"wf-engine/fleet"
	"wf-engine/global"
	"wf-engine/workflow"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func runserver(cmd *cobra.Command, args []string) error {
	log.Info("server is listening on 8000")
	if err := fleet.RunServer(); err != nil {
		return err
	}

	return nil
}

func runworkflow(cmd *cobra.Command, args []string) error {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go global.State.Activate(ctx)

	R := workflow.NewRoot("root")
	A := workflow.NewJob([]workflow.Node{R}, "sending freight1 to (10, 10)", "freight1")
	B := workflow.NewJob([]workflow.Node{R}, "sending freight2 to (10, 10)", "freight2")
	C := workflow.NewJob([]workflow.Node{R}, "sending freight3 to (10, 10)", "freight3")
	T1 := workflow.NewTerminal([]workflow.Node{A, B, C}, "all robots have started moving")
	D := workflow.NewConditional([]workflow.Node{A, B, C}, "are all robots at (10, 10)?")
	T2 := workflow.NewTerminal([]workflow.Node{D}, "all robots have reached (10, 10)")

	for _, n := range []workflow.Node{R, A, B, C, T1, T2} {
		fmt.Printf("%s -> %s\n", n.Name(), n.ID())
	}

	err := workflow.Run(R)
	if err != nil {
		return err
	}

	log.Info("Graph is completed")
	return nil
}

// Execute starts the program.
func Execute() {
	root := &cobra.Command{
		Use:   "wf-engine",
		Short: "Workflow engine that coordinates robots",
	}

	server := &cobra.Command{
		Use:     "runserver",
		Short:   "Run server",
		Example: "wf-engine runserver",
		RunE:    runserver,
	}

	workflow := &cobra.Command{
		Use:     "runworkflow",
		Short:   "Run workflow",
		Example: "wf-engine runworkflow",
		RunE:    runworkflow,
	}

	root.AddCommand(server, workflow)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

package cmd

import (
	"context"
	"fmt"
	"os"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go global.State.Activate(ctx)

	R := workflow.NewRoot("Start")
	A := workflow.NewJob([]workflow.Node{R}, "A")
	B := workflow.NewJob([]workflow.Node{R}, "B")
	C := workflow.NewJob([]workflow.Node{A, B}, "C")
	D := workflow.NewConditional([]workflow.Node{C}, "Conditional")
	T1 := workflow.NewTerminal([]workflow.Node{D}, "Ending One")
	T2 := workflow.NewTerminal([]workflow.Node{D}, "Ending Two")

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
		fmt.Println(err)
		os.Exit(1)
	}
}

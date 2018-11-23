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
	"github.com/spf13/viper"
)

func init() {
	viper.Reset()
	viper.AddConfigPath("conf")
	viper.SetConfigName("application")
	viper.SetConfigType("toml")
}

func runserver() {
	log.Infof("server is listening on %d", viper.GetInt("http.port"))
	if err := fleet.RunServer(); err != nil {
		log.Error(err)
	}
}

func runworkflow(cmd *cobra.Command, args []string) error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Make sure global state can poll robots correctly before starting workflow
	done := make(chan struct{})
	go runserver()
	go global.State.Activate(ctx, done)
	<-done

	R := workflow.NewRoot("root")
	A := workflow.NewJob([]workflow.Node{R}, "sending freight1 to (10, 10)", "freight1")
	B := workflow.NewJob([]workflow.Node{R}, "sending freight2 to (10, 10)", "freight2")
	C := workflow.NewJob([]workflow.Node{R}, "sending freight3 to (10, 10)", "freight3")
	D := workflow.NewConditional([]workflow.Node{A, B, C}, "are all robots at (10, 10)?")

	workflow.NewTerminal([]workflow.Node{A, B, C}, "all robots have started moving")
	workflow.NewTerminal([]workflow.Node{D}, "all robots have reached (10, 10)")

	err := workflow.Run(R)
	if err != nil {
		return err
	}

	log.Info("Graph is completed")

	var input string
	for {
		fmt.Print("Enter q to terminate program: ")
		fmt.Scanln(&input)
		if input == "q" {
			break
		}
	}

	return nil
}

// Execute starts the program.
func Execute() {
	root := &cobra.Command{
		Use:   "wf-engine",
		Short: "Workflow engine that coordinates robots",
	}

	workflow := &cobra.Command{
		Use:     "runworkflow",
		Short:   "Run workflow",
		Example: "wf-engine runworkflow",
		RunE:    runworkflow,
	}

	root.AddCommand(workflow)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

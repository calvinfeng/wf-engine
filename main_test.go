package main

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"wf-engine/fleet"
	"wf-engine/global"
	wf "wf-engine/workflow"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	log.SetLevel(log.ErrorLevel)

	viper.Reset()
	viper.AddConfigPath("conf")
	viper.SetConfigName("application")
	viper.SetConfigType("toml")
}

func TestWorkflow(t *testing.T) {
	if err := viper.ReadInConfig(); err != nil {
		t.Error(err)
		return
	}

	testserver := httptest.NewServer(fleet.LoadRoutes())
	defer testserver.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := viper.ReadInConfig(); err != nil {
		t.Error(err)
		return
	}

	viper.Set("http.port", strings.Split(testserver.URL, ":")[2])

	// Make sure global state can poll robots correctly before starting workflow
	done := make(chan struct{})
	go global.State.Activate(ctx, done)
	<-done

	root := wf.NewRoot("start")
	A := wf.NewJob([]wf.Node{root}, "navigate to (10, 10)", "freight1")
	B := wf.NewJob([]wf.Node{root}, "navigate to (10, 10)", "freight2")
	C := wf.NewJob([]wf.Node{A, B}, "navigate to (10, 10)", "freight3")
	wf.NewTerminal([]wf.Node{C}, "all robots have started moving")

	D := wf.NewConditional([]wf.Node{C}, "are all robots at (10, 10)?")
	wf.NewTerminal([]wf.Node{D}, "all robots have reached (10, 10)")

	err := wf.Run(root)
	if err != nil {
		t.Error(err)
	}
}

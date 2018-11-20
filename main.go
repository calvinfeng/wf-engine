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
	r := workflow.NewRoot("start of all")
	fmt.Println(r.ID())
}

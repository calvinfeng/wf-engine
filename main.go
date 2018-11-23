package main

import (
	"wf-engine/cmd"

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
	cmd.Execute()
}

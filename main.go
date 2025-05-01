package main

import (
	"github.com/kasia-kittel/rekor-verifier/cmd"
	"github.com/kasia-kittel/rekor-verifier/pkg/log"
)


func main() {
	log.StdLogger.Println("Starting")

	cmd.Execute()
}

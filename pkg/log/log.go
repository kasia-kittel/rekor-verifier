package log

import (
	"log"
	"os"
)

var StdOutLogger = createStdOutLogger()

func createStdOutLogger() *log.Logger {
	log := log.New(os.Stdout, "rekor-verifier ", log.LstdFlags)
	return log
}
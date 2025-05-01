package log

import (
	"log"
	"os"
)

var StdLogger = createStdLogger()

func createStdLogger() *log.Logger {
	log := log.New(os.Stdout, "rekor-verifier ", log.LstdFlags)
	return log
}
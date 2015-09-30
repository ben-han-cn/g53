package g53

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[debug]", 0)

func debug(fmt string, args ...interface{}) {
	logger.Printf(fmt, args...)
}

package logger

import (
	"log"
	"os"
)

// Logger Instance to logs the statements
var Log = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
package main

import (
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "[nugo] ", 0)

func LogDebugf(format string, a ...interface{}) {
	if os.Getenv("NUGO_DEBUG") == "1" {
		logger.Printf(format, a...)
	}
}

func LogInfof(format string, a ...interface{}) {
	logger.Printf(format, a...)
}

package main

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

func LogDebugf(format string, a ...interface{}) {
	if os.Getenv("NUGO_DEBUG") == "1" {
		logger.Printf("[debug] %s", fmt.Sprintf(format, a...))
	}
}

func LogInfof(format string, a ...interface{}) {
	logger.Printf("[info] %s", fmt.Sprintf(format, a...))
}

func LogErrorf(format string, a ...interface{}) {
	logger.Printf("[error] %s", fmt.Sprintf(format, a...))
}

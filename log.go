package main

import (
	"github.com/codegangsta/negroni"
	"log"
	"net/http"
	"os"
	"time"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

var logger *log.Logger = log.New(os.Stdout, "[nugo] ", 0)

func LogDebugf(format string, a ...interface{}) {
	if os.Getenv("NUGO_DEBUG") == "1" {
		logger.Printf(format, a...)
	}
}

func LogInfof(format string, a ...interface{}) {
	logger.Printf(format, a...)
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{logger}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Printf("Started %s %s", r.Method, r.URL)

	next(rw, r)

	res := rw.(negroni.ResponseWriter)
	l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
}

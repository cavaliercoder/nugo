package main

import (
	"github.com/codegangsta/negroni"
	"net/http"
	"time"
)

func NewHandler() http.Handler {
	n := negroni.New(negroni.NewRecovery())
	n.Use(negroni.HandlerFunc(Logger))
	n.Use(negroni.HandlerFunc(DefaultHeaders))
	n.Use(negroni.HandlerFunc(Mux))

	return n
}

// PanicOn triggers a failed server response if the given error is not nil.
func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}

// Logger logs the beginning and end of every transaction.
func Logger(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	LogInfof("Started %s %s", req.Method, req.URL)

	next(res, req)

	nres := res.(negroni.ResponseWriter)
	LogInfof("Completed %v %s in %v", nres.Status(), http.StatusText(nres.Status()), time.Since(start))
}

// DefaultHeaders applies HTTP headers that are valid for all server responses.
func DefaultHeaders(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	res.Header().Set("Server", "nugo")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Content-Type", "application/atom+xml;charset=utf-8")

	next(res, req)
}

// Mux route client requests to appropriate handler.
func Mux(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	switch req.URL.Path {
	case "/":
		GetRoot(res, req)
		break

	case "/$metadata":
		http.ServeFile(res, req, "metadata.xml")
		break

	case "/Search()":
		GetSearch(res, req)
		break
	}

	next(res, req)
}

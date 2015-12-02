package main

import (
	"github.com/codegangsta/negroni"
	"net/http"
	"time"
)

// NewHandler creates a http.Handler which hosts all middleware and routing
// functions for the HTTP server.
func NewHandler(config *Configuration) http.Handler {
	n := negroni.New(negroni.NewRecovery())
	n.Use(negroni.HandlerFunc(Logger))
	n.Use(negroni.HandlerFunc(DefaultHeaders))

	// register routes for each repository
	mux := http.NewServeMux()
	for _, r := range config.Repositories {
		LogDebugf("Registering route %s to repository: '%s'", r.RemotePath, r)
		mux.HandleFunc(r.RemotePath, RepoRouter(r))
	}

	n.UseHandler(mux)

	return n
}

// PanicOn triggers a failed server response if the given error is not nil.
func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}

// Logger is a Negroni middleware which logs the beginning and end of every
// transaction.
func Logger(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	LogInfof("Started %s %s", req.Method, req.URL)

	next(res, req)

	nres := res.(negroni.ResponseWriter)
	LogInfof("Completed %v %s in %v", nres.Status(), http.StatusText(nres.Status()), time.Since(start))
}

// DefaultHeaders is a Negroni middleware which applies HTTP headers that are
// valid for all server responses.
func DefaultHeaders(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	res.Header().Set("Server", "nugo")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Content-Type", "application/atom+xml;charset=utf-8")

	next(res, req)
}

// RepoRouter creates a http.HandlerFunc to route requests for each hosted
// package repository.
func RepoRouter(repo *Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		LogDebugf("Routing request: %s to repo '%s' registered at %s", req.URL.Path, repo.Name, repo.RemotePath)

		path := req.URL.Path[len(repo.RemotePath)-1:]

		// TODO: repos with a nested path break
		switch path {
		case "/":
			GetRoot(res, req)
			break

		case "/$metadata":
			http.ServeFile(res, req, "metadata.xml")
			break

		case "/Search()":
			GetSearch(res, req, repo)
			break
		}
	}
}

package main

import (
	"encoding/xml"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// load config
	config := GetConfig()

	// init router
	r := mux.NewRouter()
	r.HandleFunc("/", ListPackages)

	// init server
	n := negroni.Classic()
	n.UseHandler(r)

	n.Run(config.ListenPort)
}

func RenderXml(res http.ResponseWriter, req *http.Request, status int, v interface{}) {
	if v == nil {
		v = new(struct{})
	}

	var err error
	var data []byte
	if req.URL.Query().Get("pretty") == "true" {
		data, err = xml.MarshalIndent(v, "", "    ")
	} else {
		data, err = xml.Marshal(v)
	}

	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/xml")
	res.WriteHeader(status)
	res.Write([]byte(xml.Header))
	res.Write(data)
}

func ListPackages(res http.ResponseWriter, req *http.Request) {
	files, err := LoadPackageFiles()
	if err != nil {
		panic(err)
	}

	RenderXml(res, req, http.StatusOK, files)
}

func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}

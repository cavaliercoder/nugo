package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// load config
	config := GetConfig()

	// init router
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		log.Printf("Get: %s", req.URL)
		SetDefaultHeaders(res)

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

	})

	// serve
	http.ListenAndServe(config.ListenPort, nil)
}

func PanicOn(err error) {
	if err != nil {
		panic(err)
	}
}

func XMLEscape(s string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(s))

	return string(b.Bytes())
}

func SetDefaultHeaders(res http.ResponseWriter) {
	res.Header().Set("Server", "nugo")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Content-Type", "application/atom+xml;charset=utf-8")
}

func GetRoot(res http.ResponseWriter, req *http.Request) {

	// use buffered output so when can measure the content length
	var b bytes.Buffer

	fmt.Fprintf(&b, `<?xml version="1.0" encoding="utf-8" standalone="yes"?>%s`, "\n")
	fmt.Fprintf(&b, `<service xml:base="%s" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:app="http://www.w3.org/2007/app" xmlns="http://www.w3.org/2007/app">`, config.BaseURL)
	fmt.Fprint(&b, `<workspace>`)
	fmt.Fprint(&b, `<atom:title>Default</atom:title>`)
	fmt.Fprint(&b, `<collection href="Packages">`)
	fmt.Fprint(&b, `<atom:title>Packages</atom:title>`)
	fmt.Fprint(&b, `</collection>`)
	fmt.Fprint(&b, `</workspace>`)
	fmt.Fprint(&b, `</service>`)

	// flush buffer to client
	res.Header().Set("Content-Length", fmt.Sprintf("%d", b.Len()))
	res.Write(b.Bytes())
}

func printStringProperty(w io.Writer, tag string, value string) {
	if value == "" {
		fmt.Fprintf(w, `<d:%s m:null="true"></d:%s>`, tag, tag)
	} else {
		fmt.Fprintf(w, `<d:%s>%s</d:%s>`, tag, XMLEscape(value), tag)
	}
}

func printBoolProperty(w io.Writer, tag string, value bool) {
	if value {
		fmt.Fprintf(w, `<d:%s m:type="Edm.Boolean">true</d:%s>`, tag, tag)
	} else {
		fmt.Fprintf(w, `<d:%s m:type="Edm.Boolean">false</d:%s>`, tag, tag)
	}
}

func printIntProperty(w io.Writer, tag string, value int) {
	fmt.Fprintf(w, `<d:%s m:type="Edm.Int32">%d</d:%s>`, tag, value, tag)
}

func printDateProperty(w io.Writer, tag string, value time.Time) {
	fmt.Fprintf(w, `<d:%s m:type="Edm.DateTime">%s</d:%s>`, tag, value.Format("2006-01-02T15:04:05.999999Z"), tag)
}

func GetSearch(res http.ResponseWriter, req *http.Request) {
	config := GetConfig()

	// load packages from cache
	packages, err := GetConfig().Repositories[0].GetPackages()
	if err != nil {
		panic(err)
	}

	// use buffered output so when can measure the content length
	var b bytes.Buffer

	// print document header
	fmt.Fprintf(&b, `<?xml version="1.0" encoding="utf-8" standalone="yes"?>%s`, "\n")
	fmt.Fprintf(&b, `<feed xml:base="%s" xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata" xmlns="http://www.w3.org/2005/Atom">`, config.BaseURL)

	fmt.Fprint(&b, `<title type="text">Search</title>`)
	fmt.Fprintf(&b, `<id>%s/Search</id>`, config.BaseURL)
	fmt.Fprintf(&b, `<updated>%s</updated>`, time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	fmt.Fprint(&b, `<link rel="self" title="Search" href="Search" />`)

	// print packages
	for _, p := range packages {
		id := XMLEscape(p.Manifest.ID)
		version := XMLEscape(p.Manifest.Version)

		fmt.Fprint(&b, `<entry>`)
		fmt.Fprintf(&b, `<id>%s/Packages(Id='%s',Version='%s')</id>`, config.BaseURL, id, version)
		fmt.Fprintf(&b, `<title type="text">%s</title>`, id)
		fmt.Fprintf(&b, `<summary type="text">%s</summary>`, XMLEscape(p.Manifest.Summary))
		fmt.Fprintf(&b, `<link rel="edit-media" title="Package" href="Packages(Id='%s',Version='%s')/$value" xmlns="http://www.w3.org/2005/Atom" />`, id, version)
		fmt.Fprintf(&b, `<link rel="edit" title="Package" href="Packages(Id='%s',Version='%s')" xmlns="http://www.w3.org/2005/Atom" />`, id, version)
		fmt.Fprint(&b, `<category term="NuGet.Server.DataServices.Package" scheme="http://schemas.microsoft.com/ado/2007/08/dataservices/scheme" xmlns="http://www.w3.org/2005/Atom" />`)
		fmt.Fprintf(&b, `<content type="application/zip" src="%s/%s/%s" xmlns="http://www.w3.org/2005/Atom" />`, config.BaseURL, id, version)

		fmt.Fprint(&b, `<m:properties xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata" xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices">`)

		printStringProperty(&b, "Version", version)
		printStringProperty(&b, "Title", p.Manifest.Title)
		printStringProperty(&b, "Owners", p.Manifest.Owners)
		printStringProperty(&b, "IconUrl", p.Manifest.IconURL)
		printStringProperty(&b, "LicenseUrl", p.Manifest.LicenseURL)
		printStringProperty(&b, "ProjectUrl", p.Manifest.ProjectURL)
		printIntProperty(&b, "DownloadCount", 0)
		printBoolProperty(&b, "RequireLicenseAcceptance", p.Manifest.RequireLicenseAcceptance)
		printBoolProperty(&b, "DevelopmentDependency", false)
		printStringProperty(&b, "Description", p.Manifest.Description)
		printStringProperty(&b, "ReleaseNotes", p.Manifest.ReleaseNotes)
		printDateProperty(&b, "Published", time.Now().UTC())

		fmt.Fprint(&b, `</m:properties>`)

		fmt.Fprint(&b, `</entry>`)
	}

	// print document footer
	fmt.Fprintf(&b, "</feed>")

	// flush buffer to client
	res.Header().Set("Content-Length", fmt.Sprintf("%d", b.Len()))
	res.Write(b.Bytes())
}

package main

import (
	"encoding/xml"
	"io"
)

type Manifest struct {
	Authors                  string `xml:"metadata>authors"`
	Copyright                string `xml:"metadata>copyright"`
	Description              string `xml:"metadata>description"`
	IconURL                  string `xml:"metadata>iconUrl"`
	ID                       string `xml:"metadata>id"`
	LicenseURL               string `xml:"metadata>licenseUrl"`
	Owners                   string `xml:"metadata>owners"`
	ProjectURL               string `xml:"metadata>projectUrl"`
	ReleaseNotes             string `xml:"metadata>releaseNotes"`
	RequireLicenseAcceptance bool   `xml:"metadata>requireLicenseAcceptance"`
	Summary                  string `xml:"metadata>summary"`
	Tags                     string `xml:"metadata>tags"`
	Title                    string `xml:"metadata>title"`
	Version                  string `xml:"metadata>version"`
}

// ReadManifest decodes a .nuspec package manifest file into a Manifest struct.
func ReadManifest(r io.Reader) (*Manifest, error) {
	m := &Manifest{}

	d := xml.NewDecoder(r)
	err := d.Decode(m)
	PanicOn(err)

	return m, nil
}

package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Package struct {
	Path         string
	CreatedDate  time.Time
	ModifiedDate time.Time
	Manifest     Manifest
	IsLatest     bool
}

func (c *Package) String() string {
	return fmt.Sprintf("(Id='%s',Version='%s')", c.Manifest.ID, c.Manifest.Version)
}

func LoadPackage(path string) (*Package, error) {
	LogDebugf("Reading package: %s", path)
	p := &Package{}

	// get file info
	m, err := os.Stat(path)
	PanicOn(err)

	p.Path, err = filepath.Abs(path)
	PanicOn(err)

	p.ModifiedDate = m.ModTime()

	// read package zip file
	r, err := zip.OpenReader(p.Path)
	PanicOn(err)
	defer r.Close()

	// find .nuspec manifest file in zip file
	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".nuspec") {
			// open manifest
			mr, err := f.Open()
			PanicOn(err)
			defer mr.Close()

			// decode
			m, err := ReadManifest(mr)
			PanicOn(err)

			p.Manifest = *m

			break
		}
	}

	if p.Manifest.ID == "" {
		return nil, fmt.Errorf("Package file %s has no manifest", path)
	}

	return p, nil
}

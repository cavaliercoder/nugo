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
	return fmt.Sprintf("(id='%s', version='%s')", c.Manifest.ID, c.Manifest.Version)
}

func LoadPackage(path string) (*Package, error) {
	p := &Package{}

	m, err := os.Stat(path)
	PanicOn(err)

	p.Path, err = filepath.Abs(path)
	PanicOn(err)

	p.ModifiedDate = m.ModTime()

	// read package zip file
	r, err := zip.OpenReader(p.Path)
	PanicOn(err)

	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".nuspec") {
			mr, err := f.Open()
			PanicOn(err)

			defer mr.Close()

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

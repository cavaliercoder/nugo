package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Repository struct {
	FilePath   string
	packages   []Package
	cacheValid bool
}

func NewRepository(path string) *Repository {
	return &Repository{
		FilePath:   path,
		cacheValid: false,
	}
}

func (c *Repository) GetPackages() ([]Package, error) {
	if !c.cacheValid {
		// get a list of files in the repository directory
		files, err := ioutil.ReadDir(c.FilePath)
		PanicOn(err)

		packages := make([]*Package, 0)
		latestPackages := make(map[string]*Package, 0)

		for _, file := range files {
			// filter for .nupkg files
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".nupkg") {
				path := fmt.Sprintf("%s/%s", c.FilePath, file.Name())

				// load the package
				p, err := LoadPackage(path)
				PanicOn(err)

				packages = append(packages, p)

				id := strings.ToLower(p.Manifest.ID)

				// update latest version pointer
				if latest, ok := latestPackages[id]; ok {
					// compare version with previous
					v1, err := NewVersion(latest.Manifest.Version)
					PanicOn(err)

					v2, err := NewVersion(p.Manifest.Version)
					PanicOn(err)

					if v2.GreaterThan(v1) {
						latestPackages[id] = p
					}
				} else {
					latestPackages[id] = p
				}
			}
		}

		// set latest versions
		for _, p := range latestPackages {
			p.IsLatest = true
		}

		// copy package to repo struct
		out := make([]Package, len(packages))
		for i, p := range packages {
			out[i] = *p
		}

		c.packages = out
		c.cacheValid = true
	}

	return c.packages, nil
}

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
		if err != nil {
			return nil, err
		}

		c.packages = make([]Package, 0)
		for _, file := range files {
			// filter for .nupkg files
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".nupkg") {
				path := fmt.Sprintf("%s/%s", c.FilePath, file.Name())

				// load the package
				p, err := LoadPackage(path)
				if err != nil {
					panic(err)
				}

				c.packages = append(c.packages, *p)
			}
		}
		c.cacheValid = true
	}

	return c.packages, nil
}

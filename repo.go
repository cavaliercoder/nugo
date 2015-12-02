package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

type Repository struct {
	Name       string `yaml:"name"`
	LocalPath  string `yaml:"localPath"`
	RemotePath string `yaml:"remotePath"`
	packages   []Package
	cacheValid bool
}

type RepositorySearchParams struct {
	ByID       string
	SearchTerm string
	Limit      int
	Skip       int
	Filter     struct {
		LatestOnly        bool
		IncludePrerelease bool
		TargetFramework   string
	}
}

func NewRepositorySearchParams(v url.Values) *RepositorySearchParams {
	params := RepositorySearchParams{}

	params.SearchTerm = strings.Trim(v.Get("searchTerm"), " '")
	params.Filter.LatestOnly = v.Get("$filter") == "IsLatestVersion"
	params.Filter.IncludePrerelease = v.Get("includePrerelease") == "true"
	params.Filter.TargetFramework = strings.Trim(v.Get("targetFramework"), " '")

	return &params
}

func NewRepository(path string) *Repository {
	return &Repository{
		LocalPath:  path,
		cacheValid: false,
	}
}

func (c *Repository) String() string {
	return c.Name
}

func (c *Repository) RefreshCache() error {
	LogDebugf("Refreshing package cache for repo: %s", c)
	start := time.Now()

	// get a list of files in the repository directory
	files, err := ioutil.ReadDir(c.LocalPath)
	PanicOn(err)

	// fan out and load each package
	packageCount := 0
	packageChannel := make(chan *Package, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".nupkg") {
			packageCount++

			// load the package in parallel
			go func(filename string) {
				path := fmt.Sprintf("%s/%s", c.LocalPath, filename)
				p, err := LoadPackage(path)
				// TODO: Better error handling for erroneous packages
				PanicOn(err)

				packageChannel <- p
			}(file.Name())
		}
	}

	// parse packages as they are loaded
	packages := make([]*Package, packageCount)
	latestPackages := make(map[string]*Package, 0)
	for i := 0; i < packageCount; i++ {
		p := <-packageChannel
		packages[i] = p

		// update latest version pointer
		id := strings.ToLower(p.Manifest.ID)
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

	// set latest version property
	for _, p := range latestPackages {
		p.IsLatest = true
	}

	// dereference package slice
	out := make([]Package, len(packages))
	for i, p := range packages {
		out[i] = *p
	}

	// update cache
	// TODO: make cache updates atomic
	c.packages = out
	c.cacheValid = true

	LogInfof("Updated repository '%s' with %d packages in %v", c, packageCount, time.Since(start))

	return nil
}

func (c *Repository) GetPackages(params *RepositorySearchParams) ([]Package, error) {
	LogDebugf("Starting search: %#v", params)

	// inclusive or exclusive search?
	skipByDefault := params.ByID != "" || params.SearchTerm != ""
	if skipByDefault {
		LogDebugf("All packages will be excluded unless explicitely included")
	} else {
		LogDebugf("All packages will be included unless explicitely excluded")
	}

	// apply search params
	out := make([]Package, 0)
	for _, p := range c.packages {
		skip := skipByDefault

		// filter by ID
		if params.ByID != "" && (strings.ToLower(p.Manifest.ID) == strings.ToLower(params.ByID)) {
			LogDebugf("Package matches ID: %s", &p)
			skip = false
		}

		// filter by term
		if params.SearchTerm != "" && stringInString(params.SearchTerm, p.Manifest.ID, p.Manifest.Tags, p.Manifest.Description) {
			LogDebugf("Package matches search term: %s", &p)
			skip = false
		}

		// filter by latest only
		if params.Filter.LatestOnly && !p.IsLatest {
			LogDebugf("Package is superceded: %s", &p)
			skip = true
		}

		if !skip {
			out = append(out, p)
		}
	}

	// TODO: Implement skip and limit filters

	LogDebugf("Repo search yeilded %d results", len(out))
	return out, nil
}

func stringInString(needle string, haystacks ...string) bool {
	needle = strings.ToLower(needle)

	for _, haystack := range haystacks {
		if strings.Contains(strings.ToLower(haystack), needle) {
			return true
		}
	}

	return false
}

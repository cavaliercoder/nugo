package main

import (
	"testing"
)

func TestLoadPackages(t *testing.T) {
	// load package files
	packages, err := LoadPackageFiles()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(packages) == 0 {
		t.Errorf("0 packages loaded")
	}

	// validate name of each package file
	for _, p := range packages {
		if p.Manifest.ID == "" {
			t.Errorf("Package %s has no ID", p.Path)
		}

		if p.Manifest.Version == "" {
			t.Errorf("Package %s has no version", p.Path)
		}

		if p.Manifest.Description == "" {
			t.Errorf("Package %s has no description", p.Path)
		}

		if p.Manifest.Summary == "" {
			t.Errorf("Package %s has no summary", p.Path)
		}
	}

	t.Logf("Tested %d packages", len(packages))
}

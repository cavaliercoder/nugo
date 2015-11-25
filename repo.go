package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func LoadPackageFiles() ([]Package, error) {
	config := GetConfig()

	files, err := ioutil.ReadDir(config.PackagePath)
	if err != nil {
		return nil, err
	}

	packages := make([]Package, 0)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".nupkg") {
			path := fmt.Sprintf("%s/%s", config.PackagePath, file.Name())

			p, err := LoadPackage(path)
			if err != nil {
				panic(err)
			}

			packages = append(packages, *p)
		}
	}

	return packages, nil
}

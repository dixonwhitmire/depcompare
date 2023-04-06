// Package depcompare provides an API to load and compare dependencies.
package depcompare

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	GradleTextDependency  = "gradletxt"
	GradleBuildDependency = "gradlebuild"
)

var (
	DependencyTypes    = []string{GradleBuildDependency, GradleTextDependency}
	loadDependencyFunc func(depFile *os.File) map[string]string
)

// isValidDepType returns true if the depType is contained within the depTypes slice.
func isValidDepType(depType string, depTypes []string) bool {
	for _, s := range depTypes {
		if strings.EqualFold(s, depType) {
			return true
		}
	}
	return false
}

func loadGradleTextFile(depFile *os.File) map[string]string {
	m := make(map[string]string)

	s := bufio.NewScanner(depFile)
	for s.Scan() {
		tokens := strings.Split(s.Text(), ":")
		if len(tokens) < 2 {
			log.Println("depcompare.loadGradleTextFile: dependency format error - skipping")
			continue
		}

		dep := strings.Join(tokens[0:2], ":")
		version := ""
		if len(tokens) == 3 {
			version = tokens[2]
		}
		m[dep] = version
	}
	return m
}

func loadGradleBuildFile(depFile *os.File) map[string]string {
	return nil
}

// Load loads a dependency list from file, returning a map of dependency names and versions.
func Load(depType string, depPath string) (map[string]string, error) {
	if !isValidDepType(depType, DependencyTypes) {
		return nil, errors.New("depcompare.Load: depType is not valid")
	}

	f, err := os.Open(depPath)
	if err != nil {
		return nil, fmt.Errorf("depcompare.Load: error opening dependency file %w", err)
	}
	defer f.Close()

	switch depType {
	case GradleBuildDependency:
		loadDependencyFunc = loadGradleBuildFile
	case GradleTextDependency:
		loadDependencyFunc = loadGradleTextFile
	}

	m := loadDependencyFunc(f)
	return m, nil
}

// Compare compares one dependency list to another, returning comparison results in a map.
func Compare(depsA map[string]string, depsB map[string]string) map[string][]string {
	m := make(map[string][]string)
	return m
}

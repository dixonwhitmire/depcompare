// Package depcompare provides an API to load and compare dependencies.
package depcompare

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	// dependency type values
	gradleTextDependency  = "gradletxt"
	gradleBuildDependency = "gradlebuild"
	// keys for compare map
	BaseOnlyKey  = "base-only"
	DepOnlyKey   = "dep-only"
	IntersectKey = "intersect"
)

var (
	dependencyTypes    = []string{gradleBuildDependency, gradleTextDependency}
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

// loadGradleTextFile parses dependencies from a text file, where each depdency is written on a
// separate line.
func loadGradleTextFile(depFile *os.File) map[string]string {
	m := make(map[string]string)

	s := bufio.NewScanner(depFile)
	lineNum := 0
	for s.Scan() {
		tokens := strings.Split(s.Text(), ":")
		lineNum++
		if len(tokens) < 2 {
			log.Printf("depcompare.loadGradleTextFile: skipping line %d due to format error", lineNum)
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

// loadGradleBuildFile parses dependencies from a gradle.build file
func loadGradleBuildFile(depFile *os.File) map[string]string {
	m := make(map[string]string)
	s := bufio.NewScanner(depFile)
	lineNum := 0
	var inDependencyBlock = false

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		lineNum++

		// do not evaulate comments
		if strings.HasPrefix(line, "/*") ||
			strings.HasPrefix(line, "*") ||
			strings.HasPrefix(line, "/*") ||
			strings.HasPrefix(line, "//") {
			continue
		}

		if strings.Contains(strings.ToLower(line), "dependencies") {
			inDependencyBlock = true
		}

		if inDependencyBlock && strings.Contains("}", line) {
			inDependencyBlock = false
		}

		if !strings.Contains(strings.ToLower(line), "dependencies") && inDependencyBlock {
			// dependency information is formatted as:
			// [dependency type] `[dependency coordinates]`
			// example:
			// implementation 'org.apache.commons:commons-lang3:3.12.0'
			// parse dependency using single quotes as boundaries
			startIdx := strings.Index(line, "'")
			endIdx := strings.LastIndex(line, "'")
			if startIdx == -1 || endIdx == -1 {
				log.Printf("depcompare.loadGradleBuildFile: dependency not found on line %d", lineNum)
				continue
			}

			depEntry := line[startIdx+1 : endIdx]

			// split dependency into tokens
			tokens := strings.Split(depEntry, ":")
			if len(tokens) < 2 {
				log.Printf("depcompare.loadGradleBuildFile: skipping line %d due to format error", lineNum)
				continue
			}

			depGroupArtifact := strings.Join(tokens[0:2], ":")
			// version is an optional 3rd component
			depVersion := ""
			if len(tokens) == 3 {
				depVersion = tokens[2]
			}
			m[depGroupArtifact] = depVersion
		}

	}
	return m
}

// Load loads a dependency list from file, returning a map of dependency names and versions.
func Load(depType string, depPath string) (map[string]string, error) {
	if !isValidDepType(depType, dependencyTypes) {
		return nil, errors.New("depcompare.Load: depType is not valid")
	}

	f, err := os.Open(depPath)
	if err != nil {
		return nil, fmt.Errorf("depcompare.Load: error opening dependency file %w", err)
	}
	defer f.Close()

	switch depType {
	case gradleBuildDependency:
		loadDependencyFunc = loadGradleBuildFile
	case gradleTextDependency:
		loadDependencyFunc = loadGradleTextFile
	}

	m := loadDependencyFunc(f)
	return m, nil
}

// Compare compares a dependency list to a base dependency list, returning comparison results in a map.
// Map keys include:
// intersect: lists dependencies in both deps and baseDeps
// deps-only: lists dependencies in deps but not in baseDeps
// base-only: lists dependencies in base but not in deps
func Compare(deps map[string]string, baseDeps map[string]string) map[string][]string {
	m := map[string][]string{
		IntersectKey: make([]string, 0, len(baseDeps)),
		DepOnlyKey:   make([]string, 0, len(deps)),
		BaseOnlyKey:  make([]string, 0, len(baseDeps)),
	}

	var targetKey = ""

	// intersect and dep-only cases
	for k := range deps {

		if _, ok := baseDeps[k]; ok {
			targetKey = IntersectKey
		} else {
			targetKey = DepOnlyKey
		}
		m[targetKey] = append(m[targetKey], k)
	}

	// base-only case
	for k := range baseDeps {
		if _, ok := deps[k]; !ok {
			m[BaseOnlyKey] = append(m[BaseOnlyKey], k)
		}
	}
	// sort results
	sort.Strings(m[BaseOnlyKey])
	sort.Strings(m[IntersectKey])
	sort.Strings(m[DepOnlyKey])

	return m
}

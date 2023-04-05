// Package main contains the CLI entrypoint for the depcompare utility.
// depcompare compares a dependency list against a target dependency list.
// Comparison results display dependencies which exist in both the provided and target list.
//
// Usage:
// depcompare --type=[gradleb|gradlet] --mode=[intersect|diff] [dep file path] [target dep file path]
//
// The --type flag supports the dependency management type, such as Gradle, Maven, PyPi, etc.
// Valid --type values include:
//	- gradleb: gradle build file aka "build.gradle"
//	- gradlet: text file containing "gradle" compatible dependencies in short form. Example: org.apache.commons:commons-collections4:4.4.
//
// The --mode flag specifies the type of comparision used.
// --mode=intersect displays common dependencies in both lists, regardless of version.
// --mode=diff displays dependencies in the "target" (second) dep file path which are not in the first dep file path.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const depCompareUsage = `
Usage:
	depcompare --type=[gradleb|gradlet] --mode [diff|intersect] [dependency path] [target path]`

// valid values for CLI flags
var validDepTypeFlags = []string{"gradleb", "gradlet"}
var validModeFlags = []string{"diff", "intersect"}

// isValidFlag returns true if a flag value is valid.
func isValidFlag(flagValue string, validValues []string) bool {
	for _, s := range validValues {
		if strings.EqualFold(flagValue, s) {
			return true
		}
	}
	return false
}

// validate validates the provided CLI flags and arguments.
func validate(depTypeFlag string, modeFlag string, args []string) error {

	if !isValidFlag(depTypeFlag, validDepTypeFlags) {
		return fmt.Errorf("main.validate: invalid type flag. Expecting one of %v", validDepTypeFlags)
	}

	if !isValidFlag(modeFlag, validModeFlags) {
		return fmt.Errorf("main.validate: invalid mode flag. Expecting one of %v", validModeFlags)
	}

	if len(args) < 2 {
		return errors.New("main.validate: required arguments missing. Expecting [dependency path] [target path]")
	}

	// CLI only processes the first two arguments
	for _, a := range args[0:2] {
		f, err := os.Open(a)
		if err != nil {
			return fmt.Errorf("main.validate: invalid dependency path: %w", err)
		}
		defer f.Close()
	}

	return nil
}

// main is the entrypoint for the depcompare CLI.
func main() {
	var depTypeFlag string
	flag.StringVar(&depTypeFlag, "type", "gradleb", "Specifies the type of dependency compared. Defaults to gradleb.")

	var modeFlag string
	flag.StringVar(&modeFlag, "mode", "intersect", "Specifies the type of comparision made. Defaults to intersect")
	flag.Parse()

	err := validate(depTypeFlag, modeFlag, flag.Args())
	if err != nil {
		log.Fatalf("main: arguments are invalid. error:%v\n%s", err, depCompareUsage)
	}
}

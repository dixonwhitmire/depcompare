// Package main provides a CLI entrypoint to the depcompare API.
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dixonwhitmire/depcompare/internal/depcompare"
)

const usage = `
depcompare --type=[gradleb|gradlet] [dep path] [base dep path]
`

// printResult prints a single result set/key to stdout
func printResult(bannerText string, results []string) {
	fmt.Println("=======================================")
	fmt.Println(bannerText)
	fmt.Println("=======================================")
	for _, r := range results {
		fmt.Println(r)
	}
	fmt.Println("=======================================")

}

// displayResults prints comparision results to stdout.
func displayResults(results map[string][]string) {
	printResult("Common Dependencies", results[depcompare.IntersectKey])
	printResult("Base Only Dependencies", results[depcompare.BaseOnlyKey])
	printResult("Deps Only Dependencies", results[depcompare.DepOnlyKey])
}

func main() {
	var depTypeFlag string
	flag.StringVar(&depTypeFlag, "type", "the dependency type", "--type=gradleb")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalf("invalid number of arguments\n%s\n", usage)
	}

	log.Printf("main: dependency type set to %s", depTypeFlag)
	depPath := args[0]
	baseDepPath := args[1]

	deps, err := depcompare.Load(depTypeFlag, depPath)
	if err != nil {
		log.Fatalf("main: error loading dependency file %v", err)
	}
	log.Printf("main: loaded dependency file %s", depPath)

	baseDeps, err := depcompare.Load(depTypeFlag, baseDepPath)
	if err != nil {
		log.Fatalf("main: error loading base dependency file %v", err)
	}
	log.Printf("main: loaded base dependency file %s", baseDepPath)

	results := depcompare.Compare(deps, baseDeps)
	displayResults(results)
}

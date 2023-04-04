package main

import (
	"fmt"
	"math/rand"

	"github.com/dixonwhitmire/gotemplate/internal/message"
)

func main() {

	names := []string{"Homer", "Marge", "Bart", "Lisa", "Maggie", "Santa's Little Helper", "Snowball II"}
	i := rand.Intn(len(names))
	n := names[i]
	fmt.Println(message.Greeting(n))
}

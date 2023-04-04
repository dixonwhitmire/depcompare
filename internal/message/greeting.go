package message

import "fmt"

func Greeting(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

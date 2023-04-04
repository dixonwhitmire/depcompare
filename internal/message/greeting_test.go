package message

import (
	"testing"
)

func TestGreeting(t *testing.T) {

	tests := map[string]string{
		"Homer":  "Hello, Homer!",
		"Marge":  "Hello, Marge!",
		"Bart":   "Hello, Bart!",
		"Lisa":   "Hello, Lisa!",
		"Maggie": "Hello, Maggie!",
	}

	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			got := Greeting(name)

			if want != got {
				t.Errorf("want: %s, got: %s", want, got)
			}
		})
	}
}

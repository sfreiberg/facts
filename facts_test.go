package facts

import (
	"fmt"
	"testing"
)

func TestFacts(t *testing.T) {
	// Gather facts
	f := FindFacts()

	// Print as json
	json, err := f.ToPrettyJson()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", json)
}

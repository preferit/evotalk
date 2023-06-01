// Command rebel generates rebelious statement.
package main

import (
	"fmt"
	"os"

	"github.com/preferit/rebel"
)

func main() {
	rebel.MinimizeRepetition = true
	rebel.Shout(os.Stdout)
	fmt.Println()
}

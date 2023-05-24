// Command rebel generates rebelious statement.
package main

import (
	"fmt"
	"os"

	"github.com/preferit/rebel/phrase"
)

func main() {
	phrase.Shout(os.Stdout)
	fmt.Println()
}

// Command rebel generates rebelious statement.
package main

import (
	"os"

	"github.com/preferit/rebel/phrase"
)

func main() {
	phrase.Shout(os.Stdout)
}

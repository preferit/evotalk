// Command rebel generates rebelious statement.
package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/preferit/rebel/phrase"
)

func main() {
	// cache of last statement
	cache := "/tmp/lastphrase"
	last, _ := os.ReadFile(cache)
	// ignore error, we don't care
	var buf bytes.Buffer
	for {
		buf.Reset()
		phrase.Shout(&buf)
		if bytes.Equal(last, buf.Bytes()) {
			continue
		}
		break
	}
	// save
	_ = os.WriteFile(cache, buf.Bytes(), 0644)
	fmt.Println(buf.String())
}


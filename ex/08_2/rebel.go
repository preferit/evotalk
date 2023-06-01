package rebel

import (
	"io"
	"os"
)

// Shout writes a statement to os.Stdout. Affected by
// MinimizeRepetition flag.
func Shout(w io.Writer) {
	next := randValue(phrases)

	if MinimizeRepetition {
		cache := "/tmp/lastphrase"
		last, _ := os.ReadFile(cache)
		for {
			if string(last) == next {
				next = randValue(phrases)
				continue
			}
			break
		}
		_ = os.WriteFile(cache, []byte(next), 0644)
	}
	w.Write([]byte(next))
}

var MinimizeRepetition bool

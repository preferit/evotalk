package phrase

import (
	"fmt"
	"math/rand"
)

// Shout writes a statement to os.Stdout
func Shout() {
	fmt.Println(RandValue(phrases))
}

// RandValue returns a random element of the given slice or the zero value.
func RandValue[T any](slice []T) (v T) {
	if len(slice) == 0 {
		return
	}
	i := rand.Intn(len(slice))
	return slice[i]
}

var phrases = []string{
	"We are Rebelz!",
	"You may take our lives, but you will never take our freeeedooooom!",
	"We ain't gonna take it, No!, we ain't gonna take it!",
}

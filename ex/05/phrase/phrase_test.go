package phrase

import (
	"bytes"
	"strings"
	"testing"
)

// TestShout asserts that statements do Not ends with a new line
func TestShout(t *testing.T) {
	var buf bytes.Buffer
	Shout(&buf)
	if strings.HasSuffix(buf.String(), "\n") {
		t.Error("Shout should not end with new line")
	}
}

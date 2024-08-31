package scanner

import (
	"testing"
)

func TestAdvance(t *testing.T) {
	src := []byte("hello")
	current := 0
	for i := 0; i < len(src); i++ {
		c := advance(src, &current)
		if c != src[i] {
			t.Errorf("expected %c, got %c", src[i], c)
		}
	}
}

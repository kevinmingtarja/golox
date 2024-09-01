package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdvance(t *testing.T) {
	scanner := New([]byte("// hello"), nil)
	for i := 0; i < len(scanner.src); i++ {
		c := scanner.advance()
		if c != scanner.src[i] {
			t.Errorf("expected %c, got %c", scanner.src[i], c)
		}
	}
}

func TestMatch(t *testing.T) {
	scanner := New([]byte("!="), nil)

	c := scanner.advance()
	assert.Equal(t, '!', rune(c))
	assert.Equal(t, 1, scanner.current)

	assert.True(t, scanner.match('='))
	assert.Equal(t, 2, scanner.current)
}

func TestNoMatch(t *testing.T) {
	scanner := New([]byte("!true"), nil)

	c := scanner.advance()
	assert.Equal(t, '!', rune(c))
	assert.Equal(t, 1, scanner.current)

	assert.False(t, scanner.match('='))
	assert.Equal(t, 1, scanner.current)

	c = scanner.advance()
	assert.Equal(t, 't', rune(c))
	assert.Equal(t, 2, scanner.current)
}

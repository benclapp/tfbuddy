package rtfp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	s := []string{"a", "b", "c"}
	e := "b"

	assert.Equal(t, true, contains(s, e))
}

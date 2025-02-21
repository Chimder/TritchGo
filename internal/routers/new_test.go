package routers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	assert.Equal(t, 5, result, "Add(2,3) must return 5")
}

func Add(a, d int) int {
	return a + d
}

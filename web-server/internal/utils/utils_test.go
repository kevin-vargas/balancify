package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	actual := Identity(1)
	expect := 1
	assert.Equal(t, expect, actual)
}

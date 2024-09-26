package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	names := map[string]bool{}

	count := 10000000

	for i := 0; i < count; i++ {
		names[randomString()] = true
	}

	assert.Equal(t, len(names), count)

}

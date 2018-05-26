package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	strSlice := []string{"hello", "world", "key", "value"}
	assert.True(t, contains(strSlice, "hello"))
	assert.True(t, contains(strSlice, "world"))
	assert.True(t, contains(strSlice, "value"))
	assert.True(t, contains(strSlice, "key"))
	assert.False(t, contains(strSlice, "john"))
	assert.False(t, contains(strSlice, "smith"))
}

package golang_starter_test

import (
	"testing"

	"golang_starter"

	"github.com/stretchr/testify/assert"
)

// this file is just for the sake of coverage

func TestConst(t *testing.T) {
	constants := golang_starter.Const()
	assert.NotNil(t, constants)
}

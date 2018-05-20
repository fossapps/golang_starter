package starter_test

import (
	"testing"

	"github.com/fossapps/starter"

	"github.com/stretchr/testify/assert"
)

// this file is just for the sake of coverage

func TestConst(t *testing.T) {
	constants := starter.Const()
	assert.NotNil(t, constants)
}

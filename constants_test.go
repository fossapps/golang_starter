package crazy_nl_backend_test

import (
	"testing"

	"crazy_nl_backend"

	"github.com/stretchr/testify/assert"
)

// this file is just for the sake of coverage

func TestConst(t *testing.T) {
	constants := crazy_nl_backend.Const()
	assert.NotNil(t, constants)
}

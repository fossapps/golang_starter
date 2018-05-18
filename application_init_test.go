// +build integration

package golang_starter_test

import (
	"golang_starter/config"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestApplicationInit(t *testing.T) {
	session, err := mgo.Dial(config.GetMongoConfig().Connection)
	assert.Nil(t, err)
	assert.NotNil(t, session)
	defer session.Close()
}

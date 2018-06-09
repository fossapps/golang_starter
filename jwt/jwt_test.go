package jwt_test

import (
	"testing"
	"github.com/fossapps/starter/jwt"
	"time"
	"github.com/fossapps/starter/db"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"strings"
	"net/http/httptest"
)

func TestClient_CreateForUser(t *testing.T) {
	expect := assert.New(t)
	config := jwt.Config{
		Expiry: 1 * time.Second,
		Secret: "random secret",
	}
	manager := jwt.Client{
		Config: config,
	}
	token, err := manager.CreateForUser(&db.User{
		Email: "test",
		ID: bson.NewObjectId(),
		Permissions: []string{"sudo"},
	})
	expect.Nil(err)
	expect.NotNil(token)
	expect.Equal(3, len(strings.Split(token, ".")))
}

func TestClient_GetJwtDataFromRequestOnlyWorksWithSameSecret(t *testing.T) {
	expect := assert.New(t)
	config := jwt.Config{
		Expiry: 1 * time.Second,
		Secret: "random secret",
	}
	manager := jwt.Client{
		Config: config,
	}
	token, _ := manager.CreateForUser(&db.User{
		Email: "test",
		ID: bson.NewObjectId(),
		Permissions: []string{"sudo"},
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "bearer " + token)
	claims, err := manager.GetJwtDataFromRequest(request)
	expect.Nil(err)
	expect.NotNil(claims) // same manager should work.
	newManager := jwt.Client{
		Config: jwt.Config{
			Secret: "new secret",
			Expiry: 1 * time.Second,
		},
	}
	claims, err = newManager.GetJwtDataFromRequest(request)
	expect.Nil(claims)
	expect.NotNil(err)
}

func TestClient_GetJwtDataFromRequestReturnsErrorForExpiredJwts(t *testing.T) {
	expect := assert.New(t)
	config := jwt.Config{
		Expiry: - 1 * time.Second,
		Secret: "random secret",
	}
	manager := jwt.Client{
		Config: config,
	}
	token, err := manager.CreateForUser(&db.User{
		Email: "test",
		ID: bson.NewObjectId(),
		Permissions: []string{"sudo"},
	})
	request := httptest.NewRequest("GET", "/", nil)
	request.Header.Add("Authorization", "bearer " + token)

	claims, err := manager.GetJwtDataFromRequest(request)
	expect.NotNil(err)
	expect.Nil(claims)
}

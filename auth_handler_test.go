package starter_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fossapps/starter"
	"github.com/fossapps/starter/db"
	"github.com/fossapps/starter/mock"

	"errors"
	"github.com/fossapps/starter/adapter"
	"github.com/fossapps/starter/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func getLogger() starter.Logger {
	logger := logrus.New()
	logger.Out = httptest.NewRecorder()
	return logger
}

// region LoginHandler

func TestServer_LoginHandlerRespondsWithUnauthorizedIfNoHeader(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := starter.Server{
		Logger: getLogger(),
	}
	server.LoginHandler()(responseRecorder, request)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithUnauthorizedIfWrongPassword(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(&user)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := starter.Server{
		Logger: getLogger(),
		Db:     mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, "wrong_password")
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}
func TestServer_LoginHandlerRespondsWithBadRequestIfNoUser(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := starter.Server{
		Logger: getLogger(),
		Db:     mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, "wrong_password")
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithOkOnCorrectCredentials(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockRefreshTokenManager.EXPECT().Add(gomock.Any(), gomock.Any())
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(&user)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := starter.Server{
		Db: mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, pass)
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	res := new(starter.LoginResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.NotNil(res.RefreshToken)
	expect.NotNil(res.JWT)
	expect.True(len(res.RefreshToken) >= 128)
	expect.True(strings.Count(res.JWT, ".") == 2)
}

// endregion

// region RefreshTokenHandler
func TestServer_RefreshTokenHandlerStoresRefreshTokenInDb(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	mockRefreshTokenManager.EXPECT().Add(gomock.Any(), gomock.Any())
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail(email).Return(&user)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)
	server := starter.Server{
		Db: mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, pass)
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	res := new(starter.LoginResponse)
	json.NewDecoder(responseRecorder.Body).Decode(&res)
	expect.NotNil(res.RefreshToken)
	expect.NotNil(res.JWT)
	expect.True(len(res.RefreshToken) >= 128)
	expect.True(strings.Count(res.JWT, ".") == 2)
}

func TestServer_RefreshTokenHandlerRespondsWithStatusBadRequestIfNoAuthToken(t *testing.T) {
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	server := starter.Server{}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerRespondsWithStatusUnauthorizedIfRefreshTokenInvalid(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(nil)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")
	server := starter.Server{
		Db:     mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerRefreshTokenNotLinkedToUserRespondsWithStatusUnauthorized(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager.EXPECT().FindByID("some_user").Times(1).Return(nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockToken := &db.RefreshToken{
		Token: "auth_token",
		User:  "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := starter.Server{
		Db:     mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusUnauthorized, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerReturnsJWT(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	refreshTokenCtrl := gomock.NewController(t)
	mockUsersCtrl := gomock.NewController(t)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUser := &db.User{
		ID:          "random",
		Email:       "random",
		Permissions: []string{"sudo"},
	}
	mockUserManager.EXPECT().FindByID("some_user").Times(1).Return(mockUser)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockToken := &db.RefreshToken{
		Token: "auth_token",
		User:  "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := starter.Server{
		Db:     mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
}

// endregion

// region Sessions.List
func TestServer_RefreshTokensListReturnsBadRequestWhenTokenWrong(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRequestHelper := mock.NewMockRequestHelper(ctrl)
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Return(nil, errors.New("jwt error"))
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	server := starter.Server{ReqHelper: mockRequestHelper}
	server.RefreshTokensList()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
}
func getJwtForUser(id string, email string, permission []string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":          id,
		"email":       email,
		"permissions": permission,
		"exp":         time.Now().Add(config.GetApplicationConfig().JWTExpiryTime).Unix(),
	})
	signedString, _ := token.SignedString([]byte(config.GetApplicationConfig().JWTSecret))
	return signedString
}

func TestServer_RefreshTokensListReturnsInternalServerIfDbError(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	requestHelperCtrl := gomock.NewController(t)
	defer requestHelperCtrl.Finish()
	mockRequestHelper := mock.NewMockRequestHelper(requestHelperCtrl)
	userID := "some_random_id"
	claims := adapter.Claims{
		ID:          userID,
		Email:       "admin@example.com",
		Permissions: []string{"sudo"},
	}
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Return(&claims, nil)
	mockDb := mock.NewMockDB(dbCtrl)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	mockRefreshTokenManager.EXPECT().List(userID).Return(nil, errors.New("dbError"))
	mockDb.EXPECT().RefreshTokens().AnyTimes().Return(mockRefreshTokenManager)
	token := getJwtForUser(userID, "admin@example.com", []string{"sudo"})
	// that's valid jwt
	mockRequest := httptest.NewRequest("GET", "/", nil)
	mockRequest.Header.Add("Authorization", "Bearer "+token)
	responseRecorder := httptest.NewRecorder()
	server := starter.Server{
		Db:        mockDb,
		ReqHelper: mockRequestHelper,
	}
	server.RefreshTokensList()(responseRecorder, mockRequest)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_RefreshTokensListReturnsRefreshTokenList(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockDb := mock.NewMockDB(dbCtrl)
	requestHelperCtrl := gomock.NewController(t)
	defer requestHelperCtrl.Finish()
	mockRequestHelper := mock.NewMockRequestHelper(requestHelperCtrl)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	userID := "some_random_id"
	claims := adapter.Claims{
		ID:          userID,
		Email:       "admin@example.com",
		Permissions: []string{"sudo"},
	}
	mockRequestHelper.EXPECT().GetJwtData(gomock.Any()).Return(&claims, nil)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	refreshTokens := []db.RefreshToken{
		{Token: "token1", User: "some_random_id"},
		{Token: "token2", User: "some_random_id"},
	}
	// get list and return
	mockRefreshTokenManager.EXPECT().List(userID).Return(refreshTokens, nil)
	mockDb.EXPECT().RefreshTokens().AnyTimes().Return(mockRefreshTokenManager)
	token := getJwtForUser(userID, "admin@example.com", []string{"sudo"})
	// that's valid jwt
	mockRequest := httptest.NewRequest("GET", "/", nil)
	mockRequest.Header.Add("Authorization", "Bearer "+token)
	responseRecorder := httptest.NewRecorder()
	server := starter.Server{
		Db:        mockDb,
		ReqHelper: mockRequestHelper,
	}
	server.RefreshTokensList()(responseRecorder, mockRequest)
	expect.Equal(http.StatusOK, responseRecorder.Code)
	var responseTokens []db.RefreshToken
	json.NewDecoder(responseRecorder.Body).Decode(&responseTokens)
	expect.Equal(refreshTokens, responseTokens)
}

// endregion

// region Sessions.Delete

func TestServer_DeleteSessionReturnsNotFoundIfNoToken(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockDb := mock.NewMockDB(dbCtrl)
	requestHelperCtrl := gomock.NewController(t)
	defer requestHelperCtrl.Finish()
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockRefreshTokenManager.EXPECT().Delete("token").Return(mgo.ErrNotFound)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	mockDb.EXPECT().RefreshTokens().Return(mockRefreshTokenManager)
	server := starter.Server{
		Db: mockDb,
	}
	router := mux.NewRouter()
	router.HandleFunc("/session/{token}", server.DeleteSession())
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("DELETE", "/session/token", nil)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusNotFound, responseRecorder.Code)
}

func TestServer_DeleteSessionReturnsInternalServerErrorIfDbError(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockDb := mock.NewMockDB(dbCtrl)
	requestHelperCtrl := gomock.NewController(t)
	defer requestHelperCtrl.Finish()
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockRefreshTokenManager.EXPECT().Delete("token").Return(errors.New("invalid query"))
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	mockDb.EXPECT().RefreshTokens().Return(mockRefreshTokenManager)
	server := starter.Server{
		Db: mockDb,
	}
	router := mux.NewRouter()
	router.HandleFunc("/session/{token}", server.DeleteSession())
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("DELETE", "/session/token", nil)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_DeleteSessionReturnsNoContent(t *testing.T) {
	expect := assert.New(t)
	dbCtrl := gomock.NewController(t)
	defer dbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockDb := mock.NewMockDB(dbCtrl)
	requestHelperCtrl := gomock.NewController(t)
	defer requestHelperCtrl.Finish()
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockRefreshTokenManager.EXPECT().Delete("token").Return(nil)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	mockDb.EXPECT().RefreshTokens().Return(mockRefreshTokenManager)
	server := starter.Server{
		Db: mockDb,
	}
	router := mux.NewRouter()
	router.HandleFunc("/session/{token}", server.DeleteSession())
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("DELETE", "/session/token", nil)
	router.ServeHTTP(responseRecorder, request)
	expect.Equal(http.StatusNoContent, responseRecorder.Code)
}

// endregion

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
	"github.com/globalsign/mgo"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"github.com/fossapps/starter/logger"
	"github.com/fossapps/starter/jwt"
)

func getLogger() logger.Client {
	client := logrus.New()
	client.Out = httptest.NewRecorder()
	return client
}

func getJwtForUser(permission []string) string {
	return "token: " + strings.Join(permission, ",")
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
	defer mockDbCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(&user, nil)
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

func TestServer_LoginHandlerHandlesFindByEmailDbError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(nil, errors.New("db error"))
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := starter.Server{
		Logger: getLogger(),
		Db:     mockDb,
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, "wrong_password")
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithBadRequestIfNoUser(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(nil, nil)
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

func TestServer_LoginHandlerHandlesJwtGenerationError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	JwtCtrl := gomock.NewController(t)
	defer JwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(JwtCtrl)
	mockJwt.EXPECT().CreateForUser(gomock.Any()).Times(1).Return("", errors.New("error"))
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(&user, nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := starter.Server{
		Db:  mockDb,
		Jwt: mockJwt,
		Logger:getLogger(),
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.SetBasicAuth(email, pass)
	server.LoginHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_LoginHandlerRespondsWithOkOnCorrectCredentials(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	JwtCtrl := gomock.NewController(t)
	defer JwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(JwtCtrl)
	mockJwt.EXPECT().CreateForUser(gomock.Any()).Times(1).Return("sample_token", nil)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockRefreshTokenManager.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail("admin@example.com").Return(&user, nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	server := starter.Server{
		Db:  mockDb,
		Jwt: mockJwt,
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
	expect.Equal("sample_token", res.JWT)
}

// endregion

// region RefreshTokenHandler
func TestServer_RefreshTokenHandlerStoresRefreshTokenInDb(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	jwtCtrl := gomock.NewController(t)
	defer jwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(jwtCtrl)
	mockJwt.EXPECT().CreateForUser(gomock.Any()).Times(1).Return("sample_token", nil)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	email := "admin@example.com"
	pass := "pass"
	mockRefreshTokenManager.EXPECT().Add(gomock.Any(), gomock.Any())
	hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	user := db.User{Email: email, Password: string(hash)}
	mockUserManager.EXPECT().FindByEmail(email).Return(&user, nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)
	server := starter.Server{
		Db:  mockDb,
		Jwt: mockJwt,
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
	expect.Equal("sample_token", res.JWT)
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
	defer mockDbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(nil, nil)
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

func TestServer_RefreshTokenHandlerHandlesRefreshToken_FindOneDbError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(nil, errors.New("db error"))
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := starter.Server{
		Db:     mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerHandlesUser_FindByIdDbError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	userCtrl := gomock.NewController(t)
	defer userCtrl.Finish()
	mockUser := mock.NewMockUserManager(userCtrl)
	mockUser.EXPECT().FindByID("some_user").AnyTimes().Return(nil, errors.New("db error"))
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockToken := &db.RefreshToken{
		Token: "auth_token",
		User:  "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken, nil)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)
	mockDb.EXPECT().Users().AnyTimes().Return(mockUser)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := starter.Server{
		Db:     mockDb,
		Logger: getLogger(),
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerRefreshTokenNotLinkedToUserRespondsWithStatusUnauthorized(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	mockUserManager := mock.NewMockUserManager(mockUsersCtrl)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb := mock.NewMockDB(mockDbCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().Times(1)
	mockUserManager.EXPECT().FindByID("some_user").Times(1).Return(nil, nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockToken := &db.RefreshToken{
		Token: "auth_token",
		User:  "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken, nil)
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

func TestServer_RefreshTokenHandlerHandlesJwtError(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	jwtCtrl := gomock.NewController(t)
	defer jwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(jwtCtrl)
	mockJwt.EXPECT().CreateForUser(gomock.Any()).Times(1).Return("", errors.New("error"))
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
	mockUserManager.EXPECT().FindByID("some_user").Times(1).Return(mockUser, nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockToken := &db.RefreshToken{
		Token: "auth_token",
		User:  "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken, nil)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := starter.Server{
		Db:     mockDb,
		Logger: getLogger(),
		Jwt:    mockJwt,
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusInternalServerError, responseRecorder.Code)
}

func TestServer_RefreshTokenHandlerReturnsJWT(t *testing.T) {
	expect := assert.New(t)
	mockDbCtrl := gomock.NewController(t)
	defer mockDbCtrl.Finish()
	refreshTokenCtrl := gomock.NewController(t)
	defer refreshTokenCtrl.Finish()
	mockUsersCtrl := gomock.NewController(t)
	defer mockUsersCtrl.Finish()
	jwtCtrl := gomock.NewController(t)
	defer jwtCtrl.Finish()
	mockJwt := mock.NewMockJwtManager(jwtCtrl)
	mockJwt.EXPECT().CreateForUser(gomock.Any()).Times(1).Return("sample_token", nil)
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
	mockUserManager.EXPECT().FindByID("some_user").Times(1).Return(mockUser, nil)
	mockDb.EXPECT().Users().Times(1).Return(mockUserManager)
	mockToken := &db.RefreshToken{
		Token: "auth_token",
		User:  "some_user",
	}
	mockRefreshTokenManager.EXPECT().FindOne("auth_token").Times(1).Return(mockToken, nil)
	mockDb.EXPECT().RefreshTokens().Times(1).Return(mockRefreshTokenManager)

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", nil)
	request.Header.Add("Authorization", "Bearer auth_token")

	server := starter.Server{
		Db:     mockDb,
		Logger: getLogger(),
		Jwt:    mockJwt,
	}
	server.RefreshTokenHandler()(responseRecorder, request)
	expect.Equal(http.StatusOK, responseRecorder.Code)
}

// endregion

// region Sessions.List
func TestServer_RefreshTokensListReturnsBadRequestWhenTokenWrong(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jwtManagerCtrl := gomock.NewController(t)
	defer jwtManagerCtrl.Finish()
	mockJwtManager := mock.NewMockJwtManager(jwtManagerCtrl)
	mockRequestHelper := mock.NewMockRequestHelper(ctrl)
	mockJwtManager.EXPECT().GetJwtDataFromRequest(gomock.Any()).Return(nil, errors.New("jwt error"))
	expect := assert.New(t)
	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)
	server := starter.Server{ReqHelper: mockRequestHelper, Jwt: mockJwtManager}
	server.RefreshTokensList()(responseRecorder, request)
	expect.Equal(http.StatusBadRequest, responseRecorder.Code)
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
	jwtHelperCtrl := gomock.NewController(t)
	defer jwtHelperCtrl.Finish()
	mockJwtHelper := mock.NewMockJwtManager(jwtHelperCtrl)
	userID := "some_random_id"
	claims := jwt.Claims{
		ID:          userID,
		Email:       "admin@example.com",
		Permissions: []string{"sudo"},
	}
	mockJwtHelper.EXPECT().GetJwtDataFromRequest(gomock.Any()).Return(&claims, nil)
	mockDb := mock.NewMockDB(dbCtrl)
	mockRefreshTokenManager := mock.NewMockRefreshTokenManager(refreshTokenCtrl)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	mockRefreshTokenManager.EXPECT().List(userID).Return(nil, errors.New("dbError"))
	mockDb.EXPECT().RefreshTokens().AnyTimes().Return(mockRefreshTokenManager)
	token := getJwtForUser([]string{"sudo"})
	// that's valid jwt
	mockRequest := httptest.NewRequest("GET", "/", nil)
	mockRequest.Header.Add("Authorization", "Bearer "+token)
	responseRecorder := httptest.NewRecorder()
	server := starter.Server{
		Db:        mockDb,
		ReqHelper: mockRequestHelper,
		Jwt:       mockJwtHelper,
		Logger:    getLogger(),
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
	jwtManagerCtrl := gomock.NewController(t)
	defer jwtManagerCtrl.Finish()
	mockJwtManager := mock.NewMockJwtManager(jwtManagerCtrl)
	userID := "some_random_id"
	claims := jwt.Claims{
		ID:          userID,
		Email:       "admin@example.com",
		Permissions: []string{"sudo"},
	}
	mockJwtManager.EXPECT().GetJwtDataFromRequest(gomock.Any()).Return(&claims, nil)
	mockDb.EXPECT().Clone().AnyTimes().Return(mockDb)
	mockDb.EXPECT().Close().AnyTimes()
	refreshTokens := []db.RefreshToken{
		{Token: "token1", User: "some_random_id"},
		{Token: "token2", User: "some_random_id"},
	}
	// get list and return
	mockRefreshTokenManager.EXPECT().List(userID).Return(refreshTokens, nil)
	mockDb.EXPECT().RefreshTokens().AnyTimes().Return(mockRefreshTokenManager)
	token := getJwtForUser([]string{"sudo"})
	// that's valid jwt
	mockRequest := httptest.NewRequest("GET", "/", nil)
	mockRequest.Header.Add("Authorization", "Bearer "+token)
	responseRecorder := httptest.NewRecorder()
	server := starter.Server{
		Db:        mockDb,
		ReqHelper: mockRequestHelper,
		Jwt:       mockJwtManager,
		Logger:    getLogger(),
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

package api

import (
	"api-demo/domain"
	mockDomain "api-demo/mock/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func setup(t *testing.T) (UserApi, *fiber.App, *mockDomain.MockUserService) {
	ctrl := gomock.NewController(t)
	userService := mockDomain.NewMockUserService(ctrl)

	app := fiber.New()

	userApi, err := NewUserApi(userService)
	assert.NoError(t, err, "user api creation cannot fail")

	userApi.AddRoutes(app)

	return userApi, app, userService
}

type createUserMatcher struct {
	expectedUser domain.User
}

func (c createUserMatcher) Matches(i interface{}) bool {
	user, ok := i.(domain.User)
	if !ok {
		return false
	}
	c.expectedUser.Id = user.Id
	return c.expectedUser == user
}

func (createUserMatcher) String() string {
	return "create user params"
}

func Test_CreateUser(t *testing.T) {
	user := domain.NewUser("Shashank Pachava", "admin")

	_, app, userService := setup(t)

	userService.EXPECT().CreateUser(context.Background(), createUserMatcher{user}).Return(user, nil)

	// http.Request
	req := httptest.NewRequest(http.MethodPost, "http://acme.com/users", nil).
		WithContext(context.Background())

	createUser := CreateUserDto{
		Name: user.Name,
		Role: user.Role,
	}
	b, err := json.Marshal(createUser)
	assert.NoError(t, err, "cannot fail marshalling")
	req.Body = &nopCloser{bytes.NewReader(b)}
	req.ContentLength = int64(len(b))

	// http.Response
	resp, err := app.Test(req, -1)
	assert.NoError(t, err, "create user api failed")

	createdUser := UserDto{
		Id:   user.Id.String(),
		Name: user.Name,
		Role: user.Role,
	}

	expectedB, err := json.Marshal(createdUser)
	assert.NoError(t, err, "expected no error here")

	respB, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "expected no error here")

	assert.Equal(t, expectedB, respB, "expected different resp body")
}

func Test_GetUserById(t *testing.T) {
	user := domain.NewUser("Shashank Pachava", "admin")

	_, app, userService := setup(t)

	userService.EXPECT().GetUserById(context.Background(), user.Id).Return(user, nil)

	// http.Request
	req := httptest.NewRequest(http.MethodGet, "http://acme.com/users/"+user.Id.String(), nil).
		WithContext(context.Background())

	// http.Response
	resp, err := app.Test(req, -1)
	assert.NoError(t, err, "get user api failed")

	expectedDto := userToDto(user)

	expectedB, err := json.Marshal(expectedDto)
	assert.NoError(t, err, "expected no error here")

	respB, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "expected no error here")

	assert.Equal(t, expectedB, respB, "expected different resp body")
}

type userPropertiesMatcher struct {
	expectedProperties domain.UserProperties
}

func (u userPropertiesMatcher) Matches(i interface{}) bool {
	properties, ok := i.(*domain.UserProperties)
	if !ok {
		return false
	}
	if u.expectedProperties.Name != properties.Name {
		if u.expectedProperties.Name == nil {
			return false
		}
		if properties.Name == nil {
			return false
		}
		if *u.expectedProperties.Name != *properties.Name {
			return false
		}
	}
	if u.expectedProperties.Role != properties.Role {
		if u.expectedProperties.Role == nil {
			return false
		}
		if properties.Role == nil {
			return false
		}
		if *u.expectedProperties.Role != *properties.Role {
			return false
		}
	}

	return true
}

func (u userPropertiesMatcher) String() string {
	return fmt.Sprintf("%#v", u.expectedProperties)
}

func Test_GetUsersByProperty(t *testing.T) {
	user := domain.NewUser("Shashank Pachava", "admin")

	_, app, userService := setup(t)

	name := "Shashank Pachava"

	expectedProperties := domain.UserProperties{
		Name: &name,
	}

	userService.EXPECT().
		GetByProperty(context.Background(), userPropertiesMatcher{expectedProperties: expectedProperties}).
		Return([]domain.User{user}, nil)

	// http.Request
	queries := make(url.Values)
	queries["name"] = []string{name}
	requestUrl, err := url.Parse("http://acme.com/users?" + queries.Encode())
	assert.NoError(t, err, "could not parse url")
	req := httptest.NewRequest(http.MethodGet, requestUrl.String(), nil).
		WithContext(context.Background())

	// http.Response
	resp, err := app.Test(req, -1)
	assert.NoError(t, err, "filter users api failed")

	expectedDto := userToDto(user)

	expectedB, err := json.Marshal([]UserDto{expectedDto})
	assert.NoError(t, err, "expected no error here")

	respB, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "expected no error here")

	assert.Equal(t, string(expectedB), string(respB), "expected different resp body")
}

func Test_GetUsersByProperty_NoUsers(t *testing.T) {
	_, app, userService := setup(t)

	userService.EXPECT().
		GetByProperty(context.Background(), userPropertiesMatcher{expectedProperties: domain.UserProperties{}}).
		Return(nil, nil)

	// http.Request
	req := httptest.NewRequest(http.MethodGet, "http://acme.com/users", nil).
		WithContext(context.Background())

	// http.Response
	resp, err := app.Test(req, -1)
	assert.NoError(t, err, "filter users api failed")

	expectedB, err := json.Marshal([]UserDto{})
	assert.NoError(t, err, "expected no error here")

	respB, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "expected no error here")

	assert.Equal(t, string(expectedB), string(respB), "expected different resp body")
}

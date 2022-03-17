package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	dq "doublequote/pkg/domain"
	"doublequote/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_handleRegister(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()
		s.UserService.On(
			"FindUsers",
			mock.Anything,
			dq.UserFilter{Email: utils.StringPtr("test@example.com"),
				Limit: 1,
			},
			dq.UserInclude{},
		).Return([]*dq.User{}, 0, nil)

		req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"email": "test@example.com", "password": "password"}`))
		assert.Nil(t, err)

		expect := dq.User{
			ID:       0,
			Email:    "test@example.com",
			Password: "password",
		}
		ret := dq.User{
			ID:        1,
			Email:     "test@example.com",
			Password:  "password",
			CreatedAt: utils.MustParseTime(t, "2021-06-06"),
		}
		s.UserService.On("CreateUser", req.Context(), &expect).Return(&ret, nil)

		rr := MakeRequest(req, s.handleRegister)
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, `{"data": {"require_email_verification": false} }`, rr.Body.String())
	})

	t.Run("CreateReturnsError", func(t *testing.T) {
		s := NewTestServer()
		s.UserService.On("FindUsers", mock.Anything, mock.Anything, dq.UserInclude{}).Return([]*dq.User{}, 0, nil)

		req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"email": "test@example.com", "password": "password"}`))
		assert.Nil(t, err)

		s.UserService.
			On("CreateUser", mock.Anything, mock.Anything).
			Return(&dq.User{}, errors.New("sql: you should have used mongo"))

		rr := MakeRequest(req, s.handleRegister)
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"title": "Internal error.", "type": "about:blank"}`, rr.Body.String())
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		s := NewTestServer()
		s.UserService.On(
			"FindUsers",
			mock.Anything,
			dq.UserFilter{
				Email: utils.StringPtr("test"),
				Limit: 1,
			},
			dq.UserInclude{},
		).Return([]*dq.User{}, 0, nil)

		req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"email": "test", "password": "password"}`))
		assert.Nil(t, err)

		rr := MakeRequest(req, s.handleRegister)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"title": "Validation failed.", "type": "about:blank", "invalid_params": [{"name": "email", "reason": "Email is required."}]}`, rr.Body.String())
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		s := NewTestServer()
		s.UserService.On("FindUsers", mock.Anything, mock.Anything, dq.UserInclude{}).Return([]*dq.User{}, 0, nil)

		req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"email": "test@example.com", "password": "short"}`))
		assert.Nil(t, err)

		rr := MakeRequest(req, s.handleRegister)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"title": "Validation failed.", "type": "about:blank", "invalid_params": [{"name": "password", "reason": "Password must be greater than 6 characters and less than 64."}]}`, rr.Body.String())
	})

	t.Run("EmailTaken", func(t *testing.T) {
		s := NewTestServer()

		s.UserService.On("FindUsers", mock.Anything, dq.UserFilter{Email: utils.StringPtr("test@example.com"), Limit: 1}, dq.UserInclude{}).Return([]*dq.User{{}}, 1, nil)

		req, err := http.NewRequest("POST", "/register", strings.NewReader(`{"email": "test@example.com", "password": "password"}`))
		assert.Nil(t, err)

		rr := MakeRequest(req, s.handleRegister)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"title": "Validation failed.", "type": "about:blank", "invalid_params": [{"name": "email", "reason": "Email is taken."}]}`, rr.Body.String())
	})
}

func TestServer_handleEmailVerification(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()
		s.now = func() time.Time {
			return utils.MustParseTime(t, "2021-06-06")
		}

		s.CryptoService.On("VerifyToken", "ABC").Return(map[string]interface{}{"id": 1}, nil)

		req, err := http.NewRequest("POST", "/verify-email", strings.NewReader(`{"token": "ABC"}`))
		assert.Nil(t, err)

		s.UserService.
			On(
				"UpdateUser",
				req.Context(),
				1,
				dq.UserUpdate{EmailVerifiedAt: utils.TimePtr(utils.MustParseTime(t, "2021-06-06"))},
			).
			Return(&dq.User{}, nil)

		rr := MakeRequest(req, s.handleEmailVerification)
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	})
}

func TestServer_handleLogin(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email": "test@example.com", "password": "password"}`))
		assert.Nil(t, err)

		s.UserService.
			On(
				"FindUser",
				req.Context(),
				dq.UserFilter{Email: utils.StringPtr("test@example.com")},
				dq.UserInclude{},
			).
			Return(&dq.User{ID: 1, Email: "test@example.com", Password: "hashed-password"}, nil)

		s.CryptoService.On("VerifyPassword", "hashed-password", "password").Return(true)

		s.SessionService.On("Create", rr, req, 1).Return(nil, nil)

		h := http.HandlerFunc(s.handleLogin)
		h.ServeHTTP(rr, req)

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		s := NewTestServer()
		s.now = func() time.Time {
			return utils.MustParseTime(t, "2021-06-06")
		}

		req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email": "test", "password": "password"}`))
		assert.Nil(t, err)

		rr := MakeRequest(req, s.handleLogin)

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"title": "Validation failed.", "type": "about:blank", "invalid_params": [{"name": "email", "reason": "Email is required."}]}`, rr.Body.String())
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		s := NewTestServer()
		s.now = func() time.Time {
			return utils.MustParseTime(t, "2021-06-06")
		}

		req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email": "test@example.com", "password": "2shrt"}`))
		assert.Nil(t, err)

		rr := MakeRequest(req, s.handleLogin)

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"title": "Validation failed.", "type": "about:blank", "invalid_params": [{"name": "password", "reason": "Password must be greater than 6 characters and less than 64."}]}`, rr.Body.String())
	})

	t.Run("IncorrectPassword", func(t *testing.T) {
		s := NewTestServer()

		req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email": "test@example.com", "password": "password"}`))
		assert.Nil(t, err)

		s.UserService.
			On(
				"FindUser",
				req.Context(),
				dq.UserFilter{Email: utils.StringPtr("test@example.com")},
				dq.UserInclude{},
			).
			Return(&dq.User{ID: 1, Email: "test@example.com", Password: "hashed-password"}, nil)

		s.CryptoService.On("VerifyPassword", "hashed-password", "password").Return(false)

		rr := MakeRequest(req, s.handleLogin)

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"title": "A user with these credentials could not be found.", "type": "about:blank"}`, rr.Body.String())
	})

	t.Run("IncorrectEmail", func(t *testing.T) {
		s := NewTestServer()

		req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"email": "notfound@example.com", "password": "password"}`))
		assert.Nil(t, err)

		s.UserService.
			On(
				"FindUser",
				req.Context(),
				dq.UserFilter{Email: utils.StringPtr("notfound@example.com")},
				dq.UserInclude{},
			).
			Return(nil, nil)

		rr := MakeRequest(req, s.handleLogin)

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"title": "A user with these credentials could not be found.", "type": "about:blank"}`, rr.Body.String())
	})
}

func TestServer_handleProfile(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		req, err := http.NewRequest("GET", "/profile", nil)
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleProfile, &dq.User{Email: "test@example.com"})

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"data": {"email": "test@example.com"} }`, rr.Body.String())
	})
}

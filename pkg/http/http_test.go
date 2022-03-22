package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"doublequote/pkg/domain"
	mock2 "doublequote/pkg/mock"
	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	*Server

	UserService       *mock2.UserService
	CryptoService     *mock2.CryptoService
	SessionService    *mock2.SessionService
	CollectionService *mock2.CollectionService
	FeedService       *mock2.FeedService
}

func NewTestServer() *TestServer {
	s := TestServer{Server: NewServer()}

	s.CryptoService = &mock2.CryptoService{}
	s.Server.cryptoService = s.CryptoService

	s.UserService = &mock2.UserService{}
	s.Server.userService = s.UserService

	s.SessionService = &mock2.SessionService{}
	s.Server.sessionService = s.SessionService

	s.CollectionService = &mock2.CollectionService{}
	s.Server.collectionService = s.CollectionService

	s.FeedService = &mock2.FeedService{}
	s.Server.feedService = s.FeedService

	return &s
}

// MakeRequest is a utility function for making a request to a handler and returning its result
func MakeRequest(r *http.Request, h http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, r)

	return rr
}

// MakeAuthenticatedRequest is a utility function for making an authenticated request to a handler and returning its result
func MakeAuthenticatedRequest(r *http.Request, h http.HandlerFunc, u *domain.User) *httptest.ResponseRecorder {
	r = r.WithContext(domain.NewContextWithUser(r.Context(), u))

	return MakeRequest(r, h)
}

// MakeMiddlewareRequest is a utility function for making a request to a middleware
// It returns the next handlers http.ResponseWriter, http.Request and whether the middleware actually accepted the request
func MakeMiddlewareRequest(r *http.Request, mw func(next http.Handler) http.Handler) (mwWriter http.ResponseWriter, mwReq *http.Request, next bool) {
	var route http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mwWriter = w
		mwReq = r
		next = true
	})

	h := mw(route)

	h.ServeHTTP(httptest.NewRecorder(), r)

	return
}

func TestValidationError(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		w := httptest.NewRecorder()

		errs := map[string][]string{
			"password": {"Password is taken by another user. Try again."}, // Don't do this!
		}

		ValidationError(w, errs)

		assert.JSONEq(t, `{"title": "Validation failed.", "type": "about:blank", "invalid_params": [{"name": "password", "reason": "Password is taken by another user. Try again."}]}`, w.Body.String())
	})
}

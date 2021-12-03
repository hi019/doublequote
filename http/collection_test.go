package http

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	dq "doublequote"
	"doublequote/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_handleCreateCollection(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("CreateCollection", mock.Anything, &dq.Collection{Name: "Test", UserID: 0}).
			Return(&dq.Collection{
				ID:   1,
				Name: "Test",
			}, nil)

		req, err := http.NewRequest("POST", "", strings.NewReader(`{"name": "Test"}`))
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleCreateCollection, &dq.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, `{"data": {"id": 1, "name": "Test"} }`, rr.Body.String())
	})

	t.Run("DbErr", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("CreateCollection", mock.Anything, &dq.Collection{Name: "Test", UserID: 0}).
			Return(nil, fmt.Errorf("sqlite: /dev/null does not support sqlite"))

		req, err := http.NewRequest("POST", "", strings.NewReader(`{"name": "Test"}`))
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleCreateCollection, &dq.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"title": "Internal error.", "type": "about:blank"}`, rr.Body.String())
	})
}

func TestServer_handleListCollections(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("FindCollections", mock.Anything, dq.CollectionFilter{UserID: utils.IntPtr(0), Limit: 100}, dq.CollectionInclude{}).
			Return([]*dq.Collection{{ID: 1, Name: "Test"}}, 1, nil)

		req, err := http.NewRequest("GET", "", nil)
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleListCollections, &dq.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{ "data": { "collections": [{"id": 1, "name": "Test"}] } }`, rr.Body.String())
	})

	t.Run("DbErr", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("FindCollections", mock.Anything, dq.CollectionFilter{UserID: utils.IntPtr(0), Limit: 100}, dq.CollectionInclude{}).
			Return(nil, 0, fmt.Errorf("mongo: unexpected query 'SELECT'"))

		req, err := http.NewRequest("GET", "", nil)
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleListCollections, &dq.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"title": "Internal error.", "type": "about:blank"}`, rr.Body.String())
	})
}

func TestServer_handleGetCollectionsFeeds(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		s.FeedService.
			On("FindFeeds", mock.Anything, dq.FeedFilter{CollectionID: utils.IntPtr(1), Limit: 500}, dq.FeedInclude{}).
			Return([]*dq.Feed{{ID: 1, Name: "Test", Domain: "test.com"}}, 1, nil)

		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, dq.CollectionInclude{}).
			Return(&dq.Collection{}, nil)

		req, err := http.NewRequest("GET", "", nil)
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleGetCollectionFeeds, &dq.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{ "data": { "feeds": [{"id": 1, "name": "Test", "domain": "test.com"}] } }`, rr.Body.String())
	})

	t.Run("DbErr", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("FindCollections", mock.Anything, dq.CollectionFilter{UserID: utils.IntPtr(0), Limit: 100}, dq.CollectionInclude{}).
			Return(nil, 0, fmt.Errorf("mongo: unexpected query 'SELECT'"))

		req, err := http.NewRequest("GET", "", nil)
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleListCollections, &dq.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"title": "Internal error.", "type": "about:blank"}`, rr.Body.String())
	})

	t.Run("CollectionNotFound", func(t *testing.T) {
		s := NewTestServer()

		s.FeedService.
			On("FindFeeds", mock.Anything, dq.FeedFilter{CollectionID: utils.IntPtr(1), Limit: 500}, dq.FeedInclude{}).
			Return([]*dq.Feed{{ID: 1, Name: "Test", Domain: "test.com"}}, 1, nil)

		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, dq.CollectionInclude{}).
			Return(&dq.Collection{}, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Collection"))

		req, err := http.NewRequest("GET", "", nil)
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleGetCollectionFeeds, &dq.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"title":"Collection not found.", "type":"about:blank"}`, rr.Body.String())
	})

	t.Run("CollectionOwnedByOtherUser", func(t *testing.T) {
		s := NewTestServer()

		s.FeedService.
			On("FindFeeds", mock.Anything, dq.FeedFilter{CollectionID: utils.IntPtr(1), Limit: 500}, dq.FeedInclude{}).
			Return([]*dq.Feed{{ID: 1, Name: "Test", Domain: "test.com"}}, 1, nil)

		// Return collection owned by other user
		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, dq.CollectionInclude{}).
			Return(&dq.Collection{UserID: 2}, nil)

		req, err := http.NewRequest("GET", "", nil)
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleGetCollectionFeeds, &dq.User{ID: 1})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"title":"Collection not found.", "type":"about:blank"}`, rr.Body.String())
	})
}

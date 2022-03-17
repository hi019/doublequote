package http

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"doublequote/pkg/domain"
	"doublequote/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_handleCreateCollection(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("CreateCollection", mock.Anything, &domain.Collection{Name: "Test", UserID: 0}).
			Return(&domain.Collection{
				ID:   1,
				Name: "Test",
			}, nil)

		req, err := http.NewRequest("POST", "", strings.NewReader(`{"name": "Test"}`))
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleCreateCollection, &domain.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, `{"data": {"id": 1, "name": "Test"} }`, rr.Body.String())
	})

	t.Run("DbErr", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("CreateCollection", mock.Anything, &domain.Collection{Name: "Test", UserID: 0}).
			Return(nil, fmt.Errorf("sqlite: /dev/null does not support sqlite"))

		req, err := http.NewRequest("POST", "", strings.NewReader(`{"name": "Test"}`))
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleCreateCollection, &domain.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"title": "Internal error.", "type": "about:blank"}`, rr.Body.String())
	})
}

func TestServer_handleListCollections(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("FindCollections", mock.Anything, domain.CollectionFilter{UserID: utils.IntPtr(0), Limit: 100}, domain.CollectionInclude{}).
			Return([]*domain.Collection{{ID: 1, Name: "Test"}}, 1, nil)

		req, err := http.NewRequest("GET", "", nil)
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleListCollections, &domain.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{ "data": { "collections": [{"id": 1, "name": "Test"}] } }`, rr.Body.String())
	})

	t.Run("DbErr", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("FindCollections", mock.Anything, domain.CollectionFilter{UserID: utils.IntPtr(0), Limit: 100}, domain.CollectionInclude{}).
			Return(nil, 0, fmt.Errorf("mongo: unexpected query 'SELECT'"))

		req, err := http.NewRequest("GET", "", nil)
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleListCollections, &domain.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"title": "Internal error.", "type": "about:blank"}`, rr.Body.String())
	})
}

func TestServer_handleGetCollectionsFeeds(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		s.FeedService.
			On("FindFeeds", mock.Anything, domain.FeedFilter{CollectionID: utils.IntPtr(1), Limit: 500}, domain.FeedInclude{}).
			Return([]*domain.Feed{{ID: 1, Name: "Test", Domain: "test.com"}}, 1, nil)

		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, domain.CollectionInclude{}).
			Return(&domain.Collection{}, nil)

		req, err := http.NewRequest("GET", "", nil)
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleGetCollectionFeeds, &domain.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{ "data": { "feeds": [{"id": 1, "name": "Test", "domain": "test.com"}] } }`, rr.Body.String())
	})

	t.Run("DbErr", func(t *testing.T) {
		s := NewTestServer()

		s.CollectionService.
			On("FindCollections", mock.Anything, domain.CollectionFilter{UserID: utils.IntPtr(0), Limit: 100}, domain.CollectionInclude{}).
			Return(nil, 0, fmt.Errorf("mongo: unexpected query 'SELECT'"))

		req, err := http.NewRequest("GET", "", nil)
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleListCollections, &domain.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"title": "Internal error.", "type": "about:blank"}`, rr.Body.String())
	})

	t.Run("CollectionNotFound", func(t *testing.T) {
		s := NewTestServer()

		s.FeedService.
			On("FindFeeds", mock.Anything, domain.FeedFilter{CollectionID: utils.IntPtr(1), Limit: 500}, domain.FeedInclude{}).
			Return([]*domain.Feed{{ID: 1, Name: "Test", Domain: "test.com"}}, 1, nil)

		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, domain.CollectionInclude{}).
			Return(&domain.Collection{}, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection"))

		req, err := http.NewRequest("GET", "", nil)
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleGetCollectionFeeds, &domain.User{Email: "test@example.com"})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"title":"Collection not found.", "type":"about:blank"}`, rr.Body.String())
	})

	t.Run("CollectionOwnedByOtherUser", func(t *testing.T) {
		s := NewTestServer()

		s.FeedService.
			On("FindFeeds", mock.Anything, domain.FeedFilter{CollectionID: utils.IntPtr(1), Limit: 500}, domain.FeedInclude{}).
			Return([]*domain.Feed{{ID: 1, Name: "Test", Domain: "test.com"}}, 1, nil)

		// Return collection owned by other user
		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, domain.CollectionInclude{}).
			Return(&domain.Collection{UserID: 2}, nil)

		req, err := http.NewRequest("GET", "", nil)
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		rr := MakeAuthenticatedRequest(req, s.handleGetCollectionFeeds, &domain.User{ID: 1})
		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"title":"Collection not found.", "type":"about:blank"}`, rr.Body.String())
	})
}

func TestServer_handlePutCollectionFeeds(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		s := NewTestServer()

		// Setup mocks
		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, domain.CollectionInclude{}).
			Return(&domain.Collection{}, nil)

		s.CollectionService.
			On("UpdateCollection", mock.Anything, 1, domain.CollectionUpdate{FeedsIDs: &[]int{1, 2}}).
			Return(&domain.Collection{}, nil)

		// Setup request
		req, err := http.NewRequest("PUT", "", strings.NewReader(`{"feeds": [1, 2]}`))
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		// Make request
		rr := MakeAuthenticatedRequest(req, s.handlePutCollectionFeeds, &domain.User{})

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{ "data": { "feeds": [1,2] } }`, rr.Body.String())
	})

	t.Run("OtherUsersCollection", func(t *testing.T) {
		s := NewTestServer()

		// Setup mocks
		s.CollectionService.
			On("FindCollectionByID", mock.Anything, 1, domain.CollectionInclude{}).
			Return(&domain.Collection{}, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection"))

		// Setup request
		req, err := http.NewRequest("PUT", "", strings.NewReader(`{"feeds": [1, 2]}`))
		utils.AddParamToContext(t, req, "collectionID", "1")
		assert.Nil(t, err)

		// Make request
		rr := MakeAuthenticatedRequest(req, s.handlePutCollectionFeeds, &domain.User{})

		s.UserService.AssertExpectations(t)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{ "title": "Collection not found.", "type": "about:blank" }`, rr.Body.String())
	})
}

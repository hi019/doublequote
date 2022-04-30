package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"doublequote/pkg/domain"
	"doublequote/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

func (s *Server) registerCollectionRoutes(r chi.Router) {
	r.Get("/collections/{collectionID}/feeds", s.handleGetCollectionFeeds)
	r.Put("/collections/{collectionID}/feeds", s.handlePutCollectionFeeds)
	r.Get("/collections/{collectionID}/entries", s.handleGetCollectionEntries)
	r.Get("/collections", s.handleListCollections)
	r.Post("/collections", s.handleCreateCollection)
}

type createCollectionRequest struct {
	Name string `json:"name"`
}

type createCollectionResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func (s *Server) handleCreateCollection(w http.ResponseWriter, r *http.Request) {
	var req createCollectionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.ErrInvalidJSONBody))
		return
	}

	errors := validate.Validate(
		&validators.StringLengthInRange{
			Name:    "name",
			Field:   req.Name,
			Min:     1,
			Max:     64,
			Message: fmt.Sprintf(domain.ErrFieldGTEAndLTE, "Name", 1, 64),
		},
	)
	if errors.HasAny() {
		ValidationError(w, errors.Errors)
		return
	}

	created, err := s.collectionService.CreateCollection(r.Context(), &domain.Collection{
		Name:   req.Name,
		UserID: domain.UserIDFromContext(r.Context()),
	})
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	sendJSON(w, r, createCollectionResponse{Name: req.Name, ID: created.ID})
}

type listCollectionsResponse struct {
	Collections []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"collections"`
}

func (s *Server) handleListCollections(w http.ResponseWriter, r *http.Request) {
	filter := domain.CollectionFilter{
		UserID: utils.IntPtr(domain.UserIDFromContext(r.Context())),
		Limit:  100,
	}
	c, _, err := s.collectionService.FindCollections(r.Context(), filter, domain.CollectionInclude{})
	if err != nil {
		Error(w, r, err)
		return
	}

	var res listCollectionsResponse
	for _, col := range c {
		s := struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{ID: col.ID, Name: col.Name}
		res.Collections = append(res.Collections, s)
	}

	w.WriteHeader(http.StatusOK)
	sendJSON(w, r, res)
}

type feedResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

type getCollectionFeedsResponse struct {
	Feeds []feedResponse `json:"feeds"`
}

func (s *Server) handleGetCollectionFeeds(w http.ResponseWriter, r *http.Request) {
	colID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		Error(w, r, err)
		return
	}

	// Make sure the requesting user owns the collection
	if col, err := s.collectionService.FindCollectionByID(r.Context(), colID, domain.CollectionInclude{}); err != nil {
		Error(w, r, err)
		return
	} else if col.UserID != domain.UserIDFromContext(r.Context()) {
		Error(w, r, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection"))
		return
	}

	filter := domain.FeedFilter{
		CollectionID: utils.IntPtr(colID),
		Limit:        500,
	}
	feeds, _, err := s.feedService.FindFeeds(r.Context(), filter, domain.FeedInclude{})
	if err != nil {
		Error(w, r, err)
		return
	}

	var res getCollectionFeedsResponse
	res.Feeds = []feedResponse{}
	for _, feed := range feeds {
		s := feedResponse{
			ID:     feed.ID,
			Name:   feed.Name,
			Domain: feed.Domain,
		}
		res.Feeds = append(res.Feeds, s)
	}

	w.WriteHeader(http.StatusOK)
	sendJSON(w, r, res)
}

type putCollectionFeedsRequest struct {
	Feeds []int `json:"feeds"`
}

func (s *Server) handlePutCollectionFeeds(w http.ResponseWriter, r *http.Request) {
	colID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		Error(w, r, err)
		return
	}

	// Make sure the requesting user owns the collection
	if col, err := s.collectionService.FindCollectionByID(r.Context(), colID, domain.CollectionInclude{}); err != nil {
		Error(w, r, err)
		return
	} else if col.UserID != domain.UserIDFromContext(r.Context()) {
		Error(w, r, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection"))
		return
	}

	// Parse and validate request body
	var req putCollectionFeedsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.ErrInvalidJSONBody))
		return
	}

	// Update collection
	_, err = s.collectionService.UpdateCollection(r.Context(), colID, domain.CollectionUpdate{FeedsIDs: &req.Feeds})
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	sendJSON(w, r, req)
}

type getCollectionEntriesResponse struct {
	CollectionEntries []*domain.CollectionEntry `json:"collection_entries"`
}

func (s *Server) handleGetCollectionEntries(w http.ResponseWriter, r *http.Request) {
	colID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		Error(w, r, err)
		return
	}

	// Make sure the requesting user owns the collection
	if col, err := s.collectionService.FindCollectionByID(r.Context(), colID, domain.CollectionInclude{}); err != nil {
		Error(w, r, err)
		return
	} else if col.UserID != domain.UserIDFromContext(r.Context()) {
		Error(w, r, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection"))
		return
	}

	entries, _, err := s.collectionEntryService.FindCollectionEntries(
		r.Context(),
		domain.CollectionEntryFilter{CollectionID: utils.Ptr(colID)},
		domain.CollectionEntryInclude{Entry: true},
	)
	if err != nil {
		Error(w, r, err)
		return
	}

	// So that the JSON array in the response isn't null if entries is empty
	res := append([]*domain.CollectionEntry{}, entries...)
	sendJSON(w, r, getCollectionEntriesResponse{res})
}

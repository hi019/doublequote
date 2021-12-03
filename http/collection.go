package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dq "doublequote"
	"doublequote/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

func (s *Server) registerCollectionRoutes(r chi.Router) {
	r.Get("/collection", s.handleListCollections)
	r.Post("/collection", s.handleCreateCollection)
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
		Error(w, r, dq.Errorf(dq.EINVALID, dq.ErrInvalidJSONBody))
		return
	}

	errors := validate.Validate(
		&validators.StringLengthInRange{
			Name:    "name",
			Field:   req.Name,
			Min:     1,
			Max:     64,
			Message: fmt.Sprintf(dq.ErrFieldGTEAndLTE, "Name", 1, 64),
		},
	)
	if errors.HasAny() {
		ValidationError(w, errors.Errors)
		return
	}

	created, err := s.CollectionService.CreateCollection(r.Context(), &dq.Collection{
		Name:   req.Name,
		UserID: dq.UserIDFromContext(r.Context()),
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
	filter := dq.CollectionFilter{
		UserID: utils.IntPtr(dq.UserIDFromContext(r.Context())),
		Limit:  100,
	}
	c, _, err := s.CollectionService.FindCollections(r.Context(), filter, dq.CollectionInclude{})
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

type getCollectionFeedsResponse struct {
	Feeds []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Domain string `json:"domain"`
	} `json:"feeds"`
}

func (s *Server) handleGetCollectionFeeds(w http.ResponseWriter, r *http.Request) {
	colID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		Error(w, r, err)
		return
	}

	// Make sure the requesting user owns the collection
	if col, err := s.CollectionService.FindCollectionByID(r.Context(), colID, dq.CollectionInclude{}); err != nil {
		Error(w, r, err)
		return
	} else if col.UserID != dq.UserIDFromContext(r.Context()) {
		Error(w, r, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Collection"))
		return
	}

	filter := dq.FeedFilter{
		CollectionID: utils.IntPtr(colID),
		Limit:        500,
	}
	feeds, _, err := s.FeedService.FindFeeds(r.Context(), filter, dq.FeedInclude{})
	if err != nil {
		Error(w, r, err)
		return
	}

	var res getCollectionFeedsResponse
	for _, feed := range feeds {
		s := struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Domain string `json:"domain"`
		}{
			ID:     feed.ID,
			Name:   feed.Name,
			Domain: feed.Domain,
		}
		res.Feeds = append(res.Feeds, s)
	}

	w.WriteHeader(http.StatusOK)
	sendJSON(w, r, res)
}

//func (s *Server) handlePutCollectionFeeds(w http.ResponseWriter, r *http.Request) {
//	colID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
//	if err != nil {
//		Error(w, r, err)
//		return
//	}
//
//	//s.CollectionService.UpdateCollection(r.Context(), colID, )
//
//	w.WriteHeader(http.StatusOK)
//	sendJSON(w, r, res)
//}

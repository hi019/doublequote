package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"doublequote/pkg/domain"

	"github.com/go-chi/chi/v5"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

func (s *Server) registerFeedRoutes(r chi.Router) {
	//r.Post("/feeds/{feedID}/ingest", s.handleIngestFeed)
	r.Post("/feeds", s.handleCreateFeed)
}

type createFeedRequest struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	RssURL string `json:"rss_url"`
}

type createFeedResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	RssURL string `json:"rss_url"`
}

func (s *Server) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req createFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.ErrInvalidJSONBody))
		return
	}

	// Validate request
	errors := validate.Validate(
		&validators.StringLengthInRange{
			Name:    "name",
			Field:   req.Name,
			Message: fmt.Sprintf(domain.ErrFieldRequired, "Name"),
		},
		&validators.StringLengthInRange{
			Name:    "domain",
			Field:   req.Domain,
			Min:     4,
			Max:     256,
			Message: fmt.Sprintf(domain.ErrFieldGTEAndLTE, "Domain", 4, 256),
		},
		&validators.StringLengthInRange{
			Name:    "rssURL",
			Field:   req.Domain,
			Min:     4,
			Max:     256,
			Message: fmt.Sprintf(domain.ErrFieldGTEAndLTE, "rssURL", 4, 256),
		},
	)
	if errors.HasAny() {
		ValidationError(w, errors.Errors)
		return
	}

	// Create feed
	f := domain.Feed{
		Name:   req.Name,
		RssURL: req.RssURL,
		Domain: req.Domain,
	}
	created, err := s.feedService.CreateFeed(r.Context(), &f)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Send response
	w.WriteHeader(http.StatusCreated)
	sendJSON(w, r, createFeedResponse{
		ID:     created.ID,
		Name:   created.Name,
		RssURL: created.RssURL,
		Domain: created.Domain,
	})
}

//func (s *Server) handleIngestFeed(w http.ResponseWriter, r *http.Request) {
//	feedID, err := strconv.Atoi(chi.URLParam(r, "feedID"))
//	if err != nil {
//		Error(w, r, err)
//		return
//	}
//
//	feed, err := s.feedService.FindFeedByID(r.Context(), feedID, domain.FeedInclude{})
//	if err != nil {
//		Error(w, r, err)
//		return
//	}
//
//	err = s.ingestService.Ingest(*feed)
//}

// TODO delete, update for feeds

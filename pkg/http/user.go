package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"doublequote/pkg/domain"
	"doublequote/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

func (s *Server) registerPublicUserRoutes(r chi.Router) {
	r.Post("/register", s.handleRegister)
	r.Post("/verify-email", s.handleEmailVerification)
	r.Post("/login", s.handleLogin)
}

func (s *Server) registerUserRoutes(r chi.Router) {
	r.Get("/me", s.handleProfile)
	r.Get("/authcheck", s.handleAuthCheck)
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	RequireEmailVerification bool `json:"require_email_verification"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.ErrInvalidJSONBody))
		return
	}

	errors := validate.Validate(
		&validators.EmailIsPresent{
			Name:    "email",
			Field:   req.Email,
			Message: fmt.Sprintf(domain.ErrFieldRequired, "Email"),
		},
		&validators.StringLengthInRange{
			Name:    "password",
			Field:   req.Password,
			Min:     6,
			Max:     64,
			Message: fmt.Sprintf(domain.ErrFieldGTEAndLTE, "Password", 6, 64),
		},
		&EmailIsUnique{
			Email:       req.Email,
			UserService: s.userService,
		},
	)
	if errors.HasAny() {
		ValidationError(w, errors.Errors)
		return
	}

	u := domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	_, err = s.userService.CreateUser(r.Context(), &u)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	sendJSON(w, r, registerResponse{s.config.App.RequireEmailVerification})
}

type emailVerificationRequest struct {
	Token string `json:"token"`
}

func (s *Server) handleEmailVerification(w http.ResponseWriter, r *http.Request) {
	var req emailVerificationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.ErrInvalidJSONBody))
		return
	}

	errors := validate.Validate(
		&validators.StringIsPresent{
			Name:    "token",
			Field:   req.Token,
			Message: fmt.Sprintf(domain.ErrFieldRequired, "Token"),
		},
	)
	if errors.HasAny() {
		ValidationError(w, errors.Errors)
		return
	}

	// Validate JWT
	data, err := s.cryptoService.VerifyToken(req.Token)
	if err != nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.ErrInvalidJSONBody))
		return
	}

	// Update user
	_, err = s.userService.UpdateUser(r.Context(), data["id"].(int), domain.UserUpdate{EmailVerifiedAt: utils.TimePtr(s.now())})
	if err != nil {
		Error(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.ErrInvalidJSONBody))
		return
	}

	errors := validate.Validate(
		&validators.EmailIsPresent{
			Name:    "email",
			Field:   req.Email,
			Message: fmt.Sprintf(domain.ErrFieldRequired, "Email"),
		},
		&validators.StringLengthInRange{
			Name:    "password",
			Field:   req.Password,
			Min:     6,
			Max:     64,
			Message: fmt.Sprintf(domain.ErrFieldGTEAndLTE, "Password", 6, 64),
		},
	)
	if errors.HasAny() {
		ValidationError(w, errors.Errors)
		return
	}

	u, err := s.userService.FindUser(r.Context(), domain.UserFilter{Email: utils.StringPtr(req.Email)}, domain.UserInclude{})
	if err != nil {
		Error(w, r, err)
		return
	}
	// TODO an attacker could determine whether an email is registered or not through a timing attack
	if u == nil {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.UserNotFound))
		return
	}

	correct := s.cryptoService.VerifyPassword(u.Password, req.Password)
	if !correct {
		Error(w, r, domain.Errorf(domain.EINVALID, domain.UserNotFound))
		return
	}

	_, err = s.sessionService.Create(w, r, u.ID)
	if err != nil {
		Error(w, r, err)
		return
	}

	r = r.WithContext(domain.NewContextWithUser(r.Context(), u))
}

type profileResponse struct {
	Email string `json:"email"`
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	u := domain.UserFromContext(r.Context())

	w.WriteHeader(http.StatusOK)
	sendJSON(w, r, profileResponse{u.Email})
}

// TODO test
// If the user is not authenticated, it will be caught by the auth middleware and will result in a 401.
func (s *Server) handleAuthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

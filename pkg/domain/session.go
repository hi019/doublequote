package domain

import "net/http"

type Session interface {
	Delete() error
	UserID() int
}

type SessionService interface {
	Get(r *http.Request) (Session, error)
	Create(w http.ResponseWriter, r *http.Request, uid int) (Session, error)
}

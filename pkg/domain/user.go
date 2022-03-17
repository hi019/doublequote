package domain

import (
	"context"
	"time"
)

// User represents a user.
type User struct {
	ID int

	Email    string
	Password string

	Collections []*Collection

	EmailVerifiedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserService represents a service for managing users.
type UserService interface {
	// FindUserByID retrieves a user by ID along with their associated auth objects.
	// Returns ENOTFOUND if user does not exist.
	FindUserByID(ctx context.Context, id int, include UserInclude) (*User, error)

	// FindUsers retrieves a list of users by filter. Also returns total count of matching
	// users which may differ from returned results if filter.Limit is specified.
	FindUsers(ctx context.Context, filter UserFilter, include UserInclude) ([]*User, int, error)

	// FindUser retrieves a single user based on a filter
	FindUser(ctx context.Context, filter UserFilter, include UserInclude) (*User, error)

	// CreateUser creates a new user.
	CreateUser(ctx context.Context, user *User) (*User, error)

	// UpdateUser updates a user object. Returns EUNAUTHORIZED if current user is not
	// the user that is being updated. Returns ENOTFOUND if user does not exist.
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)

	// DeleteUser permanently deletes a user and all owned dials. Returns EUNAUTHORIZED
	// if current user is not the user being deleted. Returns ENOTFOUND if
	// user does not exist.
	DeleteUser(ctx context.Context, id int) error
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	// Filtering fields.
	ID    *int
	Email *string

	// Restrict to subset of results.
	Offset int
	Limit  int
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Email           *string
	Password        *string
	EmailVerifiedAt *time.Time
}

type UserInclude struct {
	Collections bool
}

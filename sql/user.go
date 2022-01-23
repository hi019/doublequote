package sql

import (
	"context"

	dq "doublequote"
	"doublequote/ent"
	"doublequote/ent/user"
)

// Ensure service implements interface.
var _ dq.UserService = (*UserService)(nil)

// UserService represents a service for managing users.
type UserService struct {
	sql           *SQL
	eventService  dq.EventService
	cryptoService dq.CryptoService
}

// NewUserService returns a new instance of UserService.
func NewUserService(db *SQL, es dq.EventService, cr dq.CryptoService) *UserService {
	return &UserService{
		sql:           db,
		eventService:  es,
		cryptoService: cr,
	}
}

func (s *UserService) FindUserByID(ctx context.Context, id int, include dq.UserInclude) (*dq.User, error) {
	dbU, err := s.sql.client.User.
		Query().
		With(withUserInclude(include)).
		Where(user.IDEQ(id)).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "User")
	}
	if err != nil {
		return nil, err
	}

	u := sqlUserToDQUser(dbU)
	return u, nil
}

func (s *UserService) FindUsers(ctx context.Context, filter dq.UserFilter, include dq.UserInclude) ([]*dq.User, int, error) {
	u, err := s.sql.client.User.Query().
		With(withUserInclude(include)).
		Where(ifPresent(user.IDEQ, filter.ID), ifPresent(user.EmailEQ, filter.Email)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		All(ctx)

	return sqlUserSliceToDQUserSlice(u), len(u), err
}

func (s *UserService) FindUser(ctx context.Context, filter dq.UserFilter, include dq.UserInclude) (*dq.User, error) {
	u, err := s.sql.client.User.Query().
		With(withUserInclude(include)).
		Where(ifPresent(user.IDEQ, filter.ID), ifPresent(user.EmailEQ, filter.Email)).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "User")
	}
	if err != nil {
		return nil, err
	}

	return sqlUserToDQUser(u), nil
}

func (s *UserService) CreateUser(ctx context.Context, u *dq.User) (*dq.User, error) {
	hash, err := s.cryptoService.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	dbU, err := s.sql.client.User.
		Create().
		SetEmail(u.Email).
		SetPassword(hash).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	dqU := sqlUserToDQUser(dbU)

	err = s.eventService.Publish(dq.EventTopicUserCreated, dq.UserCreatedPayload{User: dqU})
	if err != nil {
		return nil, err
	}

	return dqU, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, upd dq.UserUpdate) (*dq.User, error) {
	q := s.sql.client.User.Update().
		Where(user.IDEQ(id))
	if upd.Email != nil {
		q.SetEmail(*upd.Email)
	}
	if upd.Password != nil {
		q.SetPassword(*upd.Password)
	}
	if upd.EmailVerifiedAt != nil {
		q.SetEmailVerifiedAt(*upd.EmailVerifiedAt)
	}

	_, err := q.Save(ctx)
	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "User")
	}
	if err != nil {
		return nil, err
	}

	// Refetch user
	return s.FindUserByID(ctx, id, dq.UserInclude{})
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.sql.client.User.DeleteOneID(id).Exec(ctx)
}

func sqlUserToDQUser(u *ent.User) *dq.User {
	verifiedAt := u.EmailVerifiedAt
	n := &dq.User{
		ID:              u.ID,
		Email:           u.Email,
		Password:        u.Password,
		EmailVerifiedAt: verifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}

	return n
}

func sqlUserSliceToDQUserSlice(us ent.Users) (out []*dq.User) {
	for _, u := range us {
		out = append(out, sqlUserToDQUser(u))
	}

	return out
}

func withUserInclude(include dq.UserInclude) func(q *ent.UserQuery) {
	return func(q *ent.UserQuery) {
		if include.Collections {
			q.WithCollections()
		}
	}
}

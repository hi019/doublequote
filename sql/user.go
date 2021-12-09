package sql

import (
	"context"

	"doublequote"
	"doublequote/prisma"
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
	q := s.sql.prisma.
		User.
		FindFirst(prisma.User.ID.Equals(id)).
		With(buildUserInclude(include)...)

	dbU, err := q.Exec(ctx)
	if err == prisma.ErrNotFound {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "User")
	}
	if err != nil {
		return nil, err
	}

	u := sqlUserToDQUser(dbU)
	return u, nil
}

func (s *UserService) FindUsers(ctx context.Context, filter dq.UserFilter, include dq.UserInclude) ([]*dq.User, int, error) {
	u, err := s.sql.prisma.User.FindMany(
		prisma.User.ID.EqualsIfPresent(filter.ID),
		prisma.User.Email.EqualsIfPresent(filter.Email),
	).
		With(buildUserInclude(include)...).
		Skip(filter.Offset).
		Take(filter.Limit).
		Exec(ctx)

	// TODO implement Count when available https://github.com/prisma/prisma-client-go/issues/229

	return sqlUserSliceToDQUserSlice(u), len(u), err
}

func (s *UserService) FindUser(ctx context.Context, filter dq.UserFilter, include dq.UserInclude) (*dq.User, error) {
	u, err := s.sql.prisma.User.FindFirst(
		prisma.User.ID.EqualsIfPresent(filter.ID),
		prisma.User.Email.EqualsIfPresent(filter.Email),
	).
		With(buildUserInclude(include)...).
		Exec(ctx)

	if err == prisma.ErrNotFound {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Collection")
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

	dbU, err := s.sql.prisma.User.CreateOne(
		prisma.User.Email.Set(u.Email),
		prisma.User.Password.Set(hash),
	).Exec(ctx)
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
	dbU, err := s.sql.prisma.User.FindUnique(prisma.User.ID.Equals(id)).
		Update(
			prisma.User.Email.SetIfPresent(upd.Email),
			prisma.User.Password.SetIfPresent(upd.Password),
			prisma.User.EmailVerifiedAt.SetIfPresent(upd.EmailVerifiedAt),
		).
		Exec(ctx)
	if err == prisma.ErrNotFound {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "User")
	}
	if err != nil {
		return nil, err
	}

	return sqlUserToDQUser(dbU), err
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	_, err := s.sql.prisma.User.FindUnique(prisma.User.ID.Equals(id)).
		Delete().
		Exec(ctx)
	return err
}

func sqlUserToDQUser(u *prisma.UserModel) *dq.User {
	verifiedAt, _ := u.EmailVerifiedAt()
	n := &dq.User{
		ID:              u.ID,
		Email:           u.Email,
		Password:        u.Password,
		EmailVerifiedAt: verifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}

	if u.RelationsUser.Collections != nil {
		n.Collections = sqlColSliceToDQColSlice(u.Collections())
	}

	return n
}

func sqlUserSliceToDQUserSlice(us []prisma.UserModel) (out []*dq.User) {
	for _, u := range us {
		out = append(out, sqlUserToDQUser(&u))
	}

	return out
}

func buildUserInclude(include dq.UserInclude) (filters []prisma.IUserRelationWith) {
	if include.Collections {
		filters = append(filters, prisma.User.Collections.Fetch())
	}

	return
}

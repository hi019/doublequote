package sql

import (
	"context"
	"testing"
	"time"

	"doublequote"
	dqmock "doublequote/mock"
	"doublequote/prisma"
	"doublequote/utils"
	"github.com/stretchr/testify/assert"
	smock "github.com/stretchr/testify/mock"
)

type TestUserService struct {
	svc UserService

	eventService  *dqmock.EventService
	cryptoService *dqmock.CryptoService

	db struct {
		db     *SQL
		client *prisma.PrismaClient
		mock   *prisma.Mock
		ensure func(t *testing.T)
	}
}

func NewTestUserService() *TestUserService {
	ts := &TestUserService{}

	db, client, mock, ensure := NewTestDB()
	ts.db.db = db
	ts.db.client = client
	ts.db.mock = mock
	ts.db.ensure = ensure

	ts.eventService = &dqmock.EventService{}
	ts.cryptoService = &dqmock.CryptoService{}

	ts.svc.eventService = ts.eventService
	ts.svc.cryptoService = ts.cryptoService
	ts.svc.sql = ts.db.db

	return ts
}

func TestUserService_FindUserByID(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		result := prisma.UserModel{
			InnerUser: prisma.InnerUser{
				ID:        1,
				Email:     "test@example.com",
				Password:  "hashedpassword",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.User.Expect(
			s.db.client.User.FindFirst(
				prisma.User.ID.Equals(1),
			),
		).Returns(result)

		found, err := s.svc.FindUserByID(context.Background(), 1, dq.UserInclude{})

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, "test@example.com", found.Email)
		assert.Equal(t, "hashedpassword", found.Password)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found.CreatedAt)
		assert.Equal(t, time.Time{}, found.UpdatedAt)
	})

	t.Run("Include", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		result := prisma.UserModel{
			InnerUser: prisma.InnerUser{
				ID: 1,
			},
			RelationsUser: prisma.RelationsUser{Collections: []prisma.CollectionModel{
				{
					InnerCollection: prisma.InnerCollection{
						ID: 1,
					},
				},
			}},
		}
		s.db.mock.User.Expect(
			s.db.client.User.FindFirst(
				prisma.User.ID.Equals(1),
			).With(prisma.User.Collections.Fetch()),
		).Returns(result)

		found, err := s.svc.FindUserByID(context.Background(), 1, dq.UserInclude{Collections: true})

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, 1, found.Collections[0].ID)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		s.db.mock.User.Expect(
			s.db.client.User.FindFirst(
				prisma.User.ID.Equals(1),
			),
		).Errors(prisma.ErrNotFound)

		found, err := s.svc.FindUserByID(context.Background(), 1, dq.UserInclude{})

		assert.Nil(t, found)
		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "User not found.", dq.ErrorMessage(err))
	})
}

func TestUserService_FindUsers(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		result := []prisma.UserModel{
			{
				InnerUser: prisma.InnerUser{
					ID:        1,
					Email:     "test@example.com",
					CreatedAt: utils.MustParseTime(t, "2021-06-06"),
					UpdatedAt: time.Time{},
				},
			},
		}
		s.db.mock.User.Expect(
			s.db.client.User.FindMany(
				prisma.User.ID.EqualsIfPresent(nil),
				prisma.User.Email.EqualsIfPresent(utils.StringPtr("test@example.com")),
			).Skip(0).Take(1),
		).ReturnsMany(result)

		found, count, err := s.svc.FindUsers(context.Background(), dq.UserFilter{Email: utils.StringPtr("test@example.com"), Offset: 0, Limit: 1}, dq.UserInclude{})

		assert.Nil(t, err)
		assert.Equal(t, 1, count)
		assert.Equal(t, 1, found[0].ID)
		assert.Equal(t, "test@example.com", found[0].Email)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found[0].CreatedAt)
		assert.Equal(t, time.Time{}, found[0].UpdatedAt)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		db, client, mock, ensure := NewTestDB()
		defer ensure(t)

		es := &dqmock.EventService{}
		cr := &dqmock.CryptoService{}
		s := NewUserService(db, es, cr)

		mock.User.Expect(
			client.User.FindFirst(
				prisma.User.ID.Equals(1),
			),
		).Errors(prisma.ErrNotFound)

		found, err := s.FindUserByID(context.Background(), 1, dq.UserInclude{})

		assert.Nil(t, found)
		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "User not found.", dq.ErrorMessage(err))
	})
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		s.eventService.On("Publish", smock.Anything, smock.Anything).Return(nil)
		s.cryptoService.On("HashPassword", "password").Return("hashed-password", nil)

		toCreate := &dq.User{
			Email:    "test@example.com",
			Password: "password",
		}

		result := prisma.UserModel{
			InnerUser: prisma.InnerUser{
				ID:        1,
				Email:     "test@example.com",
				Password:  "hashed-password",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.User.Expect(
			s.db.client.User.CreateOne(
				prisma.User.Email.Set("test@example.com"),
				prisma.User.Password.Set("hashed-password"),
			),
		).Returns(result)

		created, err := s.svc.CreateUser(context.Background(), toCreate)

		assert.Equal(t, err, nil)
		assert.Equal(t, 1, created.ID)
		assert.Equal(t, "test@example.com", created.Email)
		assert.Equal(t, "hashed-password", created.Password)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), created.CreatedAt)
		assert.Equal(t, time.Time{}, created.UpdatedAt)

		s.cryptoService.AssertExpectations(t)
	})

	t.Run("PublishesEvent", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		s.eventService.On("Publish", dq.EventTopicUserCreated, smock.Anything).Return()

		s.cryptoService.On("HashPassword", "password").Return("hashed-password", nil)

		toCreate := &dq.User{
			Email:    "test@example.com",
			Password: "password",
		}

		result := prisma.UserModel{
			InnerUser: prisma.InnerUser{
				ID:        1,
				Email:     "test@example.com",
				Password:  "hashed-password",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.User.Expect(
			s.db.client.User.CreateOne(
				prisma.User.Email.Set("test@example.com"),
				prisma.User.Password.Set("hashed-password"),
			),
		).Returns(result)

		created, err := s.svc.CreateUser(context.Background(), toCreate)
		assert.Equal(t, err, nil)

		s.eventService.AssertCalled(t, "Publish", dq.EventTopicUserCreated, dq.UserCreatedPayload{User: created})
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		result := prisma.UserModel{
			InnerUser: prisma.InnerUser{
				ID:        1,
				Email:     "test@example.com",
				Password:  "hashed-password",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.User.Expect(
			s.db.client.User.FindUnique(prisma.User.ID.Equals(1)).
				Update(
					prisma.User.Email.SetIfPresent(utils.StringPtr("test@example.com")),
					prisma.User.Password.SetIfPresent(utils.StringPtr("hashed-password")),
					prisma.User.EmailVerifiedAt.SetIfPresent(nil),
				),
		).Returns(result)

		updated, err := s.svc.UpdateUser(context.Background(), 1, dq.UserUpdate{Email: utils.StringPtr("test@example.com"), Password: utils.StringPtr("hashed-password")})

		assert.Equal(t, err, nil)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "test@example.com", updated.Email)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), updated.CreatedAt)
		assert.Equal(t, time.Time{}, updated.UpdatedAt)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		s.db.mock.User.Expect(
			s.db.client.User.FindUnique(prisma.User.ID.Equals(1)).
				Update(
					prisma.User.Email.SetIfPresent(utils.StringPtr("test@example.com")),
					prisma.User.Password.SetIfPresent(utils.StringPtr("hashed-password")),
					prisma.User.EmailVerifiedAt.SetIfPresent(nil),
				),
		).Errors(prisma.ErrNotFound)

		updated, err := s.svc.UpdateUser(context.Background(), 1, dq.UserUpdate{Email: utils.StringPtr("test@example.com"), Password: utils.StringPtr("hashed-password")})

		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "User not found.", dq.ErrorMessage(err))
		assert.Nil(t, updated)
	})

}

func TestUserService_DeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestUserService()
		defer s.db.ensure(t)

		result := prisma.UserModel{
			InnerUser: prisma.InnerUser{
				ID:        1,
				Email:     "test@example.com",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.User.Expect(
			s.db.client.User.FindUnique(prisma.User.ID.Equals(1)).
				Delete(),
		).Returns(result)

		err := s.svc.DeleteUser(context.Background(), 1)

		assert.Equal(t, err, nil)
	})
}

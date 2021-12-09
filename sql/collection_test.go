package sql

import (
	"context"
	"testing"
	"time"

	dq "doublequote"
	"doublequote/prisma"
	"doublequote/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCollectionService struct {
	svc CollectionService

	db struct {
		db     *SQL
		client *prisma.PrismaClient
		mock   *prisma.Mock
		ensure func(t *testing.T)
	}
}

func NewTestCollectionService() *TestCollectionService {
	ts := &TestCollectionService{}

	db, client, mock, ensure := NewTestDB()
	ts.db.db = db
	ts.db.client = client
	ts.db.mock = mock
	ts.db.ensure = ensure

	ts.svc.sql = ts.db.db

	return ts
}

func TestCollectionService_FindCollectionByID(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		result := prisma.CollectionModel{
			InnerCollection: prisma.InnerCollection{
				ID:        1,
				Name:      "Tech News",
				UserID:    1,
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
			RelationsCollection: prisma.RelationsCollection{
				User: &prisma.UserModel{
					InnerUser: prisma.InnerUser{
						ID:    1,
						Email: "test@example.com",
					},
				},
			},
		}
		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindFirst(
				prisma.Collection.ID.Equals(1),
			).With(
				prisma.Collection.User.Fetch(),
			),
		).Returns(result)

		found, err := s.svc.FindCollectionByID(context.Background(), 1, dq.CollectionInclude{})

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, "Tech News", found.Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found.CreatedAt)
		assert.Equal(t, time.Time{}, found.UpdatedAt)

		assert.Equal(t, "test@example.com", found.User.Email)
		assert.Equal(t, 1, found.User.ID)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindFirst(
				prisma.Collection.ID.Equals(1),
			).With(
				prisma.Collection.User.Fetch(),
			),
		).Errors(prisma.ErrNotFound)

		found, err := s.svc.FindCollectionByID(context.Background(), 1, dq.CollectionInclude{})

		assert.Nil(t, found)
		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "Collection not found.", dq.ErrorMessage(err))
	})

	t.Run("WithFeeds", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		result := prisma.CollectionModel{
			InnerCollection: prisma.InnerCollection{
				ID:        1,
				Name:      "Tech News",
				UserID:    1,
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
			RelationsCollection: prisma.RelationsCollection{
				Feeds: []prisma.FeedModel{
					{
						InnerFeed: prisma.InnerFeed{
							ID:   1,
							Name: "Test",
						},
					},
				},
			},
		}
		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindFirst(
				prisma.Collection.ID.Equals(1),
			).With(
				prisma.Collection.User.Fetch(),
				prisma.Collection.Feeds.Fetch(),
			),
		).Returns(result)

		found, err := s.svc.FindCollectionByID(context.Background(), 1, dq.CollectionInclude{Feeds: true})

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, "Tech News", found.Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found.CreatedAt)
		assert.Equal(t, time.Time{}, found.UpdatedAt)

		assert.Equal(t, 1, result.Feeds()[0].ID)
		assert.Equal(t, "Test", result.Feeds()[0].Name)
	})
}

func TestCollectionService_FindCollections(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		result := []prisma.CollectionModel{
			{
				InnerCollection: prisma.InnerCollection{
					ID:        1,
					Name:      "News",
					CreatedAt: utils.MustParseTime(t, "2021-06-06"),
					UpdatedAt: time.Time{},
				},
				RelationsCollection: prisma.RelationsCollection{
					User: &prisma.UserModel{
						InnerUser: prisma.InnerUser{
							ID:    1,
							Email: "test@example.com",
						},
					},
				},
			},
			{
				InnerCollection: prisma.InnerCollection{
					ID:        2,
					Name:      "News",
					CreatedAt: utils.MustParseTime(t, "2021-08-06"),
					UpdatedAt: time.Time{},
				},
				RelationsCollection: prisma.RelationsCollection{
					User: &prisma.UserModel{
						InnerUser: prisma.InnerUser{
							ID:    2,
							Email: "test2@example.com",
						},
					},
				},
			},
		}
		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindMany(
				prisma.Collection.ID.EqualsIfPresent(nil),
				prisma.Collection.Name.EqualsIfPresent(utils.StringPtr("News")),
				prisma.Collection.UserID.EqualsIfPresent(nil),
			).With(
				prisma.Collection.User.Fetch(),
			).
				Skip(0).
				Take(2),
		).ReturnsMany(result)

		found, count, err := s.svc.FindCollections(
			context.Background(),
			dq.CollectionFilter{Name: utils.StringPtr("News"), Offset: 0, Limit: 2},
			dq.CollectionInclude{},
		)

		assert.Nil(t, err)
		require.Equal(t, 2, count)

		// Collection 0
		assert.Equal(t, 1, found[0].ID)
		assert.Equal(t, "News", found[0].Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found[0].CreatedAt)
		assert.Equal(t, time.Time{}, found[0].UpdatedAt)
		assert.Equal(t, "test@example.com", found[0].User.Email)
		assert.Equal(t, 1, found[0].User.ID)

		// Collection 1
		assert.Equal(t, 2, found[1].ID)
		assert.Equal(t, "News", found[1].Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-08-06"), found[1].CreatedAt)
		assert.Equal(t, time.Time{}, found[1].UpdatedAt)
		assert.Equal(t, "test2@example.com", found[1].User.Email)
		assert.Equal(t, 2, found[1].User.ID)
	})

	t.Run("WithFeeds", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		result := []prisma.CollectionModel{
			{
				InnerCollection: prisma.InnerCollection{
					ID:   1,
					Name: "News",
				},
				RelationsCollection: prisma.RelationsCollection{
					Feeds: []prisma.FeedModel{{
						InnerFeed: prisma.InnerFeed{ID: 1, Name: "The Verge"},
					}},
				},
			},
		}
		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindMany(
				prisma.Collection.ID.EqualsIfPresent(nil),
				prisma.Collection.Name.EqualsIfPresent(utils.StringPtr("News")),
				prisma.Collection.UserID.EqualsIfPresent(nil),
			).With(
				prisma.Collection.User.Fetch(),
				prisma.Collection.Feeds.Fetch(),
			).
				Skip(0).
				Take(1),
		).ReturnsMany(result)

		found, count, err := s.svc.FindCollections(
			context.Background(),
			dq.CollectionFilter{Name: utils.StringPtr("News"), Offset: 0, Limit: 1},
			dq.CollectionInclude{Feeds: true},
		)

		assert.Nil(t, err)
		require.Equal(t, 1, count)

		assert.Equal(t, 1, found[0].ID)
		assert.Equal(t, "News", found[0].Name)

		assert.Equal(t, 1, found[0].Feeds[0].ID)
	})
}

func TestCollectionService_CreateCollection(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		toCreate := &dq.Collection{
			Name:   "Tech News",
			UserID: 1,
		}

		result := prisma.CollectionModel{
			InnerCollection: prisma.InnerCollection{
				ID:        1,
				Name:      "Tech News",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
			RelationsCollection: prisma.RelationsCollection{
				User: &prisma.UserModel{
					InnerUser: prisma.InnerUser{
						ID:    2,
						Email: "test@example.com",
					},
				},
			},
		}
		s.db.mock.Collection.Expect(
			s.db.client.Collection.CreateOne(
				prisma.Collection.Name.Set("Tech News"),
				prisma.Collection.User.Link(
					prisma.User.ID.Equals(1),
				),
			).With(
				prisma.Collection.User.Fetch(),
			),
		).Returns(result)

		created, err := s.svc.CreateCollection(context.Background(), toCreate)

		assert.Equal(t, err, nil)
		assert.Equal(t, 1, created.ID)
		assert.Equal(t, "Tech News", created.Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), created.CreatedAt)
		assert.Equal(t, time.Time{}, created.UpdatedAt)

		assert.Equal(t, "test@example.com", created.User.Email)
		assert.Equal(t, 2, created.User.ID)
	})
}

func TestCollectionService_UpdateCollection(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		result := prisma.CollectionModel{
			InnerCollection: prisma.InnerCollection{
				ID:        1,
				Name:      "Tech News",
				UserID:    2,
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
			RelationsCollection: prisma.RelationsCollection{
				User: &prisma.UserModel{
					InnerUser: prisma.InnerUser{
						ID:    2,
						Email: "test@example.com",
					},
				},
			},
		}
		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindUnique(prisma.Collection.ID.Equals(1)).
				Update(
					prisma.Collection.Name.SetIfPresent(utils.StringPtr("Tech News")),
					prisma.Collection.UserID.SetIfPresent(utils.IntPtr(2)),
					prisma.Collection.Feeds.Link(prisma.Feed.ID.InIfPresent(nil)),
				),
		).Returns(result)

		updated, err := s.svc.UpdateCollection(
			context.Background(),
			1,
			dq.CollectionUpdate{
				UserID: utils.IntPtr(2),
				Name:   utils.StringPtr("Tech News"),
			},
		)

		assert.Equal(t, err, nil)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "Tech News", updated.Name)
		assert.Equal(t, 2, updated.UserID)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), updated.CreatedAt)
		assert.Equal(t, time.Time{}, updated.UpdatedAt)

		assert.Equal(t, "test@example.com", updated.User.Email)
		assert.Equal(t, 2, updated.User.ID)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindUnique(prisma.Collection.ID.Equals(1)).
				Update(
					prisma.Collection.Name.SetIfPresent(utils.StringPtr("Tech News 2")),
					prisma.Collection.UserID.SetIfPresent(nil),
					prisma.Collection.Feeds.Link(prisma.Feed.ID.InIfPresent(nil)),
				),
		).Errors(prisma.ErrNotFound)

		updated, err := s.svc.
			UpdateCollection(
				context.Background(),
				1,
				dq.CollectionUpdate{Name: utils.StringPtr("Tech News 2")},
			)

		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "Collection not found.", dq.ErrorMessage(err))
		assert.Nil(t, updated)
	})
}

func TestCollectionService_DeleteCollection(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestCollectionService()
		defer s.db.ensure(t)

		result := prisma.CollectionModel{
			InnerCollection: prisma.InnerCollection{
				ID:        1,
				Name:      "News",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.Collection.Expect(
			s.db.client.Collection.FindUnique(prisma.Collection.ID.Equals(1)).
				Delete(),
		).Returns(result)

		err := s.svc.DeleteCollection(context.Background(), 1)

		assert.Equal(t, err, nil)
	})
}

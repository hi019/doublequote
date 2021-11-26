package sql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	dq "doublequote"
	"doublequote/prisma"
	"doublequote/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestFeedService struct {
	svc FeedService

	db struct {
		sql    *SQL
		client *prisma.PrismaClient
		mock   *prisma.Mock
		ensure func(t *testing.T)
	}
}

func NewTestFeedService() *TestFeedService {
	ts := &TestFeedService{}

	db, client, mock, ensure := NewTestDB()
	ts.db.sql = db
	ts.db.client = client
	ts.db.mock = mock
	ts.db.ensure = ensure

	ts.svc.sql = ts.db.sql

	return ts
}

func TestFeedService_FindFeedByID(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:        1,
				Name:      "The Verge",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
			},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindFirst(
				prisma.Feed.ID.Equals(1),
			),
		).Returns(result)

		found, err := s.svc.FindFeedByID(context.Background(), 1, dq.FeedInclude{})

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, "The Verge", found.Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found.CreatedAt)
		assert.Equal(t, time.Time{}, found.UpdatedAt)
	})

	t.Run("WithCollection", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:        1,
				Name:      "The Verge",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
			},
			RelationsFeed: prisma.RelationsFeed{Collections: []prisma.CollectionModel{
				{
					InnerCollection: prisma.InnerCollection{ID: 1},
				},
			}},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindFirst(
				prisma.Feed.ID.Equals(1),
			).With(
				prisma.Feed.Collections.Fetch(),
			),
		).Returns(result)

		found, err := s.svc.FindFeedByID(context.Background(), 1, dq.FeedInclude{Collections: true})

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, "The Verge", found.Name)

		assert.Equal(t, 1, found.Collections[0].ID)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindFirst(
				prisma.Feed.ID.Equals(1),
			),
		).Errors(sql.ErrNoRows)

		found, err := s.svc.FindFeedByID(context.Background(), 1, dq.FeedInclude{})

		assert.Nil(t, found)
		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "Feed not found.", dq.ErrorMessage(err))
	})
}

func TestFeedService_FindFeed(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:        1,
				Name:      "The Verge",
				RssURL:    "https://theverge.com/rss",
				Domain:    "https://theverge.com",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
			},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindFirst(
				prisma.Feed.ID.EqualsIfPresent(nil),
				prisma.Feed.Name.EqualsIfPresent(utils.StringPtr("The Verge")),
				prisma.Feed.RssURL.EqualsIfPresent(utils.StringPtr("https://theverge.com/rss")),
				prisma.Feed.Domain.EqualsIfPresent(utils.StringPtr("theverge.com")),
				prisma.Feed.Collections.Every(
					prisma.Collection.ID.EqualsIfPresent(nil),
				),
			),
		).Returns(result)

		found, err := s.svc.FindFeed(context.Background(), dq.FeedFilter{
			Name:   utils.StringPtr("The Verge"),
			RssURL: utils.StringPtr("https://theverge.com/rss"),
			Domain: utils.StringPtr("theverge.com"),
		},
			dq.FeedInclude{},
		)

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, "The Verge", found.Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found.CreatedAt)
		assert.Equal(t, time.Time{}, found.UpdatedAt)
	})

	t.Run("WithCollections", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:        1,
				Name:      "The Verge",
				RssURL:    "https://theverge.com/rss",
				Domain:    "https://theverge.com",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
			},
			RelationsFeed: prisma.RelationsFeed{Collections: []prisma.CollectionModel{
				{
					InnerCollection: prisma.InnerCollection{
						ID:   1,
						Name: "News!",
					},
				},
			}},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindFirst(
				prisma.Feed.ID.EqualsIfPresent(nil),
				prisma.Feed.Name.EqualsIfPresent(utils.StringPtr("The Verge")),
				prisma.Feed.RssURL.EqualsIfPresent(utils.StringPtr("https://theverge.com/rss")),
				prisma.Feed.Domain.EqualsIfPresent(utils.StringPtr("theverge.com")),
				prisma.Feed.Collections.Every(
					prisma.Collection.ID.EqualsIfPresent(nil),
				),
			).With(prisma.Feed.Collections.Fetch()),
		).Returns(result)

		found, err := s.svc.FindFeed(context.Background(), dq.FeedFilter{
			Name:   utils.StringPtr("The Verge"),
			RssURL: utils.StringPtr("https://theverge.com/rss"),
			Domain: utils.StringPtr("theverge.com"),
		},
			dq.FeedInclude{Collections: true},
		)

		assert.Nil(t, err)
		assert.Equal(t, 1, found.ID)
		assert.Equal(t, "The Verge", found.Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found.CreatedAt)
		assert.Equal(t, time.Time{}, found.UpdatedAt)

		assert.Len(t, found.Collections, 1)
		assert.Equal(t, 1, found.Collections[0].ID)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindFirst(
				prisma.Feed.ID.Equals(1),
			),
		).Errors(sql.ErrNoRows)

		found, err := s.svc.FindFeedByID(context.Background(), 1, dq.FeedInclude{})

		assert.Nil(t, found)
		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "Feed not found.", dq.ErrorMessage(err))
	})
}

func TestFeedService_FindFeeds(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		result := []prisma.FeedModel{
			{
				InnerFeed: prisma.InnerFeed{
					ID:        1,
					Name:      "News",
					CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				},
			},
			{
				InnerFeed: prisma.InnerFeed{
					ID:        2,
					Name:      "News",
					CreatedAt: utils.MustParseTime(t, "2021-08-06"),
				},
			},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindMany(
				prisma.Feed.ID.EqualsIfPresent(nil),
				prisma.Feed.Name.EqualsIfPresent(utils.StringPtr("News")),
				prisma.Feed.Domain.EqualsIfPresent(nil),
				prisma.Feed.RssURL.EqualsIfPresent(nil),
				prisma.Feed.Collections.Every(
					prisma.Collection.ID.EqualsIfPresent(nil),
				),
			).
				Skip(0).
				Take(2),
		).ReturnsMany(result)

		found, count, err := s.svc.FindFeeds(
			context.Background(),
			dq.FeedFilter{Name: utils.StringPtr("News"), Offset: 0, Limit: 2},
			dq.FeedInclude{},
		)

		assert.Nil(t, err)
		require.Equal(t, 2, count)

		// Feed 0
		assert.Equal(t, 1, found[0].ID)
		assert.Equal(t, "News", found[0].Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), found[0].CreatedAt)
		assert.Equal(t, time.Time{}, found[0].UpdatedAt)

		// Feed 1
		assert.Equal(t, 2, found[1].ID)
		assert.Equal(t, "News", found[1].Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-08-06"), found[1].CreatedAt)
		assert.Equal(t, time.Time{}, found[1].UpdatedAt)
	})

	t.Run("WithCollection", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		collection := prisma.InnerCollection{
			ID:   1,
			Name: "Just News",
		}
		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:   1,
				Name: "News",
			},
			RelationsFeed: prisma.RelationsFeed{
				Collections: []prisma.CollectionModel{{InnerCollection: collection}},
			},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindMany(
				prisma.Feed.ID.EqualsIfPresent(nil),
				prisma.Feed.Name.EqualsIfPresent(nil),
				prisma.Feed.Domain.EqualsIfPresent(nil),
				prisma.Feed.RssURL.EqualsIfPresent(nil),
				prisma.Feed.Collections.Every(
					prisma.Collection.ID.EqualsIfPresent(utils.IntPtr(1)),
				),
			).
				With(
					prisma.Feed.Collections.Fetch(),
				).
				Skip(0).
				Take(2),
		).ReturnsMany([]prisma.FeedModel{result})

		found, count, err := s.svc.FindFeeds(
			context.Background(),
			dq.FeedFilter{CollectionID: utils.IntPtr(1), Offset: 0, Limit: 2},
			dq.FeedInclude{Collections: true},
		)

		assert.Nil(t, err)
		require.Equal(t, 1, count)
		require.Len(t, found[0].Collections, 1)

		assert.Equal(t, 1, found[0].Collections[0].ID)
		assert.Equal(t, "Just News", found[0].Collections[0].Name)
	})
}

func TestFeedService_CreateFeed(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		toCreate := &dq.Feed{
			Name:   "The Verge",
			Domain: "theverge.com",
			RssURL: "https://theverge.com/rss",
		}

		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:        1,
				Name:      "The Verge",
				Domain:    "theverge.com",
				RssURL:    "https://theverge.com/rss",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.CreateOne(
				prisma.Feed.Name.Set("The Verge"),
				prisma.Feed.RssURL.Set("https://theverge.com/rss"),
				prisma.Feed.Domain.Set("theverge.com"),
			),
		).Returns(result)

		created, err := s.svc.CreateFeed(context.Background(), toCreate)

		assert.Equal(t, err, nil)
		assert.Equal(t, 1, created.ID)
		assert.Equal(t, "The Verge", created.Name)
		assert.Equal(t, "https://theverge.com/rss", created.RssURL)
		assert.Equal(t, "theverge.com", created.Domain)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), created.CreatedAt)
		assert.Equal(t, time.Time{}, created.UpdatedAt)
	})
}

func TestFeedService_UpdateFeed(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:        1,
				Name:      "Tech News",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: utils.MustParseTime(t, "2021-06-07"),
			},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindUnique(prisma.Feed.ID.Equals(1)).
				Update(
					prisma.Feed.Name.SetIfPresent(utils.StringPtr("Tech News")),
					prisma.Feed.RssURL.SetIfPresent(nil),
					prisma.Feed.Domain.SetIfPresent(nil),
				),
		).Returns(result)

		updated, err := s.svc.UpdateFeed(
			context.Background(),
			1,
			dq.FeedUpdate{
				Name: utils.StringPtr("Tech News"),
			},
		)

		assert.Equal(t, err, nil)
		assert.Equal(t, 1, updated.ID)
		assert.Equal(t, "Tech News", updated.Name)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-06"), updated.CreatedAt)
		assert.Equal(t, utils.MustParseTime(t, "2021-06-07"), updated.UpdatedAt)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindUnique(prisma.Feed.ID.Equals(1)).
				Update(
					prisma.Feed.Name.SetIfPresent(utils.StringPtr("Tech News 2")),
					prisma.Feed.RssURL.SetIfPresent(nil),
					prisma.Feed.Domain.SetIfPresent(nil),
				),
		).Errors(sql.ErrNoRows)

		updated, err := s.svc.
			UpdateFeed(
				context.Background(),
				1,
				dq.FeedUpdate{Name: utils.StringPtr("Tech News 2")},
			)

		assert.Equal(t, dq.ENOTFOUND, dq.ErrorCode(err))
		assert.Equal(t, "Feed not found.", dq.ErrorMessage(err))
		assert.Nil(t, updated)
	})
}

func TestFeedService_DeleteFeed(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		s := NewTestFeedService()
		defer s.db.ensure(t)

		result := prisma.FeedModel{
			InnerFeed: prisma.InnerFeed{
				ID:        1,
				Name:      "News",
				CreatedAt: utils.MustParseTime(t, "2021-06-06"),
				UpdatedAt: time.Time{},
			},
		}
		s.db.mock.Feed.Expect(
			s.db.client.Feed.FindUnique(prisma.Feed.ID.Equals(1)).
				Delete(),
		).Returns(result)

		err := s.svc.DeleteFeed(context.Background(), 1)

		assert.Equal(t, err, nil)
	})
}

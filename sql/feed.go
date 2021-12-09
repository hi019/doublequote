package sql

import (
	"context"

	dq "doublequote"
	"doublequote/prisma"
)

// Ensure service implements interface.
var _ dq.FeedService = (*FeedService)(nil)

type FeedService struct {
	sql *SQL
}

func NewFeedService(sql *SQL) *FeedService {
	return &FeedService{sql: sql}
}

func (s *FeedService) FindFeedByID(ctx context.Context, id int, include dq.FeedInclude) (*dq.Feed, error) {
	col, err := s.sql.prisma.Feed.
		FindFirst(prisma.Feed.ID.Equals(id)).
		With(buildFeedInclude(include)...).
		Exec(ctx)
	if err == prisma.ErrNotFound {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Feed")
	}
	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(col), err
}

func (s *FeedService) FindFeeds(ctx context.Context, filter dq.FeedFilter, include dq.FeedInclude) ([]*dq.Feed, int, error) {
	cols, err := s.sql.prisma.Feed.FindMany(
		prisma.Feed.ID.EqualsIfPresent(filter.ID),
		prisma.Feed.Name.EqualsIfPresent(filter.Name),
		prisma.Feed.Domain.EqualsIfPresent(filter.Domain),
		prisma.Feed.RssURL.EqualsIfPresent(filter.RssURL),
		prisma.Feed.Collections.Every(
			prisma.Collection.ID.EqualsIfPresent(filter.CollectionID),
		),
	).
		With(buildFeedInclude(include)...).
		Skip(filter.Offset).
		Take(filter.Limit).
		Exec(ctx)

	// TODO implement Count when available https://github.com/prisma/prisma-client-go/issues/229

	return sqlFeedSliceToFeedSlice(cols), len(cols), err
}

func (s *FeedService) FindFeed(ctx context.Context, filter dq.FeedFilter, include dq.FeedInclude) (*dq.Feed, error) {
	c, err := s.sql.prisma.Feed.FindFirst(
		prisma.Feed.ID.EqualsIfPresent(filter.ID),
		prisma.Feed.Name.EqualsIfPresent(filter.Name),
		prisma.Feed.RssURL.EqualsIfPresent(filter.RssURL),
		prisma.Feed.Domain.EqualsIfPresent(filter.Domain),
		prisma.Feed.Collections.Every(
			prisma.Collection.ID.EqualsIfPresent(filter.CollectionID),
		),
	).
		With(buildFeedInclude(include)...).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(c), err
}

func (s *FeedService) CreateFeed(ctx context.Context, feed *dq.Feed) (*dq.Feed, error) {
	c, err := s.sql.prisma.Feed.CreateOne(
		prisma.Feed.Name.Set(feed.Name),
		prisma.Feed.RssURL.Set(feed.RssURL),
		prisma.Feed.Domain.Set(feed.Domain),
	).
		Exec(ctx)

	return sqlFeedToFeed(c), err
}

func (s *FeedService) UpdateFeed(ctx context.Context, id int, upd dq.FeedUpdate) (*dq.Feed, error) {
	dbU, err := s.sql.prisma.Feed.FindUnique(prisma.Feed.ID.Equals(id)).
		Update(
			prisma.Feed.Name.SetIfPresent(upd.Name),
			prisma.Feed.RssURL.SetIfPresent(upd.RssURL),
			prisma.Feed.Domain.SetIfPresent(upd.Domain),
		).
		Exec(ctx)
	if err == prisma.ErrNotFound {
		err = dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Feed")
	}
	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(dbU), err
}

func (s *FeedService) DeleteFeed(ctx context.Context, id int) error {
	_, err := s.sql.prisma.Feed.FindUnique(prisma.Feed.ID.Equals(id)).
		Delete().
		Exec(ctx)
	return err
}

func sqlFeedToFeed(c *prisma.FeedModel) *dq.Feed {
	f := &dq.Feed{
		ID:        c.ID,
		Name:      c.Name,
		Domain:    c.Domain,
		RssURL:    c.RssURL,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.RelationsFeed.Collections != nil {
		for _, c := range c.Collections() {
			f.Collections = append(f.Collections, sqlColToDQCol(&c))
		}
	}

	return f
}

func sqlFeedSliceToFeedSlice(cs []prisma.FeedModel) (out []*dq.Feed) {
	for _, u := range cs {
		out = append(out, sqlFeedToFeed(&u))
	}

	return
}

func buildFeedInclude(include dq.FeedInclude) (filters []prisma.IFeedRelationWith) {
	if include.Collections {
		filters = append(filters, prisma.Feed.Collections.Fetch())
	}

	return
}

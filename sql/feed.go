package sql

import (
	"context"

	dq "doublequote"
	"doublequote/ent"
	"doublequote/ent/collection"
	"doublequote/ent/feed"
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
	feed, err := s.sql.client.Feed.Query().
		With(withFeedInclude(include)).
		Where(feed.IDEQ(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Feed")
	}
	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(feed), err
}

func (s *FeedService) FindFeeds(ctx context.Context, filter dq.FeedFilter, include dq.FeedInclude) ([]*dq.Feed, int, error) {
	feeds, err := s.sql.client.Feed.Query().
		Where(
			ifPresent(feed.IDEQ, filter.ID),
			ifPresent(feed.NameEQ, filter.Name),
			ifPresent(feed.DomainEQ, filter.Domain),
			ifPresent(feed.RssURLEQ, filter.RssURL),
			feed.HasCollectionsWith(ifPresent(collection.IDEQ, filter.CollectionID)),
		).
		With(withFeedInclude(include)).
		Offset(filter.Offset).
		Limit(filter.Limit).
		All(ctx)

	return sqlFeedSliceToFeedSlice(feeds), len(feeds), err
}

func (s *FeedService) FindFeed(ctx context.Context, filter dq.FeedFilter, include dq.FeedInclude) (*dq.Feed, error) {
	f, err := s.sql.client.Feed.Query().
		Where(
			ifPresent(feed.IDEQ, filter.ID),
			ifPresent(feed.NameEQ, filter.Name),
			ifPresent(feed.RssURLEQ, filter.RssURL),
			ifPresent(feed.DomainEQ, filter.Domain),
			feed.HasCollectionsWith(ifPresent(collection.IDEQ, filter.CollectionID)),
		).
		With(withFeedInclude(include)).
		First(ctx)

	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(f), err
}

func (s *FeedService) CreateFeed(ctx context.Context, feed *dq.Feed) (*dq.Feed, error) {
	f, err := s.sql.client.Feed.Create().
		SetName(feed.Name).
		SetDomain(feed.Domain).
		SetRssURL(feed.RssURL).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// TODO access f.User.ID without another query?
	return s.FindFeedByID(ctx, f.ID, dq.FeedInclude{})
}

func (s *FeedService) UpdateFeed(ctx context.Context, id int, upd dq.FeedUpdate) (*dq.Feed, error) {
	q := s.sql.client.Feed.UpdateOneID(id)

	if upd.Name != nil {
		q.SetName(*upd.Name)
	}
	if upd.RssURL != nil {
		q.SetRssURL(*upd.RssURL)
	}
	if upd.Domain != nil {
		q.SetDomain(*upd.Domain)
	}

	dbU, err := q.Save(ctx)

	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Feed")
	}
	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(dbU), err
}

func (s *FeedService) DeleteFeed(ctx context.Context, id int) error {
	return s.sql.client.Feed.DeleteOneID(id).Exec(ctx)
}

func sqlFeedToFeed(c *ent.Feed) *dq.Feed {
	f := &dq.Feed{
		ID:        c.ID,
		Name:      c.Name,
		Domain:    c.Domain,
		RssURL:    c.RssURL,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if cs, err := c.Edges.CollectionsOrErr(); err != nil {
		for _, c := range cs {
			f.Collections = append(f.Collections, sqlColToDQCol(c))
		}
	}

	return f
}

func sqlFeedSliceToFeedSlice(cs []*ent.Feed) (out []*dq.Feed) {
	for _, u := range cs {
		out = append(out, sqlFeedToFeed(u))
	}

	return
}

func withFeedInclude(include dq.FeedInclude) func(q *ent.FeedQuery) {
	return func(q *ent.FeedQuery) {
		if include.Collections {
			q.WithCollections()
		}
	}
}

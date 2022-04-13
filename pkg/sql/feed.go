package sql

import (
	"context"

	"doublequote/ent"
	"doublequote/ent/collection"
	"doublequote/ent/feed"
	"doublequote/pkg/domain"
)

// Ensure service implements interface.
var _ domain.FeedService = (*FeedService)(nil)

type FeedService struct {
	sql *SQL
}

func NewFeedService(sql *SQL) *FeedService {
	return &FeedService{sql: sql}
}

func (s *FeedService) FindFeedByID(ctx context.Context, id int, include domain.FeedInclude) (*domain.Feed, error) {
	f, err := s.sql.client.Feed.Query().
		With(withFeedInclude(include)).
		Where(feed.IDEQ(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Feed")
	}
	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(f), err
}

func (s *FeedService) FindFeeds(ctx context.Context, filter domain.FeedFilter, include domain.FeedInclude) ([]*domain.Feed, int, error) {
	feeds, err := s.sql.client.Feed.Query().
		Where(
			ifPresent(feed.IDEQ, filter.ID),
			ifPresent(feed.NameEQ, filter.Name),
			ifPresent(feed.DomainEQ, filter.Domain),
			ifPresent(feed.RssURLEQ, filter.RssURL),
			// TODO add this to to other FindXs
			maybeHasRelation(filter.CollectionID, feed.HasCollectionsWith, ifPresent(collection.IDEQ, filter.CollectionID)),
		).
		With(withFeedInclude(include)).
		Offset(filter.Offset).
		Limit(filter.Limit).
		All(ctx)

	return sqlFeedSliceToFeedSlice(feeds), len(feeds), err
}

func (s *FeedService) FindFeed(ctx context.Context, filter domain.FeedFilter, include domain.FeedInclude) (*domain.Feed, error) {
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

func (s *FeedService) CreateFeed(ctx context.Context, feed *domain.Feed) (*domain.Feed, error) {
	f, err := s.sql.client.Feed.Create().
		SetName(feed.Name).
		SetDomain(feed.Domain).
		SetRssURL(feed.RssURL).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// TODO access f.User.ID without another query?
	return s.FindFeedByID(ctx, f.ID, domain.FeedInclude{})
}

func (s *FeedService) UpdateFeed(ctx context.Context, id int, upd domain.FeedUpdate) (*domain.Feed, error) {
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
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Feed")
	}
	if err != nil {
		return nil, err
	}

	return sqlFeedToFeed(dbU), err
}

func (s *FeedService) DeleteFeed(ctx context.Context, id int) error {
	return s.sql.client.Feed.DeleteOneID(id).Exec(ctx)
}

func sqlFeedToFeed(c *ent.Feed) *domain.Feed {
	f := &domain.Feed{
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

func sqlFeedSliceToFeedSlice(cs []*ent.Feed) (out []*domain.Feed) {
	for _, u := range cs {
		out = append(out, sqlFeedToFeed(u))
	}

	return
}

func withFeedInclude(include domain.FeedInclude) func(q *ent.FeedQuery) {
	return func(q *ent.FeedQuery) {
		if include.Collections {
			q.WithCollections()
		}
	}
}

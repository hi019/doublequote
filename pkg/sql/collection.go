package sql

import (
	"context"

	"doublequote/ent"
	"doublequote/ent/collection"
	"doublequote/ent/feed"
	"doublequote/ent/user"
	"doublequote/pkg/domain"
)

// Ensure service implements interface.
var _ domain.CollectionService = (*CollectionService)(nil)

type CollectionService struct {
	sql *SQL
}

func NewCollectionService(sql *SQL) *CollectionService {
	return &CollectionService{sql: sql}
}

func (s *CollectionService) FindCollectionByID(ctx context.Context, id int, include domain.CollectionInclude) (*domain.Collection, error) {
	c, err := s.sql.client.Collection.
		Query().
		With(withCollectionInclude(include)).
		Where(collection.IDEQ(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection")
	}
	if err != nil {
		return nil, err
	}

	return sqlColToDQCol(c), err
}

func (s *CollectionService) FindCollections(ctx context.Context, filter domain.CollectionFilter, include domain.CollectionInclude) ([]*domain.Collection, int, error) {
	cols, err := s.sql.client.Collection.Query().
		Where(
			ifPresent(collection.IDEQ, filter.ID),
			ifPresent(collection.NameEQ, filter.Name),
			collection.HasUserWith(ifPresent(user.IDEQ, filter.UserID)),
			collection.HasFeedsWith(ifPresent(feed.IDEQ, filter.FeedID)),
		).
		With(withCollectionInclude(include)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		All(ctx)

	return sqlColSliceToDQColSlice(cols), len(cols), err
}

func (s *CollectionService) FindCollection(ctx context.Context, filter domain.CollectionFilter, include domain.CollectionInclude) (*domain.Collection, error) {
	c, err := s.sql.client.Collection.Query().
		Where(
			ifPresent(collection.IDEQ, filter.ID),
			ifPresent(collection.NameEQ, filter.Name),
			collection.HasUserWith(ifPresent(user.IDEQ, filter.UserID)),
			collection.HasFeedsWith(ifPresent(feed.IDEQ, filter.FeedID)),
		).
		First(ctx)

	if ent.IsNotFound(err) {
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection")
	}
	if err != nil {
		return nil, err
	}

	return sqlColToDQCol(c), err
}

func (s *CollectionService) CreateCollection(ctx context.Context, col *domain.Collection) (*domain.Collection, error) {
	c, err := s.sql.client.Collection.Create().
		SetName(col.Name).
		SetUserID(col.UserID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// TODO access c.User.ID without another query?
	return s.FindCollectionByID(ctx, c.ID, domain.CollectionInclude{})
}

func (s *CollectionService) UpdateCollection(ctx context.Context, id int, upd domain.CollectionUpdate) (c *domain.Collection, err error) {
	q := s.sql.client.Collection.UpdateOneID(id)

	if upd.Name != nil {
		q.SetName(*upd.Name)
	}
	if upd.UserID != nil {
		q.SetUserID(*upd.UserID)
	}
	if upd.FeedsIDs != nil {
		q.ClearFeeds().AddFeedIDs(*upd.FeedsIDs...)
	}

	u, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	// TODO access s.User.ID without another query?
	return s.FindCollectionByID(ctx, u.ID, domain.CollectionInclude{})
}

func (s *CollectionService) DeleteCollection(ctx context.Context, id int) error {
	return s.sql.client.Collection.DeleteOneID(id).Exec(ctx)
}

func sqlColToDQCol(c *ent.Collection) *domain.Collection {
	n := &domain.Collection{
		ID:        c.ID,
		Name:      c.Name,
		UserID:    c.Edges.User.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if u, err := c.Edges.UserOrErr(); err != nil {
		n.User = *sqlUserToDQUser(u)
	}

	return n
}

func sqlColSliceToDQColSlice(cs []*ent.Collection) (out []*domain.Collection) {
	for _, u := range cs {
		out = append(out, sqlColToDQCol(u))
	}

	return
}

func withCollectionInclude(include domain.CollectionInclude) func(q *ent.CollectionQuery) {
	return func(q *ent.CollectionQuery) {
		q.WithUser()

		if include.Feeds {
			q.WithFeeds()
		}
	}
}

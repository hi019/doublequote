package sql

import (
	"context"

	dq "doublequote"
	"doublequote/ent"
	"doublequote/ent/collection"
	"doublequote/ent/user"
)

// Ensure service implements interface.
var _ dq.CollectionService = (*CollectionService)(nil)

type CollectionService struct {
	sql *SQL
}

func NewCollectionService(sql *SQL) *CollectionService {
	return &CollectionService{sql: sql}
}

func (s *CollectionService) FindCollectionByID(ctx context.Context, id int, include dq.CollectionInclude) (*dq.Collection, error) {
	c, err := s.sql.client.Collection.
		Query().
		With(withCollectionInclude(include)).
		Where(collection.IDEQ(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Collection")
	}
	if err != nil {
		return nil, err
	}

	return sqlColToDQCol(c), err
}

func (s *CollectionService) FindCollections(ctx context.Context, filter dq.CollectionFilter, include dq.CollectionInclude) ([]*dq.Collection, int, error) {
	cols, err := s.sql.client.Collection.Query().
		Where(
			ifPresent(collection.IDEQ, filter.ID),
			ifPresent(collection.NameEQ, filter.Name),
			ifPresent(collection.IDEQ, filter.UserID),
		).
		With(withCollectionInclude(include)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		All(ctx)

	return sqlColSliceToDQColSlice(cols), len(cols), err
}

func (s *CollectionService) FindCollection(ctx context.Context, filter dq.CollectionFilter, include dq.CollectionInclude) (*dq.Collection, error) {
	c, err := s.sql.client.Collection.Query().
		Where(
			ifPresent(collection.IDEQ, filter.ID),
			ifPresent(collection.NameEQ, filter.Name),
			collection.HasUserWith(ifPresent(user.IDEQ, filter.UserID)),
		).
		First(ctx)

	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, "Collection not found.")
	}
	if err != nil {
		return nil, err
	}

	return sqlColToDQCol(c), err
}

func (s *CollectionService) CreateCollection(ctx context.Context, col *dq.Collection) (*dq.Collection, error) {
	c, err := s.sql.client.Collection.Create().
		SetName(col.Name).
		SetUserID(col.UserID).
		Save(ctx)

	return sqlColToDQCol(c), err
}

func (s *CollectionService) UpdateCollection(ctx context.Context, id int, upd dq.CollectionUpdate) (c *dq.Collection, err error) {
	q := s.sql.client.Collection.UpdateOneID(id)

	if upd.Name != nil {
		q.SetName(*upd.Name)
	}
	if upd.UserID != nil {
		q.SetUserID(*upd.UserID)
	}
	if upd.FeedsIDs != nil {
		// TODO feeds
	}

	u, err := q.Save(ctx)

	return sqlColToDQCol(u), nil
}

func (s *CollectionService) DeleteCollection(ctx context.Context, id int) error {
	return s.sql.client.Collection.DeleteOneID(id).Exec(ctx)
}

func sqlColToDQCol(c *ent.Collection) *dq.Collection {
	n := &dq.Collection{
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

func sqlColSliceToDQColSlice(cs []*ent.Collection) (out []*dq.Collection) {
	for _, u := range cs {
		out = append(out, sqlColToDQCol(u))
	}

	return
}

func withCollectionInclude(include dq.CollectionInclude) func(q *ent.CollectionQuery) {
	return func(q *ent.CollectionQuery) {
		q.WithUser()

		if include.Feeds {
			q.WithFeeds()
		}
	}
}

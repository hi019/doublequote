package sql

import (
	"context"

	dq "doublequote"
	"doublequote/prisma"
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
	q, err := s.sql.prisma.Collection.
		FindFirst(prisma.Collection.ID.Equals(id)).
		With(
			buildCollectionInclude(include)...,
		).
		Exec(ctx)
	if err == prisma.ErrNotFound {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Collection")
	}
	if err != nil {
		return nil, err
	}

	return sqlColToDQCol(q), err
}

func (s *CollectionService) FindCollections(ctx context.Context, filter dq.CollectionFilter, include dq.CollectionInclude) ([]*dq.Collection, int, error) {
	cols, err := s.sql.prisma.Collection.FindMany(
		prisma.Collection.ID.EqualsIfPresent(filter.ID),
		prisma.Collection.Name.EqualsIfPresent(filter.Name),
		prisma.Collection.UserID.EqualsIfPresent(filter.UserID),
	).
		With(buildCollectionInclude(include)...).
		Skip(filter.Offset).
		Take(filter.Limit).
		Exec(ctx)

	// TODO implement Count when available https://github.com/prisma/prisma-client-go/issues/229

	return sqlColSliceToDQColSlice(cols), len(cols), err
}

func (s *CollectionService) FindCollection(ctx context.Context, filter dq.CollectionFilter, include dq.CollectionInclude) (*dq.Collection, error) {
	c, err := s.sql.prisma.Collection.FindFirst(
		prisma.Collection.ID.EqualsIfPresent(filter.ID),
		prisma.Collection.Name.EqualsIfPresent(filter.Name),
	).
		With(buildCollectionInclude(include)...).
		Exec(ctx)

	if err == prisma.ErrNotFound {
		return nil, dq.Errorf(dq.ENOTFOUND, "Collection not found.")
	}
	if err != nil {
		return nil, err
	}

	return sqlColToDQCol(c), err
}

func (s *CollectionService) CreateCollection(ctx context.Context, col *dq.Collection) (*dq.Collection, error) {
	c, err := s.sql.prisma.Collection.CreateOne(
		prisma.Collection.Name.Set(col.Name),
		prisma.Collection.User.Link(
			prisma.User.ID.Equals(col.UserID),
		),
	).
		With(
			prisma.Collection.User.Fetch(),
		).
		Exec(ctx)

	return sqlColToDQCol(c), err
}

func (s *CollectionService) UpdateCollection(ctx context.Context, id int, upd dq.CollectionUpdate) (c *dq.Collection, err error) {
	// TODO https://github.com/prisma/prisma-client-go/issues/699
	var dbU *prisma.CollectionModel

	// Update Collection attributes
	dbU, err = s.sql.prisma.Collection.FindUnique(prisma.Collection.ID.Equals(id)).
		Update(
			prisma.Collection.Name.SetIfPresent(upd.Name),
			prisma.Collection.UserID.SetIfPresent(upd.UserID),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Unlink Feeds from Collection
	_, err = s.sql.prisma.Prisma.ExecuteRaw(`DELETE FROM _CollectionToFeed WHERE A = ?`, id).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Link specified Feeds to the Collection
	for _, fID := range upd.FeedsIDs {
		dbU, err = s.sql.prisma.Collection.FindUnique(prisma.Collection.ID.Equals(id)).
			Update(
				prisma.Collection.Feeds.Link(prisma.Feed.ID.Equals(fID)),
			).
			Exec(ctx)
		if err == prisma.ErrNotFound {
			err = dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Collection")
		}
		if err != nil {
			return nil, err
		}
	}

	return sqlColToDQCol(dbU), nil
}

func (s *CollectionService) DeleteCollection(ctx context.Context, id int) error {
	_, err := s.sql.prisma.Collection.FindUnique(prisma.Collection.ID.Equals(id)).
		Delete().
		Exec(ctx)
	return err
}

func sqlColToDQCol(c *prisma.CollectionModel) *dq.Collection {
	n := &dq.Collection{
		ID:        c.ID,
		Name:      c.Name,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.RelationsCollection.User != nil {
		n.User = *sqlUserToDQUser(c.User())
	}

	if c.RelationsCollection.Feeds != nil {
		n.Feeds = sqlFeedSliceToFeedSlice(c.Feeds())
	}

	return n
}

func sqlColSliceToDQColSlice(cs []prisma.CollectionModel) (out []*dq.Collection) {
	for _, u := range cs {
		out = append(out, sqlColToDQCol(&u))
	}

	return
}

func buildCollectionInclude(include dq.CollectionInclude) (filters []prisma.ICollectionRelationWith) {
	filters = append(filters, prisma.Collection.User.Fetch())

	if include.Feeds {
		filters = append(filters, prisma.Collection.Feeds.Fetch())
	}

	return
}

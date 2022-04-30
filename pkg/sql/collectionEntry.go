package sql

import (
	"context"

	"doublequote/ent"
	"doublequote/ent/collectionentry"
	"doublequote/pkg/domain"
)

var _ domain.CollectionEntryService = (*CollectionEntryService)(nil)

type CollectionEntryService struct {
	sql *SQL
}

func NewCollectionEntryService(sql *SQL) *CollectionEntryService {
	return &CollectionEntryService{sql: sql}
}

func (s *CollectionEntryService) FindCollectionEntryByID(ctx context.Context, id int, include domain.CollectionEntryInclude) (*domain.CollectionEntry, error) {
	c, err := s.sql.client.CollectionEntry.
		Query().
		With(withCollectionEntryInclude(include)).
		Where(collectionentry.IDEQ(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection entry")
	}
	if err != nil {
		return nil, err
	}

	return transformColEntry(c), err
}

func (s *CollectionEntryService) FindCollectionEntries(ctx context.Context, filter domain.CollectionEntryFilter, include domain.CollectionEntryInclude) ([]*domain.CollectionEntry, int, error) {
	cols, err := s.sql.client.CollectionEntry.Query().
		Where(
			ifPresent(collectionentry.IDEQ, filter.ID),
			ifPresent(collectionentry.CollectionID, filter.CollectionID),
			ifPresent(collectionentry.EntryID, filter.EntryID),
		).
		With(withCollectionEntryInclude(include)).
		Limit(filter.Limit).
		Offset(filter.Offset).
		All(ctx)

	return transformColEntryArray(cols), len(cols), err
}

func (s *CollectionEntryService) FindCollectionEntry(ctx context.Context, filter domain.CollectionEntryFilter, include domain.CollectionEntryInclude) (*domain.CollectionEntry, error) {
	c, err := s.sql.client.CollectionEntry.Query().
		Where(
			ifPresent(collectionentry.IDEQ, filter.ID),
			ifPresent(collectionentry.CollectionID, filter.CollectionID),
			ifPresent(collectionentry.EntryID, filter.EntryID),
		).
		First(ctx)

	if ent.IsNotFound(err) {
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection entry")
	}
	if err != nil {
		return nil, err
	}

	return transformColEntry(c), err
}

func (s *CollectionEntryService) CreateCollectionEntry(ctx context.Context, ce *domain.CollectionEntry) (*domain.CollectionEntry, error) {
	c, err := s.sql.client.CollectionEntry.Create().
		SetCollectionID(ce.CollectionID).
		SetEntryID(ce.EntryID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return transformColEntry(c), err
}

func (s *CollectionEntryService) UpdateCollectionEntry(ctx context.Context, id int, upd domain.CollectionEntryUpdate) (c *domain.CollectionEntry, err error) {
	q := s.sql.client.CollectionEntry.UpdateOneID(id)

	if upd.EntryID != nil {
		q.SetEntryID(*upd.EntryID)
	}
	if upd.CollectionID != nil {
		q.SetCollectionID(*upd.CollectionID)
	}

	ce, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return transformColEntry(ce), err
}

func (s *CollectionEntryService) DeleteCollectionEntry(ctx context.Context, id int) error {
	return s.sql.client.Collection.DeleteOneID(id).Exec(ctx)
}

// TODO rename other transformers
// transformColEntry transforms a *ent.CollectionEntry to a *domain.CollectionEntry
func transformColEntry(ce *ent.CollectionEntry) *domain.CollectionEntry {
	n := &domain.CollectionEntry{
		ID:        ce.ID,
		CreatedAt: ce.CreatedAt,
		UpdatedAt: ce.UpdatedAt,
	}

	if u, err := ce.Edges.EntryOrErr(); err != nil {
		n.Entry = *sqlEntryToEntry(u)
	}

	return n
}

// TODO rename other transformers
// transformColEntry transforms an array of Ent collection entry to an array of domain collection entry
func transformColEntryArray(cs []*ent.CollectionEntry) (out []*domain.CollectionEntry) {
	for _, c := range cs {
		out = append(out, transformColEntry(c))
	}

	return
}

func withCollectionEntryInclude(include domain.CollectionEntryInclude) func(q *ent.CollectionEntryQuery) {
	return func(q *ent.CollectionEntryQuery) {
		if include.Collection {
			q.WithCollection()
		}

		if include.Entry {
			q.WithEntry()
		}
	}
}

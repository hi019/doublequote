package sql

//
//import (
//	"context"
//
//	"doublequote/ent"
//	"doublequote/ent/collection"
//	"doublequote/ent/collectionEntry"
//	"doublequote/ent/feed"
//	"doublequote/ent/user"
//	"doublequote/pkg/domain"
//)
//
//// Ensure service implements interface.
//var _ domain.CollectionEntryService = (*CollectionEntryService)(nil)
//
//type CollectionEntryService struct {
//	sql *SQL
//}
//
//func NewCollectionEntryService(sql *SQL) *CollectionEntryService {
//	return &CollectionEntryService{sql: sql}
//}
//
//func (s *CollectionEntryService) FindCollectionEntryByID(ctx context.Context, id int, include domain.CollectionEntryInclude) (*domain.CollectionEntry, error) {
//	c, err := s.sql.client.CollectionEntry.
//		Query().
//		With(withCollectionInclude(include)).
//		Where(collectionentry.IDEQ(id)).
//		Only(ctx)
//	if ent.IsNotFound(err) {
//		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection entry")
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	return sqlColToDQCol(c), err
//}
//
//func (s *CollectionEntryService) FindCollectionEntries(ctx context.Context, filter domain.CollectionEntryFilter, include domain.CollectionEntryInclude) ([]*domain.CollectionEntry, int, error) {
//	cols, err := s.sql.client.Collection.Query().
//		Where(
//			ifPresent(collection.IDEQ, filter.ID),
//			ifPresent(collection.NameEQ, filter.Name),
//			ifPresent(collection.IDEQ, filter.UserID),
//		).
//		With(withCollectionInclude(include)).
//		Limit(filter.Limit).
//		Offset(filter.Offset).
//		All(ctx)
//
//	return sqlColSliceToDQColSlice(cols), len(cols), err
//}
//
//func (s *CollectionEntryService) FindCollectionEntry(ctx context.Context, filter domain.CollectionEntryFilter, include domain.CollectionEntryInclude) (*domain.CollectionEntry, error) {
//	c, err := s.sql.client.Collection.Query().
//		Where(
//			ifPresent(collection.IDEQ, filter.ID),
//			ifPresent(collection.NameEQ, filter.Name),
//			collection.HasUserWith(ifPresent(user.IDEQ, filter.UserID)),
//			collection.HasFeedsWith(ifPresent(feed.IDEQ, filter.FeedID)),
//		).
//		First(ctx)
//
//	if ent.IsNotFound(err) {
//		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Collection")
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	return sqlColToDQCol(c), err
//}
//
//func (s *CollectionEntryService) CreateCollectionEntry(ctx context.Context, col *domain.CollectionEntry) (*domain.CollectionEntry, error) {
//	c, err := s.sql.client.Collection.Create().
//		SetName(col.Name).
//		SetUserID(col.UserID).
//		Save(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	// TODO access c.User.ID without another query?
//	return s.FindCollectionByID(ctx, c.ID, domain.CollectionInclude{})
//}
//
//func (s *CollectionEntryService) UpdateCollectionEntry(ctx context.Context, id int, upd domain.CollectionEntryUpdate) (c *domain.CollectionEntry, err error) {
//	q := s.sql.client.Collection.UpdateOneID(id)
//
//	if upd.Name != nil {
//		q.SetName(*upd.Name)
//	}
//	if upd.UserID != nil {
//		q.SetUserID(*upd.UserID)
//	}
//	if upd.FeedsIDs != nil {
//		q.ClearFeeds().AddFeedIDs(*upd.FeedsIDs...)
//	}
//
//	u, err := q.Save(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	// TODO access s.User.ID without another query?
//	return s.FindCollectionByID(ctx, u.ID, domain.CollectionInclude{})
//}
//
//func (s *CollectionEntryService) DeleteCollectionEntry(ctx context.Context, id int) error {
//	return s.sql.client.Collection.DeleteOneID(id).Exec(ctx)
//}
//
//// TODO rename other transformers
//// transformColEntry transforms a *ent.CollectionEntry to a *domain.CollectionEntry
//func transformColEntry(ce *ent.CollectionEntry) *domain.CollectionEntry {
//	n := &domain.CollectionEntry{
//		ID:        ce.ID,
//		CreatedAt: ce.CreatedAt,
//		UpdatedAt: ce.UpdatedAt,
//	}
//
//	if u, err := ce.Edges.EntryOrErr(); err != nil {
//		n.Entry = *sqlEntryToEntry(u)
//	}
//
//	return n
//}
//
//// TODO rename other transformers
//// transformColEntry transforms an array of Ent collection entry to an array of domain collection entry
//func transformColEntryArray(cs []*ent.CollectionEntry) (out []*domain.CollectionEntry) {
//	for _, u := range cs {
//		out = append(out, sqlColToDQCol(u))
//	}
//
//	return
//}
//
//func withCollectionInclude(include domain.CollectionInclude) func(q *ent.CollectionQuery) {
//	return func(q *ent.CollectionQuery) {
//		q.WithUser()
//
//		if include.Feeds {
//			q.WithFeeds()
//		}
//	}
//}

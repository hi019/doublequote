package sql

import (
	"context"

	"doublequote/ent"
	"doublequote/ent/entry"
	"doublequote/ent/feed"
	"doublequote/pkg/domain"
)

var _ domain.EntryService = (*EntryService)(nil)

type EntryService struct {
	sql *SQL
}

func (s *EntryService) CreateManyEntry(ctx context.Context, entries []domain.Entry) ([]*domain.Entry, error) {
	bulk := make([]*ent.EntryCreate, len(entries))
	for i, e := range entries {
		bulk[i] = s.sql.client.Entry.Create().
			SetTitle(e.Title).
			SetURL(e.URL).
			SetAuthor(e.Author).
			SetContentKey(e.ContentKey).
			SetFeedID(e.FeedID)
	}
	created, err := s.sql.client.Entry.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, err
	}

	return sqlEntrySliceToEntrySlice(created), nil
}

func NewEntryService(sql *SQL) *EntryService {
	return &EntryService{sql: sql}
}

func (s *EntryService) FindEntryByID(ctx context.Context, id int, include domain.EntryInclude) (*domain.Entry, error) {
	e, err := s.sql.client.Entry.
		Query().
		With(withEntryInclude(include)).
		Where(entry.IDEQ(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Entry")
	}
	if err != nil {
		return nil, err
	}

	return sqlEntryToEntry(e), nil
}

func (s *EntryService) FindEntry(ctx context.Context, filter domain.EntryFilter, include domain.EntryInclude) (*domain.Entry, error) {
	c, err := s.sql.client.Entry.Query().
		Where(
			ifPresent(entry.IDEQ, filter.ID),
			ifPresent(entry.TitleEQ, filter.Title),
			entry.HasFeedWith(ifPresent(feed.IDEQ, filter.FeedID)),
		).
		First(ctx)

	if ent.IsNotFound(err) {
		return nil, domain.Errorf(domain.ENOTFOUND, domain.ErrNotFound, "Entry")
	}
	if err != nil {
		return nil, err
	}

	return sqlEntryToEntry(c), err
}

func (s *EntryService) CreateEntry(ctx context.Context, entry domain.Entry) (*domain.Entry, error) {
	e, err := s.sql.client.Entry.Create().
		SetTitle(entry.Title).
		SetURL(entry.URL).
		SetAuthor(entry.Author).
		SetContentKey(entry.ContentKey).
		SetFeedID(entry.FeedID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.FindEntryByID(ctx, e.ID, domain.EntryInclude{})
}

func (s *EntryService) UpdateEntry(ctx context.Context, id int, upd domain.EntryUpdate) (*domain.Entry, error) {
	q := s.sql.client.Entry.UpdateOneID(id)

	if upd.Author != nil {
		q.SetAuthor(*upd.Author)
	}
	if upd.FeedID != nil {
		q.SetFeedID(*upd.FeedID)
	}
	if upd.Title != nil {
		q.SetTitle(*upd.Title)
	}
	if upd.URL != nil {
		q.SetURL(*upd.URL)
	}
	if upd.ContentKey != nil {
		q.SetContentKey(*upd.ContentKey)
	}

	u, err := q.Save(ctx)
	if err != nil {
		return nil, err
	}

	return s.FindEntryByID(ctx, u.ID, domain.EntryInclude{})
}

func sqlEntryToEntry(c *ent.Entry) *domain.Entry {
	e := &domain.Entry{
		ID:        c.ID,
		Title:     c.Title,
		URL:       c.URL,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if f, err := c.Edges.FeedOrErr(); err == nil {
		e.Feed = *sqlFeedToFeed(f)
	}

	return e
}

func sqlEntrySliceToEntrySlice(cs []*ent.Entry) (out []*domain.Entry) {
	for _, u := range cs {
		out = append(out, sqlEntryToEntry(u))
	}

	return
}

func withEntryInclude(include domain.EntryInclude) func(q *ent.EntryQuery) {
	return func(q *ent.EntryQuery) {
		q.WithFeed()
	}
}

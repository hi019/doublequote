package sql

import (
	"context"

	dq "doublequote"
	"doublequote/ent"
	"doublequote/ent/entry"
	"doublequote/ent/feed"
)

var _ dq.EntryService = (*EntryService)(nil)

type EntryService struct {
	sql *SQL
}

func (s *EntryService) CreateManyEntry(ctx context.Context, entries []dq.Entry) ([]*dq.Entry, error) {
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

func (s *EntryService) FindEntryByID(ctx context.Context, id int, include dq.EntryInclude) (*dq.Entry, error) {
	e, err := s.sql.client.Entry.
		Query().
		With(withEntryInclude(include)).
		Where(entry.IDEQ(id)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Entry")
	}
	if err != nil {
		return nil, err
	}

	return sqlEntryToEntry(e), nil
}

func (s *EntryService) FindEntry(ctx context.Context, filter dq.EntryFilter, include dq.EntryInclude) (*dq.Entry, error) {
	c, err := s.sql.client.Entry.Query().
		Where(
			ifPresent(entry.IDEQ, filter.ID),
			ifPresent(entry.TitleEQ, filter.Title),
			entry.HasFeedWith(ifPresent(feed.IDEQ, filter.FeedID)),
		).
		First(ctx)

	if ent.IsNotFound(err) {
		return nil, dq.Errorf(dq.ENOTFOUND, dq.ErrNotFound, "Entry")
	}
	if err != nil {
		return nil, err
	}

	return sqlEntryToEntry(c), err
}

func (s *EntryService) CreateEntry(ctx context.Context, entry dq.Entry) (*dq.Entry, error) {
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

	return s.FindEntryByID(ctx, e.ID, dq.EntryInclude{})
}

func (s *EntryService) UpdateEntry(ctx context.Context, id int, upd dq.EntryUpdate) (*dq.Entry, error) {
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

	return s.FindEntryByID(ctx, u.ID, dq.EntryInclude{})
}

func sqlEntryToEntry(c *ent.Entry) *dq.Entry {
	e := &dq.Entry{
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

func sqlEntrySliceToEntrySlice(cs []*ent.Entry) (out []*dq.Entry) {
	for _, u := range cs {
		out = append(out, sqlEntryToEntry(u))
	}

	return
}

func withEntryInclude(include dq.EntryInclude) func(q *ent.EntryQuery) {
	return func(q *ent.EntryQuery) {
		q.WithFeed()
	}
}

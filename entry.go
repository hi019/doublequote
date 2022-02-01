package dq

import (
	"context"
	"time"
)

type Entry struct {
	ID int

	Title      string
	URL        string
	Author     string
	ContentKey string

	Feed   Feed
	FeedID int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type EntryService interface {
	FindEntryByID(ctx context.Context, id int, include EntryInclude) (*Entry, error)
	FindEntry(ctx context.Context, filter EntryFilter, include EntryInclude) (*Entry, error)
	CreateEntry(ctx context.Context, entry Entry) (*Entry, error)
	CreateManyEntry(ctx context.Context, entry []Entry) ([]*Entry, error)
	UpdateEntry(ctx context.Context, id int, upd EntryUpdate) (*Entry, error)
}

type EntryFilter struct {
	ID         *int
	Title      *string
	URL        *string
	Author     *string
	ContentKey *string
	FeedID     *int

	Offset int
	Limit  int
}

type EntryUpdate struct {
	Title      *string
	URL        *string
	Author     *string
	ContentKey *string
	FeedID     *int
}

type EntryInclude struct {
}

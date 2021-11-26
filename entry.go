package dq

import (
	"context"
	"time"
)

type Entry struct {
	ID int

	Title  string
	URL    string
	Author string

	Feed   Feed
	FeedID int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type EntryService interface {
	FindEntryByID(ctx context.Context, id int, include EntryInclude) (*Entry, error)
}

type EntryInclude struct {
}

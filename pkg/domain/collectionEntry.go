package domain

import (
	"context"
	"time"
)

type CollectionEntry struct {
	ID int

	IsRead bool

	Collection   Collection
	CollectionID int

	Entry   Entry
	EntryID int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CollectionEntryService interface {
	FindCollectionEntryByID(ctx context.Context, id int, include CollectionEntryInclude) (*CollectionEntry, error)
	FindCollectionEntries(ctx context.Context, filter CollectionEntryFilter, include CollectionEntryInclude) ([]*CollectionEntry, int, error)
	FindCollectionEntry(ctx context.Context, filter CollectionEntryFilter, include CollectionEntryInclude) (*CollectionEntry, error)
	CreateCollectionEntry(ctx context.Context, col *CollectionEntry) (*CollectionEntry, error)
	CreateManyCollectionEntry(ctx context.Context, col []CollectionEntry) ([]CollectionEntry, error)
	UpdateCollectionEntry(ctx context.Context, id int, upd CollectionEntryUpdate) (*CollectionEntry, error)
	DeleteCollectionEntry(ctx context.Context, id int) error
}

type CollectionEntryFilter struct {
	ID *int

	CollectionID *int
	EntryID      *int

	Offset int
	Limit  int
}

type CollectionEntryUpdate struct {
	IsRead       *bool
	CollectionID *int
	EntryID      *int
}

type CollectionEntryInclude struct {
	Collection bool
	Entry      bool
}

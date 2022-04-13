package domain

import (
	"context"
	"time"
)

type Collection struct {
	ID int

	Name string

	User   User
	UserID int

	Feeds []*Feed

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CollectionService interface {
	FindCollectionByID(ctx context.Context, id int, include CollectionInclude) (*Collection, error)

	FindCollections(ctx context.Context, filter CollectionFilter, include CollectionInclude) ([]*Collection, int, error)

	FindCollection(ctx context.Context, filter CollectionFilter, include CollectionInclude) (*Collection, error)

	CreateCollection(ctx context.Context, col *Collection) (*Collection, error)

	UpdateCollection(ctx context.Context, id int, upd CollectionUpdate) (*Collection, error)

	DeleteCollection(ctx context.Context, id int) error
}

type CollectionFilter struct {
	ID     *int
	Name   *string
	UserID *int

	FeedID *int

	Offset int
	Limit  int
}

type CollectionUpdate struct {
	Name     *string
	UserID   *int
	FeedsIDs *[]int
}

type CollectionInclude struct {
	Feeds bool
}

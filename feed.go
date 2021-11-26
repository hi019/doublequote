package dq

import (
	"context"
	"time"
)

type Feed struct {
	ID int

	Name   string
	RssURL string
	Domain string

	Collections []*Collection

	CreatedAt time.Time
	UpdatedAt time.Time
}

type FeedService interface {
	FindFeedByID(ctx context.Context, id int, include FeedInclude) (*Feed, error)

	FindFeeds(ctx context.Context, filter FeedFilter, include FeedInclude) ([]*Feed, int, error)

	FindFeed(ctx context.Context, filter FeedFilter, include FeedInclude) (*Feed, error)

	CreateFeed(ctx context.Context, feed *Feed) (*Feed, error)

	UpdateFeed(ctx context.Context, id int, upd FeedUpdate) (*Feed, error)

	DeleteFeed(ctx context.Context, id int) error
}

type FeedFilter struct {
	ID *int

	Name         *string
	RssURL       *string
	Domain       *string
	CollectionID *int

	Offset int
	Limit  int
}

type FeedUpdate struct {
	Name   *string
	RssURL *string
	Domain *string
}

type FeedInclude struct {
	Collections bool
}

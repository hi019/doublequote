package ingest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	dq "doublequote"
	"github.com/google/uuid"

	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
	"golang.org/x/sync/errgroup"
)

type IngestService struct {
	feedService    dq.FeedService
	entryService   dq.EntryService
	storageService dq.StorageService
	fp             *gofeed.Parser
	http           *http.Client
}

func NewService(feedService dq.FeedService, entryService dq.EntryService, storageService dq.StorageService) *IngestService {
	fp := gofeed.NewParser()
	// TODO rate limiting
	h := http.Client{
		Timeout: 10 * time.Second,
	}

	return &IngestService{
		feedService:    feedService,
		entryService:   entryService,
		storageService: storageService,
		fp:             fp,
		http:           &h,
	}
}

func (s *IngestService) Ingest(feed dq.Feed) error {
	items, err := s.getItems(feed)
	if err != nil {
		return err
	}

	fetchItem := func(item *gofeed.Item) (*dq.Entry, error) {
		e := parseItem(item, feed)

		e, err = s.saveEntryContent(e)
		if err != nil {
			return e, err
		}

		return e, nil
	}

	var out []dq.Entry

	g := errgroup.Group{}
	for _, item := range items {
		i := item
		g.Go(func() error {
			e, err := fetchItem(i)
			if err != nil {
				return err
			}
			out = append(out, *e)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	_, err = s.entryService.CreateManyEntry(context.Background(), out)
	return err
}

func (s *IngestService) getItems(feed dq.Feed) ([]*gofeed.Item, error) {
	res, err := http.Get(feed.RssURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err := s.fp.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	return f.Items, nil
}

func (s *IngestService) saveEntryContent(entry *dq.Entry) (*dq.Entry, error) {
	res, err := http.Get(entry.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	parsedURL, err := url.Parse(entry.URL)
	if err != nil {
		return nil, err
	}
	content, err := readability.FromReader(res.Body, parsedURL)

	key := fmt.Sprintf("content_%s", uuid.New().String())
	if err := s.storageService.Set(context.Background(), key, []byte(content.Content)); err != nil {
		return nil, err
	}

	entry.ContentKey = key

	return entry, nil

}

func parseItem(i *gofeed.Item, f dq.Feed) *dq.Entry {
	author := "Unknown"
	if len(i.Authors) > 0 {
		author = i.Authors[0].Name
	}

	return &dq.Entry{
		Title:  i.Title,
		URL:    i.Link,
		Author: author,
		FeedID: f.ID,
	}
}

package ingest

// TODO rename to parser?

import (
	"net/http"
	"net/url"
	"time"

	"doublequote/pkg/domain"
	"github.com/go-shiori/go-readability"
	"github.com/mmcdole/gofeed"
)

type Service struct {
	feedService    domain.FeedService
	entryService   domain.EntryService
	storageService domain.StorageService
	fp             *gofeed.Parser
	http           *http.Client
}

func NewService(feedService domain.FeedService, entryService domain.EntryService, storageService domain.StorageService) *Service {
	fp := gofeed.NewParser()
	// TODO rate limiting
	h := http.Client{
		Timeout: 10 * time.Second,
	}

	return &Service{
		feedService:    feedService,
		entryService:   entryService,
		storageService: storageService,
		fp:             fp,
		http:           &h,
	}
}

func (s *Service) GetEntries(feed *domain.Feed) (entries []*domain.Entry, err error) {
	tmpItems, err := s.getItems(feed)
	if err != nil {
		return nil, err
	}

	for _, item := range tmpItems {
		entries = append(entries, parseItem(item, feed))
	}

	return
}

func (s *Service) GetEntryContent(entry *domain.Entry) (content string, err error) {
	return s.getEntryContent(entry)
}

func (s *Service) getItems(feed *domain.Feed) ([]*gofeed.Item, error) {
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

func (s *Service) getEntryContent(entry *domain.Entry) (string, error) {
	res, err := http.Get(entry.URL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	parsedURL, err := url.Parse(entry.URL)
	if err != nil {
		return "", err
	}
	content, err := readability.FromReader(res.Body, parsedURL)
	if err != nil {
		return "", err
	}

	return content.Content, nil

}

func parseItem(i *gofeed.Item, f *domain.Feed) *domain.Entry {
	author := "Unknown"
	if len(i.Authors) > 0 {
		author = i.Authors[0].Name
	}

	return &domain.Entry{
		Title:  i.Title,
		URL:    i.Link,
		Author: author,
		FeedID: f.ID,
	}
}

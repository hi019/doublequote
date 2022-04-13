package job

import (
	"context"
	"fmt"

	"doublequote/pkg/domain"
	"github.com/fatih/semgroup"
	"github.com/google/uuid"
)

type IngestJob struct {
	svc               domain.IngestService
	feedService       domain.FeedService
	storageService    domain.StorageService
	entryService      domain.EntryService
	collectionService domain.CollectionService
}

func NewIngestJob(
	svc domain.IngestService,
	feedService domain.FeedService,
	storageService domain.StorageService,
	entryService domain.EntryService,
	collectionService domain.CollectionService,
) *IngestJob {
	return &IngestJob{svc: svc, feedService: feedService, storageService: storageService, entryService: entryService, collectionService: collectionService}
}

func (j *IngestJob) Run() error {
	feeds, _, err := j.feedService.FindFeeds(context.Background(), domain.FeedFilter{}, domain.FeedInclude{})
	if err != nil {
		return err
	}

	s := semgroup.NewGroup(context.Background(), 5)

	for _, feed := range feeds {
		s.Go(func() error {
			return j.ingestFeed(feed)
		})
	}

	if err := s.Wait(); err != nil {
		return err
	}

	return nil
}

func (j *IngestJob) ingestFeed(feed *domain.Feed) error {
	var toSave []domain.Entry

	entries, err := j.svc.GetEntries(feed)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !j.shouldSave(entry) {
			continue
		}

		content, err := j.svc.GetEntryContent(entry)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("content_%s", uuid.New().String())
		if err := j.storageService.Set(context.Background(), key, []byte(content)); err != nil {
			return err
		}

		entry.ContentKey = key

		toSave = append(toSave, *entry)
	}

	_, err = j.entryService.CreateManyEntry(context.TODO(), toSave)
	return err
}

//func (j *IngestJob) createCollectionEntries(feed domain.Feed, entry domain.Entry) error {
//	collection, err := j.collectionService.FindCollection(
//		context.Background(),
//		domain.CollectionFilter{FeedID: utils.Ptr(feed.ID)},
//		domain.CollectionInclude{},
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (j *IngestJob) shouldSave(entry *domain.Entry) bool {
	// TODO maybe we shouldn't ignore the error here..
	_, err := j.entryService.FindEntry(context.TODO(), domain.EntryFilter{URL: &entry.URL}, domain.EntryInclude{})

	return domain.ErrorCode(err) == domain.ENOTFOUND
}

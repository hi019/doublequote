package job

import (
	"context"

	"doublequote/pkg/domain"
	"golang.org/x/sync/errgroup"
)

type IngestJob struct {
	svc         domain.IngestService
	feedService domain.FeedService
}

func NewIngestJob(
	svc domain.IngestService,
	feedService domain.FeedService,
) *IngestJob {
	return &IngestJob{svc, feedService}
}

func (j *IngestJob) Run() error {
	feeds, _, err := j.feedService.FindFeeds(context.Background(), domain.FeedFilter{}, domain.FeedInclude{})
	if err != nil {
		return err
	}

	g := errgroup.Group{}

	for _, feed := range feeds {
		g.Go(func() error {
			return j.svc.Ingest(*feed)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (j *IngestJob) shouldIngest(entry *domain.Entry) {

}

package domain

type IngestService interface {
	Ingest(feed Feed) error
}

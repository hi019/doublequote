package dq

type IngestService interface {
	Ingest(feed Feed) error
}

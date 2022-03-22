package domain

type IngestService interface {
	GetEntries(feed *Feed) (entries []*Entry, err error)
	GetEntryContent(entry *Entry) (content string)
}

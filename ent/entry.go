// Code generated by entc, DO NOT EDIT.

package ent

import (
	"doublequote/ent/entry"
	"doublequote/ent/feed"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
)

// Entry is the model entity for the Entry schema.
type Entry struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// URL holds the value of the "url" field.
	URL string `json:"url,omitempty"`
	// Author holds the value of the "author" field.
	Author string `json:"author,omitempty"`
	// ContentKey holds the value of the "content_key" field.
	ContentKey string `json:"content_key,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the EntryQuery when eager-loading is set.
	Edges        EntryEdges `json:"edges"`
	feed_entries *int
}

// EntryEdges holds the relations/edges for other nodes in the graph.
type EntryEdges struct {
	// Feed holds the value of the feed edge.
	Feed *Feed `json:"feed,omitempty"`
	// CollectionEntries holds the value of the collection_entries edge.
	CollectionEntries []*CollectionEntry `json:"collection_entries,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// FeedOrErr returns the Feed value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e EntryEdges) FeedOrErr() (*Feed, error) {
	if e.loadedTypes[0] {
		if e.Feed == nil {
			// The edge feed was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: feed.Label}
		}
		return e.Feed, nil
	}
	return nil, &NotLoadedError{edge: "feed"}
}

// CollectionEntriesOrErr returns the CollectionEntries value or an error if the edge
// was not loaded in eager-loading.
func (e EntryEdges) CollectionEntriesOrErr() ([]*CollectionEntry, error) {
	if e.loadedTypes[1] {
		return e.CollectionEntries, nil
	}
	return nil, &NotLoadedError{edge: "collection_entries"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Entry) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case entry.FieldID:
			values[i] = new(sql.NullInt64)
		case entry.FieldTitle, entry.FieldURL, entry.FieldAuthor, entry.FieldContentKey:
			values[i] = new(sql.NullString)
		case entry.FieldCreatedAt, entry.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case entry.ForeignKeys[0]: // feed_entries
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Entry", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Entry fields.
func (e *Entry) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case entry.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			e.ID = int(value.Int64)
		case entry.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				e.Title = value.String
			}
		case entry.FieldURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field url", values[i])
			} else if value.Valid {
				e.URL = value.String
			}
		case entry.FieldAuthor:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field author", values[i])
			} else if value.Valid {
				e.Author = value.String
			}
		case entry.FieldContentKey:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content_key", values[i])
			} else if value.Valid {
				e.ContentKey = value.String
			}
		case entry.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				e.CreatedAt = value.Time
			}
		case entry.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				e.UpdatedAt = value.Time
			}
		case entry.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field feed_entries", value)
			} else if value.Valid {
				e.feed_entries = new(int)
				*e.feed_entries = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryFeed queries the "feed" edge of the Entry entity.
func (e *Entry) QueryFeed() *FeedQuery {
	return (&EntryClient{config: e.config}).QueryFeed(e)
}

// QueryCollectionEntries queries the "collection_entries" edge of the Entry entity.
func (e *Entry) QueryCollectionEntries() *CollectionEntryQuery {
	return (&EntryClient{config: e.config}).QueryCollectionEntries(e)
}

// Update returns a builder for updating this Entry.
// Note that you need to call Entry.Unwrap() before calling this method if this Entry
// was returned from a transaction, and the transaction was committed or rolled back.
func (e *Entry) Update() *EntryUpdateOne {
	return (&EntryClient{config: e.config}).UpdateOne(e)
}

// Unwrap unwraps the Entry entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (e *Entry) Unwrap() *Entry {
	tx, ok := e.config.driver.(*txDriver)
	if !ok {
		panic("ent: Entry is not a transactional entity")
	}
	e.config.driver = tx.drv
	return e
}

// String implements the fmt.Stringer.
func (e *Entry) String() string {
	var builder strings.Builder
	builder.WriteString("Entry(")
	builder.WriteString(fmt.Sprintf("id=%v", e.ID))
	builder.WriteString(", title=")
	builder.WriteString(e.Title)
	builder.WriteString(", url=")
	builder.WriteString(e.URL)
	builder.WriteString(", author=")
	builder.WriteString(e.Author)
	builder.WriteString(", content_key=")
	builder.WriteString(e.ContentKey)
	builder.WriteString(", created_at=")
	builder.WriteString(e.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(e.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Entries is a parsable slice of Entry.
type Entries []*Entry

func (e Entries) config(cfg config) {
	for _i := range e {
		e[_i].config = cfg
	}
}

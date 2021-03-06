// Code generated by entc, DO NOT EDIT.

package ent

import (
	"doublequote/ent/feed"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
)

// Feed is the model entity for the Feed schema.
type Feed struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// RssURL holds the value of the "rssURL" field.
	RssURL string `json:"rssURL,omitempty"`
	// Domain holds the value of the "domain" field.
	Domain string `json:"domain,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the FeedQuery when eager-loading is set.
	Edges FeedEdges `json:"edges"`
}

// FeedEdges holds the relations/edges for other nodes in the graph.
type FeedEdges struct {
	// Collections holds the value of the collections edge.
	Collections []*Collection `json:"collections,omitempty"`
	// Entries holds the value of the entries edge.
	Entries []*Entry `json:"entries,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// CollectionsOrErr returns the Collections value or an error if the edge
// was not loaded in eager-loading.
func (e FeedEdges) CollectionsOrErr() ([]*Collection, error) {
	if e.loadedTypes[0] {
		return e.Collections, nil
	}
	return nil, &NotLoadedError{edge: "collections"}
}

// EntriesOrErr returns the Entries value or an error if the edge
// was not loaded in eager-loading.
func (e FeedEdges) EntriesOrErr() ([]*Entry, error) {
	if e.loadedTypes[1] {
		return e.Entries, nil
	}
	return nil, &NotLoadedError{edge: "entries"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Feed) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case feed.FieldID:
			values[i] = new(sql.NullInt64)
		case feed.FieldName, feed.FieldRssURL, feed.FieldDomain:
			values[i] = new(sql.NullString)
		case feed.FieldCreatedAt, feed.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Feed", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Feed fields.
func (f *Feed) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case feed.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			f.ID = int(value.Int64)
		case feed.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				f.Name = value.String
			}
		case feed.FieldRssURL:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field rssURL", values[i])
			} else if value.Valid {
				f.RssURL = value.String
			}
		case feed.FieldDomain:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field domain", values[i])
			} else if value.Valid {
				f.Domain = value.String
			}
		case feed.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				f.CreatedAt = value.Time
			}
		case feed.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				f.UpdatedAt = value.Time
			}
		}
	}
	return nil
}

// QueryCollections queries the "collections" edge of the Feed entity.
func (f *Feed) QueryCollections() *CollectionQuery {
	return (&FeedClient{config: f.config}).QueryCollections(f)
}

// QueryEntries queries the "entries" edge of the Feed entity.
func (f *Feed) QueryEntries() *EntryQuery {
	return (&FeedClient{config: f.config}).QueryEntries(f)
}

// Update returns a builder for updating this Feed.
// Note that you need to call Feed.Unwrap() before calling this method if this Feed
// was returned from a transaction, and the transaction was committed or rolled back.
func (f *Feed) Update() *FeedUpdateOne {
	return (&FeedClient{config: f.config}).UpdateOne(f)
}

// Unwrap unwraps the Feed entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (f *Feed) Unwrap() *Feed {
	tx, ok := f.config.driver.(*txDriver)
	if !ok {
		panic("ent: Feed is not a transactional entity")
	}
	f.config.driver = tx.drv
	return f
}

// String implements the fmt.Stringer.
func (f *Feed) String() string {
	var builder strings.Builder
	builder.WriteString("Feed(")
	builder.WriteString(fmt.Sprintf("id=%v", f.ID))
	builder.WriteString(", name=")
	builder.WriteString(f.Name)
	builder.WriteString(", rssURL=")
	builder.WriteString(f.RssURL)
	builder.WriteString(", domain=")
	builder.WriteString(f.Domain)
	builder.WriteString(", created_at=")
	builder.WriteString(f.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(f.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Feeds is a parsable slice of Feed.
type Feeds []*Feed

func (f Feeds) config(cfg config) {
	for _i := range f {
		f[_i].config = cfg
	}
}

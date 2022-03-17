package sql

import (
	"context"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"

	"doublequote/ent"

	_ "doublequote/ent/entps"
)

// SQL represents the database connection
type SQL struct {
	client *ent.Client
	dbUrl  string

	// Returns the current time. Defaults to time.Now()
	// Can be mocked for tests.
	Now func() time.Time
}

// NewSQL returns a new instance of SQL associated with the given datasource name
func NewSQL(dbUrl string) *SQL {
	return &SQL{
		dbUrl: dbUrl,
		Now:   time.Now,
	}
}

// Open opens the database connection
func (sql *SQL) Open() error {
	c, err := createClient(sql.dbUrl)
	if err != nil {
		return err
	}

	// Migrate database
	err = c.Schema.Create(context.Background(), schema.WithAtlas(true))
	if err != nil {
		return err
	}

	sql.client = c

	return nil
}

// Close closes the database connection
func (sql *SQL) Close() error {
	return sql.client.Close()
}

// TODO add postgres support
func createClient(fileName string) (*ent.Client, error) {
	return ent.Open("sqlite3", fileName)
}

// Utility function. Receives a function and a pointer to an argument. If the argument is nil, an empty selector
// is returned. If it's not, the argument is passed to the function and the result is returned
func ifPresent[V any, P ~func(*entsql.Selector)](f func(V) P, v *V) P {
	if v != nil {
		return f(*v)
	} else {
		return func(*entsql.Selector) {}
	}
}

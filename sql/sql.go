package sql

import (
	"database/sql"
	"time"

	entsql "entgo.io/ent/dialect/sql"

	"doublequote/ent"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

// SQL represents the database connection.
type SQL struct {
	client *ent.Client
	dbUrl  string

	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now func() time.Time
}

// NewSQL returns a new instance of SQL associated with the given datasource name.
func NewSQL(dbUrl string) *SQL {
	return &SQL{
		dbUrl: dbUrl,
		Now:   time.Now,
	}
}

// Open opens the database connection.
func (sql *SQL) Open() error {
	c, err := createClient(sql.dbUrl)
	if err != nil {
		return err
	}

	sql.client = c

	return nil
}

// Close closes the database connection.
func (sql *SQL) Close() error {
	return sql.client.Close()
}

// TODO add postgres support
func createClient(fileName string) (*ent.Client, error) {
	sqlDB, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	drv := entsql.OpenDB("sqlite3", sqlDB)
	return ent.NewClient(ent.Driver(drv)), nil
}

// Utility function. Receives a function and a pointer to an argument. If the argument is nil, an empty selector
// is returned. If it's not, the argument is passed to the function and the result is returned.
func ifPresent[V any, P ~func(*entsql.Selector)](f func(V) P, v *V) P {
	if v != nil {
		return f(*v)
	} else {
		return func(*entsql.Selector) {}
	}
}
package sql

import (
	"context"
	"os"
	"time"

	"doublequote/prisma"
	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

//go:generate go run github.com/prisma/prisma-client-go generate --schema=../schema.prisma

// SQL represents the database connection.
type SQL struct {
	prisma *prisma.PrismaClient
	ctx    context.Context
	cancel func()

	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now func() time.Time
}

// NewSQL returns a new instance of SQL associated with the given datasource name.
func NewSQL(dbUrl string) *SQL {
	os.Setenv("DATABASE_URL", dbUrl)

	db := &SQL{
		prisma: prisma.NewClient(),
		Now:    time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

// Open opens the database connection.
func (sql *SQL) Open() error {
	return sql.prisma.Prisma.Connect()
}

// Close closes the database connection.
func (sql *SQL) Close() error {
	// Cancel background context.
	sql.cancel()

	return sql.prisma.Prisma.Disconnect()
}

package sql

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenCloseDB(t *testing.T) {
	err := godotenv.Load("../.env.testing")
	require.Nil(t, err)

	t.Run("OK", func(t *testing.T) {
		db := NewSQL(os.Getenv("DATABASE_URL"))

		err := db.Open()
		assert.Nil(t, err)

		err = db.Close()
		assert.Nil(t, err)
	})
}

// NewTestDB returns a sql instance with a mock Prisma client, and other utilities for testing.
// Closing the connection is not needed.
// func NewTestDB() (*SQL, *prisma.PrismaClient, *prisma.Mock, func(t *testing.T)) {
// 	client, mock, ensure := prisma.NewMock()

// 	db := &SQL{
// 		prisma: client,
// 		Now:    time.Now,
// 	}
// 	db.ctx, db.cancel = context.WithCancel(context.Background())

// 	return db, client, mock, ensure
// }

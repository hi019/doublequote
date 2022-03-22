package sql

// func TestOpenCloseDB(t *testing.T) {
// 	err := godotenv.Load("../.env.testing")
// 	require.Nil(t, err)

// 	t.Run("OK", func(t *testing.T) {
// 		db := NewSQL(os.Getenv("DATABASE_URL"))

// 		err := db.open()
// 		assert.Nil(t, err)

// 		err = db.close()
// 		assert.Nil(t, err)
// 	})
// }

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

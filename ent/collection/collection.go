// Code generated by entc, DO NOT EDIT.

package collection

import (
	"time"
)

const (
	// Label holds the string label denoting the collection type in the database.
	Label = "collection"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeFeeds holds the string denoting the feeds edge name in mutations.
	EdgeFeeds = "feeds"
	// Table holds the table name of the collection in the database.
	Table = "collections"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "collections"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_collections"
	// FeedsTable is the table that holds the feeds relation/edge. The primary key declared below.
	FeedsTable = "collection_feeds"
	// FeedsInverseTable is the table name for the Feed entity.
	// It exists in this package in order to avoid circular dependency with the "feed" package.
	FeedsInverseTable = "feeds"
)

// Columns holds all SQL columns for collection fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "collections"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_collections",
}

var (
	// FeedsPrimaryKey and FeedsColumn2 are the table columns denoting the
	// primary key for the feeds relation (M2M).
	FeedsPrimaryKey = []string{"collection_id", "feed_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)
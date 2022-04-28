package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CollectionEntry holds the schema definition for the CollectionEntry entity.
type CollectionEntry struct {
	ent.Schema
}

// Fields of the CollectionEntry.
func (CollectionEntry) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_read").Default(false),

		field.Int("collection_id"),
		field.Int("entry_id"),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the CollectionEntry.
func (CollectionEntry) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("collection", Collection.Type).
			Ref("collection_entries").
			Field("collection_id").
			Required().
			Unique(),
		edge.From("entry", Entry.Type).
			Ref("collection_entries").
			Field("entry_id").
			Required().
			Unique(),
	}
}

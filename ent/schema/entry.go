package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Entry holds the schema definition for the Entry entity.
type Entry struct {
	ent.Schema
}

// Fields of the Entry.
func (Entry) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("url"),
		field.String("author"),
		field.String("content_key"),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Entry.
func (Entry) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("feed", Feed.Type).
			Ref("entries").
			Unique(),
	}
}

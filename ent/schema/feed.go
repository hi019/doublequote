package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Feed holds the schema definition for the Feed entity.
type Feed struct {
	ent.Schema
}

// Fields of the Feed.
func (Feed) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("rssURL"),
		field.String("domain"),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Feed.
func (Feed) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("collections", Collection.Type).
			Ref("feeds"),
		edge.To("entries", Entry.Type),
	}
}

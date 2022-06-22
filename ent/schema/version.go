package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Version holds the schema definition for the Version entity.
type Version struct {
	ent.Schema
}

// Fields of the Version.
func (Version) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Default("beta"),
		field.Time("created_at").Optional(),
		field.Time("updated_at").Optional(),
	}
}

// Edges of the Version.
func (Version) Edges() []ent.Edge {
	return nil
}

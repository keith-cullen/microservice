package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Thing holds the schema definition for the Thing entity.
type Thing struct {
	ent.Schema
}

// Fields of the Thing.
func (Thing) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
	}
}

// Edges of the Thing.
func (Thing) Edges() []ent.Edge {
	return nil
}

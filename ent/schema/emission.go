package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Emission holds the schema definition for the Emission entity.
type Emission struct {
	ent.Schema
}

// Fields of the Emission.
func (Emission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.Enum("emission_type").
			Values("GNRE"),
		field.UUID("client_id", uuid.New()),
		field.Text("message").Optional(),
		field.Enum("status").
			Values("PROCESSING", "FAILED", "FINISHED", "EXCEPTION"),
		field.UUID("user_id", uuid.New()),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()),
		field.Time("deleted_at").Optional(),
	}
}

// Edges of the Emission.
func (Emission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("gnre_emission", GnreEmission.Type).
			Unique(),
	}
}

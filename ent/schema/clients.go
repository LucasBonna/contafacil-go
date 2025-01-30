package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Clients holds the schema definition for the Clients entity.
type Clients struct {
	ent.Schema
}

// Fields of the Clients.
func (Clients) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("name").MaxLen(255),
		field.String("cnpj").MaxLen(20).Unique(),
		field.Enum("role").
			Values("ADMIN", "USER"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.Time("deleted_at").Optional(),
	}
}

// Edges of the Clients.
func (Clients) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}

func (Clients) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Table("clients"),
	}
}

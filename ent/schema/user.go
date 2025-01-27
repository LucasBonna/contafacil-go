package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.String("username"),
		field.String("password"),
		field.String("api_key").Unique(),
		field.Enum("role").Values("USER", "ADMIN"),
		field.UUID("client_id", uuid.New()),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()),
		field.Time("deleted_at").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("clients", Clients.Type).
			Ref("users").
			Unique().
			Required().
			Field("client_id"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username", "client_id").Unique(),
	}
}

package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// AccessLog holds the schema definition for the AccessLog entity.
type AccessLog struct {
	ent.Schema
}

// Fields of the AccessLog.
func (AccessLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).
			Unique(),
		field.String("ip"),
		field.String("method"),
		field.String("endpoint"),
		field.Text("request_body").Optional(),
		field.Text("request_headers").Optional(),
		field.Text("request_params").Optional(),
		field.Text("request_query").Optional(),
		field.Text("response_body").Optional(),
		field.Text("response_headers").Optional(),
		field.String("response_time").MaxLen(50).Optional(),
		field.Int("status_code").Optional(),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at").Default(time.Now()),
	}
}

// Edges of the AccessLog.
func (AccessLog) Edges() []ent.Edge {
	return nil
}

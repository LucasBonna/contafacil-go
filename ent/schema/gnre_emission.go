package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Gnre_emissions holds the schema definition for the Gnre_emissions entity.
type GnreEmission struct {
	ent.Schema
}

// Fields of the Gnre_emissions.
func (GnreEmission) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Unique().Immutable(),
		field.UUID("xml", uuid.New()),
		field.UUID("pdf", uuid.New()).Optional(),
		field.UUID("comprovante_pdf", uuid.New()).Optional(),
		field.Float("guia_amount").Positive(),
		field.String("numero_recibo").Optional(),
		field.String("chave_nota"),
		field.String("cod_barras_guia").Optional(),
		field.String("num_nota"),
		field.String("destinatario"),
		field.String("cpf_cnpj"),
	}
}

// Edges of the Gnre_emissions.
func (GnreEmission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("emission", Emission.Type).
			Ref("gnre_emission").
			Unique().
			Required(),
	}
}

func (GnreEmission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Table("gnre_emission"),
	}
}

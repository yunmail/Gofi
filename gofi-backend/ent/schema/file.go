package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			StructTag(`json:"name"`).
			NotEmpty(),
		field.Bool("isDirectory").
			StructTag(`json:"isDirectory"`).
			Default(false),
		field.Int("size").
			Default(0),
		field.String("extension").
			StructTag(`json:"extension"`).
			NotEmpty(),
		field.String("mime").
			StructTag(`json:"mime"`).
			NotEmpty(),
		field.String("path").
			StructTag(`json:"path"`).
			NotEmpty(),
		field.Int64("lastModified").
			StructTag(`json:"lastModified"`).
			Default(0),
		field.String("content").
			StructTag(`json:"content"`).
			NotEmpty(),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return nil
}

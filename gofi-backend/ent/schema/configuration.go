package schema

import (
	"github.com/facebookincubator/ent"
	"github.com/facebookincubator/ent/schema/field"
)

// Configuration holds the schema definition for the Configuration entity.
type Configuration struct {
	ent.Schema
}

// Fields of the Configuration.
func (Configuration) Fields() []ent.Field {
	return []ent.Field{
		field.String("version"). // 动态字段,无需持久化到数据库
			StructTag(`json:"version"`).
			Optional().
			Immutable(),
		field.String("appPath"). // 动态字段,无需持久化到数据库
			StructTag(`json:"appPath"`).
			Optional().
			Immutable(),
		field.String("defStoragePath"). // 动态字段,无需持久化到数据库,由于ent在defaultStoragePath会生成重复关键字"defaultStoragePath"，所以用defStoragePath代替
			StructTag(`json:"defaultStoragePath"`).
			Optional().
			Immutable(),
		field.String("themeStyle").
			StructTag(`json:"themeStyle"`).
			NotEmpty(),
		field.String("themeColor").
			StructTag(`json:"themeColor"`).
			NotEmpty(),
		field.String("navMode").
			StructTag(`json:"navMode"`).
			NotEmpty(),
		field.String("databaseFilePath").
			StructTag(`json:"databaseFilePath"`).
			NotEmpty(),
		field.String("customStoragePath").
			StructTag(`json:"customStoragePath"`).
			Default(""),
		field.String("logDirectoryPath").
			StructTag(`json:"logDirectoryPath"`).
			NotEmpty(),
		field.Bool("initialized").
			StructTag(`json:"initialized"`).
			Default(false),
		field.Time("created").
			Immutable().
			StructTag(`json:"created"`),
		field.Time("updated").
			StructTag(`json:"updated"`),
	}
}

// Edges of the Configuration.
func (Configuration) Edges() []ent.Edge {
	return nil
}

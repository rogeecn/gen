package gen

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"

	"go.ipao.vip/gen/field"
	"go.ipao.vip/gen/helper"
)

func DefaultConfig() Config {
	cfg := Config{
		OutPath:          "./database/query",
		Mode:             WithDefaultQuery,
		OutFile:          "query.gen.go",
		FieldSignable:    true,
		FieldWithTypeTag: true,
	}
	cfg.WithImportPkgPath("go.ipao.vip/gen/types")

	return cfg
}

type ConfigOptRelation struct {
	Relation string `yaml:"relation"`
	Table    string `yaml:"table"`
	Options  *struct {
		RelatePointer      bool          `yaml:"relate_pointer"`
		RelateSlice        bool          `yaml:"relate_slice"`
		RelateSlicePointer bool          `yaml:"relate_slice_pointer"`
		JSONTag            string        `yaml:"json_tag"`
		GORMTag            field.GormTag `yaml:"gorm_tag"`
		Tag                field.Tag     `yaml:"tag"`
		OverwriteTag       field.Tag     `yaml:"overwrite_tag"`
	} `yaml:"options"`
}

type ConfigOpt struct {
	Imports     []string                                `yaml:"imports"`
	FieldType   map[string]map[string]string            `yaml:"field_type"`
	FieldRelate map[string]map[string]ConfigOptRelation `yaml:"field_relate"`
}

func GenerateWithDefault(db *gorm.DB, transformConfigFile string) {
	g := NewGenerator(DefaultConfig())
	g.UseDB(db)

	g.WithTableNameStrategy(func(tableName string) string {
		if strings.HasPrefix(tableName, "_") {
			return ""
		}
		if tableName == "migrations" {
			return ""
		}
		return tableName
	})

	if transformConfigFile == "" {
		g.ApplyBasic(g.GenerateAllTable()...)
		g.Execute()
		return
	}

	conf, err := os.ReadFile(transformConfigFile)
	if err != nil {
		panic(fmt.Errorf("read transform config file %q fail: %w", transformConfigFile, err))
	}

	// yaml config parse to ConfigOpt
	var cfgOpt ConfigOpt
	if err := helper.UnmarshalYAML([]byte(conf), &cfgOpt); err != nil {
		panic(fmt.Errorf("parse yaml config fail: %w", err))
	}

	g.WithImportPkgPath(cfgOpt.Imports...)

	tables, err := db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("get all tables fail: %w", err))
	}

	mapTables := make(map[string][]ModelOpt)
	for _, table := range tables {
		opts := []ModelOpt{}
		if fieldTypes, ok := cfgOpt.FieldType[table]; ok {
			for f, typ := range fieldTypes {
				opts = append(opts, FieldType(f, typ))
			}
		}
		if fieldTypes, ok := cfgOpt.FieldRelate[table]; ok {
			for f, relation := range fieldTypes {
				r := field.RelationshipType(relation.Relation)

				switch r {
				case field.HasOne, field.BelongsTo, field.HasMany, field.Many2Many:
				default:
					panic("unsupported relationship type: " + relation.Relation)
				}

				opts = append(opts, FieldRelateModel(r, f, g.GenerateModel(relation.Table), &field.RelateConfig{
					RelatePointer:      relation.Options.RelatePointer,
					RelateSlice:        relation.Options.RelateSlice,
					RelateSlicePointer: relation.Options.RelateSlicePointer,
					JSONTag:            relation.Options.JSONTag,
					GORMTag:            relation.Options.GORMTag,
					Tag:                relation.Options.Tag,
					OverwriteTag:       relation.Options.OverwriteTag,
				}))
			}
		}
		mapTables[table] = opts
	}

	models := []interface{}{}
	for tbl, opts := range mapTables {
		models = append(models, g.GenerateModel(tbl, opts...))
	}
	g.ApplyBasic(models...)

	// Generate
	g.Execute()
}

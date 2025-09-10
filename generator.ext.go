package gen

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go.ipao.vip/gen/field"
	"go.ipao.vip/gen/helper"
)

func DefaultConfig() Config {
    cfg := Config{
		OutPath:          "./database",
		Mode:             WithDefaultQuery,
		OutFile:          "query.gen.go",
		FieldSignable:    true,
		FieldWithTypeTag: true,
    }
    // Co-locate models and queries in the same directory by default
    cfg.ModelPkgPath = cfg.OutPath
    cfg.WithImportPkgPath("go.ipao.vip/gen/types")

    return cfg
}

type ConfigOptRelation struct {
	Relation       string `yaml:"relation"`
	Table          string `yaml:"table"`
	Pivot          string `yaml:"pivot"`
	ForeignKey     string `yaml:"foreign_key"`
	References     string `yaml:"references"`
	JoinForeignKey string `yaml:"join_foreign_key"`
	JoinReferences string `yaml:"join_references"`
	Json           string `yaml:"json"`

	Options *struct {
		RelatePointer      bool `yaml:"relate_pointer"`
		RelateSlice        bool `yaml:"relate_slice"`
		RelateSlicePointer bool `yaml:"relate_slice_pointer"`
	} `yaml:"options"`
}

func (c *ConfigOptRelation) Config(db *gorm.DB) *field.RelateConfig {
	if c.Relation != string(field.Many2Many) && (len(c.ForeignKey) == 0 || len(c.References) == 0) {
		panic(fmt.Errorf("foreign_key and references must be set for relation %s", c.Relation))
	}
	opt := &field.RelateConfig{}

	var f func(string) string
	if ns, ok := db.NamingStrategy.(schema.NamingStrategy); ok {
		ns.SingularTable = true
		f = ns.SchemaName
	} else if db.NamingStrategy != nil {
		f = db.NamingStrategy.SchemaName
	} else {
		panic("no valid NamingStrategy")
	}

	opt.GORMTag = field.GormTag(map[string][]string{})
	if c.Relation == string(field.Many2Many) {
		opt.GORMTag["many2many"] = []string{c.Pivot}
	}

	if c.ForeignKey != "" {
		opt.GORMTag["foreignKey"] = []string{f(c.ForeignKey)}
	}

	if c.References != "" {
		opt.GORMTag["references"] = []string{f(c.References)}
	}

	if c.JoinForeignKey != "" {
		opt.GORMTag["joinForeignKey"] = []string{f(c.JoinForeignKey)}
	}

	if c.JoinReferences != "" {
		opt.GORMTag["joinReferences"] = []string{f(c.JoinReferences)}
	}

	if c.Json != "" {
		opt.JSONTag = c.Json
	}
	opt.RelatePointer = true

	return opt
}

type ConfigOpt struct {
	Ignores     []string                                `yaml:"ignores"`
	Imports     []string                                `yaml:"imports"`
	FieldType   map[string]map[string]string            `yaml:"field_type"`
	FieldRelate map[string]map[string]ConfigOptRelation `yaml:"field_relate"`
}

func GenerateWithDefault(db *gorm.DB, transformConfigFile string) {
	g := NewGenerator(DefaultConfig())
	g.UseDB(db)

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

	g.WithTableNameStrategy(func(tableName string) string {
		if strings.HasPrefix(tableName, "_") {
			return ""
		}

		// ignores table
		for _, ignore := range cfgOpt.Ignores {
			if strings.EqualFold(ignore, tableName) {
				return ""
			}
		}

		return tableName
	})

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

				opts = append(opts, FieldRelate(r, f, g.GenerateModel(relation.Table), relation.Config(db)))
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

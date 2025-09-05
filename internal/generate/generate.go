package generate

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"go.ipao.vip/gen/internal/model"
)

/*
** The feature of mapping table from database server to Golang struct
** Provided by @qqxhb
 */

func getFields(db *gorm.DB, conf *model.Config, columns []*model.Column) (fields []*model.Field) {
	for _, col := range columns {
		col.SetDataTypeMap(conf.DataTypeMap)
		col.WithNS(conf.FieldJSONTagNS)

		m := col.ToField(conf.FieldNullable, conf.FieldCoverable, conf.FieldSignable)

		// Prefer precise field wrapper type based on database type name when generating query fields.
		// This enables strongly-typed helpers like field.Money, field.Inet, field.JSONB, etc.
		if genType := fieldWrapperForDBType(col.DatabaseTypeName()); genType != "" {
			m.CustomGenType = genType
		}

		if filterField(m, conf.FilterOpts) == nil {
			continue
		}
		if _, ok := col.ColumnType.ColumnType(); ok &&
			!conf.FieldWithTypeTag { // remove type tag if FieldWithTypeTag == false
			m.GORMTag.Remove("type")
		}

		m = modifyField(m, conf.ModifyOpts)
		if ns, ok := db.NamingStrategy.(schema.NamingStrategy); ok {
			ns.SingularTable = true
			m.Name = ns.SchemaName(ns.TablePrefix + m.Name)
		} else if db.NamingStrategy != nil {
			m.Name = db.NamingStrategy.SchemaName(m.Name)
		}

		fields = append(fields, m)
	}
	for _, create := range conf.CreateOpts {
		m := create.Operator()(nil)
		if m.Relation != nil {
			if m.Relation.Model() != nil {
				stmt := gorm.Statement{DB: db}
				_ = stmt.Parse(m.Relation.Model())
				if stmt.Schema != nil {
					m.Relation.AppendChildRelation(ParseStructRelationShip(&stmt.Schema.Relationships)...)
				}
			}
			m.Type = strings.ReplaceAll(
				m.Type,
				conf.ModelPkg+".",
				"",
			) // remove modelPkg in field's Type, avoid import error
		}

		fields = append(fields, m)
	}
	return fields
}

// fieldWrapperForDBType maps PostgreSQL column types to field wrapper names.
func fieldWrapperForDBType(dbType string) string {
	switch strings.ToLower(dbType) {
	// JSON/XML/Money
	case "jsonb":
		return "JSONB"
	case "json":
		return "JSON"
	case "xml":
		return "XML"
	case "money":
		return "Money"

	// Network
	case "inet":
		return "Inet"
	case "cidr":
		return "CIDR"
	case "macaddr":
		return "MACAddr"

	// Geometry
	case "point":
		return "Point"
	case "box":
		return "Box"
	case "path":
		return "Path"
	case "polygon":
		return "Polygon"
	case "circle":
		return "Circle"

	// Bit string
	case "bit", "varbit":
		return "BitString"

	// Ranges
	case "int4range":
		return "Int4Range"
	case "int8range":
		return "Int8Range"
	case "numrange":
		return "NumRange"
	case "tsrange":
		return "TsRange"
	case "tstzrange":
		return "TstzRange"
	case "daterange":
		return "DateRange"

	// Full text
	case "tsvector":
		return "TSVector"
	case "tsquery":
		return "TSQuery"
	}
	return ""
}

func filterField(m *model.Field, opts []model.FieldOption) *model.Field {
	for _, opt := range opts {
		if opt.Operator()(m) == nil {
			return nil
		}
	}
	return m
}

func modifyField(m *model.Field, opts []model.FieldOption) *model.Field {
	for _, opt := range opts {
		m = opt.Operator()(m)
	}
	return m
}

// get mysql db' name
var modelNameReg = regexp.MustCompile(`^\w+$`)

func checkStructName(name string) error {
	if name == "" {
		return nil
	}
	if !modelNameReg.MatchString(name) {
		return fmt.Errorf("model name cannot contains invalid character")
	}
	if name[0] < 'A' || name[0] > 'Z' {
		return fmt.Errorf("model name must be initial capital")
	}
	return nil
}

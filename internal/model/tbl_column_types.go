package model

import (
	"gorm.io/gorm"
)

// defaultDataTypeMap holds framework default mappings from DB types to Go types for generated models.
// These mappings prioritize PostgreSQL types and align with field helpers in go.ipao.vip/gen/field.
var defaultDataTypeMap = map[string]func(columnType gorm.ColumnType) (dataType string){
	// primitives and bytes
	"bytea":   func(gorm.ColumnType) string { return "[]byte" },
	"boolean": func(gorm.ColumnType) string { return "bool" },

	// time/date/timestamp (including verbose names)
	"date":      func(gorm.ColumnType) string { return "types.Date" },
	"time":      func(gorm.ColumnType) string { return "types.Time" },
	"timestamp": func(gorm.ColumnType) string { return "time.Time" },

	// JSON/XML/MONEY
	"json":  func(gorm.ColumnType) string { return "types.JSON" },
	"jsonb": func(gorm.ColumnType) string { return "types.JSON" },
	"xml":   func(gorm.ColumnType) string { return "types.XML" },
	"money": func(gorm.ColumnType) string { return "types.Money" },

	// UUID
	"uuid": func(gorm.ColumnType) string { return "types.UUID" },

	// Network
	"inet":    func(gorm.ColumnType) string { return "types.Inet" },
	"cidr":    func(gorm.ColumnType) string { return "types.CIDR" },
	"macaddr": func(gorm.ColumnType) string { return "types.MACAddr" },

	// Geometry
	"point":   func(gorm.ColumnType) string { return "types.Point" },
	"box":     func(gorm.ColumnType) string { return "types.Box" },
	"path":    func(gorm.ColumnType) string { return "types.Path" },
	"polygon": func(gorm.ColumnType) string { return "types.Polygon" },
	"circle":  func(gorm.ColumnType) string { return "types.Circle" },

	// Bit string
	"bit":    func(gorm.ColumnType) string { return "types.BitString" },
	"varbit": func(gorm.ColumnType) string { return "types.BitString" },

	// Arrays (underscore element type names)
	"_text":    func(gorm.ColumnType) string { return "[]string" },
	"_varchar": func(gorm.ColumnType) string { return "[]string" },
	"_bpchar":  func(gorm.ColumnType) string { return "[]string" },
	"_char":    func(gorm.ColumnType) string { return "[]string" },
	"_int2":    func(gorm.ColumnType) string { return "[]int32" },
	"_int4":    func(gorm.ColumnType) string { return "[]int32" },
	"_int8":    func(gorm.ColumnType) string { return "[]int64" },
	"_float4":  func(gorm.ColumnType) string { return "[]float32" },
	"_float8":  func(gorm.ColumnType) string { return "[]float64" },
	"_bool":    func(gorm.ColumnType) string { return "[]bool" },
	"_uuid":    func(gorm.ColumnType) string { return "[]string" },
	"_numeric": func(gorm.ColumnType) string { return "[]float64" },

	// Arrays via generic "array" with detail type info, e.g., TEXT[] / INTEGER[]
	"text[]":             func(gorm.ColumnType) string { return "[]string" },
	"varchar[]":          func(gorm.ColumnType) string { return "[]string" },
	"char[]":             func(gorm.ColumnType) string { return "[]string" },
	"int2[]":             func(gorm.ColumnType) string { return "[]int32" },
	"int4[]":             func(gorm.ColumnType) string { return "[]int32" },
	"integer[]":          func(gorm.ColumnType) string { return "[]int32" },
	"int8[]":             func(gorm.ColumnType) string { return "[]int64" },
	"bigint[]":           func(gorm.ColumnType) string { return "[]int64" },
	"float4[]":           func(gorm.ColumnType) string { return "[]float32" },
	"real[]":             func(gorm.ColumnType) string { return "[]float32" },
	"float8[]":           func(gorm.ColumnType) string { return "[]float64" },
	"double precision[]": func(gorm.ColumnType) string { return "[]float64" },
	"boolean[]":          func(gorm.ColumnType) string { return "[]bool" },
	"uuid[]":             func(gorm.ColumnType) string { return "[]string" },
	"numeric[]":          func(gorm.ColumnType) string { return "[]float64" },

	// Ranges
	"int4range": func(gorm.ColumnType) string { return "types.Int4Range" },
	"int8range": func(gorm.ColumnType) string { return "types.Int8Range" },
	"numrange":  func(gorm.ColumnType) string { return "types.NumRange" },
	"tsrange":   func(gorm.ColumnType) string { return "types.TsRange" },
	"tstzrange": func(gorm.ColumnType) string { return "types.TstzRange" },
	"daterange": func(gorm.ColumnType) string { return "types.DateRange" },

	// Full text search
	"tsvector": func(gorm.ColumnType) string { return "types.TSVector" },
	"tsquery":  func(gorm.ColumnType) string { return "types.TSQuery" },
}

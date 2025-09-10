package template

// Model used as a variable because it cannot load template file after packed, params still can pass file
const Model = NotEditMark + `
package {{.StructInfo.Package}}

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"go.ipao.vip/gen"
	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

{{if .TableName -}}const TableName{{.ModelStructName}} = "{{.TableName}}"{{- end}}

// {{.ModelStructName}} {{.StructComment}}
type {{.ModelStructName}} struct {
    {{range .Fields}}
    {{if .MultilineComment -}}
	/*
{{.ColumnComment}}
    */
	{{end -}}
    {{.Name}} {{.Type}} ` + "`{{.Tags}}` " +
	"{{if not .MultilineComment}}{{if .ColumnComment}}// {{.ColumnComment}}{{end}}{{end}}" +
	`{{end}}
}

// Quick operations without importing query package
// Update applies changed fields to the database using the default DB.
func (m *{{.ModelStructName}}) Update(ctx context.Context) (gen.ResultInfo, error) { return Q.{{.ModelStructName}}.WithContext(ctx).Updates(m) }

// Save upserts the model using the default DB.
func (m *{{.ModelStructName}}) Save(ctx context.Context) error { return Q.{{.ModelStructName}}.WithContext(ctx).Save(m) }

// Create inserts the model using the default DB.
func (m *{{.ModelStructName}}) Create(ctx context.Context) error { return Q.{{.ModelStructName}}.WithContext(ctx).Create(m) }

// Delete removes the row represented by the model using the default DB.
func (m *{{.ModelStructName}}) Delete(ctx context.Context) (gen.ResultInfo, error) { return Q.{{.ModelStructName}}.WithContext(ctx).Delete(m) }

`

// ModelMethod model struct DIY method
const ModelMethod = `

{{if .Doc -}}// {{.DocComment -}}{{end}}
func ({{.GetBaseStructTmpl}}){{.MethodName}}({{.GetParamInTmpl}})({{.GetResultParamInTmpl}}){{.Body}}
`

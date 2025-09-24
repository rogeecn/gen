package parser

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
)

// Param parameters in method
type Param struct { // (user model.User)
	PkgPath   string // package's path: internal/model
	Package   string // package's name: model
	Name      string // param's name: user
	Type      string // param's type: User
	IsArray   bool   // is array or not
	IsPointer bool   // is pointer or not
}

// Eq if param equal to another
func (p *Param) Eq(q Param) bool {
	return p.Package == q.Package && p.Type == q.Type
}

// IsError ...
func (p *Param) IsError() bool {
	return p.Type == "error"
}

// IsGenM ...
func (p *Param) IsGenM() bool {
	return p.Package == "gen" && p.Type == "M"
}

// IsGenRowsAffected ...
func (p *Param) IsGenRowsAffected() bool {
	return p.Package == "gen" && p.Type == "RowsAffected"
}

// IsMap ...
func (p *Param) IsMap() bool {
	return strings.HasPrefix(p.Type, "map[")
}

// IsGenT ...
func (p *Param) IsGenT() bool {
	return p.Package == "gen" && p.Type == "T"
}

// IsInterface ...
func (p *Param) IsInterface() bool {
	return p.Type == "interface{}"
}

// IsNull ...
func (p *Param) IsNull() bool {
	return p.Package == "" && p.Type == "" && p.Name == ""
}

// InMainPkg ...
func (p *Param) InMainPkg() bool {
	return p.Package == "main"
}

// IsTime ...
func (p *Param) IsTime() bool {
	return p.Package == "time" && p.Type == "Time"
}

// IsSQLResult ...
func (p *Param) IsSQLResult() bool {
	return (p.Package == "sql" && p.Type == "Result") || (p.Package == "gen" && p.Type == "SQLResult")
}

// IsSQLRow ...
func (p *Param) IsSQLRow() bool {
	return (p.Package == "sql" && p.Type == "Row") || (p.Package == "gen" && p.Type == "SQLRow")
}

// IsSQLRows ...
func (p *Param) IsSQLRows() bool {
	return (p.Package == "sql" && p.Type == "Rows") || (p.Package == "gen" && p.Type == "SQLRows")
}

// SetName ...
func (p *Param) SetName(name string) {
	p.Name = name
}

// TypeName ...
func (p *Param) TypeName() string {
	if p.IsArray {
		return "[]" + p.Type
	}
	return p.Type
}

// TmplString param to string in tmpl
func (p *Param) TmplString() string {
	var res strings.Builder
	if p.Name != "" {
		res.WriteString(p.Name)
		res.WriteString(" ")
	}

	if p.IsArray {
		res.WriteString("[]")
	}
	if p.IsPointer {
		res.WriteString("*")
	}
	if p.Package != "" {
		res.WriteString(p.Package)
		res.WriteString(".")
	}
	res.WriteString(p.Type)
	return res.String()
}

// IsBaseType judge whether the param type is basic type
func (p *Param) IsBaseType() bool {
	switch p.Type {
	case "string", "byte":
		return true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	case "float64", "float32":
		return true
	case "bool":
		return true
	case "time.Time":
		return true
	default:
		return false
	}
}

func (p *Param) astGetParamType(param *ast.Field) {
	switch v := param.Type.(type) {
	case *ast.Ident:
		p.Type = v.Name
		if v.Obj != nil {
			p.Package = "UNDEFINED" // set a placeholder
		}
	case *ast.SelectorExpr:
		p.astGetEltType(v)
	case *ast.ArrayType:
		p.astGetEltType(v.Elt)
		p.IsArray = true
	case *ast.Ellipsis:
		p.astGetEltType(v.Elt)
		p.IsArray = true
	case *ast.MapType:
		p.astGetMapType(v)
	case *ast.InterfaceType:
		p.Type = "interface{}"
	case *ast.StarExpr:
		p.IsPointer = true
		p.astGetEltType(v.X)
	case *ast.IndexExpr:
		p.astGetEltType(v.X)
	case *ast.IndexListExpr:
		p.astGetEltType(v.X)
	default:
		log.Printf("Unsupported param type: %+v", v)
	}
}

func (p *Param) astGetEltType(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.Type = v.Name
		if v.Obj != nil {
			p.Package = "UNDEFINED"
		}
	case *ast.SelectorExpr:
		p.Type = v.Sel.Name
		p.astGetPackageName(v.X)
	case *ast.MapType:
		p.astGetMapType(v)
	case *ast.StarExpr:
		p.IsPointer = true
		p.astGetEltType(v.X)
	case *ast.InterfaceType:
		p.Type = "interface{}"
	case *ast.ArrayType:
		p.astGetEltType(v.Elt)
		p.Type = "[]" + p.Type
	case *ast.IndexExpr:
		p.astGetEltType(v.X)
	default:
		log.Printf("Unsupported param type: %+v", v)
	}
}

func (p *Param) astGetPackageName(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.Package = v.Name
	}
}

func (p *Param) astGetMapType(expr *ast.MapType) {
	p.Type = fmt.Sprintf("map[%s]%s", astGetType(expr.Key), astGetType(expr.Value))
}

func astGetType(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.InterfaceType:
		return "interface{}"
	}
	return ""
}

// Method Apply to query struct and base struct custom method
type Method struct {
	Receiver   Param
	MethodName string
	Doc        string
	Params     []Param
	Result     []Param
	Body       string
}

// FuncSign function signature
func (m Method) FuncSign() string {
	return fmt.Sprintf("%s(%s) (%s)", m.MethodName, m.GetParamInTmpl(), m.GetResultParamInTmpl())
}

// GetBaseStructTmpl return method bind info string
func (m *Method) GetBaseStructTmpl() string {
	return m.Receiver.TmplString()
}

// GetParamInTmpl return param list
func (m *Method) GetParamInTmpl() string {
	return paramToString(m.Params)
}

// GetResultParamInTmpl return result list
func (m *Method) GetResultParamInTmpl() string {
	return paramToString(m.Result)
}

// paramToString param list to string used in tmpl
func paramToString(params []Param) string {
	res := make([]string, len(params))
	for i, param := range params {
		res[i] = param.TmplString()
	}
	return strings.Join(res, ",")
}

// DocComment return comment sql add "//" every line
func (m *Method) DocComment() string {
	return strings.Replace(strings.TrimSpace(m.Doc), "\n", "\n//", -1)
}

func DefaultMethodTableName(structName string) *Method {
	return &Method{
		Receiver:   Param{IsPointer: true, Type: structName},
		MethodName: "TableName",
		Doc:        fmt.Sprint("TableName ", structName, "'s table name "),
		Result:     []Param{{Type: "string"}},
		Body:       fmt.Sprintf("{\n\treturn TableName%s\n} ", structName),
	}
}

func getParamList(expr *ast.FieldList) []Param {
	if expr == nil {
		return nil
	}
	params := make([]Param, 0, expr.NumFields())
	for _, param := range expr.List {
		p := Param{}
		p.astGetParamType(param)
		if len(param.Names) == 0 {
			params = append(params, p)
			continue
		}
		for _, name := range param.Names {
			newParam := p
			newParam.Name = name.Name
			params = append(params, newParam)
		}
	}
	return params
}

func fixParamPackagePath(imports map[string]string, params []Param) {
	for i := range params {
		if params[i].Package == "UNDEFINED" {
			if importPath, ok := imports[params[i].Type]; ok {
				params[i].Package = ""
				params[i].PkgPath = strings.Trim(importPath, `"`)
			}
		}
	}
}
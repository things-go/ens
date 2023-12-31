package codegen

import (
	"fmt"
	"strings"

	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

// Deprecated: Don't use this, deprecated next major version.
func (g *CodeGen) GenAssist(modelImportPath string) *CodeGen {
	pkgQualifierPrefix := ""
	if p := ens.PkgName(modelImportPath); p != "" {
		pkgQualifierPrefix = p + "."
	}
	if !g.disableDocComment {
		g.Printf("// Code generated by %s. DO NOT EDIT.\n", g.byName)
		g.Printf("// version: %s\n", g.version)
		g.Println()
	}
	g.Printf("package %s\n", g.packageName)
	g.Println()

	//* import
	g.Println("import (")
	if pkgQualifierPrefix != "" {
		g.Printf("\"%s\"\n", modelImportPath)
		g.Println()
	}
	g.Println(`assist "github.com/things-go/gorm-assist"`)
	g.Println(`"gorm.io/gorm"`)
	g.Println(")")

	constField := func(structName, fieldName string) string {
		return fmt.Sprintf(`xx_%s_%s`, structName, fieldName)
	}
	//* struct
	for _, et := range g.entities {
		structName := utils.CamelCase(et.Name)
		tableName := et.Name

		constTableName := fmt.Sprintf("xx_%s_TableName", structName)
		{ //* const field
			g.Println("const (")
			g.Printf("// hold model `%s` table name\n", structName)
			g.Printf("%s = \"%s\"\n", constTableName, tableName)
			g.Printf("// hold model `%s` column name\n", structName)
			for _, field := range et.Fields {
				g.Printf("%s = \"%s\"\n", constField(structName, utils.CamelCase(field.Name)), field.Name)
			}
			g.Println(")")
			g.Println()
		}

		varModel := fmt.Sprintf(`xxx_%s_Model`, structName)
		funcInnerNew := fmt.Sprintf(`new_%s`, structName)
		{ //* var field
			g.Printf("var %s = %s (%s)\n", varModel, funcInnerNew, constTableName)
			g.Println()
		}

		typeNative := fmt.Sprintf("%s_Native", structName)
		//* type
		{
			g.Printf("type %s struct {\n", typeNative)
			g.Println("xAlias string")
			g.Println("ALL assist.Asterisk")
			for _, field := range et.Fields {
				g.Printf("%s assist.%s\n", utils.CamelCase(field.Name), field.RapierDataType)
			}
			g.Println("}")
			g.Println()
		}
		//* function X_xxx
		{
			g.Printf("// X_%s model with TableName `%s`.\n", structName, tableName)
			g.Printf("func X_%s() %s {\n", structName, typeNative)
			g.Printf("return %s\n", varModel)
			g.Println("}")
			g.Println()
		}
		//* function new_xxx
		{
			g.Printf("func %s(xAlias string) %s {\n", funcInnerNew, typeNative)
			g.Printf("return %s {\n", typeNative)
			g.Println("xAlias: xAlias,")
			g.Println("ALL:  assist.NewAsterisk(xAlias),")
			for _, field := range et.Fields {
				fieldName := utils.CamelCase(field.Name)
				g.Printf("%s: assist.New%s(xAlias, %s),\n", fieldName, field.RapierDataType, constField(structName, fieldName))
			}
			g.Println("}")
			g.Println("}")
			g.Println()
		}
		//* function New_xxxx
		{
			g.Printf("// New_%s new instance.\n", structName)
			g.Printf("func New_%s(xAlias string) %s {\n", structName, typeNative)
			g.Printf("if xAlias == %s {\n", constTableName)
			g.Printf("return %s\n", varModel)
			g.Println("} else {")
			g.Printf("return %s(xAlias)\n", funcInnerNew)
			g.Println("}")
			g.Println("}")
			g.Println()
		}
		//* method As
		{
			g.Println("// As alias")
			g.Printf("func (*%[1]s) As(alias string) %[1]s {\n", typeNative)
			g.Printf("return New_%s(alias)\n", structName)
			g.Println("}")
			g.Println()
		}
		//* method X_Alias
		{
			g.Printf("// X_Alias hold table name when call New_%[1]s or %[1]s_Native.As that you defined.\n", structName)
			g.Printf("func (x *%s) X_Alias() string {\n", typeNative)
			g.Println("return x.xAlias")
			g.Println("}")
			g.Println()
		}
		// impl TableName interface
		{
			//* method TableName
			g.Printf("// TableName hold model `%s` table name returns `%s`.\n", structName, tableName)
			g.Printf("func (*%s) TableName() string {\n", typeNative)
			g.Printf("return %s\n", constTableName)
			g.Println("}")
			g.Println()
		}

		//* method New_Executor
		{
			modelName := pkgQualifierPrefix + structName
			g.Println("// New_Executor new entity executor which suggest use only once.")
			g.Printf("func (*%s) New_Executor(db *gorm.DB) *assist.Executor[%s] {\n", typeNative, modelName)
			g.Printf("return assist.NewExecutor[%s](db)\n", modelName)
			g.Println("}")
			g.Println()
		}
		//* method Select_Expr
		{
			g.Println("// Select_Expr select model fields")
			g.Printf("func (x *%s) Select_Expr() []assist.Expr {\n", typeNative)
			g.Println("return []assist.Expr{")
			for _, field := range et.Fields {
				g.Printf("x.%s,\n", utils.CamelCase(field.Name))
			}
			g.Println("}")
			g.Println("}")
			g.Println()
		}

		//* method Select_VariantExpr
		{
			g.Println("// Select_VariantExpr select model fields, but time.Time field convert to timestamp(int64).")
			g.Printf("func (x *%s) Select_VariantExpr(prefixes ...string) []assist.Expr {\n", typeNative)
			g.Println("if len(prefixes) > 0 && prefixes[0] != \"\" {")
			g.Println("return []assist.Expr{")
			for _, field := range et.Fields {
				g.Println(genAssist_SelectVariantExprField(structName, field, true))
			}
			g.Println("}")
			g.Println("} else {")
			g.Println("return []assist.Expr{")
			for _, field := range et.Fields {
				g.Println(genAssist_SelectVariantExprField(structName, field, false))
			}
			g.Println("}")
			g.Println("}")
			g.Println("}")
			g.Println()
		}
	}
	return g
}

func genAssist_SelectVariantExprField(structName string, field *ens.FieldDescriptor, hasPrefix bool) string {
	fieldName := utils.CamelCase(field.Name)

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString("x.")
	b.WriteString(fieldName)
	if field.Type.IsTime() {
		b.WriteString(".UnixTimestamp()")
		if field.Nullable {
			b.WriteString(".IfNull(0)")
		}
		if !hasPrefix {
			fmt.Fprintf(b, ".As(x.%s.ColumnName())", fieldName)
		}
	}
	if hasPrefix {
		fmt.Fprintf(b, ".As(x.%s.FieldName(prefixes...))", fieldName)
	}
	b.WriteString(",")
	return b.String()
}

package mysql

import (
	"fmt"
	"strings"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"

	"github.com/things-go/ens"
	"github.com/things-go/ens/internal/insql"
	"github.com/things-go/ens/rapier"
	"github.com/things-go/ens/sqlx"
	"github.com/things-go/ens/utils"
)

func intoColumnSql(col *schema.Column) string {
	nullable := col.Type.Null
	autoIncrement := autoIncrement(col.Attrs)

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString(col.Type.Raw)
	if !nullable {
		b.WriteString(" NOT NULL")
	}
	if autoIncrement {
		b.WriteString(" AUTO_INCREMENT")
	} else {
		dv, ok := insql.DefaultValue(col)
		if ok {
			fmt.Fprintf(b, " DEFAULT '%s'", strings.Trim(dv, `"`))
		} else if nullable {
			b.WriteString(" DEFAULT NULL")
		}
	}
	return b.String()
}

func intoIndexSql(index *schema.Index) string {
	fields := insql.IndexPartColumnNames(index.Parts)
	indexType := findIndexType(index.Attrs)
	fieldList := "`" + strings.Join(fields, "`,`") + "`"
	if insql.IndexEqual(index.Table.PrimaryKey, index) {
		return fmt.Sprintf("PRIMARY KEY (%s) USING %s", fieldList, indexType)
	} else if index.Unique {
		return fmt.Sprintf("UNIQUE KEY `%s` (%s) USING %s", index.Name, fieldList, indexType)
	} else {
		return fmt.Sprintf("KEY `%s` (%s) USING %s", index.Name, fieldList, indexType)
	}
}

func intoForeignKeySql(fk *schema.ForeignKey) string {
	columnNameList := "`" + strings.Join(insql.ColumnNames(fk.Columns), "`,`") + "`"
	refColumnNameList := "`" + strings.Join(insql.ColumnNames(fk.RefColumns), "`,`") + "`"
	return fmt.Sprintf(
		"CONSTRAINT `%s` FOREIGN KEY (%s) REFERENCES `%s` (%s) ON DELETE %s ON UPDATE %s",
		fk.Symbol, columnNameList, fk.RefTable.Name, refColumnNameList, fk.OnDelete, fk.OnUpdate,
	)
}

func intoTableSql(tb *schema.Table) string {
	b := &strings.Builder{}
	b.Grow(64)
	fmt.Fprintf(b, "CREATE TABLE `%s` (\n", tb.Name)

	remain := len(tb.Columns) + len(tb.Indexes) + len(tb.ForeignKeys)
	if tb.PrimaryKey != nil {
		remain++
	}
	suffixOrEmpty := func(r int) string {
		if r == 0 {
			return ""
		}
		return ","
	}
	//* columns
	for _, col := range tb.Columns {
		remain--
		suffix := suffixOrEmpty(remain)
		comment, ok := insql.Comment(col.Attrs)
		if ok {
			comment = fmt.Sprintf(" COMMENT '%s'", comment)
		}
		fmt.Fprintf(b, "  `%s` %s%s%s\n", col.Name, intoColumnSql(col), comment, suffix)
	}
	//* pk + indexes
	if tb.PrimaryKey != nil {
		remain--
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", intoIndexSql(tb.PrimaryKey), suffix)
	}
	for _, index := range tb.Indexes {
		remain--
		if insql.IndexEqual(tb.PrimaryKey, index) { // ignore primary key, maybe include
			continue
		}
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", intoIndexSql(index), suffix)
	}
	//* foreignKeys
	for _, fk := range tb.ForeignKeys {
		remain--
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", intoForeignKeySql(fk), suffix)
	}

	engine := mysql.EngineInnoDB
	charset := "utf8mb4"
	collate := ""
	comment := ""
	for _, attr := range tb.Attrs {
		switch val := attr.(type) {
		case *mysql.Engine:
			engine = val.V
		case *schema.Charset:
			charset = val.V
		case *schema.Collation:
			collate = val.V
		case *schema.Comment:
			comment = val.Text
			// case *mysql.AutoIncrement: // ignore this
		}
	}
	fmt.Fprintf(b, ") ENGINE=%s DEFAULT CHARSET=%s", engine, charset)
	if collate != "" {
		fmt.Fprintf(b, " COLLATE='%s'", collate)
	}
	if comment != "" {
		fmt.Fprintf(b, " COMMENT='%s'", comment)
	}
	return b.String()
}

// column, type, not null, authIncrement, default, [primaryKey|index], comment
func intoGormTag(tb *schema.Table, col *schema.Column) string {
	pkPriority, isPk := 0, false
	if pk := tb.PrimaryKey; pk != nil {
		pkPriority, isPk = insql.FindIndexPartSeq(pk.Parts, col)
	}
	autoIncrement := autoIncrement(col.Attrs)

	b := &strings.Builder{}
	b.Grow(64)
	fmt.Fprintf(b, `gorm:"column:%s`, col.Name)
	if !(isPk && autoIncrement) {
		fmt.Fprintf(b, ";type:%s", col.Type.Raw)
	}
	if !col.Type.Null {
		fmt.Fprintf(b, ";not null")
	}

	if isPk {
		if autoIncrement {
			fmt.Fprintf(b, ";autoIncrement:true")
		}
	} else {
		dv, ok := insql.DefaultValue(col)
		if ok {
			if dv == `""` || dv == "" {
				dv = "''"
			} else {
				dv = strings.Trim(dv, `"`) // format: `"xxx"` or `'xxx'`
			}
			fmt.Fprintf(b, ";default:%s", dv)
		} else if col.Type.Null {
			fmt.Fprintf(b, ";default:null")
		}
	}

	//* pk + indexes
	if isPk && tb.PrimaryKey != nil {
		fmt.Fprintf(b, ";primaryKey")
		if len(tb.PrimaryKey.Parts) > 1 {
			fmt.Fprintf(b, ",priority:%d", pkPriority)
		}
	}
	for _, val := range col.Indexes {
		if insql.IndexEqual(tb.PrimaryKey, val) { // ignore primary key, may be include
			continue
		}
		if val.Unique {
			fmt.Fprintf(b, ";uniqueIndex:%s", val.Name)
		} else {
			fmt.Fprintf(b, ";index:%s", val.Name)
			// 	mysql.IndexTypeFullText
			// if v.IndexType == "FULLTEXT" {
			// 	b.WriteString(",class:FULLTEXT")
			// }
		}
		if len(val.Parts) > 1 {
			priority, ok := insql.FindIndexPartSeq(val.Parts, col)
			if ok {
				fmt.Fprintf(b, ",priority:%d", priority)
			}
		}
	}
	if comment, ok := insql.Comment(col.Attrs); ok && comment != "" {
		fmt.Fprintf(b, ";comment:%s", utils.TrimFieldComment(comment))
	}
	b.WriteString(`"`)
	return b.String()
}

func intoRapier(tb *schema.Table) *rapier.Struct {
	// * columns
	fields := make([]*rapier.StructField, 0, len(tb.Columns))
	for _, col := range tb.Columns {
		goType := intoGoType(col.Type.Raw)

		t := goType.Type.IntoRapierType()

		fields = append(fields, &rapier.StructField{
			Type:       t,
			GoName:     utils.CamelCase(col.Name),
			Nullable:   col.Type.Null,
			ColumnName: col.Name,
			Comment:    insql.MustComment(col.Attrs),
		})
	}
	return &rapier.Struct{
		GoName:    utils.CamelCase(tb.Name),
		TableName: tb.Name,
		Comment:   insql.MustComment(tb.Attrs),
		Fields:    fields,
	}
}

func intoSql(tb *schema.Table) *sqlx.Table {
	return &sqlx.Table{
		Name:    tb.Name,
		Sql:     intoTableSql(tb),
		Comment: insql.MustComment(tb.Attrs),
	}
}

func intoSchema(tb *schema.Table) *ens.EntityDescriptor {
	// * columns
	fielders := make([]*ens.FieldDescriptor, 0, len(tb.Columns))
	for _, col := range tb.Columns {
		fielders = append(fielders, &ens.FieldDescriptor{
			ColumnName: col.Name,
			Comment:    insql.MustComment(col.Attrs),
			Nullable:   col.Type.Null,
			Column:     NewColumnDef(col),
			Type:       intoGoType(col.Type.Raw),
			GoName:     utils.CamelCase(col.Name),
			GoPointer:  col.Type.Null,
			Tags:       []string{intoGormTag(tb, col)},
		})
	}
	// * indexes
	indexers := make([]*ens.IndexDescriptor, 0, len(tb.Indexes))
	for _, index := range tb.Indexes {
		indexers = append(indexers, &ens.IndexDescriptor{
			Name:   index.Name,
			Fields: insql.IndexPartColumnNames(index.Parts),
			Index:  NewIndexDef(index),
		})
	}
	//* foreignKeys
	fks := make([]*ens.ForeignKeyDescriptor, 0, len(tb.ForeignKeys))
	for _, fk := range tb.ForeignKeys {
		fks = append(fks, &ens.ForeignKeyDescriptor{
			Symbol:     fk.Symbol,
			Table:      fk.Table.Name,
			Columns:    insql.ColumnNames(fk.Columns),
			RefTable:   fk.RefTable.Name,
			RefColumns: insql.ColumnNames(fk.RefColumns),
			OnUpdate:   fk.OnUpdate,
			OnDelete:   fk.OnDelete,
			ForeignKey: NewForeignKey(fk),
		})
	}

	// * table
	return &ens.EntityDescriptor{
		Name:        tb.Name,
		Comment:     insql.MustComment(tb.Attrs),
		Table:       NewTableDef(tb),
		Fields:      fielders,
		Indexes:     indexers,
		ForeignKeys: fks,
	}
}

package proto

import "google.golang.org/protobuf/reflect/protoreflect"

type MessageField struct {
	Cardinality protoreflect.Cardinality // optional, required, or repeated
	Type        protoreflect.Kind        // 类型
	TypeName    string                   // 类型名称, 如果Type=protoreflect.MessageKind时
	Name        string                   // 名称, snake or small camel case
	ColumnName  string                   // 列名, snake case
	Comment     string                   // 注释
}

type Message struct {
	Name      string          // 名称, camel case
	TableName string          // 表名, snake name
	Comment   string          // 注释
	Fields    []*MessageField // 字段
}

type Schema struct {
	Name     string
	Entities []*Message
}

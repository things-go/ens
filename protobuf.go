package ens

import "fmt"

type ProtoMessage struct {
	DataType    string   // 数据类型
	Name        string   // 名称
	Comment     string   // 注释
	Annotations []string // 注解
}

func buildProtoMessage(field *FieldDescriptor, enableGogo, enableSea bool) *ProtoMessage {
	dataType := field.Type.Type.IntoProtoDataType()
	annotations := make([]string, 0, 16)
	if field.Type.Type == TypeInt64 ||
		field.Type.Type == TypeUint64 {
		annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
	} else if field.Type.IsTime() {
		if enableGogo {
			annotations = append(annotations, `(gogoproto.stdtime) = true`, `(gogoproto.nullable) = false`)
		} else {
			dataType = "int64"
			annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
		}
	}
	comment := field.Comment
	if comment != "" {
		comment = "// " + comment
	}
	if field.Column != nil && enableSea {
		l := fmt.Sprintf(`// #[seaql(type="%s")]`, field.Column.Definition())
		if comment != "" {
			comment = comment + "\n" + l
		} else {
			comment = l
		}
	}
	return &ProtoMessage{
		DataType:    dataType,
		Name:        field.Name,
		Comment:     comment,
		Annotations: annotations,
	}
}

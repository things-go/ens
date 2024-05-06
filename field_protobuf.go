package ens

import (
	"strings"

	"github.com/things-go/ens/matcher"
)

type ProtoMessage struct {
	DataType    string   // 数据类型
	Name        string   // 名称
	Optional    bool     // 是否可选
	Comment     string   // 注释
	Annotations []string // 注解
}

func (field *FieldDescriptor) buildProtoMessage() *ProtoMessage {
	dataType := field.Type.Type.IntoProtoDataType()
	annotations := make([]string, 0, 16)
	if field.Type.Type == TypeInt64 ||
		field.Type.Type == TypeUint64 {
		annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
	} else if field.Type.IsTime() {
		dataType = "int64"
		annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
	}

	comment := strings.TrimSuffix(matcher.TrimEnumValue(field.Comment), ",")
	if comment != "" {
		comment = "// " + comment
	}

	return &ProtoMessage{
		DataType:    dataType,
		Name:        field.Name,
		Optional:    field.Optional,
		Comment:     comment,
		Annotations: annotations,
	}
}

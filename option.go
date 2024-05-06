package ens

type Option struct {
	EnableInt          bool              `yaml:"enableInt" json:"enableInt"`                   // 使能int8,uint8,int16,uint16,int32,uint32输出为int,uint
	EnableIntegerInt   bool              `yaml:"enableIntegerInt" json:"enableIntegerInt"`     // 使能int32,uint32输出为int,uint
	EnableBoolInt      bool              `yaml:"enableBoolInt" json:"enableBoolInt"`           // 使能bool输出int
	DisableNullToPoint bool              `yaml:"disableNullToPoint" json:"disableNullToPoint"` // 禁用字段为null时输出指针类型,将输出为sql.Nullxx
	EnableForeignKey   bool              `yaml:"enableForeignKey" json:"enableForeignKey"`     // 输出外键
	Tags               map[string]string `yaml:"tags" json:"tags"`                             // tags标签列表, support smallCamelCase, camelCase, snakeCase, kebab
}

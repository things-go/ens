package command

import (
	"github.com/spf13/pflag"
	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

type Config struct {
	ens.Option
	DisableCommentTag bool   `yaml:"disableCommentTag" json:"disableCommentTag"`     // 禁用注释放入tag标签中
	Package           string `yaml:"package" json:"package"`                         // 包名
	DisableDocComment bool   `yaml:"disable_doc_comment" json:"disable_doc_comment"` // 禁用文档注释
}

func InitFlagSetForConfig(s *pflag.FlagSet, cc *Config) {
	s.StringToStringVarP(&cc.Tags, "tags", "K", map[string]string{"json": utils.StyleSnakeCase}, "tags标签,类型支持[smallCamelCase,camelCase,snakeCase,kebab]")
	s.BoolVarP(&cc.EnableInt, "enableInt", "e", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	s.BoolVarP(&cc.EnableBoolInt, "enableBoolInt", "b", false, "使能bool输出int")
	s.BoolVarP(&cc.DisableNullToPoint, "disableNullToPoint", "B", false, "禁用字段为null时输出指针类型,将输出为sql.Nullxx")
	s.BoolVarP(&cc.DisableCommentTag, "disableCommentTag", "j", false, "禁用注释放入tag标签中")
	s.BoolVarP(&cc.EnableForeignKey, "enableForeignKey", "J", false, "使用外键")
	s.StringVar(&cc.Package, "package", "", "package name")
	s.BoolVarP(&cc.DisableDocComment, "disableDocComment", "d", false, "禁用文档注释")
}

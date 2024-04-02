package ens

import "github.com/things-go/ens/utils"

const (
	TagSmallCamelCase = "smallCamelCase"
	TagCamelCase      = "camelCase"
	TagSnakeCase      = "snakeCase"
	TagKebab          = "kebab"
)

func TagName(kind, name string) string {
	vv := ""
	switch kind {
	case TagSmallCamelCase:
		vv = utils.SmallCamelCase(name)
	case TagCamelCase:
		vv = utils.CamelCase(name)
	case TagSnakeCase:
		vv = utils.SnakeCase(name)
	case TagKebab:
		vv = utils.Kebab(name)
	}
	return vv
}

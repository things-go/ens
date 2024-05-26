package utils

const (
	StyleSmallCamelCase = "smallCamelCase"
	StylePascalCase     = "pascalCase"
	StyleSnakeCase      = "snakeCase"
	StyleKebab          = "kebab"
)

func StyleName(kind, name string) string {
	vv := name
	switch kind {
	case StyleSmallCamelCase:
		vv = SmallCamelCase(name)
	case StylePascalCase:
		vv = PascalCase(name)
	case StyleSnakeCase:
		vv = SnakeCase(name)
	case StyleKebab:
		vv = Kebab(name)
	}
	return vv
}

package utils

const (
	StyleSmallCamelCase = "smallCamelCase"
	StyleCamelCase      = "camelCase"
	StyleSnakeCase      = "snakeCase"
	StyleKebab          = "kebab"
)

func StyleName(kind, name string) string {
	vv := name
	switch kind {
	case StyleSmallCamelCase:
		vv = SmallCamelCase(name)
	case StyleCamelCase:
		vv = CamelCase(name)
	case StyleSnakeCase:
		vv = SnakeCase(name)
	case StyleKebab:
		vv = Kebab(name)
	}
	return vv
}

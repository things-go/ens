package matcher

import (
	"regexp"
	"strings"
)

var reJSONTag = regexp.MustCompile(`^.*\[@(?i:jsontag):\s*([^\[\]]*)\].*`)
var reAffixJSONTag = regexp.MustCompile(`^.*\[@(affix)\s*\].*`)
var reProtobufType = regexp.MustCompile(`^.*\[@(?i:pbtype):\s*([^\[\]]*)\].*`)
var reProtobufEnumValue = regexp.MustCompile(`@EnumValue\[[^\]]*\]`)

// JsonTag 匹配json标签
// [@jsontag:id,omitempty]
func JsonTag(comment string) string {
	match := reJSONTag.FindStringSubmatch(comment)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// HasAffixJSONTag 是否有 affix, 增加 json 标签 `,string`
// [@affix]
func HasAffixJSONTag(comment string) bool {
	match := reAffixJSONTag.FindStringSubmatch(comment)
	return len(match) == 2 && strings.TrimSpace(match[1]) == "affix"
}

// PbType 匹配pbtype类型值
// [@pbtype: E.Gender]
func PbType(comment string) string {
	match := reProtobufType.FindStringSubmatch(comment)
	if len(match) == 2 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// TrimEnumValue 去除 @EnumValue[xxxx] 的数据
func TrimEnumValue(value string) string {
	return strings.TrimSpace(reProtobufEnumValue.ReplaceAllString(value, ""))
}

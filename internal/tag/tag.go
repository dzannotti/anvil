package tag

import (
	"regexp"
	"strings"
)

type Tag struct {
	value string
}

var (
	whitespace  = regexp.MustCompile(`\s`)
	specialChar = regexp.MustCompile(`[@#$%\-^&*]`)
	nonASCII    = regexp.MustCompile(`[^\x00-\x7F]`)
	multipleDot = regexp.MustCompile(`\.+`)
	boundaryDot = regexp.MustCompile(`^\.|\.$`)
)

func FromString(val string) Tag {
	return Tag{value: normalize(val)}
}

func (t Tag) AsString() string {
	return t.value
}

func (t Tag) AsStrings() []string {
	return strings.Split(t.value, ".")
}

func (t Tag) MatchExact(other Tag) bool {
	return t.value == other.value
}

func (t Tag) Match(other Tag) bool {
	if t.value == other.value {
		return true
	}
	return strings.HasPrefix(t.value, other.value+".")
}

func normalize(value string) string {
	value = strings.ToLower(value)
	value = whitespace.ReplaceAllString(value, "")
	value = specialChar.ReplaceAllString(value, "")
	value = nonASCII.ReplaceAllString(value, "")
	value = multipleDot.ReplaceAllString(value, ".")
	value = boundaryDot.ReplaceAllString(value, "")
	return value
}

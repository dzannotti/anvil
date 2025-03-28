package tag

import (
	"regexp"
	"strings"
)

var (
	whitespace  = regexp.MustCompile(`\s`)
	specialChar = regexp.MustCompile(`[@#$%\-^&*]`)
	nonASCII    = regexp.MustCompile(`[^\x00-\x7F]`)
	multipleDot = regexp.MustCompile(`\.+`)
	boundaryDot = regexp.MustCompile(`^\.|\.$`)
)

type Tag struct {
	value string
}

func FromString(value string) Tag {
	return Tag{value: normalize(value)}
}

func (t Tag) MatchExact(other Tag) bool {
	return t.value == other.value
}

func (t Tag) Match(other Tag) bool {
	return strings.HasPrefix(t.value, other.value)
}

func (t Tag) String() string {
	return t.value
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

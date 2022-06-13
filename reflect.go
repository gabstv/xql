package xql

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	TableReplacer = "{t}"
)

func ExtractStructTags(v any, tag ...string) []string {
	rv := reflect.TypeOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}
	var tags []string
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		for _, t := range tag {
			if tv, tok := f.Tag.Lookup(t); tok && tv != "-" && !strings.HasPrefix(tv, "-,") {
				if tv == "" {
					tags = append(tags, ToSnakeCase(f.Name))
				} else {
					tags = append(tags, tv)
				}
			}
		}
	}
	return RemoveDuplicates(tags)
}

func Explode(sep string, v ...string) []string {
	var out []string
	for _, s := range v {
		x := strings.Split(s, sep)
		out = append(out, x...)
	}
	return out
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func RemoveDuplicates[T comparable](v []T) []T {
	var out []T
	m := make(map[T]bool)
	for _, x := range v {
		if _, ok := m[x]; !ok {
			out = append(out, x)
			m[x] = true
		}
	}
	return out
}

package db

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// FtsQuery takes a user-generated query string and formats as a
// query suitable for postgresql's ts_query().  It does its level best
// to sanitize user input.
//
// Basic Algorithm:
//  - Split query into fields where fields contain characters conforming
//    to unicode.isLetter() and unicode.isNumber().  All other
//    characters are treated as field separators and discarded.
//  - Fields which conform to [Prefix][Number] (case insensitive) are
//    striped of the prefix and any leading zeros and used as an
//    alternate to the field. ex: T01001 -> (T01001:* | 1001:*)
//  - All fields are treated as a prefix search by appending ':*'
//  - Fields are explicitly ANDed together
//  - ex: 'hsa01040 fat' -> '(hsa01040:* | 1040:*) & fat:*'
func FtsQuery(query, prefix string) string {
	var (
		s      string
		fields []string
	)

	if query == "" {
		return ""
	}

	fields = strings.FieldsFunc(query, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	for k, v := range fields {
		s = strings.TrimPrefix(strings.ToLower(v), strings.ToLower(prefix))
		s = strings.TrimLeft(s, "0")
		_, err := strconv.Atoi(s)
		if err == nil && s != v && len(s) > 0 {
			fields[k] = fmt.Sprintf("(%s:* | %s:*)", v, s)
		} else {
			fields[k] = fmt.Sprintf("%s:*", v)
		}
	}

	return strings.Join(fields, " & ")
}

package server

import (
	"html/template"
	"strings"
)

func funcMap() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
	}
}

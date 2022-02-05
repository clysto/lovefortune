package plugin

import "text/template"

type Plugin interface {
	Name() string
	Funcs() template.FuncMap
}

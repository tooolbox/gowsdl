package main

var typesTmpl = `
{{define "SimpleType"}}
	{{$name := .Name}}
	type {{.Name}} {{toGoType .Restriction.Base}}
	const (
		{{with .Restriction}}
			{{range .Enumeration}}
				{{replaceReservedWords .Value | makeFieldPublic}} {{$name}} = "{{.Value}}" {{end}}
		{{end}}
	)
{{end}}

{{define "ComplexType"}}
	{{$name := replaceReservedWords .Name}}

	type {{.Name}} struct {
		{{if ne .ComplexContent.Extension.Base ""}}
			{{$baseType := .ComplexContent.Extension.Base}}
			{{ if $baseType }}
				*{{stripns $baseType}}
			{{end}}

			{{template "Elements" .ComplexContent.Extension.Sequence.Elements}}
		{{ else }}
			{{template "Elements" .Sequence.Elements}}
		{{end}}
	}
{{end}}

{{define "Elements"}}
	{{range .}}
		{{replaceReservedWords .Name | makeFieldPublic}} {{if eq .MaxOccurs "unbounded"}}[]{{end}}{{.Type | toGoType}}{{end}}
{{end}}

package main
import (
	"encoding/xml"
	//"time"
	{{/*range .Imports*/}}
		{{/*.*/}}
	{{/*end*/}}
)

{{range .Schemas}}
	{{range .SimpleType}}
		{{template "SimpleType" .}}
	{{end}}
	{{range .Elements}}
		{{if not .Type}}

		{{end}}
	{{end}}
	{{range .ComplexTypes}}
		{{template "ComplexType" .}}
	{{end}}
{{end}}
`

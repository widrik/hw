package main

import (
	"bytes"
	"text/template"
)

type TemplateData struct {
	PackageName string
	Imports     []string
	Structs     []Struct
}

var Template = `/*
* Code generated automatically 
* THIS FILE SHOULD NOT BE EDITED.
*/
package {{.PackageName}}

{{if .Imports}}
import (
	{{range .Imports}}
	"{{.}}"{{end}}
)
{{end}}

type ValidationError struct {
	Field string
	Err string
}

{{range $structIdx, $struct := .Structs}}
func ({{$struct.VarName}} {{$struct.Name}}) Validate() ([]ValidationError, error) {
	var errs []ValidationError
	{{range $fieldIdx, $field := .Fields}}
		{{range $nameIdx, $name := $field.Names}}
			{{if eq $field.Type.Type "array"}}
	for i := range {{$struct.VarName}}.{{$name}} {
			{{end}}

			{{if eq $field.Type.VarType "string"}}
				{{range $validatorIdx, $validator := $field.Validators}}

					{{if eq $validator.Type "len"}}
	if len({{$struct.VarName}}.{{$name}}{{if eq $field.Type.Type "array"}}[i]{{end}}) != {{$validator.Value}} {
		errs = append(errs, ValidationError{
			Field: "{{$name}}",
						{{if eq $field.Type.Type "array"}}
			Err: fmt.Sprintf("field {{$name}} must be length is {{$validator.Value}} in index %d", i),
						{{else}}
			Err: "field {{$name}} must be length is {{$validator.Value}}",
						{{end}}
		})
	}
					{{end}}

					{{if eq $validator.Type "regexp"}}
	{
		matched, err := regexp.MatchString("{{$validator.Value}}", {{$struct.VarName}}.{{$name}}{{if eq $field.Type.Type "array"}}[i]{{end}})
		if err != nil {
			return errs, err
		}
		if matched == false {
			errs = append(errs, ValidationError{
				Field: "{{$name}}",
						{{if eq $field.Type.Type "array"}}
				Err: fmt.Sprintf("field {{$name}} must be in regexp {{$validator.Value}} in index %d", i),
						{{else}}
				Err: "field {{$name}} must be in regexp {{$validator.Value}}",
						{{end}}
			})
		}
	}
					{{end}}

					{{if eq $validator.Type "in"}}
	if !({{range $valueIdx, $value := $validator.Value}}{{if eq $valueIdx 0}}{{else}} ||{{end}} {{$struct.VarName}}.{{$name}}{{if eq $field.Type.Type "array"}}[i]{{end}} == "{{$value}}"{{end}}) {
		errs = append(errs, ValidationError{
			Field: "{{$name}}",
						{{if eq $field.Type.Type "array"}}
			Err: fmt.Sprintf("field {{$name}} must be in range {{$validator.Value}} in index %d", i),
						{{else}}
			Err: "field {{$name}} must be in range {{$validator.Value}}",
						{{end}}
		})
	}
					{{end}}
				{{end}}
			{{end}}

			{{if eq $field.Type.VarType "int"}}
				{{range $validatorIdx, $validator := $field.Validators}}

					{{if eq $validator.Type "min"}}
	if {{$struct.VarName}}.{{$name}}{{if eq $field.Type.Type "array"}}[i]{{end}} < {{$validator.Value}} {
		errs = append(errs, ValidationError{
			Field: "{{$name}}",
						{{if eq $field.Type.Type "array"}}
			Err: fmt.Sprintf("field {{$name}} must be min is {{$validator.Value}} in index %d", i),
						{{else}}
			Err: "field {{$name}} must be min is {{$validator.Value}}",
						{{end}}
		})
	}
					{{end}}

					{{if eq $validator.Type "max"}}
	if {{$struct.VarName}}.{{$name}}{{if eq $field.Type.Type "array"}}[i]{{end}} > {{$validator.Value}} {
		errs = append(errs, ValidationError{
			Field: "{{$name}}",
						{{if eq $field.Type.Type "array"}}
			Err: fmt.Sprintf("field {{$name}} must be max is {{$validator.Value}} in index %d", i),
						{{else}}
			Err: "field {{$name}} must be max is {{$validator.Value}}",
						{{end}}
		})
	}
					{{end}}

					{{if eq $validator.Type "in"}}
	if !({{range $valueIdx, $value := $validator.Value}}{{if eq $valueIdx 0}}{{else}} ||{{end}} {{$struct.VarName}}.{{$name}}{{if eq $field.Type.Type "array"}}[i]{{end}} == {{$value}}{{end}}) {
		errs = append(errs, ValidationError{
			Field: "{{$name}}",
						{{if eq $field.Type.Type "array"}}
			Err: fmt.Sprintf("field {{$name}} must be in range {{$validator.Value}} in index %d", i),
						{{else}}
			Err: "field {{$name}} must be in range {{$validator.Value}}",
						{{end}}
		})
	}
					{{end}}
				{{end}}
			{{end}}
			{{if eq $field.Type.Type "array"}}
	}
			{{end}}
		{{end}}
	{{end}}
	return errs, nil
}
{{end}}
`

func CreateTemplate(file GeneratatedFile) (bytes.Buffer, error) {
	newTemplate := template.Must(template.New("structValidatorTemplate").Parse(Template))

	arrayExists := false
	regexpExists := false
	for _, structs := range file.Structs {
		for _, structField := range structs.Fields {
			if structField.Type.Type == "array" {
				arrayExists = true
			}

			for _, validator := range structField.Validators {
				if validator.Type == "regexp" {
					regexpExists = true
				}
			}
		}
	}

	importsList := []string{}
	if arrayExists {
		importsList = append(importsList, "fmt")
	}
	if regexpExists {
		importsList = append(importsList, "regexp")
	}

	templateData := TemplateData{
		PackageName: file.Name,
		Imports:     importsList,
		Structs:     file.Structs,
	}

	var tplBuffer bytes.Buffer
	err := newTemplate.Execute(&tplBuffer, templateData)
	if err != nil {
		return tplBuffer, err
	}

	return tplBuffer, nil
}

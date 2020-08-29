package main

import (
	"errors"
	"go/ast"
	"regexp"
	"strconv"
	"strings"
)

var rgxp = regexp.MustCompile(`[\s]?validate:"([^\"]+)"[\s]?`)

const (
	min    string = "min"
	max    string = "max"
	regex  string = "regexp"
	in     string = "in"
	length string = "len"
)

func CreateStructsData(astFile ast.File) (GeneratatedFile, error) {
	generatedFile := GeneratatedFile{}

	if astFile.Decls == nil {
		return generatedFile, errors.New("astFile error")
	}

	if astFile.Name == nil {
		return generatedFile, errors.New("astFile error")
	}

	generatedFile.Name = astFile.Name.Name
	for _, d := range astFile.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range g.Specs {
			astTypeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			st, err := convertDataToStruct(*astTypeSpec)
			if err != nil {
				return generatedFile, nil
			}
			if st != nil {
				generatedFile.Structs = append(generatedFile.Structs, *st)
			}
		}
	}

	return generatedFile, nil
}

func convertDataToStruct(typeSpec ast.TypeSpec) (*Struct, error) {
	resStruct := Struct{}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok || typeSpec.Name == nil {
		return nil, nil
	}
	resStruct.Name = typeSpec.Name.Name
	resStruct.VarName = strings.ToLower(strings.Split(typeSpec.Name.Name, "")[0])

	fields := structType.Fields
	if fields == nil {
		return nil, nil
	}

	fieldList := fields.List
	if fieldList == nil {
		return nil, nil
	}

	structFields := []Field{}
	for _, field := range fieldList {
		if field.Tag == nil || field.Names == nil {
			continue
		}

		fieldType := prepareFieldType(field.Type)
		if fieldType == nil {
			continue
		}

		resNames := []string{}
		for _, name := range field.Names {
			if name == nil {
				continue
			}

			resNames = append(resNames, name.Name)
		}
		if len(resNames) == 0 {
			continue
		}

		validators, err := createValidators(fieldType.VarType, field.Tag.Value)
		if err != nil {
			return nil, err
		}
		if len(validators) == 0 {
			continue
		}

		newField := Field{
			Names:      resNames,
			Type:       *fieldType,
			Validators: validators,
		}
		structFields = append(structFields, newField)
	}
	resStruct.Fields = structFields

	return &resStruct, nil
}

func prepareFieldType(expression ast.Expr) *FieldType {
	switch expressionType := expression.(type) {
	case *ast.Ident:
		return &FieldType{
			VarType: expressionType.Name,
			Type:    "var",
		}
	case *ast.ArrayType:
		ident, ok := expressionType.Elt.(*ast.Ident)
		if !ok {
			return nil
		}
		return &FieldType{
			VarType: ident.Name,
			Type:    "array",
		}
	}

	return nil
}

func createValidators(fieldType string, tag string) ([]Validator, error) {
	emptyValidators := []Validator{}
	validators := []Validator{}

	params := rgxp.FindStringSubmatch(tag)
	if len(params) == 0 || len(params[1]) == 0 {
		return emptyValidators, nil
	}

	for _, validString := range strings.Split(params[1], "|") {
		validator, err := createValidator(fieldType, validString)
		if err != nil {
			return emptyValidators, err
		}
		validators = append(validators, validator)
	}

	return validators, nil
}

func createValidator(fieldType string, validatorString string) (Validator, error) {
	var err error
	validator := Validator{}

	params := strings.SplitN(validatorString, ":", 2)
	if len(params) != 2 {
		return validator, errors.New("not valid string ")
	}

	value := params[1]

	switch params[0] {
	case "in":
		validator.Type = in
		values := strings.Split(value, ",")

		if fieldType == "int" {
			var (
				intsValues []int
				num        int
			)
			for _, v := range values {
				num, err = strconv.Atoi(v)
				if err != nil {
					return validator, err
				}
				intsValues = append(intsValues, num)
			}

			validator.Value = intsValues
		} else {
			validator.Value = values
		}
	case "min":
		validator.Type = min
		validator.Value, err = strconv.Atoi(value)
		if err != nil {
			return validator, err
		}
	case "max":
		validator.Type = max
		validator.Value, err = strconv.Atoi(value)
		if err != nil {
			return validator, err
		}
	case "len":
		validator.Type = length
		validator.Value, err = strconv.Atoi(value)
		if err != nil {
			return validator, err
		}
	case "regexp":
		validator.Type = regex
		validator.Value = value
	default:
		return validator, errors.New("not valid validator type ")
	}

	return validator, nil
}

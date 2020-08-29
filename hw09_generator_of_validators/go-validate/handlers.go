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
	availableTypes := getAvailableTypes(astFile.Decls)

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

			st, err := convertDataToStruct(*astTypeSpec, availableTypes)
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

func convertDataToStruct(typeSpec ast.TypeSpec, types map[string]string) (*Struct, error) {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok || typeSpec.Name == nil {
		return nil, nil
	}
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
		detectedType, ok := types[fieldType.VarType]
		if !ok {
			continue
		}
		fieldType.VarType = detectedType

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
		structFields = append(structFields, Field{
			Names:      resNames,
			Type:       *fieldType,
			Validators: validators,
		})
	}
	if len(structFields) == 0 {
		return nil, nil
	}

	return &Struct{
		Name:    typeSpec.Name.Name,
		VarName: strings.ToLower(typeSpec.Name.Name),
		Fields:  structFields,
	}, nil
}

func getAvailableTypes(decls []ast.Decl) map[string]string {
	baseTypes := map[string]string{
		"int":    "int",
		"int64":  "int64",
		"string": "string",
	}

	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			typeIdent, ok := typeSpec.Type.(*ast.Ident)
			if !ok {
				continue
			}

			if typeSpec.Name == nil {
				continue
			}

			if _, ok := baseTypes[typeSpec.Name.Name]; !ok {
				baseTypes[typeSpec.Name.Name] = typeIdent.Name
			}
		}
	}

	for typeIndex, base := range baseTypes {
		newType, ok := baseTypes[base]
		if !ok {
			continue
		}

		baseTypes[typeIndex] = newType
	}

	return baseTypes
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
